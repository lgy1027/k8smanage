// @Author : liguoyu
// @Date: 2019/10/29 15:42
package cache

type Cache interface {
	BaseCache
	SCache
	MCache
	HCache
	ZCache
}

type BaseCache interface {
	Get(k string) (v string, exists bool, err error)
	Set(k string, v interface{}, ex int) error
	Del(k string) error
	Close() error
}

type SCache interface {
	Smembers(k string) (vList []string, err error)
	Sadd(k string, vList []string, ex int) error
}

type MCache interface {
	MGet(keys []string) (vList []string, err error)
	MSet(kv map[string]string) (err error)
}

type HCache interface {
	HGet(key string, field interface{}) (string, error)
	HSet(key string, field interface{}, value interface{}, ex int) (err error)
}

type ZCache interface {
	ZAdd(k string, score int64, mem string) (err error)
	ZRangeByScore(key string, args []interface{}) (vList []string, err error)
	ZREVRangeByScore(key string, args []interface{}) (vList []string, err error)
	ZRange(key string, min, max int64) (vmap map[string]int64, err error)
	ZREVRange(key string, min, max int64) (vmap map[string]int64, err error)
}
