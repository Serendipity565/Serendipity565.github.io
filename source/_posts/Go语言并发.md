---
type: Post
date: 2024/03/30
title: Go语言并发
tags: Go
category: 开发
category_bar: true
abbrlink: 30761
---

## 基本概念

### 串行、并发与并行

串行：依次执行多个任务。

并行：同一时刻执行多个任务。

并发：同一时间段内执行多个任务。

### **进程、线程和协程**

进程（process）：程序在操作系统中的一次执行过程，系统进行资源分配和调度的一个独立单位。

线程（thread）：操作系统基于进程开启的轻量级进程，是操作系统调度执行的最小单位。

协程（coroutine）：非操作系统提供而是由用户自行创建和控制的用户态”线程“，比线程更轻量级。

## Goroutine

Goroutine 是 Go 程序中最基本的并发执行单元。每一个 Go 程序都至少包含一个 goroutine——main goroutine，当 Go 程序启动时它会自动创建。

goroutine 是由Go运行时负责调度。例如Go运行时会智能地将 m个goroutine 合理地分配给n个操作系统线程，实现类似m:n的调度机制，不再需要我们自行在代码层面维护一个线程池。

在Go语言编程中你不需要去自己写进程、线程、协程，当你需要让某个任务并发执行的时候，你只需要把这个任务包装成一个函数，开启一个 goroutine 去执行这个函数就可以了。

### go关键字

Go语言中使用 goroutine 非常简单，只需要在函数或方法调用前加上go关键字就可以创建一个 goroutine ，从而让该函数或方法在新创建的 goroutine 中执行。

### 启动单个Goroutine

例如：

```go
package main

import (
    "fmt"
)

func hello() {
    fmt.Println("hello")
}

func main() {
    go hello()
    fmt.Println("你好")
}
```

输出：

```go
你好
```

行结果只在终端打印了"你好"，并没有打印 hello。这是为什么呢？

其实在 Go 程序启动时，Go 程序就会为 main 函数创建一个默认的 goroutine 。在上面的代码中我们在 main 函数中使用 go 关键字创建了另外一个 goroutine 去执行 hello 函数，而此时 main goroutine 还在继续往下执行，我们的程序中此时存在两个并发执行的 goroutine。当 main 函数结束时整个程序也就结束了，同时 main goroutine 也结束了，所有由 main goroutine 创建的 goroutine 也会一同退出。也就是说我们的 main 函数退出太快，另外一个 goroutine 中的函数还未执行完程序就退出了，导致未打印出“hello”。

所以我们要想办法让 main 函数等一等将在另一个 goroutine 中运行的 hello( ) 函数。其中最简单的方式就是在 main 函数中加入 time.Sleep 了（这里的1秒钟是我们根据经验而设置的一个值，在这个示例中1秒钟足够创建新的 goroutine 执行完 hello( ) 函数了）。

修改主函数：

```go
func main() {
    go hello()
    fmt.Println("你好")
    time.Sleep(time.Second)
}
```

得到结果：

```go
你好
hello
```

为什么会先打印”你好“呢？

这是因为在程序中创建 goroutine 执行函数需要一定的时间，而与此同时 main 函数所在的 goroutine 是继续执行的。

上面程序使用的 `time.Sleep(time.Second)` 虽然可以完成实现上面的功能，但无法满足更多的使用场景。Go 语言中通过sync包为我们提供了一些常用的并发原语，当你并不关心并发操作的结果或者有其它方式收集并发操作的结果时，`WaitGroup`是实现等待一组并发操作完成的好方法。

```go
package main

import (
    "fmt"
    "sync"
)

var wg sync.WaitGroup

func hello() {
    fmt.Println("hello")
    wg.Done() // 告知当前goroutine完成
}

func main() {
    wg.Add(1) // 登记1个goroutine
    go hello()
    fmt.Println("你好")
    wg.Wait() // 阻塞等待登记的goroutine完成
}
```

### 启动多个Goroutine

