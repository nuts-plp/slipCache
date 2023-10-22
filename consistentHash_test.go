package slipCache

import (
	"strconv"
	"testing"
)

func TestMap_Add(t *testing.T) {
	hashMap := New(3, func(data []byte) uint32 {
		h, _ := strconv.Atoi(string(data))
		return uint32(h)
	})
	// 5+1 5+2 5+3 10+1 10+2 10+3  15+1 15+2 15+3
	ipCases := []string{
		"6",
		"4",
		"2",
	}
	hashMap.Add(ipCases...)
	kvCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range kvCases {
		if hashMap.Get(k) != v {
			t.Errorf("查询%s,应该返回%s", k, v)
		}
	}
	hashMap.Add("8")
	kvCases["27"] = "8"
	for k, v := range kvCases {
		if hashMap.Get(k) != v {
			t.Errorf("查%s,应返回%s", k, v)
		}
	}
}
