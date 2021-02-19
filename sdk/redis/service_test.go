package redis

import "testing"

func TestRedisSet(t *testing.T) {
	res, err := RedisSet("name", "aozhu")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
