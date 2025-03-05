---
type: Post
title: Python爬虫入门
tags: Python
category: 开发
category_bar: true
abbrlink: 28509
date: 2024-01-20 20:01:21
---

## HTTP网页结构

一个网页有三大技术要素，分别为HTML，CSS和JavaScript，其中HTML定义了网页的结构和信息，CSS定义网页的样式，JavaScript定义用户和网页的交互逻辑。其中我们爬虫最关心的就是HTML。

```html
<!DOCTYPE html>
<html>
    <body>
        <h1>这是一个一级标题</h1>
        <p>这是一段文字这是一段文字这是一段文字</p>
    </body>
</html>
```

## HTTP请求和响应

HTTP请求的方法类型主要的有两种，GET方法和POST方法，其中前者用于获得数据，后者用于创建数据（例如注册账号时提交数据）。

一个完整的HTTP请求包括请求行，请求头和请求体。请求行会包含方法类型，资源路径（查询参数）和协议版本。请求头包含Host（主机域名），User-Agent（服务器客户端相关信息）和Accept（客户端像接收的数据，接收多种数据类型可以用逗号进行分割；*/*则表示任意数据类型）。请求体可以放客户端传给服务器的任意数据，但是GET方法的请求体一般是空的。

![](/img/blog/pypc/1.png)

HTTP响应也由三个部分组成，分别是状态行，响应头和响应体。状态行还包含协议版本，状态码和状态消息，状态码和状态消息一一对应：状态码2开头表示成功，，请求已处理完成；3开头表示表示重定向，需要进一步的操作；4开头表示客户端错误，比如请求里面有错误，或者请求的资源无效，等等；5开头表示服务器错误，比如出现问题或者正在维护。响应头会包含一些告知客户端的信息，比如Date：生产响应的日期和时间；Content-Type：返回内容的类型及编码格式。响应体就是服务器想给客户端的数据内容，通常与Content-Type的类型相对应。

![](/img/blog/pypc/2.png)

## **requests 模块**

Python requests 是一个常用的 HTTP 请求库，可以方便地向网站发送 HTTP 请求，并获取响应结果。requests 模块比urllib模块更简洁。

使用 requests 发送 HTTP 请求需要先导入 requests 模块。导入后就可以发送 HTTP 请求，使用 requests 提供的方法向指定 URL 发送 HTTP 请求。

我们使用requests函数时的User-Agent时自动生成的，服务器可以轻松的分辨出是否是浏览器发出的请求，不过我们可以通过手动传入一个headers的参数帮我们模拟出浏览器发出的请求。

```python
import requests

head = {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; X64)"}
response = requests.get("https://www.baidu.com/", headers=head)
if response.ok:  # 或者用response.status_code>=200 and response.status_code<400
    print(response.text)
else:
    print("请求失败")
```

## Beautiful Soup 模块

爬到了网页信息之后，我们要提取出我们想要的不分，Beautiful Soup 可以帮我们解析网页内容Beautiful Soup自动将输入文档转换为Unicode编码，输出文档转换为utf-8编码。你不需要考虑编码方式，除非文档没有指定一个编码方式。

例如用requests和Beautiful Soup爬取豆瓣电影前250名单:

```python
import requests
from bs4 import BeautifulSoup

head = {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"}
for start_num in range(0,250,25):
    response = requests.get(f"https://movie.douban.com/top250?start={start_num}", headers=head)
    html=response.text
    soup=BeautifulSoup(html,"html.parser")
    all_title=soup.find_all("span",attrs={"class":"title"})
    for title in all_title:
        title_string=title.string
        if "/" not in title_string:
            print(title_string)
```
