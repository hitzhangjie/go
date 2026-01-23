// demo1_stack_allocation.go
// 这个demo展示栈分配的情况
// 通过使用runtime.KeepAlive确保变量不会被优化掉
// 同时确保变量不会逃逸到堆上

package main

import (
	"fmt"
	"runtime"
)

//go:noinline
func stackAllocDemo() {
	// 创建一个局部变量，确保在栈上分配
	// 使用足够大的值来避免被优化
	var localInt int64 = 0x1234567890ABCDEF
	var localPtr *int64 = &localInt
	
	// 使用runtime.KeepAlive确保变量不会被优化掉
	runtime.KeepAlive(localPtr)
	
	// 打印地址，用于验证
	fmt.Printf("Stack allocated localInt address: %p\n", &localInt)
	fmt.Printf("Stack allocated localPtr value: %p\n", localPtr)
}

//go:noinline
func stackPointerDemo() {
	// 创建一个包含指针的局部变量
	type localStruct struct {
		field1 int64
		field2 *int64
	}
	
	var local localStruct
	local.field1 = 42
	local.field2 = &local.field1
	
	// 确保变量不会被优化
	runtime.KeepAlive(&local)
	
	fmt.Printf("Stack struct address: %p\n", &local)
	fmt.Printf("Stack struct field2 pointer: %p\n", local.field2)
}

//go:noinline
func stackArrayDemo() {
	// 创建一个局部数组，包含指针
	var arr [4]*int64
	for i := range arr {
		var val int64 = int64(i)
		arr[i] = &val
	}
	
	runtime.KeepAlive(arr)
	
	fmt.Printf("Stack array address: %p\n", &arr)
	for i, ptr := range arr {
		fmt.Printf("  arr[%d] = %p (value: %d)\n", i, ptr, *ptr)
	}
}

//go:noinline
func stackSliceDemo() {
	// 创建一个局部slice，包含指针
	// 注意：slice本身在栈上，但底层数组可能在堆上
	// 这里我们确保slice本身在栈上
	var slice []*int64
	slice = make([]*int64, 3)
	
	for i := range slice {
		var val int64 = int64(i * 10)
		slice[i] = &val
	}
	
	runtime.KeepAlive(slice)
	
	fmt.Printf("Stack slice address: %p\n", &slice)
	fmt.Printf("Stack slice len: %d, cap: %d\n", len(slice), cap(slice))
}

func main() {
	fmt.Println("=== Demo 1: Stack Allocation ===")
	stackAllocDemo()
	fmt.Println()
	
	fmt.Println("=== Demo 2: Stack Pointer ===")
	stackPointerDemo()
	fmt.Println()
	
	fmt.Println("=== Demo 3: Stack Array ===")
	stackArrayDemo()
	fmt.Println()
	
	fmt.Println("=== Demo 4: Stack Slice ===")
	stackSliceDemo()
	fmt.Println()
	
	// 打印一些运行时信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Heap allocations: %d\n", m.Mallocs)
	fmt.Printf("Heap size: %d bytes\n", m.Alloc)
}
