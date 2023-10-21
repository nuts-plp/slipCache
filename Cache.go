package slipCache

import (
	"container/list"
)

// 使用go内置的链表来存储数据，使用map来实现字典，map的值指向节点
type (
	//Cache 缓存对象
	Cache struct {
		maxBytes  int64                         //最大使用内存
		nBytes    int64                         //已使用内存
		ll        *list.List                    //go内置的双向链表
		dic       map[string]*list.Element      //map实现字典，用于查找元素
		OnEvicted func(key string, value Value) //callback function 用于删除时的处理
	}
	//Value 链表节点的存储单元，存储键值对
	unit struct {
		key   string
		value Value
	}
	// Value 链表节点的大小
	Value interface {
		Len() int
	}
)

//New 实例化缓存对象
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		dic:       make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//Get 查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.dic[key]; ok {
		c.ll.MoveToFront(element) //如果查到数据，把数据移到表首
		kv := element.Value.(*unit)
		return kv.value, ok
	}
	return nil, false

}

//RemoveLRU 删除最近最少使用
func (c *Cache) RemoveLRU() {
	element := c.ll.Back()
	if element != nil {
		c.ll.Remove(element)
		kv := element.Value.(*unit)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

//Add 添加/修改
//	修改数据后，把数据提到表首，同时更新占用内存
//	添加数据后，把数据添加至表头，更新内存
//	最后循环判断，更新后的内存占用是否超出，如果超出根据LRU删除，直至占用内存小于最大内存
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.dic[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*unit)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		element := c.ll.PushFront(&unit{key, value})
		c.dic[key] = element
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes < c.nBytes && c.maxBytes != 0 {
		c.RemoveLRU()
	}
}

//Len 返回链表长度
func (c *Cache) Len() int {
	return c.ll.Len()
}
