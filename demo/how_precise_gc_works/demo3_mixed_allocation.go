// demo3_mixed_allocation.go
// 这个demo展示混合分配的情况
// 包含栈分配和堆分配，用于展示GC如何扫描不同的内存区域

package main

import (
	"fmt"
	"runtime"
)

//go:noinline
func mixedAllocationDemo() {
	// 栈分配的局部变量
	var stackInt int64 = 42
	var stackPtr *int64 = &stackInt
	
	// 堆分配（通过返回指针）
	heapPtr := func() *int64 {
		val := int64(100)
		return &val
	}()
	
	// 栈分配的结构体，包含指针字段
	type localStruct struct {
		stackField int64
		heapField  *int64
	}
	
	var local localStruct
	local.stackField = 200
	local.heapField = heapPtr
	
	// 栈分配的数组，包含指针
	var stackArr [3]*int64
	stackArr[0] = &stackInt
	stackArr[1] = heapPtr
	stackArr[2] = &local.stackField
	
	// 堆分配的slice
	heapSlice := make([]*int64, 2)
	heapSlice[0] = &stackInt
	heapSlice[1] = heapPtr
	
	fmt.Printf("Stack int: %d at %p\n", stackInt, &stackInt)
	fmt.Printf("Stack ptr: %p -> %d\n", stackPtr, *stackPtr)
	fmt.Printf("Heap ptr: %p -> %d\n", heapPtr, *heapPtr)
	fmt.Printf("Local struct at %p:\n", &local)
	fmt.Printf("  stackField: %d\n", local.stackField)
	fmt.Printf("  heapField: %p -> %d\n", local.heapField, *local.heapField)
	fmt.Printf("Stack array at %p:\n", &stackArr)
	for i, ptr := range stackArr {
		if ptr != nil {
			fmt.Printf("  [%d]: %p -> %d\n", i, ptr, *ptr)
		}
	}
	fmt.Printf("Heap slice at %p:\n", &heapSlice)
	for i, ptr := range heapSlice {
		if ptr != nil {
			fmt.Printf("  [%d]: %p -> %d\n", i, ptr, *ptr)
		}
	}
	
	// 确保所有变量都被使用
	runtime.KeepAlive(stackInt)
	runtime.KeepAlive(stackPtr)
	runtime.KeepAlive(heapPtr)
	runtime.KeepAlive(local)
	runtime.KeepAlive(stackArr)
	runtime.KeepAlive(heapSlice)
}

//go:noinline
func functionWithPointers(a *int64, b *int64) *int64 {
	// 函数参数中的指针
	// 局部变量
	var local int64 = *a + *b
	localPtr := &local
	
	// 返回局部变量的地址，强制逃逸
	return localPtr
}

//go:noinline
func nestedFunctionCalls() {
	val1 := int64(10)
	val2 := int64(20)
	
	// 调用函数，传递指针参数
	result := functionWithPointers(&val1, &val2)
	
	fmt.Printf("val1: %d at %p\n", val1, &val1)
	fmt.Printf("val2: %d at %p\n", val2, &val2)
	fmt.Printf("result: %p -> %d\n", result, *result)
	
	runtime.KeepAlive(result)
}

func main() {
	fmt.Println("=== Mixed Allocation Demo ===")
	mixedAllocationDemo()
	fmt.Println()
	
	fmt.Println("=== Nested Function Calls ===")
	nestedFunctionCalls()
	fmt.Println()
	
	// 强制GC运行
	runtime.GC()
	
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("After GC:\n")
	fmt.Printf("  Heap allocations: %d\n", memStats.Mallocs)
	fmt.Printf("  Heap size: %d bytes\n", memStats.Alloc)
	fmt.Printf("  Total GC cycles: %d\n", memStats.NumGC)
}
