---
type: Post
title: 在gin框架中使用JWT
tags: 技术
category: 开发
category_bar: true
abbrlink: 44835
date: 2024-04-07 01:01:23
---

JWT全称JSON Web Token是一种跨域认证解决方案，属于一个开放的标准，它规定了一种Token实现方式，目前多用于前后端分离项目和OAuth2.0业务场景下。

## 什么是JWT？

JWT全称JSON Web Token是一种跨域认证解决方案，属于一个开放的标准，它规定了一种Token 实现方式，目前多用于前后端分离项目和 OAuth2.0 业务场景下。

## 为什么需要JWT？

在之前的一些web项目中，我们通常使用的是`Cookie-Session`模式实现用户认证。相关流程大致如下：

1. 用户在浏览器端填写用户名和密码，并发送给服务端
2. 服务端对用户名和密码校验通过后会生成一份保存当前用户相关信息的session数据和一个与之对应的标识（通常称为session_id）
3. 服务端返回响应时将上一步的session_id写入用户浏览器的Cookie
4. 后续用户来自该浏览器的每次请求都会自动携带包含session_id的Cookie
5. 服务端通过请求中的session_id就能找到之前保存的该用户那份session数据，从而获取该用户的相关信息。

这种方案依赖于客户端（浏览器）保存 Cookie，并且需要在服务端存储用户的session数据。

在移动互联网时代，我们的用户可能使用浏览器也可能使用APP来访问我们的服务，我们的web应用可能是前后端分开部署在不同的端口，有时候我们还需要支持第三方登录，这下`Cookie-Session`的模式就有些力不从心了。

JWT就是一种基于Token的轻量级认证模式，服务端认证通过后，会生成一个JSON对象，经过签名后得到一个Token（令牌）再发回给用户，用户后续请求只需要带上这个Token，服务端解密之后就能获取该用户的相关信息了。

想要了解JWT的原理，推荐大家阅读：[阮一峰的JWT入门教程](https://www.ruanyifeng.com/blog/2018/07/json_web_token-tutorial.html)

## 安装

我们使用 Go 语言社区中的 jwt 相关库来构建我们的应用，例如：[https://github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)。

```go
go get github.com/golang-jwt/jwt/v4
```

我们将使用这个库来实现生成JWT和解析JWT的功能。

## 使用

### 默认Claim

如果我们直接使用JWT中默认的字段，没有其他定制化的需求则可以直接使用这个包中的和方法快速生成和解析token。

```go
// 用于签名的字符串
var mySigningKey = []byte("serendipity565.vercel.app")

// GenRegisteredClaims 使用默认声明创建jwt
func GenRegisteredClaims() (string, error) {
    // 创建 Claims
    claims := &jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),  // 过期时间
        Issuer:    "serendipity",                                       // 签发人
    }
    // 生成token对象
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    // 生成签名字符串
    return token.SignedString(mySigningKey)
}

// ParseRegisteredClaims 解析jwt
func ValidateRegisteredClaims(tokenString string) bool {
    // 解析token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return mySigningKey, nil
    })
    if err != nil { // 解析token失败
        return false
    }
    return token.Valid
}
```

### 自定义Claims

我们需要定制自己的需求来决定JWT中保存哪些数据，比如我们规定在JWT中要存储`username`信息，那么我们就定义一个`MyClaims`结构体如下：

```go
// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
    // 可根据需要自行添加字段
    Username             string `json:"username"`
    jwt.RegisteredClaims        // 内嵌标准的声明
}
```

然后我们定义JWT的过期时间，这里以24小时为例：

```go
const TokenExpireDuration = time.Hour ***** 24
```

接下来还需要定义一个用于签名的字符串：

```go
// CustomSecret 用于签名的字符串
var CustomSecret = []byte("HelloWorld")
```

### 生成JWT

我们可以根据自己的需要封装一个生成 token 的函数。

```go
// GenToken 生成JWT
func GenToken(username string) (string, error) {
    // 创建一个我们自己的声明
    claims := CustomClaims{
        username, // 自定义字段
        jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
            Issuer:    "my-project", // 签发人
        },
    }
    // 使用指定的签名方法创建签名对象
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    // 使用指定的secret签名并获得完整的编码后的字符串token
    return token.SignedString(CustomSecret)
}
```

### 解析JWT

根据给定的 JWT 字符串，解析出数据。

```go
// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
    // 解析token
    // 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
    token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
    // 直接使用标准的Claim则可以直接使用Parse方法
    //token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
        return CustomSecret, nil
    })
    if err != nil {
        return nil, err
    }
    // 对token对象中的Claim进行类型断言
    if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
        return claims, nil
    }
    return nil, errors.New("invalid token")
}
```

## 在gin框架中使用JWT

首先我们注册一条路由`/auth`，对外提供获取Token的渠道：

`r.POST("/auth", authHandler)`

我们的`authHandler`定义如下：

```go
func authHandler(c *gin.Context) {
    // 用户发送用户名和密码过来
    var user UserInfo
    err := c.ShouldBind(&user)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{
            "code": 2001,
            "msg":  "无效的参数",
        })
        return
    }
    // 校验用户名和密码是否正确
    if user.Username == "serendipity" && user.Password == "serendipity123" {
        // 生成Token
        tokenString, _ := GenToken(user.Username)
        c.JSON(http.StatusOK, gin.H{
            "code": 2000,
            "msg":  "success",
            "data": gin.H{"token": tokenString},
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "code": 2002,
        "msg":  "鉴权失败",
    })
    return
}
```

用户通过上面的接口获取Token之后，后续就会携带着Token再来请求我们的其他接口，这个时候就需要对这些请求的Token进行校验操作了，很显然我们应该实现一个检验Token的中间件，具体实现如下：

```go
// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
    return func(c *gin.Context) {
        // 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
        // 这里假设Token放在Header的Authorization中，并使用Bearer开头
        // 这里的具体实现方式要依据实际情况决定
        authHeader := c.Request.Header.Get("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusOK, gin.H{
                "code": 2003,
                "msg":  "请求头中auth为空",
            })
            c.Abort()
            return
        }
        // 按空格分割
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(http.StatusOK, gin.H{
                "code": 2004,
                "msg":  "请求头中auth格式有误",
            })
            c.Abort()
            return
        }
        // parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
        mc, err := ParseToken(parts[1])
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": 2005,
                "msg":  "无效的Token",
            })
            c.Abort()
            return
        }
        // 将当前请求的username信息保存到请求的上下文c上
        c.Set("username", mc.Username)
        c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
    }
}
```

**说明**：`c.Abort()`通常指的是在处理HTTP请求的过程中，提前终止请求的处理，并且不再执行后续的处理流程。

注册一个`/home`路由，发个请求验证一下。

```go
r.GET("/home", JWTAuthMiddleware(), homeHandler)

func homeHandler(c *gin.Context) {
    username := c.MustGet("username").(string)
    c.JSON(http.StatusOK, gin.H{
        "code": 2000,
        "msg":  "success",
        "data": gin.H{"username": username},
    })
}
```

如果不想自己实现上述功能，你也可以使用Github上别人封装好的包，比如：[https://github.com/appleboy/gin-jwt](https://github.com/appleboy/gin-jwt)。

### refresh token

在某些场景下，我们可能还需要使用refresh token，这里可以参考 [RFC 6749 OAuth2.0中关于refresh token的介绍](https://datatracker.ietf.org/doc/html/rfc6749#section-1.5)
