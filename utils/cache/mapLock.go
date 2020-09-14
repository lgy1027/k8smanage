package cache

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"
)

type Lock struct {
	Lock sync.RWMutex
	Data map[string]interface{}
}

func (m *Lock) GetData() map[string]interface{} {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	return m.Data
}

// 获取数据
func (m *Lock) GetV(key string) interface{} {
	v, _ := m.Get(key)
	return v
}

// 获取数据
func (m *Lock) GetVString(key string) string {
	v, ok := m.Get(key)
	if ok {
		switch v.(type) {
		case string:
			return v.(string)
			break
		case int32:
			return strconv.Itoa(int(v.(int32)))
		case int64:
			return strconv.Itoa(int(v.(int64)))
		}
	}
	return ""
}

// 获取数据
func (m *Lock) Get(key string) (interface{}, bool) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	data := m.GetData()
	if _, ok := data[key]; ok {
		return data[key], true
	}
	return nil, false
}

// 添加数据
func (m *Lock) Put(k string, v interface{}) {
	if len(m.Data) < 1 {
		m.Data = make(map[string]interface{})
	}
	m.Lock.Lock()
	m.Data[k] = v
	defer m.Lock.Unlock()
}

// 将对象转字符串
func (m *Lock) String() string {
	v, err := json.Marshal(m.GetData())
	if err == nil {
		return string(v)
	}
	return err.Error()
}

// 避免频繁更新,加锁n秒后可操作
func WriteLock(key string, lock *Lock, timeout int64) bool {
	if len(lock.GetData()) > 0 {
		v, err := lock.Get(key)
		if err {
			last := v.(int64)
			if time.Now().Unix()-last < timeout {
				return false
			}
		}
	}
	lock.Put(key, time.Now().Unix())
	return true
}