```go
package main

import (
    "fmt"
    "sync"
)

var wg sync.WaitGroup

func hello(i int) {
    defer wg.Done()
    fmt.Println("hello", i)
}

func main() {
    for i := 0; i < 10; i++ {
        wg.Add(1) // 启动一个goroutine就登记+1
        go hello(i)
    }
    wg.Wait() // 等待所有登记的goroutine都结束
}

```

多次执行上面的代码会发现每次终端上打印数字的顺序都不一致。这是因为10个 goroutine 是并发执行的，而 goroutine 的调度是随机的。

### **动态栈**

操作系统的线程一般都有固定的栈内存（通常为2MB）,而 Go 语言中的 goroutine 非常轻量级，一个 goroutine 的初始栈空间很小（一般为2KB），所以在 Go 语言中一次创建数万个 goroutine 也是可能的。并且 goroutine 的栈不是固定的，可以根据需要动态地增大或缩小， Go 的 runtime 会自动为 goroutine 分配合适的栈空间。

### **goroutine调度**

操作系统内核在调度时会挂起当前正在执行的线程并将寄存器中的内容保存到内存中，然后选出接下来要执行的线程并从内存中恢复该线程的寄存器信息，然后恢复执行该线程的现场并开始执行线程。从一个线程切换到另一个线程需要完整的上下文切换。因为可能需要多次内存访问，索引这个切换上下文的操作开销较大，会增加运行的cpu周期。

区别于操作系统内核调度操作系统线程，goroutine 的调度是Go语言运行时（runtime）层面的实现，是完全由 Go 语言本身实现的一套调度系统——go scheduler。它的作用是按照一定的规则将所有的 goroutine 调度到操作系统线程上执行。

在经历数个版本的迭代之后，目前 Go 语言的调度器采用的是 GPM 调度模型。

![](/img/blog/Gobf/1.png)

说明：

- G：表示 goroutine，包含要执行的函数和上下文信息。
- 全局队列（Global Queue）：存放等待运行的 G。
- P：表示 goroutine 执行所需的资源，最多有 GOMAXPROCS 个。
- GOMAXPROCS默认值是机器上的 CPU 核心数。可以通过runtime.GOMAXPROCS函数设置当前程序并发时占用的 CPU逻辑核心数。
- P 的本地队列：同全局队列类似，存放的也是等待运行的G，存的数量有限，不超过256个。新建 G 时，G 优先加入到 P 的本地队列，如果本地队列满了会批量移动部分 G 到全局队列。
- M：线程想运行任务就得获取 P，从 P 的本地队列获取 G，当 P 的本地队列为空时，M 也会尝试从全局队列或其他 P 的本地队列获取 G。M 运行 G，G 执行之后，M 会从 P 获取下一个 G，不断重复下去。
- Goroutine 调度器和操作系统调度器是通过 M 结合起来的，每个 M 都代表了1个内核线程，操作系统调度器负责把内核线程分配到 CPU 的核上执行。

## **Channel**

单纯地将函数并发执行是没有意义的。函数与函数间需要交换数据才能体现并发执行函数的意义。

如果说 goroutine 是Go程序并发的执行体，channel就是它们之间的连接。channel是可以让一个 goroutine 发送特定值到另一个 goroutine 的通信机制。

Go 语言中的通道（channel）是一种特殊的类型。通道像一个传送带或者队列，总是遵循先入先出的规则，保证收发数据的顺序。每一个通道都是一个具体类型的导管，也就是声明channel的时候需要为其指定元素类型。

### channel类型

```go
var 变量名称 chan 元素类型
```

### channel零值

未初始化的通道类型变量其默认零值是nil。

```go
var ch chan int
fmt.Println(ch) //输出：<nil>
```

### 初始化channel

```go
make(chan 元素类型, [缓冲大小])
```

### channel操作

通道共有发送（send）、接收（receive）和关闭（close）三种操作。而发送和接收操作都使用`<-`符号。

#### 发送

将一个值发送到通道中。

```go
ch <- 10  // 把10发送到ch中
```

#### 接收

从一个通道中接收值。

