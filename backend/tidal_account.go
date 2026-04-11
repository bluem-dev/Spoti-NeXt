package backend

// TidalAccount implements a direct Tidal API login using the OAuth2 device
// authorization flow (the same flow used by tidal-dl-ng). Once a session is
// established the account can be used to fetch a real playbackinfopostpaywall
// manifest — the only reliable way to get DOLBY_ATMOS content.
//
// The account credentials (access_token + refresh_token) are persisted inside
// the app settings so the user only needs to link the device once.

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// ── OAuth constants ────────────────────────────────────────────────────────
//
// Two client IDs are needed:
//   - tidalLoginClientID  — supports the Device Authorization Flow (limited_input_device).
//                           Used for login + token refresh only.
//   - tidalAtmosClientID  — supports Atmos streams in playbackinfopostpaywall.
//                           Used only as X-Tidal-Token in manifest requests.
//
// The access token is generated with tidalLoginClientID; the X-Tidal-Token
// header on the manifest request is what gates Atmos availability.

const (
	tidalLoginClientID = "fX2JxdmntZWK0ixT"
	tidalLoginSecret   = "1Nn9AfDAjxrgJFJbKNWLeAyKGVGmINuXPPLHVXAvxAg="

	// TV client — used to re-authenticate with an existing refresh_token to
	// obtain an access_token scoped for Dolby Atmos streams (E-AC3/EAC-3).
	// The access_token itself must be issued by this client; a header-only
	// approach (X-Tidal-Token) is insufficient.
	// Source: orpheusdl-tidal — tv_atmos_token / tv_atmos_secret
	tidalTVAtmosClientID = "cgiF7TQuB97BUIu3"
	tidalTVAtmosSecret   = "1nqpgx8uvBdZigrx4hUPDV2hOwgYAAAG5DYXOr6uNf8="

	tidalDeviceAuthURL = "https://auth.tidal.com/v1/oauth2/device_authorization"
	tidalTokenURL      = "https://auth.tidal.com/v1/oauth2/token"
	tidalAPIBase       = "https://api.tidal.com/v1/"
)

// ── Persisted session ──────────────────────────────────────────────────────

// TidalAccountSession is the data that gets saved in settings.
type TidalAccountSession struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CountryCode  string    `json:"country_code"`
	UserID       int64     `json:"user_id"`
	ClientID     string    `json:"client_id,omitempty"` // tracks which client generated this token
}

// IsCompatible returns true if the session was generated with the current
// tidalLoginClientID. Tokens are bound to the client that issued them — a
// session created with a different client ID must be discarded and re-created.
func (s *TidalAccountSession) IsCompatible() bool {
	// Sessions saved before this field existed have ClientID == "".
	// Treat them as incompatible so the user is prompted to re-login once.
	return s != nil && s.ClientID == tidalLoginClientID
}

func (s *TidalAccountSession) Valid() bool {
	return s != nil && s.AccessToken != "" && time.Now().Before(s.ExpiresAt)
}

// ── Device-link login state ────────────────────────────────────────────────

// TidalDeviceLink is returned to the frontend so the user can open the URL.
type TidalDeviceLink struct {
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	DeviceCode              string `json:"device_code"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// TidalAccountManager holds a live (or pending) session and serializes access.
type TidalAccountManager struct {
	mu      sync.Mutex
	session *TidalAccountSession
	client  *http.Client
}

var globalTidalAccount = &TidalAccountManager{
	client: &http.Client{Timeout: 15 * time.Second},
}

// SetSession replaces the stored session (called from app.go on settings load
// and after a successful login).
func (m *TidalAccountManager) SetSession(s *TidalAccountSession) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.session = s
}

// GetSession returns a copy of the current session, refreshing if needed.
func (m *TidalAccountManager) GetSession() (*TidalAccountSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil || m.session.AccessToken == "" {
		log.Println("⚠ Tidal account: no session in memory")
		return nil, fmt.Errorf("no Tidal account session — please log in first")
	}

	// Refresh if expiring within 5 minutes.
	remaining := time.Until(m.session.ExpiresAt)
	log.Printf("Tidal account session: token present, expires in %v (at %s)", remaining.Round(time.Second), m.session.ExpiresAt.Format(time.RFC3339))
	if time.Now().Add(5 * time.Minute).After(m.session.ExpiresAt) {
		log.Println("Tidal account token expiring, refreshing...")
		if err := m.refreshLocked(); err != nil {
			log.Printf("⚠ Tidal account refresh failed: %v", err)
			return nil, fmt.Errorf("token refresh failed: %w", err)
		}
		log.Println("✓ Tidal account token refreshed successfully")
	}

	copy := *m.session
	return &copy, nil
}

func (m *TidalAccountManager) refreshLocked() error {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", m.session.RefreshToken)
	data.Set("client_id", tidalLoginClientID)
	data.Set("client_secret", tidalLoginSecret)

	req, err := http.NewRequest("POST", tidalTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("refresh HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("refresh parse error: %w", err)
	}

	m.session.AccessToken = result.AccessToken
	m.session.ExpiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return nil
}

// ── Device-link flow ───────────────────────────────────────────────────────

// StartDeviceLogin initiates the device authorization flow. Returns the link
// the user must visit to approve the login.
func StartDeviceLogin() (*TidalDeviceLink, error) {
	data := url.Values{}
	data.Set("client_id", tidalLoginClientID)
	data.Set("scope", "r_usr w_usr w_sub")

	req, err := http.NewRequest("POST", tidalDeviceAuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("device auth request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("device auth HTTP %d: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		DeviceCode              string `json:"deviceCode"`
		UserCode                string `json:"userCode"`
		VerificationURI         string `json:"verificationUri"`
		VerificationURIComplete string `json:"verificationUriComplete"`
		ExpiresIn               int    `json:"expiresIn"`
		Interval                int    `json:"interval"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("device auth parse error: %w", err)
	}

	return &TidalDeviceLink{
		VerificationURI:         raw.VerificationURI,
		VerificationURIComplete: raw.VerificationURIComplete,
		DeviceCode:              raw.DeviceCode,
		ExpiresIn:               raw.ExpiresIn,
		Interval:                raw.Interval,
	}, nil
}

