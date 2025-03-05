---
type: Post
title: Go语言net/http
tags: Go
category: 开发
category_bar: true
abbrlink: 63938
date: 2024-03-19 01:32:29
---

## 初识net/http包

我们先初步介绍以下net/http包的使用，通过http.HandleFunc()和http.ListenAndServe()两个函数就可以轻松创建一个简单的Go web服务器，示例代码如下:

```go
package main

import (
    "fmt"
    "net/http"
    "strings"
)

func hello(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()       // 解析参数，默认是不会解析的
    fmt.Println(r.Form)   // 这些信息是输出到服务器端的打印信息
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }

    // 这个写入到 w 的是输出到客户端的
    fmt.Fprintf(w, "Hello!")
}

func main() {
    http.HandleFunc("/hello", hello)
    http.ListenAndServe(":8000", nil)
}

```

访问界面如下：

![](/img/blog/Gonethttp/1.png)

首先，我们调用http.HandleFunc("/hello", hello)注册路径处理函数，这里将路径/hello的处理函数设置为hello。即引号里的hello是访问地址，逗号后面的hello是我们创建的函数，用来响应这个地址的访问。处理函数的类型必须是：

```go
func (http.ResponseWriter, *http.Request)
```

其中*http.Request表示HTTP请求对象，该对象包含请求的所有信息，如URL、首部、表单内容、请求的其他内容等。

http.ResponseWriter是一个接口类型：

```go
// net/http/server.go
type ResponseWriter interface {
    Header() Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int)
}
```

main()函数通过代码http.ListenAndServe(":8000“,nil)启动一个8000端口的服务器。

ListenAndServe()函数有两个参数，当前监听的端口号和事件处理器Handler。如果ListenAndServe()传入的第一个参数地址为空，则服务器在启动后默认使用`http://localhost:8080`地址进行访问；如果这个函数传入的第二个参数为nil，则服务器在启动后将使用默认的多路复用器DefaultServeMux。

要想结束这个服务器，只需要在终端输入`Ctrl+c`。

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
)

type students struct {
    Name   string  `json:"name"`
    Age    int     `json:"age"`
    Id     string  `json:"id"`
    Height float64 `json:"height"`
}

func hello(w http.ResponseWriter, r *http.Request) {
    stu1 := students{
        Name:   "serendipity",
        Age:    21,
        Id:     "202300001",
        Height: 172.0,
    }
    res, err := json.Marshal(stu1)
    if err != nil {
        fmt.Printf("json error is %s", err)
        return
    }

    m1 := make(map[string]interface{})
    m1["name"] = "xiaoming"
    m1["age"] = 18
    m1["id"] = "202300002"
    m1["height"] = 168.2
    res1, err1 := json.Marshal(m1)
    if err1 != nil {
        fmt.Printf("json error is %s", err1)
        return
    }

    fmt.Fprintf(w, string(res))
    fmt.Fprintf(w, string(res1))
}

var count int
var mu sync.Mutex

func index(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    count++
    mu.Unlock()
    fmt.Fprintf(w, "count is %d", count)
}

func main() {
    http.HandleFunc("/index", index)
    http.HandleFunc("/main", hello)
    log.Fatal(http.ListenAndServe(":8000", nil))
}
```

我们来逐步看这个代码。首先处理了两个url，对于/main这个url我们传输了两个json数据。有两种方法可以将数据转换为json格式，一种是利用结构体,另外一种是使用map。map转换成json格式输出是的顺序并不是按我们定的顺序，而是按照字符串的从小到大的顺序。

![](/img/blog/Gonethttp/2.png)

对于/index这个网页，我们使用count这个函数，服务器每一次接收请求处理时都会另起一个goroutine，这样服务器就可以同一时间处理多个请求。然而在并发情况下，假如真的有两个请求同一时刻去更新count，那么这个值可能并不会被正确地增加；这个程序可能会引发一个严重的bug：竞态条件。为了避免这个问题，我们必须保证每次修改变量的最多只能有一个goroutine，这也就是代码里的mu.Lock()和mu.Unlock()调用将修改count的所有行为包在中间的目的。

![](/img/blog/Gonethttp/3.png)

![](/img/blog/Gonethttp/4.png)
