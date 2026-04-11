package backend

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	resolverCacheDBFile = "resolver_cache.db"
	resolverCacheBucket = "ResolvedLinks"

	// resolverCacheTTL controls how long resolved streaming URLs are cached.
	// 30 days is a reasonable balance: URLs are stable but platforms
	// occasionally retire or reassign track IDs, so we don't cache forever.
	resolverCacheTTL = 30 * 24 * time.Hour
)

type resolverCacheEntry struct {
	ISRC      string `json:"isrc"`
	TidalURL  string `json:"tidal_url,omitempty"`
	AmazonURL string `json:"amazon_url,omitempty"`
	DeezerURL string `json:"deezer_url,omitempty"`
	QobuzURL  string `json:"qobuz_url,omitempty"`
	UpdatedAt int64  `json:"updated_at"`
}

var (
	resolverCacheDB   *bolt.DB
	resolverCacheDBMu sync.Mutex
)

func InitResolverCacheDB() error {
	resolverCacheDBMu.Lock()
	defer resolverCacheDBMu.Unlock()

	if resolverCacheDB != nil {
		return nil
	}

	appDir, err := EnsureAppDir()
	if err != nil {
		return err
	}

	dbPath := filepath.Join(appDir, resolverCacheDBFile)
	db, err := bolt.Open(dbPath, 0o600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(resolverCacheBucket))
		return err
	}); err != nil {
		db.Close()
		return err
	}

	resolverCacheDB = db
	return nil
}

func CloseResolverCacheDB() {
	resolverCacheDBMu.Lock()
	defer resolverCacheDBMu.Unlock()

	if resolverCacheDB != nil {
		_ = resolverCacheDB.Close()
		resolverCacheDB = nil
	}
}

// GetCachedResolvedLinks returns cached streaming URLs for the given ISRC,
// or nil if no valid cache entry exists.
func GetCachedResolvedLinks(isrc string) *resolverCacheEntry {
	normalizedISRC := strings.ToUpper(strings.TrimSpace(isrc))
	if normalizedISRC == "" {
		return nil
	}

	if err := InitResolverCacheDB(); err != nil {
		return nil
	}

	var entry *resolverCacheEntry
	_ = resolverCacheDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(resolverCacheBucket))
		if bucket == nil {
			return nil
		}

		value := bucket.Get([]byte(normalizedISRC))
		if len(value) == 0 {
			return nil
		}

		var e resolverCacheEntry
		if err := json.Unmarshal(value, &e); err != nil {
			return nil
		}

		if e.UpdatedAt > 0 && time.Since(time.Unix(e.UpdatedAt, 0)) > resolverCacheTTL {
			return nil // expired
		}

		entry = &e
		return nil
	})

	return entry
}

// PutCachedResolvedLinks stores resolved streaming URLs for the given ISRC.
// Only entries that contain at least one URL are written.
func PutCachedResolvedLinks(isrc string, links *resolvedTrackLinks) error {
	normalizedISRC := strings.ToUpper(strings.TrimSpace(isrc))
	if normalizedISRC == "" || links == nil {
		return nil
	}

	// Don't cache if we have nothing useful.
	if links.TidalURL == "" && links.AmazonURL == "" && links.DeezerURL == "" && links.QobuzURL == "" {
		return nil
	}

	if err := InitResolverCacheDB(); err != nil {
		return err
	}

	entry := resolverCacheEntry{
		ISRC:      normalizedISRC,
		TidalURL:  links.TidalURL,
		AmazonURL: links.AmazonURL,
		DeezerURL: links.DeezerURL,
		QobuzURL:  links.QobuzURL,
		UpdatedAt: time.Now().Unix(),
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to encode resolver cache entry: %w", err)
	}

	return resolverCacheDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(resolverCacheBucket))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(normalizedISRC), payload)
	})
}
