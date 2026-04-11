package backend

import (
	"errors"
	"fmt"
	"strings"
)

type resolvedTrackLinks struct {
	TidalURL  string
	AmazonURL string
	DeezerURL string
	QobuzURL  string
	ISRC      string
}

const (
	linkResolverProviderSongstats      = "songstats"
	linkResolverProviderDeezerSongLink = "deezer-songlink"
)

func (s *SongLinkClient) resolveSpotifyTrackLinks(spotifyTrackID string, region string) (*resolvedTrackLinks, error) {
	links := &resolvedTrackLinks{}
	var attempts []string

	isrc, err := s.lookupSpotifyISRC(spotifyTrackID)
	if err != nil {
		attempts = append(attempts, fmt.Sprintf("spotify isrc: %v", err))
	} else {
		links.ISRC = isrc
	}

	// Check resolver cache before hitting external resolvers.
	if links.ISRC != "" {
		if cached := GetCachedResolvedLinks(links.ISRC); cached != nil {
			links.TidalURL = cached.TidalURL
			links.AmazonURL = cached.AmazonURL
			links.DeezerURL = cached.DeezerURL
			links.QobuzURL = cached.QobuzURL
			fmt.Printf("[Resolver] cache hit for ISRC %s\n", links.ISRC)
			return links, nil
		}
	}

	if links.ISRC != "" {
		resolvers := orderedLinkResolvers()

		for _, resolver := range resolvers {
			switch resolver {
			case linkResolverProviderSongstats:
				addedData, songstatsErr := s.resolveLinksViaSongstats(links)
				if songstatsErr != nil {
					attempts = append(attempts, fmt.Sprintf("songstats: %v", songstatsErr))
				} else if addedData {
					fmt.Println("[Resolver] using Songstats")
				}
			case linkResolverProviderDeezerSongLink:
				addedData, deezerSongLinkErr := s.resolveLinksViaDeezerSongLink(links, region)
				if deezerSongLinkErr != nil {
					attempts = append(attempts, fmt.Sprintf("deezer-songlink: %v", deezerSongLinkErr))
				} else if addedData {
					fmt.Println("[Resolver] using Songlink (deezer-songlink)")
				}
			}

			if links.TidalURL != "" && links.AmazonURL != "" {
				break
			}
		}

		// Write resolved URLs to cache if we got anything useful.
		if hasAnySongLinkData(links) {
			if cacheErr := PutCachedResolvedLinks(links.ISRC, links); cacheErr != nil {
				fmt.Printf("[Resolver] warning: failed to cache resolved links: %v\n", cacheErr)
			}
		}
	}

	if hasAnySongLinkData(links) {
		return links, nil
	}

	if len(attempts) == 0 {
		attempts = append(attempts, "no streaming URLs found")
	}

	return links, errors.New(strings.Join(attempts, " | "))
}

func orderedLinkResolvers() []string {
	preferred := GetLinkResolverSetting()
	if !GetLinkResolverAllowFallback() {
		if preferred == linkResolverProviderDeezerSongLink {
			return []string{linkResolverProviderDeezerSongLink}
		}
		return []string{linkResolverProviderSongstats}
	}

	if preferred == linkResolverProviderDeezerSongLink {
		return []string{
			linkResolverProviderDeezerSongLink,
			linkResolverProviderSongstats,
		}
	}

	return []string{
		linkResolverProviderSongstats,
		linkResolverProviderDeezerSongLink,
	}
}

func (s *SongLinkClient) resolveLinksViaSongstats(links *resolvedTrackLinks) (bool, error) {
	if links == nil || links.ISRC == "" {
		return false, fmt.Errorf("ISRC is required for Songstats resolver")
	}

	before := *links

	fmt.Printf("[Resolver] fetching Songstats links for ISRC %s\n", links.ISRC)
	if err := s.populateLinksFromSongstats(links, links.ISRC); err != nil {
		return false, err
	}

	return *links != before, nil
}

func (s *SongLinkClient) resolveLinksViaDeezerSongLink(links *resolvedTrackLinks, region string) (bool, error) {
	if links == nil || links.ISRC == "" {
		return false, fmt.Errorf("ISRC is required for Deezer song.link resolver")
	}

	before := *links
	var attempts []string

	if links.DeezerURL == "" {
		fmt.Printf("[Resolver] resolving Deezer track from ISRC %s\n", links.ISRC)
		deezerURL, err := s.lookupDeezerTrackURLByISRC(links.ISRC)
		if err != nil {
			attempts = append(attempts, fmt.Sprintf("deezer isrc: %v", err))
		} else {
			links.DeezerURL = deezerURL
			fmt.Printf("[Resolver] found Deezer URL: %s\n", links.DeezerURL)
		}
	}

	if links.DeezerURL != "" {
		fmt.Println("[Resolver] resolving streaming URLs via song.link + Deezer URL")
		deezerResp, err := s.fetchSongLinkLinksByURL(links.DeezerURL, region)
		if err != nil {
			attempts = append(attempts, fmt.Sprintf("song.link deezer: %v", err))
		} else {
			mergeSongLinkResponse(links, deezerResp)
		}

		if links.ISRC == "" {
			if resolvedISRC, deezerISRCErr := getDeezerISRC(links.DeezerURL); deezerISRCErr == nil {
				links.ISRC = resolvedISRC
			}
		}
	}

	if *links != before {
		if len(attempts) == 0 {
			return true, nil
		}
		return true, errors.New(strings.Join(attempts, " | "))
	}

	if len(attempts) == 0 {
		attempts = append(attempts, "no links found via deezer-songlink")
	}

	return false, errors.New(strings.Join(attempts, " | "))
}