```go
x := <- ch  // 第一种方式，从ch中接收值并赋值给变量x
<-ch        // 第二种方式，从ch中接收值，忽略结果
```

#### 关闭

我们通过调用内置的close函数来关闭通道。

```go
close(ch)
```

通道通常由发送方执行关闭操作，并且只有在接收方明确等待通道关闭的信号时才需要执行关闭操作。它和关闭文件不一样，通常在结束操作之后关闭文件是必须要做的，但关闭通道不是必须的。

关闭后的通道有以下特点：

1. 对一个关闭的通道再发送值就会导致 panic。
2. 对一个关闭的通道进行接收会一直获取值直到通道为空。
3. 对一个关闭的并且没有值的通道执行接收操作会得到对应类型的零值。
4. 关闭一个已经关闭的通道会导致 panic。

### 无缓冲的通道

无缓冲的通道又称为阻塞的通道。

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int)
    ch <- 10
    fmt.Println("发送成功")
}
```

结果：报错提示deadlock，即死锁。

![](/img/blog/Gobf/2.png)

我们使用 ch := make(chan int) 创建的是无缓冲的通道，无缓冲的通道只有在有接收方能够接收值的时候才能发送成功，否则会一直处于等待发送的阶段。同理，如果对一个无缓冲通道执行接收操作时，没有任何向通道中发送值的操作那么也会导致接收操作阻塞。

我们看可以通过创建一个 goroutine 去接收值来解决这个问题，例如：

```go
package main

import (
    "fmt"
    "sync"
)

var wg sync.WaitGroup

func recvive(c chan int) {
    defer wg.Done()
    ans := <-c
    fmt.Println("接收成功", ans)
}

func main() {
    ch := make(chan int)
    wg.Add(1)
    go recvive(ch)
    ch <- 10
    fmt.Println("发送成功")
    wg.Wait()
}
```

首先无缓冲通道ch上的发送操作会阻塞，直到另一个 goroutine 在该通道上执行接收操作，这时数字10才能发送成功，两个 goroutine 将继续执行。相反，如果接收操作先执行，接收方所在的 goroutine 将阻塞，直到 main goroutine 中向该通道发送数字10。

使用无缓冲通道进行通信将导致发送和接收的 goroutine 同步化。因此，无缓冲通道也被称为同步通道。

### 有缓存通道

还有另外一种解决上面死锁问题的方法，那就是使用有缓冲区的通道。我们可以在使用 make 函数初始化通道时，可以为其指定通道的容量，例如：

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int, 1)
    ch <- 10
    fmt.Println("发送成功")
}
```

只要通道的容量大于零，那么该通道就属于有缓冲的通道，通道的容量表示通道中最大能存放的元素数量。当通道内已有元素数达到最大容量后，再向通道执行发送操作就会阻塞，除非有从通道执行接收操作。

我们可以使用内置的len函数获取通道内元素的数量，使用cap函数获取通道的容量。

总结一下对通道进行操作的几种结果：

![](/img/blog/Gobf/3.png)

### 多返回值模式

当向通道中发送完数据时，我们可以通过 close( ) 函数来关闭通道。当一个通道被关闭后，再往该通道发送值会引发panic，从该通道取值的操作会先取完通道中的值。通道内的值被接收完后再对通道执行接收操作得到的值会一直都是对应元素类型的零值。那我们如何判断一个通道是否被关闭了呢？

对一个通道执行接收操作时支持使用如下多返回值模式。

```go
value, ok := <- ch
```

其中：

- value：从通道中取出的值，如果通道被关闭则返回对应类型的零值。
- ok：通道ch关闭时返回 false，否则返回 true。

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int, 2)
    ch <- 10
    ch <- 20
    close(ch)
    fmt.Println("发送成功")
    for {
        ans, ok := <-ch
        if ok != false {
            fmt.Printf("ans is %d\n", ans)
        } else {
            fmt.Println("通道已关闭")
            break
        }
    }
}
```

### for range接收值

通常我们会选择使用for range循环从通道中接收值，当通道被关闭后，会在通道内的所有值被接收完毕后会自动退出循环。上面那个示例我们使用for range改写后会很简洁。

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int, 2)
    ch <- 10
    ch <- 20
    close(ch)
    fmt.Println("发送成功")
    for v := range ch {
        fmt.Printf("ans is %d\n", v)
    }
}
```

