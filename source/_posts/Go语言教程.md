---
type: Post
title: Go语言教程
tags: Go
category: 开发
category_bar: true
abbrlink: 41611
date: 2024-01-29 23:23:24
---

## 1.Go**语言的出现**

Go语言最初由Google公司的Robert Griesemer、Ken Thompson和Rob Pike三个大牛于2007年开始设计发明，他们最终的目标是设计一种适应网络和多核时代的C语言。但是Go语言更是对C语言最彻底的一次扬弃，它舍弃了C语言中灵活但是危险的指针运算，还重新设计了C语言中部分不太合理运算符的优先级，并在很多细微的地方都做了必要的打磨和改变。

## 2.**第一个Go程序**

接下来我们来编写第一个Go程序hello.go（Go语言源文件的扩展是.go），代码如下：

```go
package main
import "fmt"
func main() {
    //输出hello，world
    fmt.Println("Hello, World!")
}
```

和C语言相似，go语言的基本组成有：

1. 第一行代码`package main`定义了包名。你必须在源文件中非注释的第一行指明这个文件属于哪个包，如：`package main`。`package main`表示一个可独立执行的程序，每个 Go 应用程序都包含一个名为`main`的包。
2. 下一行`import "fmt"`告诉Go编译器这个程序需要使用fmt包（的函数，或其他元素），fmt包实现了格式化IO（输入/输出）的函数。
3. 下一行`func main()`是程序开始执行的函数。main函数是每一个可执行程序所必须包含的，一般来说都是在启动后第一个执行的函数（如果有`init()`函数则会先执行该函数）。
4. 下一行`/*...*/`是注释，在程序执行时将被忽略。单行注释是最常见的注释形式，你可以在任何地方使用以`//`开头的单行注释。多行注释也叫块注释，均已以`/*`开头，并以`*/`结尾，且不可以嵌套使用，多行注释一般用于包的文档描述或注释成块的代码片段。

    ```go
    // 单行注释
    /*
    我是多行注释
    */
    ```

5. 下一行`fmt.Println(...)`可以将字符串输出到控制台，并在最后自动增加换行字符\n。

    使用`fmt.Print("hello, world\n")`可以得到相同的结果。

    Print和Println这两个函数也支持使用变量，如：fmt.Println(arr)。如果没有特别指定，它们会以默认的打印格式将变量arr输出到控制台。

6. 当标识符（包括常量、变量、类型、函数名、结构字段等等）以一个大写字母开头，如：Group1，那么使用这种形式的标识符的对象就可以被外部包的代码所使用（客户端程序需要先导入这个包），这被称为导出（像面向对象语言中的 public）；标识符如果以小写字母开头，则对包外是不可见的，但是他们在整个包的内部是可见并且可用的（像面向对象语言中的 private ）。

## 3.数据结构

| 类型 | 说明 |
| --- | --- |
| 布尔型 | 布尔型的值只可以是常量 true 或者 false |
| 数字类型 | 整型int和浮点型float。支持整型和浮点型数字，并且支持复数，其中位的运算采用补码 |
| 字符串类型 | Go的字符串是由单个字节连接起来的。Go语言的字符串的字节使用UTF-8编码标识unicode文本 |
| 派生类型 | (a) 指针类型(Pointer) <br> (b) 数组类型 <br> (c) 结构体类型(struct) <br> (d) 联合体类型 (union) <br> (e) 函数类型 <br> (f) 切片类型 <br> (g) 接口类型(interface) <br> (h) Map 类型 <br> (i) Channel 类型 |

### 3.1变量

声明变量的一般形式是使用 var 关键字，具体格式为：`var identifier typename`。

#### 3.1.1如果变量没有初始化

在go语言中定义了一个变量，指定变量类型，如果没有初始化，则变量默认为零值。**零值就是变量没有做初始化时系统默认设置的值**。

| 类型 | 零值 |
| --- | --- |
| 数值类型 | 0 |
| 布尔类型 | false |
| 字符串 | ""（空字符串） |

#### 3.1.2如果变量没有指定类型

在go语言中如果没有指定变量类型，可以通过变量的初始值来判断变量类型。例如：

