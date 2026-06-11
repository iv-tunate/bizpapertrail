package cache

import (
    "time"
    "github.com/jellydator/ttlcache/v3"
)

var (
	StringCache *ttlcache.Cache[string, string]
	BlacklistedTokensCache *ttlcache.Cache[string, bool]
)

func InitializeCache(){
	StringCache = ttlcache.New(ttlcache.WithTTL[string, string](5 * time.Minute))
	BlacklistedTokensCache = ttlcache.New[string, bool]()
	go StringCache.Start()
	go BlacklistedTokensCache.Start()
}

func StopCache(){
	go StringCache.Stop()
	go BlacklistedTokensCache.Stop()
}


