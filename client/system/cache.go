package client_system

import (
	"github.com/patrickmn/go-cache"
	"time"
	"strconv"
)

var cacheDriver *cache.Cache

//设置微信缓存
func SetWechatBufferCache(buffer []byte) {
	seq := strconv.Itoa(int(ReadInt(buffer, 12)))
	SetCache(seq, buffer, 10*time.Second)
}

//获取微信缓存
func GetWechatBufferCache(longHead []byte) []byte {
	seq := strconv.Itoa(int(ReadInt(longHead, 12)))
	data, has := GetCache(seq)
	if !has {
		return nil
	}
	DeleteCache(seq)
	return data.([]byte)
}

//通用设置缓存
func SetCache(key string, value interface{}, expireTime time.Duration) {
	cacheDriver.Set(key, value, expireTime)
}

//获取缓存
func GetCache(key string) (interface{}, bool) {
	return cacheDriver.Get(key)
}

//删除缓存
func DeleteCache(key string) {
	cacheDriver.Delete(key)
}
