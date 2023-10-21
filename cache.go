package slipCache

import (
	"fmt"
	"log"
	"slipCache/lru"
	"sync"
)

//cache 借助cache实现lru.Cache并发
type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return ByteView{}, false
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return ByteView{}, false
}

//Getter 设置回调函数，如果数据不存在，从源处获取数据
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

//Group 缓存数据库的核心，一个Group就是一个缓存空间
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

//NewGroup 实例化一个Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("Getter nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

//GetGroup 根据名称获取实例,如果为nil则表示没有该Group实例
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

//Get 整个最核心的方法，先判断key是否为空字符串，再向本地查询，最后向源查询，查到后添加至本地
func (g *Group) Get(key string) (ByteView, error) {
	mu.RLock()
	defer mu.RUnlock()
	if key == "" {
		return ByteView{}, fmt.Errorf("key 为空字符串")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[slipCache] cache hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.loadLocal(key)
}

func (g *Group) loadLocal(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	v := ByteView{cloneBytes(bytes)}
	g.populateCache(key, v)
	return v, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
