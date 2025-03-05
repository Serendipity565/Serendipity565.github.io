---
type: Post
title: Go语言中结构体与json映射
tags: Go
category: 开发
category_bar: true
abbrlink: 14019
date: 2024-03-17 17:23:39
---

结构体与JSON之间的互相转化：json.Marshal和json.Unmarshal函数。

Marshal(v any) ([]byte, error)：将v转成json数据，以[]byte的形式返回。

Unmarshal(data []byte, v any) error：将json解析成指定的结构体。

如果转换成功，则该函数会返回nil，表示没有出现任何错误；如果解析失败，则会返回一个非空的error对象，表示解析过程中发生了错误，具体的错误信息可以通过该error对象获取。

## struct转json

默认情况下，转化后的json中的key值和结构体中的字段名是一样的，如果我们期望转化后的json字段名和struct里的不一样的话，就得用到tag。tag在这里的用途就是提供json结构中的别名，让两者的转化更加灵活。

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    Name           string
    Age            int
    Height         float64
    Weight         *float64
    marriageStatus string
}

func main() {
    weight := 120.0
    user := User{
        Name:           "serendipity",
        Age:            21,
        Height:         172,
        Weight:         &weight,
        marriageStatus: "未婚",
    }
    jsonBytes, err := json.Marshal(user)
    if err != nil {
        fmt.Println("error: ", err)
        return
    }
    fmt.Println(string(jsonBytes))
}
```

输出：

```go
{"Name":"serendipity","Age":21,"Height":172,"Weight":120}
```

上面的代码中，由于marriageStatus是小写开头，属于private字段，不会转换为为json。

针对JSON的输出，我们在定义struct tag的时候需要注意的几点是：

- 首字母为小写时，为private字段，不会转换。这也符合Go的语法规定，以小写字母开头的变量或结构体字段等，不能在包外被访问。
- tag中带有自定义名称，那么这个自定义名称会出现在JSON的字段名中
- 字段的tag是"-"，那么这个字段不会输出到JSON
- tag中如果带有"omitempty"选项，那么如果该字段值为空，就不会输出到JSON串中，比如 false、0、nil、长度为0的 array、map、slice和string
- 如果字段类型是bool, string、int、int64等，而tag中带有",string"选项，那么这个字段在输出到JSON的时候会把该字段对应的值转换成JSON字符串

例如：

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    Name           string  `json:"name"`
    Age            int     `json:",string"`
    Height         float64 `json:"weight"`
    MarriageStatus string  `json:"-"`
}

func main() {
    user := User{
        Name:           "serendipity",
        Age:            21,
        Height:         172,
        MarriageStatus: "未婚",
    }
    jsonBytes, err := json.Marshal(user)
    if err != nil {
        fmt.Println("error: ", err)
        return
    }
    fmt.Println(string(jsonBytes))
}

```

输出：

```go
{"name":"serendipity","Age":"21","weight":172}
```

在这个例子中，我们忽视了MarriageStatus，在输出时将Height在json中改为weight，并且将int类型的age改为了string类型。

## json转struct

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    Name           string `json:"name"`
    Age            int    `json:"age"`
    Height         float64
    Child          bool `json:"-"`
    marriageStatus string
}

func main() {
    userStr := `
    {
      "name": "gopher",
      "age": 18,
      "height": 180.5,
      "child": true,
      "marriageStatus": "未婚"
    }
    `

    user := User{}
    err := json.Unmarshal([]byte(userStr), &user)
    if err != nil {
        fmt.Println("error: ", err)
        return
    }
    fmt.Printf("%#v\n", user)
}
```

输出：

```go
main.User{Name:"gopher", Age:18, Height:180.5, Child:false, marriageStatus:""}
```

- 使用Unmarshal函数时，我们需要传入结构体的指针类型，否则结构体字段的值将不会被改变，因为底层是通过指针去修改结构体字段的值。
- json解析时，json的key与结构体字段的匹配规则是：

    1.优先查找json标签值和key一样的，找到则将value赋值给对应字段。

    2.如果没有json标签值与key相匹配，则根据字段名进行匹配。

- 可以发现，如果结构体字段是非导出字段或json标签的值为”-“，将不会匹配。
