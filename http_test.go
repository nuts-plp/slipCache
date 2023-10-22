package slipCache

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestHttpPool_ServeHttp(t *testing.T) {
	db := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}
	_ = NewGroup("class", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[http search] search ", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("key 不存在")
		}))
	addr := "127.0.0.1:8080"
	peer := NewHttpPool(addr)

	log.Println("[slipCache] run at ", addr)
	log.Fatal(http.ListenAndServe(addr, peer))
}

func TestHu(t *testing.T) {
	fmt.Println(string([]byte{118, 49}))
}
