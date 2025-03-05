---
type: Post
title: Go接口断言
tags: Go
category: 开发
category_bar: true
abbrlink: 64382
date: 2024-03-15 23:23:42
---

空接口interface{}没有定义任何函数，因此Golang中所有类型都实现了空接口。当一个函数的形参interface{}，那么在函数中，需要对形参进行断言，从而得到它的真实类型。

在学习接口断言之前，先了解一下类型断言，其实接口断言也是在判断类型。

## 类型断言

类型断言可以检查i是否为nil或者为某个类型，通常有两中方式

第一种：

```go
t := i.(T)
```

这个表达式可以断言一个接口对象i里不是nil，并且接口对象i存储的值的类型是T，如果断言成功，就会返回值给t，如果断言失败，就会触发 panic。这个方式常用于 switch 结构。

第二种：

```go
t, ok:= i.(T)
```

这个表达式也是可以断言一个接口对象t里不是nil，并且接口对象t存储的值的类型是T，如果断言成功，就会返回其类型给t，并且此时ok的值为true，表示断言成功。

如果接口值的类型，并不是我们所断言的T，就会断言失败，但和第一种表达式不同的事，这个不会触发panic，而是将ok的值设为false，表示断言失败，此时t为T的零值。这几种方式常用于if else结构。

## 接口断言

### if else结构接口断言

```go
package main

import (
    "fmt"
    "math"
)

// 定义了一个shape接口
type shape interface {
    area() float64
    peri() float64
}

// 定义一个圆形结构体
type circle struct {
    r float64
}

// 定义一个三角形结构体
type triangle struct {
    a float64
    b float64
    c float64
}

// 圆形接口实现放法
func (c circle) area() float64 {
    return c.r * c.r * math.Pi
}
func (c circle) peri() float64 {
    return 2 * c.r * math.Pi
}

// 三角形接口实现方法
func (t triangle) peri() float64 {
    return t.a + t.b + t.c
}
func (t triangle) area() float64 {
    p := (t.a + t.b + t.c) / 2
    return math.Sqrt(p * (p - t.a) * (p - t.b) * (p - t.c))
}

// 定义断言函数
func check(s shape) {
    if i, ok := s.(triangle); ok {
        fmt.Printf("是三角形，三条边分别为：%f %f %f \n", i.a, i.b, i.c)
    } else if i, ok := s.(circle); ok {
        fmt.Printf("是圆形，半径为：%f \n", i.r)
    } else if i, ok := s.(*circle); ok {
        fmt.Printf("是圆形结构体指针，类型为：%T,存储的地址为：%p，指针自身的地址为：%p\n", i, &i, i)
    } else {
        fmt.Printf("无法判断类型")
    }
}

func main() {
    //初始化圆形结构体
    c1 := circle{r: 5.0}
    fmt.Printf("\n")
    fmt.Printf("下面是圆形结构体：\n")
    fmt.Printf("圆的周长为：%f \n", c1.peri())
    fmt.Printf("圆的面积为：%f \n", c1.area())

    //初始化三角形结构体
    t1 := triangle{a: 2, b: 3, c: 4}
    fmt.Printf("\n")
    fmt.Printf("下面是三角形结构体：\n")
    fmt.Printf("三角形的周长为：%f \n", t1.peri())
    fmt.Printf("三角形的面积为: %f \n", t1.area())

    //圆形结构体指针
    var c2 *circle = &circle{r: 10.0}
    fmt.Printf("\n")
    fmt.Printf("下面是圆形结构体指针：\n")
    fmt.Printf("圆的周长为：%f \n", c2.peri())
    fmt.Printf("圆的面积为：%f \n", c2.area())

    //开始接口断言
    fmt.Printf("\n")
    fmt.Printf("开始接口断言：\n")
    fmt.Printf("c1是")
    check(c1)
    fmt.Printf("t1是")
    check(t1)
    fmt.Printf("c2是")
    check(c2)
}

```

输出：

```go
下面是圆形结构体：
圆的周长为：31.415927 
圆的面积为：78.539816 

下面是三角形结构体：
三角形的周长为：9.000000
三角形的面积为: 2.904738

下面是圆形结构体指针：
圆的周长为：62.831853
圆的面积为：314.159265

开始接口断言：
c1是是圆形，半径为：5.000000
t1是是三角形，三条边分别为：2.000000 3.000000 4.000000
c2是是圆形结构体指针，类型为：*main.circle,存储的地址为：0xc000122020，指针自身的地址为：0xc0001100b8
```

可以看到，我们的断言成功了。

### switch结构接口断言

断言其实还有另一种形式，就是用在利用switch语句判断接口的类型。

每一个case会被顺序地考虑。当命中一个case时，就会执行case中的语句，因此case语句的顺序是很重要的，因为很有可能会有多个case匹配的情况。

我们再封装一个switch逻辑的接口断言函数，逻辑和之前的一模一样，只是条件语句换成了switch....case：

```go
// 使用 switch定义接口断言函数
func check(s shape) {
    switch i := s.(type) {
    case circle:
        fmt.Printf("是圆形，半径为：%f \n", i.r)
    case triangle:
        fmt.Printf("是三角形，三条边分别为：%f %f %f \n", i.a, i.b, i.c)
    case *circle:
        fmt.Printf("是圆形结构体指针，类型为：%T,存储的地址为：%p，指针自身的地址为：%p\n", i, &i, i)
    default:
        fmt.Printf("无法判断类型")
    }
}
```

输出：

```go
开始接口断言：
c1是是圆形，半径为：5.000000
t1是是三角形，三条边分别为：2.000000 3.000000 4.000000
c2是是圆形结构体指针，类型为：*main.circle,存储的地址为：0xc00008a028，指针自身的地址为：0xc00000a118
```

可以看到switch断言也正常输出了。

下面附上源代码：

[test.go](/img/blog/jieko/check/test.go)
