// demo2_heap_allocation.go
// 这个demo展示堆分配的情况
// 通过返回指针或使用interface{}来强制变量逃逸到堆上

package main

import (
	"fmt"
	"runtime"
)

//go:noinline
func heapAllocDemo1() *int64 {
	// 返回局部变量的地址，强制逃逸到堆
	var heapInt int64 = 0x1234567890ABCDEF
	return &heapInt
}

//go:noinline
func heapAllocDemo2() interface{} {
	// 使用interface{}强制逃逸
	var heapInt int64 = 0x7FFFFFFFFFFFFFFF
	return &heapInt
}

//go:noinline
func heapStructDemo() *struct {
	field1 int64
	field2 *int64
} {
	// 返回结构体指针，强制逃逸
	s := &struct {
		field1 int64
		field2 *int64
	}{
		field1: 100,
	}
	s.field2 = &s.field1
	return s
}

//go:noinline
func heapArrayDemo() *[4]*int64 {
	// 返回数组指针，强制逃逸
	arr := &[4]*int64{}
	for i := range arr {
		val := int64(i * 100)
		arr[i] = &val
	}
	return arr
}

//go:noinline
func heapSliceDemo() []*int64 {
	// 返回slice，底层数组在堆上
	slice := make([]*int64, 5)
	for i := range slice {
		val := int64(i * 200)
		slice[i] = &val
	}
	return slice
}

//go:noinline
func heapMapDemo() map[int]*int64 {
	// 返回map，键值对在堆上
	m := make(map[int]*int64)
	for i := 0; i < 3; i++ {
		val := int64(i * 300)
		m[i] = &val
	}
	return m
}

func main() {
	fmt.Println("=== Demo 1: Heap Allocation (Return Pointer) ===")
	ptr1 := heapAllocDemo1()
	fmt.Printf("Heap allocated int address: %p\n", ptr1)
	fmt.Printf("Heap allocated int value: 0x%x\n", *ptr1)
	runtime.KeepAlive(ptr1)
	fmt.Println()
	
	fmt.Println("=== Demo 2: Heap Allocation (Interface) ===")
	iface := heapAllocDemo2()
	ptr2 := iface.(*int64)
	fmt.Printf("Heap allocated int (via interface) address: %p\n", ptr2)
	fmt.Printf("Heap allocated int (via interface) value: 0x%x\n", *ptr2)
	runtime.KeepAlive(ptr2)
	fmt.Println()
	
	fmt.Println("=== Demo 3: Heap Struct ===")
	s := heapStructDemo()
	fmt.Printf("Heap struct address: %p\n", s)
	fmt.Printf("Heap struct field1: %d\n", s.field1)
	fmt.Printf("Heap struct field2 pointer: %p\n", s.field2)
	runtime.KeepAlive(s)
	fmt.Println()
	
	fmt.Println("=== Demo 4: Heap Array ===")
	arr := heapArrayDemo()
	fmt.Printf("Heap array address: %p\n", arr)
	for i, ptr := range arr {
		if ptr != nil {
			fmt.Printf("  arr[%d] = %p (value: %d)\n", i, ptr, *ptr)
		}
	}
	runtime.KeepAlive(arr)
	fmt.Println()
	
	fmt.Println("=== Demo 5: Heap Slice ===")
	slice := heapSliceDemo()
	fmt.Printf("Heap slice address: %p\n", &slice)
	fmt.Printf("Heap slice len: %d, cap: %d\n", len(slice), cap(slice))
	for i, ptr := range slice {
		if ptr != nil {
			fmt.Printf("  slice[%d] = %p (value: %d)\n", i, ptr, *ptr)
		}
	}
	runtime.KeepAlive(slice)
	fmt.Println()
	
	fmt.Println("=== Demo 6: Heap Map ===")
	m := heapMapDemo()
	fmt.Printf("Heap map address: %p\n", &m)
	for k, v := range m {
		fmt.Printf("  map[%d] = %p (value: %d)\n", k, v, *v)
	}
	runtime.KeepAlive(m)
	fmt.Println()
	
	// 打印一些运行时信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Heap allocations: %d\n", memStats.Mallocs)
	fmt.Printf("Heap size: %d bytes\n", memStats.Alloc)
	fmt.Printf("Total GC cycles: %d\n", memStats.NumGC)
}