注意：不能简单的通过len(ch)操作来判断通道是否被关闭。

### **单向通道**

```go
<- chan int// 只接收通道，只能接收不能发送
chan<- int// 只发送通道，只能发送不能接收
```

```go
package main

import (
    "fmt"
    "sync"
)

var wg sync.WaitGroup

func sent() <-chan int {
    ch := make(chan int, 2)
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i
        }
        close(ch)
    }()
    return ch
}

func receive(ch <-chan int) int {
    sum := 0
    wg.Add(1)
    go func() {
        for {
            v, ok := <-ch
            if ok != false {
                sum += v
            } else {
                break
            }
        }
        wg.Done()
    }()
    wg.Wait()
    return sum
}

func main() {
    ch1 := sent()
    ans := receive(ch1)

    fmt.Println(ans)
}
```

这就从代码层面限制了该函数返回的通道只能进行接收操作，保证了数据安全。函数可以在其他地方被其他人调用时进行发送数据而产生问题。

在函数传参及任何赋值操作中全向通道（正常通道）可以转换为单向通道，但是无法反向转换。

## select多路复用

在某些场景下我们可能需要同时从多个通道接收数据。通道在接收数据时，如果没有数据可以被接收那么当前 goroutine 将会发生阻塞。我们可以尝试使用遍历的方式来实现从多个通道中接收值。这种方式虽然可以实现从多个通道接收值的需求，但是程序的运行性能会差很多。

Go 语言内置了select关键字，使用它可以同时响应多个通道的操作。

Select 的使用方式类似于之前学到的 switch 语句，它也有一系列 case 分支和一个默认的分支。每个 case 分支会对应一个通道的通信（接收或发送）过程。select 会一直等待，直到其中的某个 case 的通信操作完成时，就会执行该 case 分支对应的语句。具体格式如下：

```go
select {
case <-ch1:
    //...
case data := <-ch2:
    //...
case ch3 <- 10:
    //...
default:
    //默认操作
}
```

Select 语句具有以下特点。

- 可处理一个或多个 channel 的发送/接收操作。
- 如果多个 case 同时满足，select 会随机选择一个执行。
- 对于没有 case 的 select 会一直阻塞，可用于阻塞 main 函数，防止退出。

下面的示例代码能够在终端打印出10以内的偶数，我们借助这个代码片段来看一下 select 的具体使用。

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 1)
    for i := 0; i <= 10; i++ {
        select {
        case x := <-ch:
            fmt.Println(x)
        case ch <- i:
        }
    }
}
```

示例中的代码首先是创建了一个缓冲区大小为1的通道 ch，进入 for 循环后：

- 第一次循环时 i = 0，select 语句中包含两个 case 分支，此时由于通道中没有值可以接收，所以x := <-ch 这个 case 分支不满足，而ch <- i这个分支可以执行，会把1发送到通道中，结束本次 for 循环；
- 第二次 for 循环时，i = 1，由于通道缓冲区已满，所以ch <- i这个分支不满足，而x := <-ch这个分支可以执行，从通道接收值1并赋值给变量 x ，所以会在终端打印出 0；

## 并发安全和锁

有时候我们的代码中可能会存在多个 goroutine 同时操作一个资源（临界区）的情况，这种情况下就会发生竞态问题（数据竞态）。

```go
package main

import (
    "fmt"
    "sync"
)

var (
    x int64
    wg sync.WaitGroup
)

func add() {
    for i := 0; i < 5000; i++ {
        x = x + 1
    }
    wg.Done()
}

