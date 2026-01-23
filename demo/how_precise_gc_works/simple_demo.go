// simple_demo.go
// 一个简单的demo，用于展示PCDATA和FUNCDATA指令

package main

//go:noinline
func stackAlloc() {
	var x int64 = 42
	var p *int64 = &x
	_ = p
}

//go:noinline
func heapAlloc() *int64 {
	x := int64(100)
	return &x
}

//go:noinline
func withPointers(a *int64, b *int64) *int64 {
	var local int64 = *a + *b
	return &local
}

func main() {
	stackAlloc()
	ptr := heapAlloc()
	_ = ptr
	val1 := int64(10)
	val2 := int64(20)
	result := withPointers(&val1, &val2)
	_ = result
}
