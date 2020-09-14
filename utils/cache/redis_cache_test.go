// @Author : liguoyu
// @Date: 2019/10/29 15:42
package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"testing"
)

var pool *redis.Pool

func init() {
	opts := make([]redis.DialOption, 0)
	//opt := redis.DialPassword("lgy")
	//opts = append(opts,opt)
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379", opts...)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

func TestRedisCache_LPush(t *testing.T) {
	cache := NewRedisCache(pool)
	if err := cache.Sadd("AA", []string{"BBB", "CCCC"}, 100); err != nil {
		t.Error(err)
	}
}

func TestRedisCache_Smembers(t *testing.T) {
	cache := NewRedisCache(pool)
	cache.Sadd("AA", []string{"BBB", "CCCC"}, 100)
	var (
		vList []string
		err   error
	)
	if vList, err = cache.Smembers("AA"); err != nil || len(vList) != 2 {
		t.Error(err)
	}
}

func TestRedisCache_MSet(t *testing.T) {
	var defaultKv = make(map[string]string)
	defaultKv["kvTest1"] = "val1"
	defaultKv["kvTest2"] = "val2"
	defaultKv["kvTest3"] = "val3"
	cache := NewRedisCache(pool)
	if err := cache.MSet(defaultKv); err != nil {
		t.Error(err)
	}
}

func TestRedisCache_MGet(t *testing.T) {
	var defaultKList = []string{"kvTest1", "kvTest2", "kvTest3", "kvTest4"}
	cache := NewRedisCache(pool)
	if vList, err := cache.MGet(defaultKList); err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", vList)
	}
}

func TestRedisCache_HSet(t *testing.T) {
	cache := NewRedisCache(pool)
	err := cache.HSet("hset", "faild", "test", 0)
	if err != nil {
		t.Error(err)
	}
}

func TestRedisCache_HGet(t *testing.T) {
	cache := NewRedisCache(pool)
	s, err := cache.HGet("hset", "faild")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf(s)
	}
}

func Test_Zadd(t *testing.T) {
	cache := NewRedisCache(pool)
	var i int64
	for i = 0; i < 10; i++ {
		err := cache.ZAdd("space", i, "node"+strconv.Itoa(int(i)))
		fmt.Println(err)
	}
}

func Test_ZRangeByScore(t *testing.T) {
	cache := NewRedisCache(pool)
	arg := []interface{}{"-inf", 10}
	vList, err := cache.ZRangeByScore("space", arg)
	fmt.Println(err)
	fmt.Println(vList)
}

func Test_ZREVRangeByScore(t *testing.T) {
	cache := NewRedisCache(pool)
	arg := []interface{}{"+inf", 0, "limit", 0, 5}
	vList, err := cache.ZREVRangeByScore("space", arg)
	fmt.Println(err)
	fmt.Println(vList)
}

func Test_Range(t *testing.T) {
	cache := NewRedisCache(pool)
	vList, err := cache.ZRange("space", 0, -1)
	fmt.Println(err)
	fmt.Println(vList)
}

func Test_ZREVRANGE(t *testing.T) {
	cache := NewRedisCache(pool)
	vList, err := cache.ZREVRange("space", 0, -1)
	fmt.Println(err)
	fmt.Println(vList)
}
