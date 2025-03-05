---
type: Post
title: Gin框架介绍及使用
tags: Go
category: 开发
category_bar: true
abbrlink: 15144
date: 2024-04-01 21:45:56
---

## Gin框架安装与使用

### 安装

下载并安装gin：

```go
go get -u [github.com/gin-gonic/gin](http://github.com/gin-gonic/gin)
```

第一个gin示例：

```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    // 创建一个默认的路由引擎
    r := gin.Default()
    // 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
    r.GET("/hello", func(c *gin.Context) {
        // c.JSON：返回JSON格式的数据
        c.JSON(200, gin.H{
            "message": "Hello world!",
        })
    })
    r.Run()
}
```

## RESTful API

REST与技术无关，代表的是一种软件架构风格，REST是Representational State Transfer的简称，中文翻译为“表征状态转移”或“表现层状态转化”。

简单来说，REST的含义就是客户端与Web服务器之间进行交互的时候，使用HTTP协议中的4个请求方法代表不同的动作。

- GET 用来获取资源
- POST 用来新建资源
- PUT 用来更新资源
- DELETE 用来删除资源。

只要API程序遵循了REST风格，那就可以称其为RESTful API。目前在前后端分离的架构中，前后端基本都是通过RESTful API来进行交互。

例如，我们现在要编写一个管理书籍的系统，我们可以查询对一本书进行查询、创建、更新和删除等操作，我们在编写程序的时候就要设计客户端浏览器与我们Web服务端交互的方式和路径。按照经验我们通常会设计成如下模式：

| 请求方法 | URL | 含义 |
| --- | --- | --- |
| GET | /book | 查询书籍信息 |
| POST | /create_book | 创建书籍记录 |
| POST | /update_book | 更新书籍信息 |
| POST | /delete_book | 删除书籍信息 |

同样的需求我们按照RESTful API设计如下：

| 请求方法 | URL | 含义 |
| --- | --- | --- |
| GET | /book | 查询书籍信息 |
| POST | /book | 创建书籍记录 |
| PUT | /book | 更新书籍信息 |
| DELETE | /book | 删除书籍信息 |

```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/book", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "GET",
        })
    })

    r.POST("/book", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "POST",
        })
    })

    r.PUT("/book", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "PUT",
        })
    })

    r.DELETE("/book", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "DELETE",
        })
    })
}
```

## Gin渲染

### HTML渲染

我们首先定义一个存放模板文件的 templates 文件夹，然后在其内部按照业务分别定义一个 posts 文件夹和一个 users 文件夹。 posts/index.html 文件的内容如下：

```html
{{define "posts/index.html"}}
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>posts/index</title>
</head>
<body>
    {{.title}}
</body>
</html>
{{end}}

```

users/index.html 文件的内容如下：

```html
{{define "users/index.html"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>users/index</title>
</head>
<body>
    {{.title}}
</body>
</html>
{{end}}

```

Gin框架中使用 `LoadHTMLGlob( )` 或者 `LoadHTMLFiles( )` 方法进行HTML模板渲染。

```go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()

    r.LoadHTMLGlob("templates/**/*")
    //或者r.LoadHTMLFiles("templates/posts/index.html", "templates/users/index.html")

    r.GET("/posts/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "posts/index.html", gin.H{
            "title": "posts/index",
        })
    })

    r.GET("users/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "users/index.html", gin.H{
            "title": "users/index",
        })
    })

    r.Run(":8080")
}
```

### 自定义模板函数

定义一个不转义相应内容的`safe`模板函数如下：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "html/template"
    "net/http"
)

func main() {
    r := gin.Default()
    r.SetFuncMap(template.FuncMap{
        "safe": func(str string) template.HTML {
            return template.HTML(str)
        },
    })
    r.LoadHTMLFiles("./index.tmpl")

    r.GET("/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", "<a href='https://serendipity565.vercel.app/'>serendipity的博客</a>")
    })

    r.Run(":8080")
}
```

在`index.tmpl`中使用定义好的`safe`模板函数：

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <title>修改模板引擎的标识符</title>
</head>
<body>
<div>{{ . | safe }}</div>
</body>
</html>
```