```go
package main
import "fmt"
func main() {
    var d = true
    fmt.Println(d)
}
```

#### 3.1.3 :=符号定义变量

当我们定义一个变量后又使用该符号初始化变量，就会产生编译错误，因为该符号其实是一个声明语句。使用格式：`typename := value`。类如：

```go
var a int = 10
var b = 10
c := 10
```

#### 3.1.4多变量声明

可以同时声明多个类型相同的变量（非全局变量），例如：

```go
var x, y int = 1, 2
var c, d  = 1, "hello"
g, h := 123, "hello"
```

关于全局变量的声明如下：`var ( vname1 v_type1 vname2 v_type2 )`，例如：

```go
var (
    a int
    b bool
    g, h := 123, "hello"
)
```

#### 3.1.5匿名变量

匿名变量的特点是一个下画线`_`，这本身就是一个特殊的标识符，被称为空白标识符。它可以像其他标识符那样用于变量的声明或赋值（任何类型都可以赋值给它），但任何赋给这个标识符的值都将被抛弃，因此这些值**不能在后续的代码中使用**，也不可以使用这个标识符作为变量对其它变量进行赋值或运算。使用匿名变量时，只需要在**变量声明的地方**使用下画线替换即可。例如：

```go
func GetData() (int, int) {
    return 10, 20
}

func main(){
    a, _ := GetData()
    _, b := GetData()
    fmt.Println(a, b)
}
```

需要注意的是匿名变量不占用内存空间，不会分配内存。匿名变量与匿名变量之间也不会因为多次声明而无法使用。

**如果你想要交换两个变量的值，则可以简单地使用`a, b = b, a`。**

### 3.2指针

变量其实是一种使用方便的占位符，用于引用计算机内存地址。Go 语言中的的取地址符是`&`，放到一个变量前使用就会返回相应变量的内存地址。指针变量其实就是用于存放某一个对象的内存地址。

和基础类型数据相同，在使用指针变量之前我们首先需要申明指针，声明格式如下：`var var_name *var-type`，其中的var-type 为指针类型，var_name 为指针变量名，* 号用于指定变量是作为一个指针。例如：

```go
var ip *int        /* 指向整型*/
var fp *float32    /* 指向浮点型 */
```

指针的初始化就是取出相对应的变量地址对指针进行赋值，具体如下：

```go
var a int= 20   /* 声明实际变量 */
var ip *int        /* 声明指针变量 */

ip = &a  /* 指针变量的存储地址 */
```

当一个指针被定义后**没有分配到任何变量**时，它的值为 **nil**，也称为空指针。它概念上和其它语言的null、NULL一样，都指代零值或空值。

### 3.3常量

常量是一个简单值的标识符，在程序运行时，不会被修改的量。

常量中的数据类型只可以是布尔型、数字型（整数型、浮点型和复数）和字符串型。

常量的定义格式：`const identifier [type] = value`

你可以省略类型说明符 [type]，因为编译器可以根据变量的值来推断其类型。

- 显式类型定义： `const b string = "abc"`
- 隐式类型定义： `const b = "abc"`

多个相同类型的声明可以简写为：`const c_name1, c_name2 = value1, value2`

#### 3.3.1特殊常量iota

iota，特殊常量，可以认为是一个可以被编译器修改的常量。在每一个const关键字出现时，被重置为0，然后再下一个const出现之前，每出现一次iota，其所代表的数字会自动增加1。iota可以被用作枚举值：

```go
const (
    a = iota
    b = iota
    c = iota
)
```

第一个 iota 等于 0，每当 iota 在新的一行被使用时，它的值都会自动加 1；所以 a=0, b=1, c=2 可以简写为如下形式：

```go
const (
    a = iota
    b
    c
)
```

#### 3.3.2iota用法

```go
package main
import "fmt"
func main() {
    const (
            a = iota   //0
            b          //1
            c          //2
            d = "ha"   //独立值，iota += 1
            e          //"ha"   iota += 1
            f = 100    //iota +=1
            g          //100  iota +=1
            h = iota   //7,恢复计数
            i          //8
    )
    fmt.Println(a,b,c,d,e,f,g,h,i)
}
```