func main() {
    wg.Add(2)
    go add()
    go add()
    wg.Wait()
    fmt.Println(x)
}
```

多次执行，发现输出如8088、7818、10000、8146等不同结果。原因是这两个 goroutine 在访问和修改全局变量 x 时就会存在数据竞争，某个 goroutine 中对全局变量 x 的修改可能会覆盖掉另一个 goroutine 中的操作，所以导致最后的结果与预期不符。

### 互斥锁

互斥锁是一种常用的控制共享资源访问的方法，它能够保证同一时间只有一个 goroutine 可以访问共享资源。Go 语言中使用sync包中提供的Mutex类型来实现互斥锁。`sync.Mutex` 提供了两个方法供我们使用：

| 方法 | 功能 |
| --- | --- |
| func (m *Mutex) Lock() | 获取互斥锁 |
| func (m *Mutex) Unlock() | 释放互斥锁 |

下面我们来解决上面那个代码出现的问题：

```go
package main

import (
    "fmt"
    "sync"
)

var (
    x  int64
    wg sync.WaitGroup
    m  sync.Mutex
)

func add() {
    for i := 0; i < 5000; i++ {
        m.Lock() //修改前加锁
        x = x + 1
        m.Unlock() //修改后解锁
    }
    wg.Done()
}

func main() {
    wg.Add(2)
    go add()
    go add()
    wg.Wait()
    fmt.Println(x)
}
```

使用互斥锁能够保证同一时间有且只有一个 goroutine 进入临界区，其他的 goroutine 则在等待锁；当互斥锁释放后，等待的 goroutine 才可以获取锁进入临界区，多个 goroutine 同时等待一个锁时，唤醒的策略是随机的。

### 读写互斥锁

互斥锁是完全互斥的，但是实际上有很多场景是读多写少的，当我们并发的去读取一个资源而不涉及资源修改的时候是没有必要加互斥锁的，这种场景下使用读写锁是更好的一种选择。读写锁在 Go 语言中使用 sync 包中的 RWMutex 类型。

`sync.RWMutex` 提供了以下5个方法。

| 方法名 | 功能 |
| --- | --- |
| func (rw *RWMutex) Lock() | 获取写锁 |
| func (rw *RWMutex) Unlock() | 释放写锁 |
| func (rw *RWMutex) RLock() | 获取读锁 |
| func (rw *RWMutex) RUnlock() | 释放读锁 |
| func (rw *RWMutex) RLocker() Locker | 返回一个实现Locker接口的读写锁 |

读写锁分为两种：读锁和写锁。当一个 goroutine 获取到读锁之后，其他的 goroutine 如果是获取读锁会继续获得锁，如果是获取写锁就会等待；而当一个 goroutine 获取写锁之后，其他的 goroutine 无论是获取读锁还是写锁都会等待。

下面我们使用代码构造一个读多写少的场景，然后分别使用互斥锁和读写锁查看它们的性能差异。

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

var (
    x       int64
    wg      sync.WaitGroup
    mutex   sync.Mutex
    rwMutex sync.RWMutex
)

// writeWithLock 使用互斥锁的写操作
func writeWithLock() {
    mutex.Lock() // 加互斥锁
    x = x + 1
    time.Sleep(10 * time.Millisecond) // 假设读操作耗时10毫秒
    mutex.Unlock()                    // 解互斥锁
    wg.Done()
}

// readWithLock 使用互斥锁的读操作
func readWithLock() {
    mutex.Lock()                 // 加互斥锁
    time.Sleep(time.Millisecond) // 假设读操作耗时1毫秒
    mutex.Unlock()               // 释放互斥锁
    wg.Done()
}

// writeWithLock 使用读写互斥锁的写操作
func writeWithRWLock() {
    rwMutex.Lock() // 加写锁
    x = x + 1
    time.Sleep(10 * time.Millisecond) // 假设读操作耗时10毫秒
    rwMutex.Unlock()                  // 释放写锁
    wg.Done()
}

// readWithRWLock 使用读写互斥锁的读操作
func readWithRWLock() {
    rwMutex.RLock()              // 加读锁
    time.Sleep(time.Millisecond) // 假设读操作耗时1毫秒
    rwMutex.RUnlock()            // 释放读锁
    wg.Done()
}

func do(wf, rf func(), wc, rc int) {
    start := time.Now()
    // wc个并发写操作
    for i := 0; i < wc; i++ {
        wg.Add(1)
        go wf()
    }

    //  rc个并发读操作
    for i := 0; i < rc; i++ {
        wg.Add(1)
        go rf()
    }

    wg.Wait()
    cost := time.Since(start)
    fmt.Printf("x:%v cost:%v\n", x, cost)
}

func main() {
    // 使用互斥锁，10并发写，1000并发读
    do(writeWithLock, readWithLock, 10, 1000)

    // 使用读写互斥锁，10并发写，1000并发读
    do(writeWithRWLock, readWithRWLock, 10, 1000)
}
```

