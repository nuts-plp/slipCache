package slipCache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(data []byte) uint32

type Map struct {
	hashFunc   HashFunc
	replicas   int
	hashCircle []int
	hashMap    map[int]string
}

//New 实例化一个全局HashMap实例
func New(replicas int, fn HashFunc) *Map {
	m := &Map{
		replicas: replicas,
		hashMap:  make(map[int]string),
		hashFunc: fn,
	}
	if fn == nil {
		m.hashFunc = crc32.ChecksumIEEE
		return m
	}
	return m
}

//Add 添加节点地址并创建对应数量的虚拟节点并使之对应真实节点
func (m *Map) Add(urls ...string) {
	for _, url := range urls {
		for i := 0; i < m.replicas; i++ {
			//节点对应的虚拟节点
			hash := int(m.hashFunc([]byte(strconv.Itoa(i) + url)))
			m.hashMap[hash] = url
			m.hashCircle = append(m.hashCircle, hash)
		}
	}
	sort.Ints(m.hashCircle)
}

func (m *Map) Get(key string) string {
	if len(m.hashCircle) == 0 {
		return ""
	}
	hash := int(m.hashFunc([]byte(key)))
	idx := sort.Search(len(m.hashCircle), func(i int) bool {
		return m.hashCircle[i] > hash
	})
	//返回虚拟节点对应的真实节点
	return m.hashMap[m.hashCircle[idx%len(m.hashCircle)]]
}