// PollDeviceLogin polls the token endpoint until the user approves or the
// code expires. Returns the completed session on success.
func PollDeviceLogin(deviceCode string, interval, expiresIn int) (*TidalAccountSession, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)
	pollInterval := time.Duration(interval) * time.Second
	if pollInterval < 2*time.Second {
		pollInterval = 5 * time.Second
	}

	for time.Now().Before(deadline) {
		time.Sleep(pollInterval)

		data := url.Values{}
		data.Set("client_id", tidalLoginClientID)
		data.Set("client_secret", tidalLoginSecret)
		data.Set("device_code", deviceCode)
		data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		data.Set("scope", "r_usr w_usr w_sub")

		req, err := http.NewRequest("POST", tidalTokenURL, strings.NewReader(data.Encode()))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 400 {
			// authorization_pending is normal — keep polling
			var errBody struct {
				Error string `json:"error"`
			}
			json.Unmarshal(body, &errBody)
			if errBody.Error == "authorization_pending" || errBody.Error == "authorization_declined" {
				continue
			}
		}

		if resp.StatusCode != 200 {
			continue
		}

		var token struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int    `json:"expires_in"`
			TokenType    string `json:"token_type"`
		}
		if err := json.Unmarshal(body, &token); err != nil || token.AccessToken == "" {
			continue
		}

		// Fetch user info for country code.
		session := &TidalAccountSession{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
			ClientID:     tidalLoginClientID,
		}

		if country, userID, err := fetchTidalUserInfo(session.AccessToken, client); err == nil {
			session.CountryCode = country
			session.UserID = userID
		} else {
			session.CountryCode = "US"
		}

		globalTidalAccount.SetSession(session)
		return session, nil
	}

	return nil, fmt.Errorf("device login timed out — code expired")
}

