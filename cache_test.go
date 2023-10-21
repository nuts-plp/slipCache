package slipCache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetterFunc_Get(t *testing.T) {
	getter := GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("123456")
	v, _ := getter.Get("123456")
	if !reflect.DeepEqual(expect, v) {
		t.Fatalf("回调函数设置失败")
	}
}

func TestGroup_Get(t *testing.T) {
	db := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	}
	loadCounts := make(map[string]int, len(db))
	slip := NewGroup("class", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := db[key]; ok {
				log.Println("[search] search ", key)
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s 不存在", key)
		}))
	for k, v := range db {
		if b, err := slip.Get(k); err != nil || b.String() != v {
			t.Fatalf("根据key取到了错误的value")
		}
		if _, err := slip.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("缓存k的值丢失")
		}
	}
	if _, err := slip.Get("unknown"); err == nil {
		t.Fatalf("缓存根据一个不存在的key取到了值")
	}
}
