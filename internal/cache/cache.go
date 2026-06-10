package cache

import (
    "time"
    "github.com/jellydator/ttlcache/v3"
)

var (
	StringCache *ttlcache.Cache[string, string]
)

func InitializeCache(){
	StringCache = ttlcache.New(ttlcache.WithTTL[string, string](5 * time.Minute))
	go StringCache.Start()
}

func StopCache(){
	go StringCache.Stop()
}