func fetchTidalUserInfo(accessToken string, client *http.Client) (string, int64, error) {
	req, err := http.NewRequest("GET", tidalAPIBase+"sessions", nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Tidal-Token", tidalLoginClientID)

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var info struct {
		SessionID   string `json:"sessionId"`
		CountryCode string `json:"countryCode"`
		UserID      int64  `json:"userId"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", 0, err
	}
	return info.CountryCode, info.UserID, nil
}

// ── TV-client re-authentication ────────────────────────────────────────────

// refreshToTVSession exchanges an existing refresh_token for a new access_token
// issued by the TV Atmos client (tidalTVAtmosClientID). The resulting token is
// scoped for Dolby Atmos playback and should be used only for manifest requests.
// It does not replace the main session — the main session keeps its own tokens.
//
// Per orpheusdl: "refresh tokens can be used with any client id — this can be
// used to switch to any client type from an existing session."
func refreshToTVSession(refreshToken string, client *http.Client) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", tidalTVAtmosClientID)
	data.Set("client_secret", tidalTVAtmosSecret)

	req, err := http.NewRequest("POST", tidalTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "TIDAL_ANDROID/1039 okhttp/3.14.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("TV session request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("TV session HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil || result.AccessToken == "" {
		return "", fmt.Errorf("TV session parse error: %w", err)
	}
	return result.AccessToken, nil
}

// ── Direct playback manifest fetch ────────────────────────────────────────

// TidalAccountGetManifest fetches a playbackinfopostpaywall manifest directly
// using the account's access token, bypassing the third-party proxies.
// Returns a "MANIFEST:<base64>" string on success — same format as the proxy
// path so it plugs directly into DownloadFromManifest.
func TidalAccountGetManifest(trackID int64, quality string) (string, string, error) {
	session, err := globalTidalAccount.GetSession()
	if err != nil {
		return "", "", err
	}

	countryCode := session.CountryCode
	if countryCode == "" {
		countryCode = "US"
	}

	DebugLog(fmt.Sprintf("[atmos-manifest] trackID=%d quality=%s countryCode=%s", trackID, quality, countryCode))
	DebugLog(fmt.Sprintf("[atmos-manifest] session clientID=%s userID=%d tokenLen=%d refreshTokenLen=%d",
		session.ClientID, session.UserID, len(session.AccessToken), len(session.RefreshToken)))

	client := &http.Client{Timeout: 20 * time.Second}

	// Obtain bearer token — for Atmos, re-auth with TV client for the correct scope.
	bearerToken := session.AccessToken
	bearerSource := "main-session"
	if quality == "DOLBY_ATMOS" {
		DebugLog(fmt.Sprintf("[atmos-tv] attempting TV re-auth with clientID=%s", tidalTVAtmosClientID))
		tvToken, tvErr := refreshToTVSession(session.RefreshToken, client)
		if tvErr != nil {
			DebugLog(fmt.Sprintf("[atmos-tv] ✗ TV re-auth failed: %v — falling back to main token", tvErr))
		} else {
			DebugLog(fmt.Sprintf("[atmos-tv] ✓ TV token obtained tokenLen=%d", len(tvToken)))
			bearerToken = tvToken
			bearerSource = "tv-session"
		}
	}

	setHeaders := func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+bearerToken)
		r.Header.Set("X-Tidal-Token", tidalTVAtmosClientID)
		r.Header.Set("User-Agent", "TIDAL_ANDROID/1039 okhttp/3.14.9")
		r.Header.Set("Connection", "Keep-Alive")
	}

	// Build the manifest endpoint URL.
	// audioMode=DOLBY_ATMOS is gated server-side by account/subscription —
	// we request HI_RES and the API upgrades to Atmos if the account has access.
	audioQuality := quality
	if quality == "DOLBY_ATMOS" {
		audioQuality = "HI_RES"
	}
	type candidate struct{ label, url string }
	candidates := []candidate{
		{"v4", fmt.Sprintf("%stracks/%d/playbackinfopostpaywall/v4?playbackmode=STREAM&assetpresentation=FULL&audioquality=%s&prefetch=false&countryCode=%s",
			tidalAPIBase, trackID, audioQuality, countryCode)},
	}

	apiURL := candidates[0].url
	DebugLog(fmt.Sprintf("[atmos-manifest] request bearer-from=%s url=%s", bearerSource, apiURL))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", "", err
	}
	setHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("Tidal API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	DebugLog(fmt.Sprintf("[atmos-manifest] HTTP=%d bodyLen=%d", resp.StatusCode, len(body)))
	if resp.StatusCode != 200 {
		DebugLog(fmt.Sprintf("[atmos-manifest] error body: %s", string(body)))
		return "", "", fmt.Errorf("Tidal API HTTP %d: %s", resp.StatusCode, string(body))
	}

	var data struct {
		TrackID          int64  `json:"trackId"`
		AudioMode        string `json:"audioMode"`
		AudioQuality     string `json:"audioQuality"`
		ManifestMimeType string `json:"manifestMimeType"`
		Manifest         string `json:"manifest"`
		BitDepth         int    `json:"bitDepth"`
		SampleRate       int    `json:"sampleRate"`
	}
	if err := json.Unmarshal(body, &data); err != nil || data.Manifest == "" {
		return "", "", fmt.Errorf("unexpected response from Tidal API: %s", string(body))
	}

	DebugLog(fmt.Sprintf("[atmos-manifest] result audioMode=%s audioQuality=%s bitDepth=%d sampleRate=%d",
		data.AudioMode, data.AudioQuality, data.BitDepth, data.SampleRate))

	if quality == "DOLBY_ATMOS" && data.AudioMode != "DOLBY_ATMOS" {
		DebugLog("[atmos-manifest] ✗ Atmos not granted — account or subscription limitation")
		return "", "", fmt.Errorf("DOLBY_ATMOS unavailable — API returned audioMode=%s audioQuality=%s", data.AudioMode, data.AudioQuality)
	}

	return "MANIFEST:" + data.Manifest, data.AudioMode, nil
}

// SetGlobalTidalSession replaces the stored session. Called from app.go on
// startup (restoring persisted credentials) and on logout (nil).
func SetGlobalTidalSession(s *TidalAccountSession) {
	globalTidalAccount.SetSession(s)
}

// CurrentTidalClientID returns the login client ID compiled into this build.
// Used by app.go to log incompatibility messages on session restore.
func CurrentTidalClientID() string {
	return tidalLoginClientID
}

// TidalAccountAvailable returns true if a valid session is loaded.
func TidalAccountAvailable() bool {
	_, err := globalTidalAccount.GetSession()
	return err == nil
}
