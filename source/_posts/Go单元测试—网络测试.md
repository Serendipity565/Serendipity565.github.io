---
type: Post
title: Go单元测试—网络测试
tags: Go
category: 开发
category_bar: true
abbrlink: 34180
date: 2024-04-15 12:32:34
---

在实际中我们遇到的场景往往会比较复杂，无论我们的代码是作为server端对外提供服务或者还是我们依赖别人提供的网络服务（调用别人提供的API接口）的场景，我们通常都不想在测试过程中真正的建立网络连接。接下来就专门介绍如何在上述两种场景下mock网络测试。

## httptest

在Web开发场景下的单元测试，如果涉及到HTTP请求推荐大家使用Go标准库 `net/http/httptest` 进行测试，能够显著提高测试效率。

在这一小节，我们以常见的gin框架为例，演示如何为http server编写单元测试。

假设我们的业务逻辑是搭建一个http server端，对外提供HTTP服务。我们编写了一个`helloHandler`函数，用来处理用户请求。

```go
package httptest

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)

type Person struct {
    Name string `json:"name"`
}

func helloHandler(c *gin.Context) {
    var p Person
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusOK, gin.H{
            "msg": "we need a name",
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "msg": fmt.Sprintf("hello %s", p.Name),
    })
}

func SetupRouter() *gin.Engine {
    router := gin.Default()
    router.POST("/hello", helloHandler)
    return router
}
```

现在我们需要为`helloHandler`函数编写单元测试，这种情况下我们就可以使用`httptest`这个工具mock一个HTTP请求和响应记录器，让我们的server端接收并处理我们mock的HTTP请求，同时使用响应记录器来记录server端返回的响应内容。

单元测试的示例代码如下：

```go
package httptest

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_helloHandler(t *testing.T) {
 // 定义两个测试用例
    tests := []struct {
        name   string
        person string
        expect string
    }{
        {"base case", `{"name": "serendipity"}`, "hello serendipity"},
        {"bad case", "", "we need a name"},
    }

    r := SetupRouter()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // mock一个HTTP请求
            req := httptest.NewRequest(
                "POST",                       // 请求方法
                "/hello",                     // 请求URL
                strings.NewReader(tt.person), // 请求参数
            )

            // mock一个响应记录器
            w := httptest.NewRecorder()

            // 让server端处理mock请求并记录返回的响应内容
            r.ServeHTTP(w, req)

            // 校验状态码是否符合预期
            assert.Equal(t, http.StatusOK, w.Code)

            // 解析并检验响应内容是否复合预期
            var resp map[string]string
            err := json.Unmarshal([]byte(w.Body.String()), &resp)
            assert.Nil(t, err)                      //断言 err 变量是否为 nil
            assert.Equal(t, tt.expect, resp["msg"]) //断言 resp["msg"] 的值是否等于 tt.expect。
        })
    }
}
```

终端运行`go test ./httptest -v`执行单元测试，查看测试结果。

```go
=== RUN   Test_helloHandler
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /hello                    --> gin/httptest.helloHandler (3 handlers)
=== RUN   Test_helloHandler/base_case
[GIN] 2024/04/15 - 14:30:08 | 200 |            0s |       192.0.2.1 | POST     "/hello"
=== RUN   Test_helloHandler/bad_case
[GIN] 2024/04/15 - 14:30:08 | 200 |            0s |       192.0.2.1 | POST     "/hello"
--- PASS: Test_helloHandler (0.01s)
    --- PASS: Test_helloHandler/base_case (0.01s)
    --- PASS: Test_helloHandler/bad_case (0.00s)
PASS
ok      gin/httptest    0.124s
```

通过这个示例我们就掌握了如何使用httptest在HTTP Server服务中为请求处理函数编写单元测试了。

## gock

上面的示例介绍了如何在HTTP Server服务类场景下为请求处理函数编写单元测试，那么如果我们是在代码中请求外部API的场景（比如通过API调用其他服务获取返回值）又该怎么编写单元测试呢？

例如，我们有以下业务逻辑代码，依赖外部API：`http://your-api.com/post` 提供的数据。

```go
package api

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

// ReqParam API请求参数
type ReqParam struct {
    X int `json:"x"`
}

// Result API返回结果
type Result struct {
    Value int `json:"value"`
}

func GetResultByAPI(x, y int) int {
    p := &ReqParam{X: x}
    b, _ := json.Marshal(p)

    // 调用其他服务的API
    resp, err := http.Post(
        "http://your-api.com/post",
        "application/json",
        bytes.NewBuffer(b),
    )
    if err != nil {
        return -1
    }
    body, _ := ioutil.ReadAll(resp.Body)
    var ret Result
    if err := json.Unmarshal(body, &ret); err != nil {
        return -1
    }
    // 对API返回的数据做一些逻辑处理
    return ret.Value + y
}
```

在对类似上述这类业务代码编写单元测试的时候，如果不想在测试过程中真正去发送请求或者依赖的外部接口还没有开发完成时，我们可以在单元测试中对依赖的API进行mock。

这里推荐使用[gock](https://github.com/h2non/gock)这个库。

### 安装

```go
go get -u gopkg.in/h2non/gock
```

### 使用

使用`gock`对外部API进行mock，即mock指定参数返回约定好的响应内容。 下面的代码中mock了两组数据，组成了两个测试用例。

```go
package api

import (
    "github.com/h2non/gock"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestGetResultByAPI(t *testing.T) {
    defer gock.Off() // 测试执行后刷新挂起的mock
    // mock 请求外部api时传参x=1返回100
    gock.New("http://your-api.com").
        Post("/post").
        MatchType("json").
        JSON(map[string]int{"x": 1}).
        Reply(200).
        JSON(map[string]int{"value": 100})
    
    res := GetResultByAPI(1, 1)
    // 校验返回结果是否符合预期
    assert.Equal(t, res, 101)

    // mock 请求外部api时传参x=2返回200
    gock.New("http://your-api.com").
        Post("/post").
        MatchType("json").
        JSON(map[string]int{"x": 2}).
        Reply(200).
        JSON(map[string]int{"value": 200})
    
    res = GetResultByAPI(2, 2)
    // 校验返回结果是否符合预期
    assert.Equal(t, res, 202)

    assert.True(t, gock.IsDone()) // 断言mock被触发
}
```

执行上面写好的单元测试，看一下测试结果。

```go
=== RUN   TestGetResultByAPI
--- PASS: TestGetResultByAPI (0.00s)
PASS
ok      gin/api 0.123s
```

测试结果和预期的完全一致。
