package slipCache

import (
	"fmt"
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
		"5",
		"10",
		"15",
	}
	hashMap.Add(ipCases...)
	kvCases := map[string]string{
		"2":  "5",
		"11": "10",
		"15": "10",
	}
	for k, v := range kvCases {
		fmt.Printf("%s -- %s--》%s\n", k, v, hashMap.Get(k))
		if hashMap.Get(k) != v {
			t.Errorf("查询%s,应该返回%s", k, v)
		}
	}
	hashMap.Add("20")
	kvCases["21"] = "20"
	for k, v := range kvCases {
		if hashMap.Get(k) != v {
			t.Errorf("查%s,应返回%s", k, v)
		}
	}
}
