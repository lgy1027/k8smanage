// @Author : liguoyu
// @Date: 2019/10/29 15:42
package protocol

import "testing"

func TestRedis_Init(t *testing.T) {
	r := &Redis{
		Address: "127.0.0.1:6379",
		Db:      1,
		pool:    nil,
	}
	if err := r.Init(); err != nil {
		t.Error(err)
	}
}
