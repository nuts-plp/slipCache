package slipCache

import (
	"reflect"
	"testing"
)

func TestByteView_String(t *testing.T) {
	str := "slice"
	b := ByteView{
		[]byte("slice"),
	}
	if !reflect.DeepEqual(str, b.String()) {
		t.Fatalf("字节切片转换为字符串失败！")
	}
}