在 Go Gin 框架中，可以使用不同的模板文件扩展名，比如 `.html` 和 `.tmpl` 等。`index.tmpl` 或 `index.tpl` 是一个模板文件，通常用于动态生成页面内容，可以包含一些动态数据和逻辑。在 Gin 框架中，通常会使用模板文件来渲染动态页面。

### 静态文件处理

当我们渲染的HTML文件中引用了静态文件时，我们只需要按照以下方式在渲染页面前调用`gin.Static`方法即可。

```go
func main() {
    r := gin.Default()

    r.Static("/static", "./static")
    r.LoadHTMLGlob("templates/**/*")
    // ...
    r.Run(":8080")
}
```

使用 `r.Static("/static", "./static")` 配置了静态文件服务。这会将 URL 路径 "/static" 映射到项目中的 "./static" 目录，客户端可以通过该路径访问静态资源，例如图片、CSS、JavaScript 等文件。

### 使用模板继承

Gin框架默认都是使用单模板，如果需要使用`block template`功能，可以通过`"github.com/gin-contrib/multitemplate"`库实现，具体示例如下：

首先，假设我们项目目录下的templates文件夹下有以下模板文件，其中`home.tmpl`和`index.tmpl`继承了`base.tmpl`：

```Plain text
templates
├── includes
│   ├── home.tmpl
│   └── index.tmpl
├── layouts
│   └── base.tmpl
└── scripts.tmpl
```

然后我们定义一个`loadTemplates`函数如下：

```go
package main

import (
    "github.com/gin-contrib/multitemplate"
    "github.com/gin-gonic/gin"
    "net/http"
    "path/filepath"
)

func loadTemplates(templatesDir string) multitemplate.Renderer {
    //创建一个 multitemplate.Renderer 实例
    r := multitemplate.NewRenderer()
    layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
    if err != nil {
        panic(err.Error())
    }
    includes, err := filepath.Glob(templatesDir + "/includes/*.tmpl")
    if err != nil {
        panic(err.Error())
    }
    // 为layouts/和includes/目录生成 templates map
    // 将文件路径赋值给include
    for _, include := range includes {
        layoutCopy := make([]string, len(layouts))
        copy(layoutCopy, layouts)
        files := append(layoutCopy, include)
        r.AddFromFiles(filepath.Base(include), files...)
    }
    return r
}

func indexFunc(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl", nil)
}

func homeFunc(c *gin.Context) {
    c.HTML(http.StatusOK, "home.tmpl", nil)
}

func main() {
    r := gin.Default()
    r.HTMLRender = loadTemplates("./templates")
    r.GET("/index", indexFunc)
    r.GET("/home", homeFunc)
    r.Run()
}
```

### 补充文件路径处理

关于模板文件和静态文件的路径，我们需要根据要求进行设置。可以使用下面的函数获取当前执行程序的路径。

```go
func getCurrentPath() string {
    if ex, err := os.Executable(); err == nil {
        return filepath.Dir(ex)
    }
    return "./"
}
```

### JSON渲染

```go
func main() {
    r := gin.Default()

    // gin.H 是map[string]interface{}的缩写
    r.GET("/someJSON", func(c *gin.Context) {
        // 方式一：自己拼接JSON
        c.JSON(http.StatusOK, gin.H{"message": "Hello world!"})
    })
    r.GET("/moreJSON", func(c *gin.Context) {
        // 方法二：使用结构体
        var msg struct {
            Name    string `json:"user"`
            Message string
            Age     int
        }
        msg.Name = "serendipity"
        msg.Message = "Hello world!"
        msg.Age = 20
        c.JSON(http.StatusOK, msg)
    })
    r.Run(":8080")
}
```

## **获取参数**

### 获取querystring参数

`querystring`指的是URL中`?`后面携带的参数。获取请求的querystring参数的方法如下：

```go
func main() {
    r := gin.Default()

    r.GET("/user/search", func(c *gin.Context) {
        username := c.DefaultQuery("username", "unknow")
        //username := c.Query("username")
        address := c.Query("address")
        //输出json结果给调用方
        c.JSON(http.StatusOK, gin.H{
            "message":  "ok",
            "username": username,
            "address":  address,
        })
    })
    r.Run()
}
```

