---
type: Post
title: Go语言内存对齐
tags: Go
category: 开发
category_bar: true
abbrlink: 28613
date: 2025-05-09 17:02:34
---

## **什么是内存对齐？**

内存对齐是一个源于CPU访问内存方式的概念。现代CPU优化为在对齐的地址上访问内存——即地址是数据大小的倍数。

例如：

- `int64`（8字节）理想情况下应该从8的倍数的内存地址开始
- `int32`（4字节）应该从4的倍数地址开始

如果变量未对齐，CPU可能需要执行多次内存读取才能获取完整数据。这会降低速度。此外，如果变量跨越两个缓存行，你将受到性能惩罚，因为CPU必须加载两个缓存行。

这就如同阅读一个分散在书中两页的句子。你翻一次页，然后再翻一次，仅仅为了获取完整信息。对齐可以让你的"句子"保持在同一页上。

简而言之：

- 对齐的数据 = 快速内存访问
- 未对齐的数据 = 慢速，可能需要多次读取

## 问题引出

思考下面这段代码的输出

```go
type S1 struct {
	num2 int8
	num1 int16
	flag bool
}

type S2 struct {
	num1 int8
	flag bool
	num2 int16
}

func main() {
	fmt.Println(unsafe.Sizeof(S1{}))
	fmt.Println(unsafe.Sizeof(S2{}))
}
```

输出结果：

```go
6
4
```

为什么仅是字段顺序不同，`S1{}` 和 `S2{}` 的大小就不一样了？

根据理论`s1`的内存结构如下：


![](/img/blog/ncdq/3.png)

如果没有内存对齐呢？`s1`的结构可能如下：


![](/img/blog/ncdq/4.png)

如果是 16 位系统的话，那么没有内存对齐的情况下，要访问 `s1.num2` 字段，就需要跨过 2 个系统字长的内存，效率就低了。具体来说，内存对齐是计算机内存分配的一种优化方式，用于确保数据结构的存储按照特定的字节边界对齐。这种对齐是为了提高计算机处理数据的效率。

## **对齐系数**

- 对齐系数：变量的内存地址必须被对齐系数整除。
- `unsafe.Alignof()`: 可以查看值在内存中的对齐系数。

### **基本类型对齐**

```go
fmt.Printf("bool size: %d, align: %d\n", unsafe.Sizeof(bool(true)), unsafe.Alignof(bool(true)))
fmt.Printf("byte size: %d, align: %d\n", unsafe.Sizeof(byte(0)), unsafe.Alignof(byte(0)))
fmt.Printf("int8 size: %d, align: %d\n", unsafe.Sizeof(int8(0)), unsafe.Alignof(int8(0)))
fmt.Printf("int16 size: %d, align: %d\n", unsafe.Sizeof(int16(0)), unsafe.Alignof(int16(0)))
fmt.Printf("int32 size: %d, align: %d\n", unsafe.Sizeof(int32(0)), unsafe.Alignof(int32(0)))
fmt.Printf("int64 size: %d, align: %d\n", unsafe.Sizeof(int64(0)), unsafe.Alignof(int64(0)))
```

输出：

```go
bool size: 1, align: 1
byte size: 1, align: 1
int8 size: 1, align: 1
int16 size: 2, align: 2
int32 size: 4, align: 4
int64 size: 8, align: 8
```

结论：基本类型的对齐系数跟它的长度一致。


![](/img/blog/ncdq/1.png)

## **结构体内部对齐**

结构体内存对齐分为内部对齐和结构体之间对齐。

我们先来看结构体内部对齐：

- 指的是结构体内部成员的相对位置（偏移量）；
- 每个成员的偏移量是 **自身大小** 和 **对齐系数** 的较小值的倍数

```go
type Demo struct {
  a bool
  b string
  c int16
}
```

假如我们定义了上面的结构体 `Demo`，如果在 64 位系统上（字长为 8 字节）通过上面的规则，可以判断出：（单位为字节）

- a: size=1, align=1
- b: size=16, align=8
- c: size=2, align=2

![](/img/blog/ncdq/2.png)

### **结构体长度填充**

上面 Demo 结构体最后还填了 6 个字节的 0，这就是结构体长度填充：

- 结构体通过填充长度，来对齐系统字长。
- 结构体长度是 **最大成员长度** 和 **系统字长** 较小值的整数倍。

## **结构体之间对齐**

- 结构体之间对齐，是为了确定结构体的第一个成员变量的内存地址，以让后面的成员地址都合法。
- 结构体的对齐系数是 **其成员的最大对齐系数**；

## **空结构体对齐**

空结构体 `struct{}`，它们的内存地址统一指向 `zerobase`，而且内存长度为 0。这也导致了它的内存对齐规则，有一些不同。具体可以分为以下 4 个情况。

### **空结构体单独存在**

空结构体单独存在时，其内存地址为 `zerobase`，不额外分配内存。

```go
package main

import (
	"fmt"
	"unsafe"
)

type TestEmpty struct {
	empty struct{}
}

func main() {
	te := TestEmpty{}
	fmt.Println("size of TestEmpty:", unsafe.Sizeof(te))
	fmt.Printf("address of te: %p\n", &te)
	fmt.Printf("address of te.empty: %p\n", &(te.empty))
	fmt.Printf("empty: size=%d, align=%d\n", unsafe.Sizeof(te.empty), unsafe.Alignof(te.empty))
}
```

输出：

```go
size of TestEmpty: 0
address of te: 0x102750360
address of te.empty: 0x102750360
empty: size=0, align=1
```

### **空结构体在结构体最前**

空结构体是结构体第一个字段时，它的地址跟结构体本身及结构体第 2 个字段一样，不占据内存空间。

