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

/*
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

*/

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