1. r.DefaultQuery(key, defaultValue)：
    - `DefaultQuery` 方法用于获取指定键的查询参数值，如果该参数不存在，则返回默认值defaultValue。
    - 如果请求中不存在指定的查询参数，则返回提供的默认值defaultValue。
    - 这个方法适用于在获取查询参数时提供一个默认值，以避免因参数不存在而引发的错误。
2. r.URL.Query( ).Get(key) 或 r.URL.Query( ) [key]：
    - `Query` 方法用于获取指定键的查询参数值。
    - 如果请求中不存在指定的查询参数，则返回空字符串。
    - 这个方法适用于在获取查询参数时直接获取其值，如果参数不存在则返回空字符串。

### 获取form参数

当前端请求的数据通过form表单提交时，例如向`/user/search`发送一个POST请求，获取请求数据的方式如下：

```go
func main() {
    r := gin.Default()

    r.POST("/user/search", func(c *gin.Context) {
        //username := c.DefaultPostForm("username", "serendipity")
        username := c.PostForm("username")
        address := c.PostForm("address")
        //输出json结果给调用方
        c.JSON(http.StatusOK, gin.H{
            "message":  "ok",
            "username": username,
            "address":  address,
        })
    })
    r.Run(":8080")
}
```

`DefaultPostForm(key,  defaultValue)`和`PostForm(key)`这两个函数用法与`DefaultQuery(key, defaultValue)`和`Query(key)`类似

### 获取JSON参数

当前端请求的数据通过JSON提交时，例如向`/json`发送一个JSON格式的POST请求，则获取请求参数的方式如下：

```go
r.POST("/json", func(c *gin.Context) {
    // 注意：下面为了举例子方便，暂时忽略了错误处理
    b, _ := c.GetRawData()  // 从c.Request.Body读取请求数据
    // 定义map或结构体
    var m map[string]interface{}
    // 反序列化
    _ = json.Unmarshal(b, &m)

    c.JSON(http.StatusOK, m)
}) 
```

更便利的获取请求参数的方式，参见下面的 **参数绑定** 小节。

### 获取path参数

请求的参数通过URL路径传递，例如：`/user/search/serendipity/China`。 获取请求URL路径中的参数的方式如下。

```go
func main() {
    r := gin.Default()

    r.GET("/user/search/:username/:address", func(c *gin.Context) {
        username := c.Param("username")
        address := c.Param("address")
        //输出json结果给调用方
        c.JSON(http.StatusOK, gin.H{
            "message":  "ok",
            "username": username,
            "address":  address,
        })
    })

    r.Run(":8080")
}
```

### 参数绑定

为了能够更方便的获取请求相关参数，提高开发效率，我们可以基于请求的`Content-Type`识别请求数据类型并利用反射机制自动提取请求中`QueryString`、`form表单`、`JSON`、`XML`等参数到结构体中。 下面的示例代码演示了`.ShouldBind()`强大的功能，它能够基于请求自动提取`JSON`、`form表单`和`QueryString`类型的数据，并把值绑定到指定的结构体对象。

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
)

type Login struct {
    User     string `form:"user" json:"user" binding:"required"`
    Password string `form:"password" json:"password" binding:"required"`
}

