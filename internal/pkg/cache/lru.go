package cache

import (
	"fmt"
	"main/config"
	"main/internal/pkg/logger"

	lru "github.com/hashicorp/golang-lru"
)

var lruCache *lru.TwoQueueCache

func init() {
	if err := initLRU(config.Conf.Svr.LRUSize); err != nil {
		panic(err)
	} else {
		fmt.Printf("lru cache init succ. size:%d\n", config.Conf.Svr.LRUSize)
	}
}

// initLRU 初始化LRU cache
func initLRU(size int) error {
	var err error
	lruCache, err = lru.New2Q(size)
	if err != nil {
		logger.ERROR(fmt.Sprintf("cache Init fail, error:%s", err.Error()))
		return err
	}
	return nil
}

// Put 写入
func LRUPut(key uint, value interface{}) {
	lruCache.Add(key, value)
}

// Get 获取
func LRUGet(key uint) (value interface{}, ok bool) {
	return lruCache.Get(key)
}
