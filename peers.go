package slipCache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

//PeerPicker 根据key获取节点
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

//PeerGetter 根据命名空间和key获取对应的value
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

type HttpGetter struct {
	baseUrl string
}

func (h *HttpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v/%v/%v",
		h.baseUrl,
		url.QueryEscape(group),
		url.QueryEscape(key))
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("路径错误，err:%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("节点服务器相应信息:%v", resp.Body)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应body失败，body:%v", resp.Body)
	}
	return bytes, nil
}

var _ PeerGetter = (*HttpGetter)(nil)

const (
	defaultRaplicas = 50
)

//HTTPPool 实例化一个PeerPicker  实现一个http peer 池
type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *Map
	httpGetters map[string]*HttpGetter
}

func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = New(defaultRaplicas, nil)
	h.peers.Add(peers...)
	h.httpGetters = make(map[string]*HttpGetter, len(peers))
	for _, peer := range peers {
		h.httpGetters[peer] = &HttpGetter{peer + h.basePath}
	}
}

func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		h.Log("Pick peer %s", peer)
		return h.httpGetters[peer], true
	}
	return nil, false
}
func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Panicf("[HTTPPool] server is %s,%s", h.self, fmt.Sprintf(format, v...))
}

var _ PeerPicker = (*HTTPPool)(nil)
