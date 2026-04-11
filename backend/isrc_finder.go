package backend

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	spotifyServerTimeURL    = "https://open.spotify.com/api/server-time"
	spotifySessionTokenURL  = "https://open.spotify.com/api/token"
	spotifyTOTPSecretsURL   = "https://git.gay/thereallo/totp-secrets/raw/branch/main/secrets/secretDict.json"
	spotifyGIDMetadataURL   = "https://spclient.wg.spotify.com/metadata/4/%s/%s?market=from_token"
	spotifyTOTPPeriod       = 30
	spotifyTOTPDigits       = 6
	spotifyBase62Alphabet   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	spotifyTokenCacheFile   = ".isrc-finder-token.json"
	spotifySecretsCacheFile = "spotify-secret-dict-cache.json"
	// spotifySecretsCacheTTL controls how long the TOTP secret dictionary is
	// cached locally before a refresh is attempted. 72 hours reduces network
	// traffic while still picking up rotations within a few days. The cache
	// falls back to stale data on network failure, so the effective resilience
	// window is longer than this value.
	spotifySecretsCacheTTL  = 72 * time.Hour
)

var spotifyAnonymousTokenMu sync.Mutex

type spotifyAnonymousToken struct {
	AccessToken                      string `json:"accessToken"`
	AccessTokenExpirationTimestampMs int64  `json:"accessTokenExpirationTimestampMs"`
}

type spotifyServerTimeResponse struct {
	ServerTime int64 `json:"serverTime"`
}

type spotifySecretsCache struct {
	FetchedAtUnix int64            `json:"fetched_at_unix"`
	Secrets       map[string][]int `json:"secrets"`
}

