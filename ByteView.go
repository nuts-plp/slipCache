package slipCache

type ByteView struct {
	b []byte
}

func (B ByteView) Len() int {
	return len(B.b)
}

func (B ByteView) ByteSlice() []byte {
	return cloneBytes(B.b)
}
func (B ByteView) String() string {
	return string(B.b)
}
func cloneBytes(bytes []byte) []byte {
	clone := make([]byte, len(bytes))
	copy(clone, bytes)
	return clone
}
