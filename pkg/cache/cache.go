package cache

import (
	"github.com/coocood/freecache"
	"log"
	"runtime/debug"
)

var TokenCache *freecache.Cache

const PoolTokenExpired = 60 * 60 * 24 * 7 // 7 天

func init() {
	var err error
	// In bytes, where 1024 * 1024 represents a single Megabyte, and 100 * 1024*1024 represents 100 Megabytes.
	cacheSize := 100 * 1024 * 1024
	TokenCache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
	if err != nil {
		log.Println("init TokenCache error ", err.Error())
		panic(err)
	}
	log.Println("init TokenCache success")
}