type spotifyTrackRawData struct {
	ExternalID []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"external_id"`
}

type spotFetchISRCResponse struct {
	Input        string   `json:"input"`
	TrackID      string   `json:"track_id"`
	GID          string   `json:"gid"`
	CanonicalURI string   `json:"canonical_uri"`
	Name         string   `json:"name"`
	Artists      []string `json:"artists"`
	AlbumName    string   `json:"album_name"`
	ReleaseDate  string   `json:"release_date"`
	Label        string   `json:"label"`
	ISRC         string   `json:"isrc"`
}

func (s *SongLinkClient) lookupSpotifyISRC(spotifyTrackID string) (string, error) {
	normalizedTrackID, err := extractSpotifyTrackID(spotifyTrackID)
	if err != nil {
		return "", err
	}

	cachedISRC, err := GetCachedISRC(normalizedTrackID)
	if err != nil {
		fmt.Printf("[ISRC] warning: failed to read ISRC cache: %v\n", err)
	} else if cachedISRC != "" {
		fmt.Printf("[ISRC] cache hit: %s\n", cachedISRC)
		return cachedISRC, nil
	}

	useSpotFetchAPI, spotFetchAPIURL := GetSpotFetchAPISettings()
	if useSpotFetchAPI {
		isrc, resolvedTrackID, err := s.lookupSpotifyISRCViaSpotFetchAPI(normalizedTrackID, spotFetchAPIURL)
		if err == nil && isrc != "" {
			fmt.Printf("[ISRC] found via SpotFetch API: %s\n", isrc)
			cacheResolvedSpotifyTrackISRC(normalizedTrackID, resolvedTrackID, isrc)
			return isrc, nil
		}
		if err != nil {
			fmt.Printf("[ISRC] warning: SpotFetch lookup failed, trying Spotify metadata: %v\n", err)
		}
	}

	payload, metadataErr := fetchSpotifyTrackRawData(s.client, normalizedTrackID)
	if metadataErr == nil {
		isrc, extractErr := extractSpotifyTrackISRC(payload)
		if extractErr == nil {
			fmt.Printf("[ISRC] found via Spotify metadata: %s\n", isrc)
			cacheResolvedSpotifyTrackISRC(normalizedTrackID, "", isrc)
			return isrc, nil
		}
		metadataErr = extractErr
	}

	if metadataErr != nil {
		fmt.Printf("[ISRC] warning: Spotify metadata failed, trying Soundplate: %v\n", metadataErr)
	}

	isrc, resolvedTrackID, soundplateErr := s.lookupSpotifyISRCViaSoundplate(normalizedTrackID)
	if soundplateErr == nil && isrc != "" {
		fmt.Printf("[ISRC] found via Soundplate: %s\n", isrc)
		cacheResolvedSpotifyTrackISRC(normalizedTrackID, resolvedTrackID, isrc)
		return isrc, nil
	}

	if metadataErr != nil && soundplateErr != nil {
		return "", fmt.Errorf("spotify metadata lookup failed: %v | soundplate lookup failed: %w", metadataErr, soundplateErr)
	}
	if soundplateErr != nil {
		return "", soundplateErr
	}
	return "", metadataErr
}

func cacheResolvedSpotifyTrackISRC(trackID string, resolvedTrackID string, isrc string) {
	if err := PutCachedISRC(trackID, isrc); err != nil {
		fmt.Printf("[ISRC] warning: failed to write ISRC cache: %v\n", err)
	}
	if resolvedTrackID != "" && resolvedTrackID != trackID {
		if err := PutCachedISRC(resolvedTrackID, isrc); err != nil {
			fmt.Printf("[ISRC] warning: failed to write ISRC cache for resolved ID: %v\n", err)
		}
	}
}

func (s *SongLinkClient) lookupSpotifyISRCViaSpotFetchAPI(spotifyTrackID string, apiBaseURL string) (string, string, error) {
	normalizedTrackID := strings.TrimSpace(spotifyTrackID)
	baseURL := strings.TrimRight(strings.TrimSpace(apiBaseURL), "/")
	if normalizedTrackID == "" {
		return "", "", fmt.Errorf("spotify track ID is required")
	}
	if baseURL == "" {
		return "", "", fmt.Errorf("spotfetch api url is required")
	}

	requestURL := fmt.Sprintf("%s/isrc/%s", baseURL, url.PathEscape(normalizedTrackID))
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create SpotFetch ISRC request: %w", err)
	}
	req.Header.Set("User-Agent", songLinkUserAgent)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("SpotFetch ISRC request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyPreview, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return "", "", fmt.Errorf("SpotFetch ISRC returned status %d (%s)", resp.StatusCode, strings.TrimSpace(string(bodyPreview)))
	}

	var payload spotFetchISRCResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", "", fmt.Errorf("failed to decode SpotFetch ISRC response: %w", err)
	}

	isrc := firstISRCMatch(payload.ISRC)
	if isrc == "" {
		return "", "", fmt.Errorf("ISRC missing in SpotFetch response")
	}

	return isrc, strings.TrimSpace(payload.TrackID), nil
}

func requestSpotifyBytes(client *http.Client, targetURL string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		details := strings.TrimSpace(string(body))
		if details == "" {
			details = resp.Status
		}
		return nil, fmt.Errorf("request failed: %s", details)
	}

	return body, nil
}

func requestSpotifyJSON(client *http.Client, targetURL string, headers map[string]string, target interface{}) error {
	body, err := requestSpotifyBytes(client, targetURL, headers)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return nil
}

func loadSpotifyCachedToken() (*spotifyAnonymousToken, error) {
	cachePath, err := spotifyTokenCachePath()
	if err != nil {
		return nil, err
	}

	body, err := os.ReadFile(cachePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read token cache: %w", err)
	}

	var token spotifyAnonymousToken
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("failed to read token cache: %w", err)
	}

	return &token, nil
}

func saveSpotifyCachedToken(token *spotifyAnonymousToken) error {
	cachePath, err := spotifyTokenCachePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return fmt.Errorf("failed to create token cache directory: %w", err)
	}

	body, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(cachePath, body, 0o644); err != nil {
		return fmt.Errorf("failed to write token cache: %w", err)
	}

	return nil
}

func loadSpotifyCachedSecrets() (*spotifySecretsCache, error) {
	cachePath, err := spotifySecretsCachePath()
	if err != nil {
		return nil, err
	}

	body, err := os.ReadFile(cachePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read secrets cache: %w", err)
	}

	var cache spotifySecretsCache
	if err := json.Unmarshal(body, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse secrets cache: %w", err)
	}

	return &cache, nil
}

func saveSpotifyCachedSecrets(cache *spotifySecretsCache) error {
	cachePath, err := spotifySecretsCachePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return fmt.Errorf("failed to create secrets cache directory: %w", err)
	}

	body, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(cachePath, body, 0o644); err != nil {
		return fmt.Errorf("failed to write secrets cache: %w", err)
	}

	return nil
}

func spotifyTokenCachePath() (string, error) {
	appDir, err := EnsureAppDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, spotifyTokenCacheFile), nil
}

func spotifySecretsCachePath() (string, error) {
	appDir, err := EnsureAppDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, spotifySecretsCacheFile), nil
}

func spotifyTokenIsValid(token *spotifyAnonymousToken) bool {
	if token == nil || token.AccessToken == "" || token.AccessTokenExpirationTimestampMs == 0 {
		return false
	}

	return time.Now().UnixMilli() < token.AccessTokenExpirationTimestampMs-30_000
}

func spotifySecretsCacheIsValid(cache *spotifySecretsCache) bool {
	if cache == nil || cache.FetchedAtUnix == 0 || len(cache.Secrets) == 0 {
		return false
	}

	return time.Since(time.Unix(cache.FetchedAtUnix, 0)) < spotifySecretsCacheTTL
}

func deriveSpotifyTOTPSecret(ciphertext []int) []byte {
	var builder strings.Builder

	for index, value := range ciphertext {
		builder.WriteString(strconv.Itoa(value ^ ((index % 33) + 9)))
	}

	return []byte(builder.String())
}

func generateSpotifyTOTP(secret []byte, timestampMs int64) string {
	counter := timestampMs / 1000 / spotifyTOTPPeriod
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, uint64(counter))

	mac := hmac.New(sha1.New, secret)
	mac.Write(counterBytes)
	digest := mac.Sum(nil)

	offset := digest[len(digest)-1] & 0x0f
	binaryCode := (int(digest[offset])&0x7f)<<24 |
		(int(digest[offset+1])&0xff)<<16 |
		(int(digest[offset+2])&0xff)<<8 |
		(int(digest[offset+3]) & 0xff)

	modulo := 1
	for i := 0; i < spotifyTOTPDigits; i++ {
		modulo *= 10
	}

	return fmt.Sprintf("%0*d", spotifyTOTPDigits, binaryCode%modulo)
}

func requestSpotifyAnonymousAccessToken(client *http.Client) (string, error) {
	spotifyAnonymousTokenMu.Lock()
	defer spotifyAnonymousTokenMu.Unlock()

	cachedToken, err := loadSpotifyCachedToken()
	if err != nil {
		return "", err
	}

	if spotifyTokenIsValid(cachedToken) {
		return cachedToken.AccessToken, nil
	}

	var serverTime spotifyServerTimeResponse
	if err := requestSpotifyJSON(client, spotifyServerTimeURL, nil, &serverTime); err != nil {
		return "", err
	}

	var secrets map[string][]int
	cachedSecrets, err := loadSpotifyCachedSecrets()
	if err != nil {
		fmt.Printf("[TOTP] warning: failed to read secrets cache: %v\n", err)
	}

	if spotifySecretsCacheIsValid(cachedSecrets) {
		secrets = cachedSecrets.Secrets
	} else {
		if err := requestSpotifyJSON(client, spotifyTOTPSecretsURL, nil, &secrets); err != nil {
			if cachedSecrets != nil && len(cachedSecrets.Secrets) > 0 {
				fmt.Printf("[TOTP] warning: failed to refresh secrets cache, using stale: %v\n", err)
				secrets = cachedSecrets.Secrets
			} else {
				return "", err
			}
		} else {
			cache := &spotifySecretsCache{
				FetchedAtUnix: time.Now().Unix(),
				Secrets:       secrets,
			}
			if err := saveSpotifyCachedSecrets(cache); err != nil {
				fmt.Printf("[TOTP] warning: failed to write secrets cache: %v\n", err)
			}
		}
	}

	version, err := latestSpotifySecretVersion(secrets)
	if err != nil {
		return "", err
	}

	secret := deriveSpotifyTOTPSecret(secrets[version])
	generatedTOTP := generateSpotifyTOTP(secret, serverTime.ServerTime*1000)

	query := url.Values{
		"reason":      {"init"},
		"productType": {"web-player"},
		"totp":        {generatedTOTP},
		"totpServer":  {generatedTOTP},
		"totpVer":     {version},
	}

	var token spotifyAnonymousToken
	if err := requestSpotifyJSON(client, spotifySessionTokenURL+"?"+query.Encode(), nil, &token); err != nil {
		return "", err
	}

	if err := saveSpotifyCachedToken(&token); err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func latestSpotifySecretVersion(secrets map[string][]int) (string, error) {
	var (
		bestVersion string
		bestNumber  int
	)

	for version := range secrets {
		number, err := strconv.Atoi(version)
		if err != nil {
			return "", fmt.Errorf("invalid secret version %q: %w", version, err)
		}
		if bestVersion == "" || number > bestNumber {
			bestVersion = version
			bestNumber = number
		}
	}

	if bestVersion == "" {
		return "", errors.New("no TOTP secret versions available")
	}

	return bestVersion, nil
}

// spotifyContentTypes lists the known Spotify content path segments.
// Any leading path segment that is not one of these is treated as a locale
// prefix (e.g. "intl-es", "es-419", "en-GB") and skipped automatically.
var spotifyContentTypes = map[string]bool{
	"track":    true,
	"album":    true,
	"playlist": true,
	"artist":   true,
}

func extractSpotifyTrackID(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("track input is required")
	}

	if strings.HasPrefix(value, "spotify:track:") {
		return value[strings.LastIndex(value, ":")+1:], nil
	}

	parsed, err := url.Parse(value)
	if err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") {
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		// Spotify sometimes prefixes the path with a locale segment such as
		// "intl-es", "es-419", or "en-GB". Skip it if it is not a known
		// content type so that locale-tagged URLs resolve correctly.
		if len(parts) >= 1 && !spotifyContentTypes[parts[0]] {
			parts = parts[1:]
		}
		if len(parts) >= 2 && parts[0] == "track" {
			return parts[1], nil
		}
		return "", errors.New("expected URL like https://open.spotify.com/track/<id>")
	}

	if len(value) == 22 {
		return value, nil
	}

	return "", errors.New("track must be a Spotify track ID, URL, or URI")
}

func spotifyTrackIDToGID(trackID string) (string, error) {
	if trackID == "" {
		return "", errors.New("track ID is empty")
	}

	value := big.NewInt(0)
	base := big.NewInt(62)

	for _, char := range trackID {
		index := strings.IndexRune(spotifyBase62Alphabet, char)
		if index < 0 {
			return "", fmt.Errorf("invalid base62 character: %q", string(char))
		}

		value.Mul(value, base)
		value.Add(value, big.NewInt(int64(index)))
	}

	hexValue := value.Text(16)
	if len(hexValue) < 32 {
		hexValue = strings.Repeat("0", 32-len(hexValue)) + hexValue
	}

	return hexValue, nil
}

func fetchSpotifyTrackRawData(client *http.Client, trackID string) ([]byte, error) {
	accessToken, err := requestSpotifyAnonymousAccessToken(client)
	if err != nil {
		return nil, err
	}

	gid, err := spotifyTrackIDToGID(trackID)
	if err != nil {
		return nil, err
	}

	return requestSpotifyBytes(
		client,
		fmt.Sprintf(spotifyGIDMetadataURL, "track", gid),
		map[string]string{
			"authorization": "Bearer " + accessToken,
			"accept":        "application/json",
		},
	)
}

func extractSpotifyTrackISRC(payload []byte) (string, error) {
	var track spotifyTrackRawData
	if err := json.Unmarshal(payload, &track); err != nil {
		return "", fmt.Errorf("failed to decode Spotify track metadata: %w", err)
	}

	for _, externalID := range track.ExternalID {
		if strings.EqualFold(strings.TrimSpace(externalID.Type), "isrc") {
			if isrc := firstISRCMatch(externalID.ID); isrc != "" {
				return isrc, nil
			}
		}
	}

	if fallbackISRC := firstISRCMatch(string(payload)); fallbackISRC != "" {
		return fallbackISRC, nil
	}

	return "", fmt.Errorf("ISRC not found in Spotify track metadata")
}

// GenerateSpotifyTOTPCode returns a fresh TOTP code and its version number
// using the dynamic secret dictionary fetched from the upstream repository.
// It is the single authoritative TOTP source for the entire backend; both the
// metadata pipeline (spotfetch.go) and the ISRC resolver (isrc_finder.go)
// must use this function so that a secret rotation is handled in one place.
//
// The client parameter is used only when a network request is needed to
// refresh the secrets cache; pass nil to use a default client.
func GenerateSpotifyTOTPCode(client *http.Client) (code string, version int, err error) {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	var serverTime spotifyServerTimeResponse
	if err = requestSpotifyJSON(client, spotifyServerTimeURL, nil, &serverTime); err != nil {
		return "", 0, fmt.Errorf("failed to fetch Spotify server time: %w", err)
	}

	var secrets map[string][]int
	cachedSecrets, cacheErr := loadSpotifyCachedSecrets()
	if cacheErr != nil {
		fmt.Printf("[TOTP] warning: failed to read secrets cache: %v\n", cacheErr)
	}

	if spotifySecretsCacheIsValid(cachedSecrets) {
		secrets = cachedSecrets.Secrets
	} else {
		if fetchErr := requestSpotifyJSON(client, spotifyTOTPSecretsURL, nil, &secrets); fetchErr != nil {
			if cachedSecrets != nil && len(cachedSecrets.Secrets) > 0 {
				fmt.Printf("[TOTP] warning: failed to refresh secrets cache, using stale: %v\n", fetchErr)
				secrets = cachedSecrets.Secrets
			} else {
				return "", 0, fmt.Errorf("failed to fetch Spotify TOTP secrets: %w", fetchErr)
			}
		} else {
			cache := &spotifySecretsCache{
				FetchedAtUnix: time.Now().Unix(),
				Secrets:       secrets,
			}
			if saveErr := saveSpotifyCachedSecrets(cache); saveErr != nil {
				fmt.Printf("[TOTP] warning: failed to write secrets cache: %v\n", saveErr)
			}
		}
	}

	versionStr, err := latestSpotifySecretVersion(secrets)
	if err != nil {
		return "", 0, err
	}

	versionInt, err := strconv.Atoi(versionStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid TOTP secret version %q: %w", versionStr, err)
	}

	secret := deriveSpotifyTOTPSecret(secrets[versionStr])
	code = generateSpotifyTOTP(secret, serverTime.ServerTime*1000)
	return code, versionInt, nil
}