输出：

![](/img/blog/Gobf/4.png)

从结果可以看出，使用读写互斥锁在读多写少的场景下能够极大地提高程序的性能。但是如果程序中的读操作和写操作数量级差别不大，那么读写互斥锁的优势就发挥不出来。

### sync.WaitGroup

Go语言中可以使用`sync.WaitGroup`来实现并发任务的同步。这个在前面我们已近提到过并简单使用过，下面就来加单介绍一下。

| 方法名 | 功能 |
| --- | --- |
| func (wg * WaitGroup) Add(delta int) | 计数器+delta |
| (wg *WaitGroup) Done() | 计数器-1 |
| (wg *WaitGroup) Wait() | 阻塞直到计数器变为0 |

`sync.WaitGroup`内部维护着一个计数器，计数器的值可以增加和减少。例如当我们启动了 N 个并发任务时，就将计数器值增加N。每个任务完成时通过调用 Done 方法将计数器减1。通过调用 Wait 来等待并发任务执行完，当计数器值为 0 时，表示所有并发任务已经完成。

需要注意`sync.WaitGroup`是一个结构体，进行参数传递的时候要传递指针。

### sync.Once

在某些场景下我们需要确保某些操作即使在高并发的场景下也只会被执行一次，例如只加载一次配置文件等。

Go语言中的sync包中提供了一个针对只执行一次场景的解决方案——`sync.Once`，`sync.Once`只有一个Do方法，其签名如下：

```go
func (o *Once) Do(f func())
```

- 注意：如果要执行的函数 f 需要传递参数就需要搭配闭包来使用。

#### 加载配置文件示例

延迟一个开销很大的初始化操作到真正用到它的时候再执行是一个很好的实践。因为预先初始化一个变量（比如在init函数中完成初始化）会增加程序的启动耗时，而且有可能实际执行过程中这个变量没有用上，那么这个初始化操作就不是必须要做的。我们来看一个例子：

```go
var icons map[string]image.Image

func loadIcons() {
    icons = map[string]image.Image{
        "left":  loadIcon("left.png"),
        "up":    loadIcon("up.png"),
        "right": loadIcon("right.png"),
        "down":  loadIcon("down.png"),
    }
}

func Icon(name string) image.Image {
    if icons == nil {
        loadIcons()
    }
    return icons[name]
}
```

多个 goroutine 并发调用Icon函数时不是并发安全的，编译器和CPU可能会在保证每个 goroutine 都满足串行一致的基础上自由地重排访问内存的顺序。loadIcons函数可能会被重排为以下结果：

```go
func loadIcons() {
    icons = make(map[string]image.Image)
    icons["left"] = loadIcon("left.png")
    icons["up"] = loadIcon("up.png")
    icons["right"] = loadIcon("right.png")
    icons["down"] = loadIcon("down.png")
}
```

在这种情况下就会出现即使判断了 icons 不是nil也不意味着变量初始化完成了。考虑到这种情况，我们能想到的办法就是添加互斥锁，保证初始化 icons 的时候不会被其他的 goroutine 操作，但是这样做又可能会引发性能问题（时间过长）。

使用`sync.Once`改造的示例代码如下：