```go
package main

import (
	"fmt"
	"unsafe"
)

type TestEmpty struct {
	empty struct{}
	a     bool
	b     string
}

func main() {
	te := TestEmpty{}
	fmt.Println("size of TestEmpty:", unsafe.Sizeof(te))
	fmt.Printf("address of te: %p\n", &te)
	fmt.Printf("address of te.empty: %p\n", &(te.empty))
	fmt.Printf("address of te.a: %p\n", &(te.a))
	fmt.Printf("address of te.b: %p\n", &(te.b))
	fmt.Printf("empty: size=%d, align=%d\n", unsafe.Sizeof(te.empty), unsafe.Alignof(te.empty))
	fmt.Printf("a: size=%d, align=%d\n", unsafe.Sizeof(te.a), unsafe.Alignof(te.a))
	fmt.Printf("b: size=%d, align=%d\n", unsafe.Sizeof(te.b), unsafe.Alignof(te.b))
}
```

输出：

```go
size of TestEmpty: 24
address of te: 0x14000136000
address of te.empty: 0x14000136000
address of te.a: 0x14000136000
address of te.b: 0x14000136008
empty: size=0, align=1
a: size=1, align=1
b: size=16, align=8
```

### **空结构体在结构体中间**

空结构体出现在结构体中时，地址跟随前一个变量。

```go
package main

import (
	"fmt"
	"unsafe"
)

type TestEmpty struct {
	a     bool
	empty struct{}
	b     string
}

func main() {
	te := TestEmpty{}
	fmt.Println("size of TestEmpty:", unsafe.Sizeof(te))
	fmt.Printf("address of te: %p\n", &te)
	fmt.Printf("address of te.a: %p\n", &(te.a))
	fmt.Printf("address of te.empty: %p\n", &(te.empty))
	fmt.Printf("address of te.b: %p\n", &(te.b))
	fmt.Printf("empty: size=%d, align=%d\n", unsafe.Sizeof(te.empty), unsafe.Alignof(te.empty))
	fmt.Printf("a: size=%d, align=%d\n", unsafe.Sizeof(te.a), unsafe.Alignof(te.a))
	fmt.Printf("b: size=%d, align=%d\n", unsafe.Sizeof(te.b), unsafe.Alignof(te.b))
}

```

输出：

```go
size of TestEmpty: 24
address of te: 0x1400012a000
address of te.a: 0x1400012a000
address of te.empty: 0x1400012a001
address of te.b: 0x1400012a008
empty: size=0, align=1
a: size=1, align=1
b: size=16, align=8

```

### **空结构体在结构体最后**

空结构体出现在结构体最后，如果开启了一个新的系统字长，则需要补零，防止与其他结构体混用地址。

```go
package main

import (
	"fmt"
	"unsafe"
)

type TestEmpty struct {
	a     bool
	b     string
	empty struct{}
}

func main() {
	te := TestEmpty{}
	fmt.Println("size of TestEmpty:", unsafe.Sizeof(te))
	fmt.Printf("address of te: %p\n", &te)
	fmt.Printf("address of te.a: %p\n", &(te.a))
	fmt.Printf("address of te.b: %p\n", &(te.b))
	fmt.Printf("address of te.empty: %p\n", &(te.empty))
	fmt.Printf("empty: size=%d, align=%d\n", unsafe.Sizeof(te.empty), unsafe.Alignof(te.empty))
	fmt.Printf("a: size=%d, align=%d\n", unsafe.Sizeof(te.a), unsafe.Alignof(te.a))
	fmt.Printf("b: size=%d, align=%d\n", unsafe.Sizeof(te.b), unsafe.Alignof(te.b))
}
```

输出：

```go
size of TestEmpty: 32
address of te: 0x14000120040
address of te.a: 0x14000120040
address of te.b: 0x14000120048
address of te.empty: 0x14000120058
empty: size=0, align=1
a: size=1, align=1
b: size=16, align=8
```

## **使用 fieldalignment -fix 工具优化结构体内存对齐**

还记得我们最开始提出的问题吗？

```go
type S1 struct {
	num2 int8
	num1 int16
	flag bool
}

type S2 struct {
	num1 int8
	flag bool
	num2 int16
}

func main() {
	fmt.Println(unsafe.Sizeof(S1{}))
	fmt.Println(unsafe.Sizeof(S2{}))
}
```

`S1` 和 `S2` 提供的程序功能是一样的，但是 `S1` 却比 `S2` 花费了更多的内存空间。所以有时候我们可以通过仅仅调整结构体内部字段的顺序就减少不少的内存空间消耗。在这个时候 `fieldalignment` 可以帮助我们自动检测并优化。

你可以运行下面命令安装 `fieldalignment` 命令：

```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
```

然后在项目根目录下运行下面命令，对我们的代码进行检查：

```bash
go vet -vettool=$(which fieldalignment) ./...
```

这里会输出：

```bash
./main.go:9:9: struct of size 6 could be 4
```

这个时候可以执行 `fieldalignment -fix 目录|文件` ，它会自动帮我们的代码进行修复，但是**强烈建议你在运行之前，备份你的代码，因为注释会被删除！**

```bash
fieldalignment -fix ./...
```

这个时候 `S1` 已经被优化好了：

```go
type S1 struct {
	num1 int16
	num2 int8
	flag bool
}
```

## **Go结构体布局的最佳实践**

1. 按从最大到最小的对齐顺序排列字段
2. 将相同大小的字段分组在一起
3. 在定义高容量或性能关键的结构体时考虑内存布局
4. 使用`go vet -fieldalignment`自动建议

---

**参考资料**

- [**Go优化指南**](https://goperf.dev/01-common-patterns/fields-alignment)