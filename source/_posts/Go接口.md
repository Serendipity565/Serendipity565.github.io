---
type: Post
title: Go接口
tags: Go
category: 开发
category_bar: true
abbrlink: 48051
date: 2024-03-14 13:12:12
---

## 什么是接口

在goland中，接口是一组方法签名。接口把所有的具有共性的方法定义在一起，任何其他类型只要实现了这些方法就是实现了这个接口。它与OOP(面向对象编程)非常相似。

接口可以让我们将不同的类型绑定到一组公共的方法上，从而实现多态和灵活的设计。

Go语言中的接口是隐式实现的，也就是说，如果一个类型实现了一个接口定义的所有方法，那么它就自动地实现了该接口。因此，我们可以通过将接口作为参数来实现对不同类型的调用，从而实现多态。

## 接口的定义及实现

### 定于接口

关键字interface用来定义接口，语法如下：

```go
type interface_name interface {
    method_name1([args ...arg_type]) [return_type]
    method_name2([args ...arg_type]) [return_type]
    method_name3([args ...arg_type]) [return_type]
    ...
    method_namen([args ...arg_type]) [return_type]
}
```

一个接口中可以定义多个方法，根据逻辑需要，自定义参数和返回值。

### 实现接口

如果一个结构体实现了某个接口的所有方法，则此结构体就实现了该接口。

我们定义一个sleeper接口，包含两个方法，即sleep()和eat()：

```go
type animal interface {
    sleep()
    eat()
}
```

声明结构体，用来实现接口：

```go
type dog struct {
    name string
    food string
}
```

声明了dog这一结构体，包含name和food两个字段。

接下来实现animal接口，我们先让dog结构体实现该接口的所有方法：

```go
func (mydog dog) sleep() {
    fmt.Printf("%s is sleeping\n", mydog.name)
}
func (mydog dog) eat() {
    fmt.Printf("%s is eating %s\n", mydog.name, mydog.food)
}
```

接下来初始化结构体，输出测试：

```go
func main() {
    xiaobai := dog{
        name: "xiaobai",
        food: "bone",
    }
    xiaobai.sleep()
    xiaobai.eat()
}

```

输出：

```go
xiaobai is sleeping
xiaobai is eating bone
```

dog结构体实现了animal接口的所有方法，因此我们认为dog实现了animal接口。

#### 没有实现接口会怎么样？

我们先定义一个结构体cat，cat只实现animal的sleep()方法，我们认为cat并没有实现animal接口：

```go
func (mycat cat) sleep() {
    fmt.Printf("%s is sleeping \n", mycat.name)
}

func main() {
    kitty := cat{
        name: "kitty",
        food: "fish",
    }
    kitty.sleep()
}
```

输出

```go
kitty is sleeping
```

输出结果显示sleep()方法并没有什么问题。

但是使用接口的方法并不只有一种，我们换一种方法来观察接口有没有实现：

```go
func main() {
    var ani animal
    kitty := cat{
        name: "kitty",
        food: "fish",
    }

    ani = kitty

    ani.sleep()
}
```

结果报错，在终端里提示`cannot use kitty (variable of type cat) as animal value in assignment: cat does not implement animal (missing method eat)`。编写代码时编译器也会出现红色下划线的错误提示，如图：

![](/img/blog/jieko/1.png)

## 接口嵌套

接口嵌套就是一个接口中包含了其他接口，如果要实现外部接口，那么就要把内部嵌套的接口对应的所有方法全实现了。

```go
type A interface {
    test1()
}

type B interface {
    test2()
}

// 定义嵌套接口
type C interface {
    A
    B
    test3()
}
```

如果想实现接口C，那不止要实现接口C的方法，还要实现接口A，B中的方法

## 空接口

空接口不包含任何的方法，正因为如此，所有的类型都实现了空接口，因此空接口可以存储任意类型的数值。

fmt包下的Print系列函数，其参数大多是空接口类型，也可以说支持任意类型：

```go
func Print(a ...interface{}) (n int, err error)
func Println(format string, a ...interface{}) (n int, err error)
func Println(a ...interface{}) (n int, err error)
```

### 示例

在需要存储不同类型数据的容器中，可以使用空接口作为容器的元素类型。这样就可以在容器中存储任意类型的值。

```go
type empty_interface interface {
}

func example(empty empty_interface) {
    fmt.Printf("example", empty)
}

func main() {
    var data []interface{}
    data = append(data, 42)
    data = append(data, "hello")

    fmt.Println("data...........", data)
}

```

data是一个空接口类型的切片，可以存储整数、字符串等不同类型的数据。

下面附上源代码：

[test.go](/img/blog/jieko/test.go)