```go
var icons map[string]image.Image

var loadIconsOnce sync.Once

func loadIcons() {
    icons = map[string]image.Image{
        "left":  loadIcon("left.png"),
        "up":    loadIcon("up.png"),
        "right": loadIcon("right.png"),
        "down":  loadIcon("down.png"),
    }
}

// Icon 是并发安全的
func Icon(name string) image.Image {
    loadIconsOnce.Do(loadIcons)
    return icons[name]
}
```

`sync.Once` 其实内部包含一个互斥锁和一个布尔值，互斥锁保证布尔值和数据的安全，而布尔值用来记录初始化是否完成。这样设计就能保证初始化操作的时候是并发安全的并且初始化操作也不会被执行多次。

### sync.Map

Go 语言中内置的 map 不是并发安全的，我们不能在多个 goroutine 中并发对内置的 map 进行读写操作，否则会存在数据竞争问题，编译时会报出`fatal error: concurrent map writes`错误。

Go语言的 sync 包中提供了一个开箱即用的并发安全版 map——`sync.Map` 。不用像内置的 map 一样使用 make 函数初始化就能直接使用。同时`sync.Map`内置了诸如`Store`、`Load`、`LoadOrStore`、`Delete`、`Range`等操作方法。

| 方法名 | 功能 |
| --- | --- |
| func (m *Map) Store(key, value interface{}) | 存储key-value数据 |
| func (m *Map) Load(key interface{}) (value interface{}, ok bool) | 查询key对应的value |
| func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) | 查询或存储key对应的value |
| func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool) | 查询并删除key |
| func (m *Map) Delete(key interface{}) | 删除key |
| func (m *Map) Range(f func(key, value interface{}) bool) | 对map中的每个key-value依次调用f |

例如：

```go
package main

import (
    "fmt"
    "strconv"
    "sync"
)

var m = sync.Map{}

func main() {
    wg := sync.WaitGroup{}
    // 对m执行20个并发的读写操作
    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func(n int) {
            key := strconv.Itoa(n)
            m.Store(key, n)
            value, _ := m.Load(key)
            fmt.Printf("k=:%v,v:=%v\n", key, value)
            wg.Done()
        }(i)
    }
    wg.Wait()
}
```

## 原子操作

针对整数数据类型（int32、uint32、int64、uint64）我们还可以使用原子操作来保证并发安全，通常直接使用原子操作比使用锁操作效率更高。Go语言中原子操作由内置的标准库 `sync/atomic` 提供。

### atomic包

| 方法 | 解释 |
| --- | --- |
| `func LoadInt32(addr *int32) (val int32)` <br> `func LoadInt64(addr *int64) (val int64)` <br> `func LoadUint32(addr *uint32) (val uint32)` <br> `func LoadUint64(addr *uint64) (val uint64)` <br> `func LoadUintptr(addr *uintptr) (val uintptr)` <br> `func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)` | 读取操作 |
| `func StoreInt32(addr *int32, val int32)` <br> `func StoreInt64(addr *int64, val int64)` <br> `func StoreUint32(addr *uint32, val uint32)` <br> `func StoreUint64(addr *uint64, val uint64)` <br> `func StoreUintptr(addr *uintptr, val uintptr)` <br> `func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)` | 写入操作 |
| `func AddInt32(addr *int32, delta int32) (new int32)` <br> `func AddInt64(addr *int64, delta int64) (new int64)` <br> `func AddUint32(addr *uint32, delta uint32) (new uint32)` <br> `func AddUint64(addr *uint64, delta uint64) (new uint64)` <br> `func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)` | 修改操作 |
| `func SwapInt32(addr *int32, new int32) (old int32)` <br> `func SwapInt64(addr *int64, new int64) (old int64)` <br> `func SwapUint32(addr *uint32, new uint32) (old uint32)` <br> `func SwapUint64(addr *uint64, new uint64) (old uint64)` <br> `func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)` <br> `func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)` | 交换操作 |
| `func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)` <br> `func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)` <br> `func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)` <br> `func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)` <br> `func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)` <br> `func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)` | 比较并交换操作 |

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