func main() {
    r := gin.Default()

    // 绑定JSON的示例 ({"user": "serendipity", "password": "123456"})
    r.POST("/loginJSON", func(c *gin.Context) {
        var login Login

        if err := c.ShouldBind(&login); err == nil {
            fmt.Printf("login info:%#v\n", login)
            c.JSON(http.StatusOK, gin.H{
                "user":     login.User,
                "password": login.Password,
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
    })

    // 绑定form表单示例 (user=serendipity&password=123456)
    r.POST("/loginForm", func(c *gin.Context) {
        var login Login
        // ShouldBind()会根据请求的Content-Type自行选择绑定器
        if err := c.ShouldBind(&login); err == nil {
            c.JSON(http.StatusOK, gin.H{
                "user":     login.User,
                "password": login.Password,
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
    })

    // 绑定QueryString示例 (/loginQuery?user=serendipity&password=123456)
    r.GET("/loginForm", func(c *gin.Context) {
        var login Login
        // ShouldBind()会根据请求的Content-Type自行选择绑定器
        if err := c.ShouldBind(&login); err == nil {
            c.JSON(http.StatusOK, gin.H{
                "user":     login.User,
                "password": login.Password,
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
    })

    r.Run(":8080")
}
```

`ShouldBind`会按照下面的顺序解析请求中的数据完成绑定：

1. 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）。
2. 如果是 `POST` 请求，首先检查 `content-type` 是否为 `JSON` 或 `XML`，然后再使用 `Form`（`form-data`）。

## 文件上传

### 单个文件上传

文件上传前端页面代码：

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <title>上传文件示例</title>
</head>
<body>
<form action="/upload" method="post" enctype="multipart/form-data">
    <input type="file" name="f1">
    <input type="submit" value="上传">
</form>
</body>
</html>
```

后端gin框架部分代码：

```go
func main() {
    r := gin.Default()

    r.POST("/upload", func(c *gin.Context) {
        file, err := c.FormFile("f1")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "message": err.Error(),
            })
            return
        }

        log.Println(file.Filename)
        dst := fmt.Sprintf("C:/tmp/%s", file.Filename)
        // 上传文件到指定的目录
        c.SaveUploadedFile(file, dst)
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("'%s' uploaded!", file.Filename),
        })
    })
    r.Run()
}
```

处理multipart forms提交文件时默认的内存限制是32 MB，我们可以通过下面的方式修改`r.MaxMultipartMemory` ，例如 r.MaxMultipartMemory = 64 << 20 将内存限制改到64MB，这里的 20 表示移动的位数，即左移 20 位，相当于将 32 乘以 2^20，即 1024 * 1024，即 1MB。

### 多个文件上传

```go
func main() {
    router := gin.Default()

    router.POST("/upload", func(c *gin.Context) {
        form, _ := c.MultipartForm()
        files := form.File["file"]

        for index, file := range files {
            log.Println(file.Filename)
            dst := fmt.Sprintf("C:/tmp/%s_%d", file.Filename, index)
            // 上传文件到指定的目录
            c.SaveUploadedFile(file, dst)
        }
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("%d files uploaded!", len(files)),
        })
    })
    router.Run()
}
```

## 重定向

### HTTP重定向

HTTP 重定向很容易。 内部、外部重定向均支持。

```go
r.GET("/test", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "https://serendipity565.vercel.app/")
})
```

### 路由重定向

路由重定向，使用`HandleContext`：

```go
r.GET("/test", func(c *gin.Context) {
    // 指定重定向的URL
    c.Request.URL.Path = "/test2"
    r.HandleContext(c)
})
r.GET("/test2", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"hello": "world"})
})
```

## Gin路由

### 普通路由

```go
r.GET("/login", func(c *gin.Context) {...})
r.POST("/login", func(c *gin.Context) {...})
```

此外，还有一个可以匹配所有请求方法的`Any`方法如下：

```go
r.Any("/test", func(c *gin.Context) {...})
```

为没有配置处理函数的路由添加处理程序，默认情况下它返回404代码，下面的代码为没有匹配到路由的请求都返回`views/404.html`页面。

```go
r.NoRoute(func(c *gin.Context) {
    c.HTML(http.StatusNotFound, "views/404.html", nil)
})
```

### 路由组

我们可以将拥有共同URL前缀的路由划分为一个路由组。习惯性一对`{}`包裹同组的路由，这只是为了看着清晰，你用不用`{}`包裹功能上没什么区别。

```go
func main() {
    r := gin.Default()
    userGroup := r.Group("/user")
    {
        userGroup.GET("/index", func(c *gin.Context) {...})
        userGroup.GET("/login", func(c *gin.Context) {...})
        userGroup.POST("/login", func(c *gin.Context) {...})
    }

    shopGroup := r.Group("/shop")
    {
        shopGroup.GET("/index", func(c *gin.Context) {...})
        shopGroup.GET("/cart", func(c *gin.Context) {...})
        shopGroup.POST("/checkout", func(c *gin.Context) {...})
    }
    r.Run()
}
```

路由组也是支持嵌套的，例如：

```go
shopGroup := r.Group("/shop")
{
    shopGroup.GET("/index", func(c *gin.Context) {...})
    shopGroup.GET("/cart", func(c *gin.Context) {...})
    shopGroup.POST("/checkout", func(c *gin.Context) {...})
    // 嵌套路由组
    xx := shopGroup.Group("xx")
    xx.GET("/oo", func(c *gin.Context) {...})
}
```

通常我们将路由分组用在划分业务逻辑或划分API版本时。

### 路由原理

Gin框架中的路由使用的是[httprouter](https://github.com/julienschmidt/httprouter)这个库，其基本原理就是构造一个路由地址的前缀树。

## Gin中间件

Gin框架允许开发者在处理请求的过程中，加入用户自己的钩子（Hook）函数。这个钩子函数就叫中间件，中间件适合处理一些公共的业务逻辑，比如登录认证、权限校验、数据分页、记录日志、耗时统计等。

### 定义中间件

Gin中的中间件必须是一个`gin.HandlerFunc`类型。

#### 记录接口耗时的中间件

例如我们像下面的代码一样定义一个统计请求耗时的中间件。

```go
// StatCost 是一个统计耗时请求耗时的中间件
func StatCost() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Set("name", "serendipity")
        // 可以通过c.Set在请求上下文中设置值，后续的处理函数能够取到该值
        // 调用该请求的剩余处理程序
        c.Next()
        // 不调用该请求的剩余处理程序
        // c.Abort()
        // 计算耗时
        cost := time.Since(start)
        log.Println(cost)
    }
}

func main() {
    r := gin.Default()

    // 使用自定义中间件
    r.Use(StatCost())

    r.GET("/hello", func(c *gin.Context) {
        name := c.MustGet("name").(string)
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello " + name,
        })
    })

    r.Run(":8080")
}
```

#### 记录响应体的中间件

我们有时候可能会想要记录下某些情况下返回给客户端的响应数据，这个时候就可以编写一个中间件来搞定。

```go
type bodyLogWriter struct {
    gin.ResponseWriter               // 嵌入gin框架ResponseWriter
    body               *bytes.Buffer // 我们记录用的response
}

// Write 写入响应体数据
func (w bodyLogWriter) Write(b []byte) (int, error) {
    w.body.Write(b)                  // 我们记录一份
    return w.ResponseWriter.Write(b) // 真正写入响应
}

// ginBodyLogMiddleware 一个记录返回给客户端响应体的中间件
// https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
func ginBodyLogMiddleware(c *gin.Context) {
    blw := &bodyLogWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
    c.Writer = blw // 使用我们自定义的类型替换默认的

    c.Next() // 执行业务逻辑

    fmt.Println("Response body: " + blw.body.String()) // 事后按需记录返回的响应
}
```

#### 跨域中间件cors

推荐使用社区的[https://github.com/gin-contrib/cors](https://github.com/gin-contrib/cors) 库，一行代码解决前后端分离架构下的跨域问题。

**注意：** 该中间件需要注册在业务处理函数前面。

这个库支持各种常用的配置项，具体使用方法如下。

```go
package main

import (
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    // CORS for https://foo.com and https://github.com origins, allowing:
    // - PUT and PATCH methods
    // - Origin header
    // - Credentials share
    // - Preflight requests cached for 12 hours
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://foo.com"},  // 允许跨域发来请求的网站
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE",  "OPTIONS"},  // 允许的请求方法
        AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        AllowOriginFunc: func(origin string) bool {  
            // 自定义过滤源站的方法
            return origin == "https://github.com"
        },
        MaxAge: 12 * time.Hour,
    }))
    r.Run()
}
```

当然你可以简单的像下面的示例代码那样使用默认配置，允许所有的跨域请求。

```go
func main() {
    r := gin.Default()
    // same as
    // config := cors.DefaultConfig()
    // config.AllowAllOrigins = true
    // router.Use(cors.New(config))
    r.Use(cors.Default())
    r.Run()
}
```

### 注册中间件

在gin框架中，我们可以为每个路由添加任意数量的中间件。

#### 为全局路由注册

```go
func main() {
    // 新建一个没有任何默认中间件的路由
    r := gin.New()
    // 注册一个全局中间件
    r.Use(StatCost())

    r.GET("/test", func(c *gin.Context) {
        name := c.MustGet("name").(string) // 从上下文取值
        log.Println(name)
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello world!",
        })
    })
    r.Run()
}
```

`gin.Default()` 和 `gin.New()` 都用于创建 Gin 框架的实例，但它们之间有一些细微的区别：

1. `gin.Default()`：
    - `gin.Default()` 方法会返回一个默认配置的 Gin 实例。
    - 默认配置包括 Logger 和 Recovery 中间件。Logger 中间件用于记录请求日志，Recovery 中间件用于处理恢复从处理程序中出现的 panic。
    - 这个方法在创建 Gin 实例时会自动使用默认的中间件，因此你无需手动添加 Logger 和 Recovery 中间件。
2. `gin.New()`：
    - `gin.New()` 方法会返回一个空白的 Gin 实例。
    - 这个方法创建的 Gin 实例不会包含任何默认的中间件。你需要手动添加所需的中间件。
    - 这个方法适用于需要完全自定义 Gin 实例的情况，你可以根据需要选择性地添加中间件。

一般来说，如果你想要快速搭建一个使用了默认 Logger 和 Recovery 中间件的 Gin 实例，可以使用 `gin.Default()`；如果你想要完全控制 Gin 实例中的中间件，可以使用 `gin.New()`。

#### 为某个路由单独注册

```go
// 给/test2路由单独注册中间件（可注册多个）
r.GET("/test2", StatCost(), func(c *gin.Context) {
    name := c.MustGet("name").(string) // 从上下文取值
    log.Println(name)
    c.JSON(http.StatusOK, gin.H{
        "message": "Hello world!",
    })
})
```

#### 为路由组注册中间件

为路由组注册中间件有以下两种写法。

写法1：

```go
shopGroup := r.Group("/shop", StatCost())
{
    shopGroup.GET("/index", func(c *gin.Context) {...})
    ...
}
```

写法2：

```go
shopGroup := r.Group("/shop")
shopGroup.Use(StatCost())
{
    shopGroup.GET("/index", func(c *gin.Context) {...})
    ...
}
```

### 中间件注意事项

#### gin中间件中使用goroutine

当在中间件或`handler`中启动新的`goroutine`时，**不能使用**原始的上下文（c *gin.Context），必须使用其只读副本，即`c.Copy()`。

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "time"
)

func SomeMiddleware(c *gin.Context) {
    // 使用只读副本
    copyOfContext := c.Copy()

    // 在 goroutine 中使用副本
    go func() {
        // 模拟耗时操作
        time.Sleep(1 * time.Second)
        // 从只读副本中获取数据
        fmt.Println(copyOfContext.Request.URL.Path)
    }()

    // 继续处理请求
    c.Next()
}

func main() {
    r := gin.Default()

    // 使用中间件
    r.Use(SomeMiddleware)

    r.GET("/test", func(c *gin.Context) {
        // 处理请求
        c.JSON(200, gin.H{
            "message": "Hello, World!",
        })
    })

    r.Run(":8080")
}

```

## 运行多个服务

我们可以在多个端口启动服务，例如：

```go
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/sync/errgroup"
)

var (
    g errgroup.Group
)

func router1() http.Handler {
    r := gin.New()
    r.Use(gin.Recovery())
    r.GET("/", func(c *gin.Context) {
        c.JSON(
            http.StatusOK,
            gin.H{
                "code":  http.StatusOK,
                "error": "Welcome server 01",
            },
        )
    })

    return r
}

func router2() http.Handler {
    r := gin.New()
    r.Use(gin.Recovery())
    r.GET("/", func(c *gin.Context) {
        c.JSON(
            http.StatusOK,
            gin.H{
                "code":  http.StatusOK,
                "error": "Welcome server 02",
            },
        )
    })

    return r
}

func main() {
    server01 := &http.Server{
        Addr:         ":8080",
        Handler:      router1(),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    server02 := &http.Server{
        Addr:         ":8081",
        Handler:      router2(),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    // 借助errgroup.Group或者自行开启两个goroutine分别启动两个服务
    g.Go(func() error {
        return server01.ListenAndServe()
    })

    g.Go(func() error {
        return server02.ListenAndServe()
    })

    if err := g.Wait(); err != nil {
        log.Fatal(err)
    }
}
```
