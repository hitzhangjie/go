# Go精确GC实现机制Demo

本目录包含用于演示Go语言精确GC实现机制的示例代码和分析文档。

## 文件说明

### 源代码文件

1. **simple_demo.go** - 最简单的示例，用于生成汇编代码分析
2. **demo1_stack_allocation.go** - 栈分配示例
3. **demo2_heap_allocation.go** - 堆分配示例
4. **demo3_mixed_allocation.go** - 混合分配示例
5. **demo4_gc_scanning.go** - GC扫描示例

### 文档

- **how_precise_gc_works.md** - 详细的分析报告，深入解析Go精确GC的实现机制

## 使用方法

### 1. 查看汇编代码

```bash
# 使用Go工具链编译并查看汇编代码
cd /path/to/go
./bin/go tool compile -S -o /tmp/demo.o demo/how_precise_gc_works/simple_demo.go

# 或者使用系统Go工具链（如果已安装）
go tool compile -S simple_demo.go
```

### 2. 运行Demo程序

注意：由于这些demo使用了标准库（fmt、runtime），需要确保Go环境正确配置。

```bash
# 设置GOROOT（如果使用源码构建的Go）
export GOROOT=/path/to/go
export PATH=$GOROOT/bin:$PATH

# 运行demo
go run simple_demo.go
go run demo1_stack_allocation.go
go run demo2_heap_allocation.go
go run demo3_mixed_allocation.go
go run demo4_gc_scanning.go
```

### 3. 分析汇编代码中的PCDATA和FUNCDATA

在生成的汇编代码中，查找以下模式：

- `PCDATA $0, $value` - 不安全点标记
- `PCDATA $1, $value` - 栈映射表索引
- `FUNCDATA $0, symbol` - 参数指针映射表
- `FUNCDATA $1, symbol` - 局部变量指针映射表

### 4. 查看栈映射表

栈映射表以符号形式出现在汇编代码中，例如：
- `gclocals·g2BeySu+wFnoycgXfElmcg==` - 空映射表（无指针）
- `gclocals·TjPuuCwdlCpTaRQGRKTrYw==` - 包含指针的映射表

## 关键概念

### PCDATA指令

- **PCDATA $0**：不安全点标记（UnsafePoint）
  - `$-1`：安全点
  - `$-2`：不安全点
- **PCDATA $1**：栈映射表索引（StackMapIndex）
  - 值对应栈映射表中的索引

### FUNCDATA指令

- **FUNCDATA $0**：参数指针映射表（ArgsPointerMaps）
- **FUNCDATA $1**：局部变量指针映射表（LocalsPointerMaps）

### 栈映射表格式

栈映射表是一个位图，每一位表示栈帧中一个指针大小的位置是否包含指针：
- `1`：该位置包含指针
- `0`：该位置不包含指针

## 注意事项

1. **逃逸分析**：Go编译器会进行逃逸分析，某些看似栈分配的变量可能被优化或逃逸到堆上
2. **编译器优化**：不同优化级别可能生成不同的代码
3. **架构差异**：不同CPU架构（x86_64、ARM64等）生成的汇编代码不同
4. **Go版本**：不同Go版本的实现可能有所差异

## 进一步学习

1. 阅读 `how_precise_gc_works.md` 了解详细实现机制
2. 查看Go源码：
   - `src/cmd/compile/internal/liveness/` - 活跃性分析
   - `src/runtime/mgcmark.go` - GC标记
   - `src/runtime/symtab.go` - pclntab使用
3. 使用工具：
   - `go tool objdump` - 反汇编二进制文件
   - `go tool compile -S` - 查看汇编代码
   - `go build -gcflags="-m"` - 查看逃逸分析结果

## 示例输出

运行`simple_demo.go`的汇编代码示例：

```assembly
main.stackAlloc STEXT size=16 args=0x0 locals=0x0 funcid=0x0 align=0x0 leaf
    0x0000 00000 (simple_demo.go:7)	TEXT	main.stackAlloc(SB), LEAF|NOFRAME|ABIInternal, $0-0
    0x0000 00000 (simple_demo.go:7)	FUNCDATA	$0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
    0x0000 00000 (simple_demo.go:7)	FUNCDATA	$1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
    0x0000 00000 (simple_demo.go:11)	RET	(R30)
```

这展示了：
- FUNCDATA指令的位置（在函数开始处）
- 参数和局部变量映射表都指向同一个符号（表示无指针）
- 这是一个LEAF函数（叶子函数，不调用其他函数）
