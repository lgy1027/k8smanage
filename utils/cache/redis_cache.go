// @Author : liguoyu
// @Date: 2019/10/29 15:42
package cache

import (
	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

type RedisCache struct {
	*redis.Pool
}

func NewRedisCache(pool *redis.Pool) *RedisCache {
	return &RedisCache{
		Pool: pool,
	}
}

func (s *RedisCache) getConn() redis.Conn {
	return s.Pool.Get()
}

func (s *RedisCache) Close() error {
	if err := s.getConn().Close(); err != nil {
		log.Error("err", err)
		return err
	}
	return nil
}

func (s *RedisCache) Set(key string, value interface{}, ex int) (err error) {
	//log.Debug("action", "Set", "key", key, "value", value, "ex", ex)
	if ex > 0 {
		_, err = s.getConn().Do("SET", key, value, "EX", ex)
	} else {
		_, err = s.getConn().Do("SET", key, value)
	}
	if err != nil {
		log.Error("err", err)
	}
	return err
}

func (s *RedisCache) Get(key string) (string, bool, error) {
	val, err := redis.String(s.getConn().Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return "", false, nil
		}
		return "", false, err
	}
	return val, true, nil
}

func (s *RedisCache) Sadd(key string, vList []string, ex int) (err error) {
	log.Debug("action", "SADD", "key", key, "vListLength", len(vList), "ex", ex)
	var args = make([]interface{}, 0, len(vList)+1)
	args = append(args, key)
	for _, v := range vList {
		args = append(args, v)
	}
	if _, err = s.getConn().Do("SADD", args...); err != nil {
		log.Error("err", err)
	}
	if ex > 0 {
		if _, err = s.getConn().Do("EXPIRE", key, ex); err != nil {
			log.Error("err", err)
		}
	}
	return
}

func (s *RedisCache) Smembers(key string) (vList []string, err error) {
	log.Debug("action", "SMEMBERS", "key", key)
	vList, err = redis.Strings(s.getConn().Do("SMEMBERS", key))
	return
}

func (s *RedisCache) Del(key string) error {
	log.Debug("action", "Del", "key", key)
	if _, err := s.getConn().Do("DEL", key); err != nil {
		log.Error("err", err)
		return err
	}
	return nil
}

func (s *RedisCache) MGet(keys []string) (vList []string, err error) {
	log.Debug("action", "MGet", "keyLength", len(keys))
	var (
		args = make([]interface{}, 0, len(keys))
	)
	for _, k := range keys {
		args = append(args, k)
	}
	vList, err = redis.Strings(s.getConn().Do("MGET", args...))
	return
}

func (s *RedisCache) MSet(kv map[string]string) (err error) {
	log.Debug("action", "MSet", "keyLength", len(kv))
	var args = make([]interface{}, 0, len(kv)*2)
	for k, v := range kv {
		args = append(args, k, v)
	}
	if _, err = s.getConn().Do("MSET", args...); err != nil {
		log.Error("err", err)
		return
	}
	return
}

func (s *RedisCache) HSet(key string, field interface{}, value interface{}, ex int) (err error) {
	log.Debug("action:", " HSet", " key:", key, " field:", field)
	if _, err = s.getConn().Do("HSet", key, field, value); err != nil {
		log.Error("err", err)
		return
	}

	if ex > 0 {
		if _, err = s.getConn().Do("EXPIRE", key, ex); err != nil {
			log.Error("err", err)
		}
	}

	return
}

func (s *RedisCache) HDel(key string, field interface{}) (err error) {
	log.Debug("action: ", " HDel", " key: ", key, " field:", field)
	if _, err = s.getConn().Do("HDEL", key, field); err != nil {
		log.Error("err", err)
		return
	}
	return
}

func (s *RedisCache) HVals(key string) ([]interface{}, bool, error) {
	log.Debug("action: ", " HVals: ", " key: ", key)
	r, err := redis.Values(s.getConn().Do("hvals", key))
	if err != nil {
		return nil, false, err
	} else if len(r) == 0 {
		return r, false, nil
	}
	return r, true, nil
}

func (s *RedisCache) HKeys(key string) ([]interface{}, bool, error) {
	log.Debug("action: ", " HKeys: ", " key: ", key)
	r, err := redis.Values(s.getConn().Do("hkeys", key))
	if err != nil {
		return nil, false, err
	} else if len(r) == 0 {
		return r, false, nil
	}
	return r, true, nil
}

func (s *RedisCache) HGet(key string, field interface{}) (string, bool, error) {
	log.Debug("action", "HGet", "key", key, "field", field)
	val, err := redis.String(s.getConn().Do("HGet", key, field))
	if err != nil {
		if err == redis.ErrNil {
			return "", false, nil
		}
		return "", false, err
	}
	return val, true, nil
}

func (s *RedisCache) ZAdd(k string, score int64, mem string) (err error) {
	args := []interface{}{k, score, mem}
	if _, err = s.getConn().Do("ZADD", args...); err != nil {
		log.Error("err", err)
		return
	}
	return
}

func (s *RedisCache) ZRangeByScore(key string, arg []interface{}) ([]string, error) {
	args := []interface{}{key}
	args = append(args, arg...)
	r, err := redis.Strings(s.getConn().Do("ZRANGEBYSCORE", args...))
	if err != nil {
		return nil, err
	}
	return r, err
}

func (s *RedisCache) ZREVRangeByScore(key string, arg []interface{}) ([]string, error) {
	args := []interface{}{key}
	args = append(args, arg...)
	r, err := redis.Strings(s.getConn().Do("ZREVRANGEBYSCORE", args...))
	if err != nil {
		return nil, err
	}
	return r, err
}

func (s *RedisCache) ZRange(key string, start, stop int64) (map[string]int64, error) {
	args := []interface{}{key, start, stop, "withscores"}
	r, err := redis.Int64Map(s.getConn().Do("ZRANGE", args...))
	if err != nil {
		return nil, err
	}
	return r, err
}

func (s *RedisCache) ZREVRange(key string, start, stop int64) (map[string]int64, error) {
	args := []interface{}{key, start, stop, "withscores"}
	r, err := redis.Int64Map(s.getConn().Do("ZREVRANGE", args...))
	if err != nil {
		return nil, err
	}
	return r, err
}
