package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const (
	// defaultSpotFetchAPIURL is the fallback SpotFetch API endpoint used when
	// the user has enabled the API but has not set a custom URL in settings.
	defaultSpotFetchAPIURL = "https://sp.afkarxyz.qzz.io/api"
)

func GetDefaultMusicPath() string {

	homeDir, err := os.UserHomeDir()
	if err != nil {

		return "C:\\Users\\Public\\Music"
	}

	return filepath.Join(homeDir, "Music")
}

func GetConfigPath() (string, error) {
	dir, err := EnsureAppDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.json"), nil
}

func LoadConfigSettings() (map[string]interface{}, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func GetSpotFetchAPISettings() (bool, string) {
	settings, err := LoadConfigSettings()
	if err != nil || settings == nil {
		return false, ""
	}

	useAPI, _ := settings["useSpotFetchAPI"].(bool)
	if !useAPI {
		return false, ""
	}

	apiURL, _ := settings["spotFetchAPIUrl"].(string)
	if apiURL == "" {
		apiURL = defaultSpotFetchAPIURL
	}

	return true, apiURL
}

func GetLinkResolverSetting() string {
	settings, err := LoadConfigSettings()
	if err != nil || settings == nil {
		return linkResolverProviderDeezerSongLink
	}

	resolver, _ := settings["linkResolver"].(string)
	switch strings.TrimSpace(strings.ToLower(resolver)) {
	case "songlink", linkResolverProviderDeezerSongLink:
		return linkResolverProviderDeezerSongLink
	case "songstats":
		return linkResolverProviderSongstats
	case "":
		return linkResolverProviderDeezerSongLink
	default:
		return linkResolverProviderDeezerSongLink
	}
}

func GetLinkResolverAllowFallback() bool {
	settings, err := LoadConfigSettings()
	if err != nil || settings == nil {
		return true
	}

	allowFallback, ok := settings["allowResolverFallback"].(bool)
	if !ok {
		return true
	}

	return allowFallback
}
