package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)

}
func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key", String("7870"))
	if v, ok := lru.Get("key"); !ok || v.(String) != "7870" {
		t.Fatal("存取失败")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("获取一个不存在的值却能够获取到！")
	}

}

func TestCache_RemoveLRU(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "123", "234", "345"
	ca := len(k1 + k2 + v1 + v2)
	lru := New(int64(ca), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("lru策略失效")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	lru := New(int64(10), func(s string, value Value) {
		keys = append(keys, s)
	})
	lru.Add("k1", String("798910"))
	lru.Add("k2", String("hfakjs"))
	expect := []string{"k1", "k2"}
	if !reflect.DeepEqual(keys, expect) {
		t.Fatalf("回调函数失败了")
	}
}
