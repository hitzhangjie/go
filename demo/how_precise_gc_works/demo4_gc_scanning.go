// demo4_gc_scanning.go
// 这个demo展示GC如何扫描栈和堆中的指针
// 通过创建复杂的指针关系来展示精确GC的工作方式

package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

// Node 表示一个链表节点，包含指针字段
type Node struct {
	Value int64
	Next  *Node
	Prev  *Node
}

//go:noinline
func createStackList() *Node {
	// 在栈上创建节点
	var n1, n2, n3 Node
	
	n1.Value = 1
	n2.Value = 2
	n3.Value = 3
	
	// 建立指针关系
	n1.Next = &n2
	n2.Next = &n3
	n2.Prev = &n1
	n3.Prev = &n2
	
	// 返回第一个节点的指针（但节点本身在栈上）
	// 注意：这里返回栈上变量的地址，在实际使用中可能会导致问题
	// 但我们可以用它来演示GC如何扫描栈上的指针
	return &n1
}

//go:noinline
func createHeapList() *Node {
	// 在堆上创建节点
	n1 := &Node{Value: 10}
	n2 := &Node{Value: 20}
	n3 := &Node{Value: 30}
	
	// 建立指针关系
	n1.Next = n2
	n2.Next = n3
	n2.Prev = n1
	n3.Prev = n2
	
	return n1
}

//go:noinline
func createMixedList() (*Node, *Node) {
	// 混合：栈上的节点指向堆上的节点
	var stackNode Node
	stackNode.Value = 100
	
	heapNode1 := &Node{Value: 200}
	heapNode2 := &Node{Value: 300}
	
	stackNode.Next = heapNode1
	heapNode1.Next = heapNode2
	heapNode1.Prev = &stackNode
	heapNode2.Prev = heapNode1
	
	return &stackNode, heapNode1
}

//go:noinline
func complexPointerStructure() {
	// 创建一个复杂的指针结构
	type Container struct {
		Ptr1 *int64
		Ptr2 *int64
		Ptr3 *Container
	}
	
	// 栈上的容器
	var stackContainer Container
	val1 := int64(1000)
	val2 := int64(2000)
	stackContainer.Ptr1 = &val1
	stackContainer.Ptr2 = &val2
	
	// 堆上的容器
	heapContainer := &Container{}
	val3 := int64(3000)
	val4 := int64(4000)
	heapContainer.Ptr1 = &val3
	heapContainer.Ptr2 = &val4
	
	// 建立交叉引用
	stackContainer.Ptr3 = heapContainer
	heapContainer.Ptr3 = &stackContainer
	
	fmt.Printf("Stack container at %p:\n", &stackContainer)
	fmt.Printf("  Ptr1: %p -> %d\n", stackContainer.Ptr1, *stackContainer.Ptr1)
	fmt.Printf("  Ptr2: %p -> %d\n", stackContainer.Ptr2, *stackContainer.Ptr2)
	fmt.Printf("  Ptr3: %p\n", stackContainer.Ptr3)
	
	fmt.Printf("Heap container at %p:\n", heapContainer)
	fmt.Printf("  Ptr1: %p -> %d\n", heapContainer.Ptr1, *heapContainer.Ptr1)
	fmt.Printf("  Ptr2: %p -> %d\n", heapContainer.Ptr2, *heapContainer.Ptr2)
	fmt.Printf("  Ptr3: %p\n", heapContainer.Ptr3)
	
	runtime.KeepAlive(&stackContainer)
	runtime.KeepAlive(heapContainer)
}

//go:noinline
func demonstrateStackMap() {
	// 这个函数用于展示栈映射表
	// 包含多个局部变量，有些是指针，有些不是
	
	var nonPtr1 int64 = 1
	var nonPtr2 int64 = 2
	var ptr1 *int64 = &nonPtr1
	var ptr2 *int64 = &nonPtr2
	var nonPtr3 int64 = 3
	
	// 创建一个包含指针的数组
	var arr [4]*int64
	arr[0] = ptr1
	arr[1] = ptr2
	arr[2] = nil
	arr[3] = &nonPtr3
	
	fmt.Printf("Stack variables:\n")
	fmt.Printf("  nonPtr1: %d at %p\n", nonPtr1, &nonPtr1)
	fmt.Printf("  nonPtr2: %d at %p\n", nonPtr2, &nonPtr2)
	fmt.Printf("  ptr1: %p -> %d\n", ptr1, *ptr1)
	fmt.Printf("  ptr2: %p -> %d\n", ptr2, *ptr2)
	fmt.Printf("  nonPtr3: %d at %p\n", nonPtr3, &nonPtr3)
	fmt.Printf("  arr: %p\n", &arr)
	
	// GC需要知道arr[0], arr[1], arr[3]是指针，而arr[2]是nil指针
	// 这些信息存储在栈映射表中
	
	runtime.KeepAlive(nonPtr1)
	runtime.KeepAlive(nonPtr2)
	runtime.KeepAlive(ptr1)
	runtime.KeepAlive(ptr2)
	runtime.KeepAlive(nonPtr3)
	runtime.KeepAlive(arr)
}

func main() {
	fmt.Println("=== Demo 1: Stack List ===")
	stackList := createStackList()
	fmt.Printf("Stack list head: %p\n", stackList)
	if stackList != nil {
		fmt.Printf("  Value: %d\n", stackList.Value)
		if stackList.Next != nil {
			fmt.Printf("  Next: %p (Value: %d)\n", stackList.Next, stackList.Next.Value)
		}
	}
	runtime.KeepAlive(stackList)
	fmt.Println()
	
	fmt.Println("=== Demo 2: Heap List ===")
	heapList := createHeapList()
	fmt.Printf("Heap list head: %p\n", heapList)
	current := heapList
	for i := 0; i < 3 && current != nil; i++ {
		fmt.Printf("  Node %d: %p (Value: %d)\n", i, current, current.Value)
		current = current.Next
	}
	runtime.KeepAlive(heapList)
	fmt.Println()
	
	fmt.Println("=== Demo 3: Mixed List ===")
	stackNode, heapNode := createMixedList()
	fmt.Printf("Stack node: %p (Value: %d)\n", stackNode, stackNode.Value)
	fmt.Printf("Heap node: %p (Value: %d)\n", heapNode, heapNode.Value)
	runtime.KeepAlive(stackNode)
	runtime.KeepAlive(heapNode)
	fmt.Println()
	
	fmt.Println("=== Demo 4: Complex Pointer Structure ===")
	complexPointerStructure()
	fmt.Println()
	
	fmt.Println("=== Demo 5: Stack Map Demonstration ===")
	demonstrateStackMap()
	fmt.Println()
	
	// 强制GC运行
	fmt.Println("=== Running GC ===")
	runtime.GC()
	
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("After GC:\n")
	fmt.Printf("  Heap allocations: %d\n", memStats.Mallocs)
	fmt.Printf("  Heap size: %d bytes\n", memStats.Alloc)
	fmt.Printf("  Total GC cycles: %d\n", memStats.NumGC)
	fmt.Printf("  GC pause total: %v\n", memStats.PauseTotalNs)
}