以上实例运行结果为：

```go
0 1 2 ha ha 100 100 7 8
```

再看个有趣的的 iota 实例：

```go
package main
import "fmt"

const (
    i=1<<iota
    j=3<<iota
    k
    l
)

func main() {
    fmt.Println("i=",i)
    fmt.Println("j=",j)
    fmt.Println("k=",k)
    fmt.Println("l=",l)
}
```

以上实例运行结果为：

```go
i= 1
j= 6
k= 12
l= 24
```

iota表示从0开始自动加1，所以i=1<<0,j=3<<1（<<表示左移的意思），即：i=1,j=6，这没问题，关键在k和l，从输出结果看，k=3<<2，l=3<<3。

简单表述:

- i=1：左移 0 位，不变仍为 1。
- j=3：左移 1 位，变为[二进制](https://www.coonote.com/linux-note/binary.html) 110，即 6。
- k=3：左移 2 位，变为二进制 1100，即 12。
- l=3：左移 3 位，变为二进制 11000，即 24。

### 3.4数组

#### 3.4.1声明数组

Go 语言数组声明需要指定元素类型及元素个数，语法格式如下：`var variable_name [SIZE] variable_type`。以上就可以定一个一维数组。

声明多为数组的格式为:`var variable_name [SIZE1][SIZE2]...[SIZEN] variable_type`。

#### 3.4.2初始化数组

1. 直接进行初始化：`var balance = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}`
2. 通过字面量在声明数组的同时快速初始化数组：`balance := [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}`
3. 数组长度不确定，编译器通过元素个数自行推断数组长度，在[ ]中填入`...`，举例如下：`var balance = [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}`和`balance := [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}`
4. 数组长度确定，指定下标进行部分初始化：`balanced := [5]float32(1:2.0, 3:7.0)`

注意：如果忽略 [] 中的数字不设置数组大小，Go 语言会根据元素的个数来设置数组的大小。

#### 3.4.3数组指针

go中一个数组变量被赋值或者被传递的时候实际上就会复制整个数组。如果数组比较大的话，这种复制往往会占有很大的开销。所以为了避免这种开销，往往需要传递一个指向数组的指针，这个数组指针并不是数组。例如：

```go
var a = [...]int{1, 2, 3} // a 是一个数组
var b = &a                // b 是指向数组的指针
```

数组指针除了可以防止数组作为参数传递的时候浪费空间，还可以利用其和`for range`来遍历数组，具体代码如下：

```go
for i, v := range b {
    // 通过数组指针迭代数组的元素
    fmt.Println(i, v)   //i表示索引，v代表具体数值
}
```

### 3.5结构体

#### **3.5.1 声明结构体**

在声明结构体之前我们首先需要定义一个结构体类型，这需要使用type和struct，type用于设定结构体的名称，struct用于定义一个新的数据类型。具体结构如下：

```go
type struct_variable_type struct {
    member definition
    member definition
    ...
    member definition
}
```

#### **3.5.2 访问结构体成员**

如果要访问结构体成员，需要使用点号`.`操作符，格式为：`结构体变量名.成员名`。举例代码如下：

```go
package main

import "fmt"

type Books struct {
    title string
    author string
}

func main() {
    var book1 Books
    Book1.title = "Go 语言入门"
    Book1.author = "mars.hao"
}
```

#### **3.5.3 结构体指针**

关于结构体指针的定义和申明同样可以套用前文中讲到的指针的相关定义，从而使用一个指针变量存放一个结构体变量的地址。定义一个结构体变量的语法：`var struct_pointer *Books`。

这种指针变量的初始化和上文指针部分的初始化方式相同`struct_pointer = &Book1`，和c语言中有所不同，使用结构体指针访问结构体成员仍然使用`.`操作符。格式如下`struct_pointer.title`。

### 3.6**切片(Slice)**

简单地说，切片就是一种简化版的**动态数组**。因为动态数组的**长度不固定**，切片的长度自然也就不能是类型的组成部分了。数组虽然有适用它们的地方，但是数组的类型和操作都不够灵活，而切片则使用得相当广泛。切片高效操作的要点是要降低内存分配的次数，尽量保证append操作（在后续的插入和删除操作中都涉及到这个函数）不会超出cap的容量，降低触发内存分配的次数和每次分配内存大小。

和数组一样，内置的len函数返回切片中有效元素的长度，内置的cap函数返回切片容量大小，容量必须大于或等于切片的长度。

#### 3.6.1定义切片

可以声明一个未指定大小的数组来定义切片：`var identifier []type`，切片不需要说明长度。

或使用make()函数来创建切片：`var slice1 []type = make([]type, len)`，也可以简写为`slice1 := make([]type, len)`。也可以指定容量，其中capacity为可选参数，`make([]T, length, capacity)`，这里len是数组的长度并且也是切片的初始长度。

#### 3.6.2**初始化切片**

```go
s :=[] int {1,2,3 }
```

直接初始化切片，[]表示是切片类型，{1,2,3}初始化值依次是1，2，3。其cap=len=3，cap表示切片最大长度。

```go
s := arr[:]
```

初始化切片s,是数组arr的引用

```go
s := arr[startIndex:endIndex]
```

将arr中从下标startIndex到endIndex-1 下的元素创建为一个新的切片，缺省endIndex时将表示一直到arr的最后一个元素，缺省startIndex时将表示从arr的第一个元素开始。

```go
s1 := s[startIndex:endIndex]
```

通过切片s初始化切片s1

```go
s :=make([]int,len,cap)
```

通过内置函数make()初始化切片s,[]int 标识为其元素类型为int的切片

## 4.循环结构

### 4.1循环语句

go语言中的循环语句只有for循环

下面是for循环的三种形式：

第一种：和C语言的for循环一样

```go
for init; condition; post { }
```

第二种：和C语言的while循环一样

```go
for condition { }
```

当陷入无限循环是，要停止无限循环，可以在命令窗口按下ctrl-c

第三种：和C语言的for(;;)一样

```go
for { }
```

- init： 一般为赋值表达式，给控制变量赋初值；
- condition： 关系表达式或逻辑表达式，循环控制条件；
- post： 一般为赋值表达式，给控制变量增量或减量。

for 循环的 range 格式可以对 slice、map、数组、字符串等进行迭代循环。格式如下：

```go
for key, value := range oldMap {
    newMap[key] = value
}
```

以上代码中的key和value是可以省略。后者用_来实现只需要key或者value的场景。

### 4.2循环控制语句

- break
- continue
- goto

下面对goto语句进行说明：

```go
package main

import "fmt"

func main() {
    var a int = 10
    LOOP: for a < 20 {
        if a == 15 {
            /* 跳过迭代 */
            a = a + 1
            goto LOOP
        }
        fmt.Printf("a的值为 : %d\n", a)
        a++    
    }  
}
```

输出：

```go
a的值为 : 10
a的值为 : 11
a的值为 : 12
a的值为 : 13
a的值为 : 14
a的值为 : 16
a的值为 : 17
a的值为 : 18
a的值为 : 19
```

## 5函数

### 5.1函数定义

Go 语言函数定义格式如下：

```go
func function_name( [parameter list] ) [return_types] {
    函数体
}
```

### 5.2函数返回多个值

go语言的函数和python相似，可以返回多个值。例如：

```go
package main

import "fmt"

func swap(x, y string) (string, string) {
    return y, x
}

func main() {
    a, b := swap("Google", "Edge")
    fmt.Println(a, b)
}
```

输出：

```go
Edge Google
```

### 5.3函数闭包（匿名函数）

匿名函数的优越性在于可以直接使用函数内的变量，不必申明。匿名函数是一种没有函数名的函数，通常用于在函数内部定义函数，或者作为函数参数进行传递。

```go
package main

import "fmt"

func getSequence() func() int {
    i:=0
    return func() int {
        i+=1
        return i  
    }
}

func main(){
    /* nextNumber 为一个函数，函数 i 为 0 */
    nextNumber := getSequence()  

    /* 调用 nextNumber 函数，i 变量自增 1 并返回 */
    fmt.Println(nextNumber())
    fmt.Println(nextNumber())
    fmt.Println(nextNumber())
   
    /* 创建新的函数 nextNumber1，并查看结果 */
    nextNumber1 := getSequence()  
    fmt.Println(nextNumber1())
    fmt.Println(nextNumber1())
}
```

输出

```go
1
2
3
1
2
```

### 5.4函数方法

一个方法就是一个包含了接受者的函数，接受者可以是命名类型或者结构体类型的一个值或者是一个指针。所有给定类型的方法属于该类型的方法集。语法格式如下：

```go
func (variable_name variable_data_type) function_name() [return_type]{
    /* 函数体*/
}
```

例如：

```go
package main

import (
    "fmt"  
)

/* 定义结构体 */
type Circle struct {
    radius float64
}

func main() {
    var c1 Circle
    c1.radius = 10.00
    fmt.Println("圆的面积 = ", c1.getArea())
}

//该 method 属于 Circle 类型对象中的方法
func (c Circle) getArea() float64 {
    //c.radius 即为 Circle 类型对象中的属性
    return 3.14 * c.radius * c.radius
}
```

这个方法会在学习接口时进一步深入理解。

## 6类型转换

与python类似，go语言的类型转换用于将一种数据类型的变量转换为另外一种类型的变量。基本语法如下：

```go
type_name(expression)
```

### 6.1数值类型转换

#### 6.1.1将整型转换为浮点型

```go
var a int = 10
var b float64 = float64(a)
```

#### 6.1.2将浮点型转换为整型：与python类似，小数点后之间舍去，不管是多少

```go
var a float64 = 10.6
var b int = int(a)
```

### 6.2字符串类型转换

#### 6.2.1将字符串类型转化为整型

```go
var str string = "10"
var num int
num, _ = strconv.Atoi(str)
```

注意：**strconv.Atoi**函数返回两个值，第一个是转换后的整型值，第二个是可能发生的错误，我们可以使用空白标识符”_“来忽略这个错误。

#### 6.2.2将整型转换为字符串类型

```go
num := 123
str := strconv.Itoa(num)
```

#### 6.2.3将浮点型转换为字符串类型

```go
str := "3.14"
num, err := strconv.ParseFloat(str, 64)
```

### 6.3接口类型转换

#### 6.3.1类型断言

类型断言用于将接口类型转换为指定类型，其语法为：

```go
value.(type) 
```

其中value是接口类型的变量，type或T是要转换成的类型。

如果类型断言成功，它将返回转换后的值和一个布尔值，表示转换是否成功。

例如：

```go
package main

import "fmt"

func main() {
    var i interface{} = "Hello, World"   //定义了一个接口类型变量i
    str, ok := i.(string)
    if ok {
        fmt.Printf("'%s' is a string\n", str)
    } else {
        fmt.Println("conversion failed")
    }
}
```

我们定义了一个接口类型变量i，并将它赋值为字符串 "Hello, World"。然后，我们使用类型断言将i转换为字符串类型，并将转换后的值赋值给变量 str。最后，我们使用ok变量检查类型转换是否成功，如果成功，我们打印转换后的字符串；否则，我们打印转换失败的消息。

#### 6.3.2类型转换

语法：

```go
T(value)
```

T是目标接口类型，value是要转换的值。

在类型转换中，我们必须保证要转换的值和目标接口类型之间是兼容的，否则编译器会报错。

```go
package main

import "fmt"

type Writer interface {
    Write([]byte) (int, error)
}

type StringWriter struct {
    str string
}

func (sw *StringWriter) Write(data []byte) (int, error) {
    sw.str += string(data)
    return len(data), nil
}

func main() {
    var w Writer = &StringWriter{}
    sw := w.(*StringWriter)
    sw.str = "Hello, World"
    fmt.Println(sw.str)
}
```

以上实例中，我们定义了一个Writer接口和一个实现了该接口的结构体StringWriter。然后，我们将StringWriter类型的指针赋值给Writer接口类型的变量w。接着，**我们使用类型转换将w转换为StringWriter类型，并将转换后的值赋值给变量sw**。最后，我们使用sw访问StringWriter 结构体中的字段str，并打印出它的值。