type Counter interface {
    Increase()
    Load() int64
}

// 普通版
type CommonCounter struct {
    counter int64
}

func (c CommonCounter) Increase() {
    c.counter++
}

func (c CommonCounter) Load() int64 {
    return c.counter
}

// 互斥锁版
type MutexCounter struct {
    counter int64
    lock    sync.Mutex
}

func (m *MutexCounter) Increase() {
    m.lock.Lock()
    defer m.lock.Unlock()
    m.counter++
}

func (m *MutexCounter) Load() int64 {
    m.lock.Lock()
    defer m.lock.Unlock()
    return m.counter
}

// 原子操作版
type AtomicCounter struct {
    counter int64
}

func (a *AtomicCounter) Increase() {
    atomic.AddInt64(&a.counter, 1)
}

func (a *AtomicCounter) Load() int64 {
    return atomic.LoadInt64(&a.counter)
}

func test(c Counter) {
    var wg sync.WaitGroup
    start := time.Now()
    for i := 0; i < 5000; i++ {
        wg.Add(1)
        go func() {
            c.Increase()
            wg.Done()
        }()
    }
    wg.Wait()
    end := time.Now()
    fmt.Println(c.Load(), end.Sub(start))
}

func main() {
    c1 := CommonCounter{} // 非并发安全
    test(c1)
    c2 := MutexCounter{} // 使用互斥锁实现并发安全
    test(&c2)
    c3 := AtomicCounter{} // 并发安全且比互斥锁效率更高
    test(&c3)
}
```

![](/img/blog/Gobf/5.png)

atomic 包提供了底层的原子级内存操作，对于同步算法的实现很有用。除了某些特殊的底层应用，使用通道或者 sync 包的函数/类型实现同步更好。

## 练习

交叉打印下面两个字符串"ABCDEFGHIJKLMNOPQRSTUVWXYZ" "0123..."

得到："AB01CD23EF34..."

仅供参考：

```go
package main

import (
    "fmt"
    "strconv"
    "sync"
)

var wg1 sync.WaitGroup
var ch1 chan int
var ch2 chan int
var s string

func sent3() {
    s1 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    i := 0
    for {
        select {
        case _, ok := <-ch2:
            if ok == false {
                s += s1[i : i+2]
                i += 2
            } else {
                s += s1[i : i+2]
                ch1 <- i
                i += 2
            }
        }
        if i > 24 {
            close(ch1)
            break
        }
    }
    wg1.Done()
}

func sent4() {
    j := 0
    ch2 <- 0
    for {
        select {
        case _, ok := <-ch1:
            if ok == false {
                s += strconv.Itoa(j) + strconv.Itoa(j+1)
                j += 2
            } else {
                s += strconv.Itoa(j) + strconv.Itoa(j+1)
                ch2 <- j
                j += 2
            }
        }
        if j >= 29 {
            close(ch2)
            break
        }
    }
    wg1.Done()
}

func main() {
    ch1 = make(chan int, 1)
    ch2 = make(chan int, 1)
    s = ""
    wg1.Add(1)
    go sent3()
    wg1.Add(1)
    go sent4()
    wg1.Wait()
    fmt.Printf("%s", s)
}

```

![](/img/blog/Gobf/6.png)

## 个人思考

### 并发与并行

我的理解是，并发更关注任务之间的切换和协调，而并行则是实打实的同时进行。并发就像一个人同时处理多个任务，比如看书时偶尔看看手机；并行则是两个人分别看不同的书，彼此互不干扰。

### Goroutine

相比于传统的线程，Goroutine 的内存占用更小，同时调度器能够动态分配合适的资源。这让我想到一个问题：如果每个任务所需的资源是极不均匀的，Goroutine 的轻量是否会成为一种负面影响，反而性能更低？

### Channel

在我看来，Channel 其实就是一个队列，只不过在这个队列里封装了一个等待拿数据和发送数据的功能，感觉这样更好理解，也更能使用好 Channel。
