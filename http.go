package slipCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

//defaultBasePath
const defaultBasePath = "/_nuts/"

type Handler interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

//HttpPool 实现一个http节点的  e.g http://example.com:8080
type HttpPool struct {
	self     string
	basePath string
}

//NewHttpPool 实力一个当前节点的地址   http://example.com:8080/_nuts/
func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

//Log 格式化服务器日志打印
func (h *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", h.self, fmt.Sprintf(format, v...))
}

func (h *HttpPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, h.basePath) {
		panic("错误的路径 ")
	}
	h.Log("%s  %s", req.Method, req.URL.Path)
	// basePath/group/key
	parts := strings.SplitN(req.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "错误的请求", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "不存在的group"+groupName, http.StatusNotFound)
		return
	}
	value, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/octet-stream")
	w.Write(value.ByteSlice())
}
