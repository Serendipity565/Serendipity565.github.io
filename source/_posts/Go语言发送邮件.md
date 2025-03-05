---
type: Post
title: Go语言发送邮件
tags: Go
category: 开发
category_bar: true
abbrlink: 13327
date: 2025-01-12 17:23:04
---

## **邮件协议**

各种事物都有一个规范，也就是协议，比如我们在浏览器里面浏览网页，需要遵循各种网络协议，我们先来简单了解一下都有哪些协议

1. SMTP
    
    `SMTP`是 简单邮件传输协议，是一组用于从源地址到目的地址传输邮件的规范，通过它来控制邮件的中转方式。它通常在 25、465、587 端口上运行。
    
    另外 `SMTP` 协议属于`TCP/IP`协议簇
    
2. POP3
    
    邮局协议的第**3**个版本，是因特网电子邮件的第一个离线协议标准(邮件服务器下载邮件到本地计算机后，可以断开网络连接继续查看邮件内容)。下载后邮件会从服务器删除。
    
3. IMAP
    
    是一种优于`POP`的新协议，与`POP`不同的是，他是典型的在线协议。和`POP`一样，`IMAP`也能下载邮件、从服务器中删除邮件或询问是否有新邮件。
    
    `IMAP`可让用户在服务器上创建并管理邮件文件夹或邮箱、删除邮件、查询某封信的一部分或全部内容。
    
    最终完成所有这些工作都不需要把邮件从服务器下载到用户的个人计算机上。
    

## 一些基础的配置

- QQ邮箱的设置:需要打开`POP3/SMTP`服务

![](/img/blog/GoSendEmail/1.png)

![](/img/blog/GoSendEmail/2.png)

![](/img/blog/GoSendEmail/3.png)

温馨提示：在使用 QQ 邮箱发送邮件的时候，需要使用授权码，而不是 QQ 密码！

## 开始编码

### 发送第一个简单的邮件

首先从一个简单的代码开始

```go
func sendEmailBYQQEmailTest(to string) error {
	from := "2493325754@qq.com"
	password := "kfpjhmkeiykmebec" // 邮箱授权码
	smtpServer := "smtp.qq.com:465"
	// 邮件内容
	msg := []byte("From: Sender Name <" + from + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: test email\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		"this is a test email")

	// 设置 PlainAuth
	// 第一个 "" 可以看作一个可选参数，多数情况下不需要设置，传空即可。
	// 它的存在是为了满足 SMTP 标准协议中的扩展需求，但实际应用中很少需要自定义。
	auth := smtp.PlainAuth("", from, password, "smtp.qq.com")

	// 创建 tls 配置
	// InsecureSkipVerify: true：表示跳过对服务器证书的验证。这在生产环境中是不安全的，通常只在开发或测试环境中使用。
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.qq.com",
	}

	// 连接到 SMTP 服务器
	conn, err := tls.Dial("tcp", smtpServer, tlsconfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, "smtp.qq.com")
	if err != nil {
		return fmt.Errorf("SMTP 客户端创建失败: %v", err)
	}
	defer client.Quit()

	// 使用 auth 进行认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("认证失败: %v", err)
	}

	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("发件人设置失败: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("收件人设置失败: %v", err)
	}

	// 写入邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("数据写入失败: %v", err)
	}
	defer wc.Close()

	_, err = wc.Write(msg)
	if err != nil {
		return fmt.Errorf("消息发送失败: %v", err)
	}

	return nil
}
```

结果如下：

![](/img/blog/GoSendEmail/4.png)

#### 解释一下：

1. **创建 SMTP 客户端**：使用 `smtp.NewClient` 创建一个新的 SMTP 客户端。
2. **认证**：使用 `client.Auth` 方法进行认证。
3. **设置发件人和收件人**：使用 `client.Mail` 和 `client.Rcpt` 方法设置发件人和收件人。
4. **写入邮件内容**：使用 `client.Data` 方法获取一个 `io.WriteCloser`，然后使用 `Write` 方法将邮件内容写入。

在这些步骤中，`client.Data`方法实际上已经开始了邮件的发送过程，而`Write`方法则将邮件内容写入到 SMTP 服务器中。

### 怎么调整格式？

我们希望我们发的邮件不仅仅是几个文字，我们希望这些文字有一定的格式，例如：一级标题，二级标题，居中等等。

所以我们采用将正文调整成 html 的样式，这样渲染的时候就可以实现我们想要的格式。

```go
func sendEmailBYQQEmailAndFormat(to string) error {
	from := "2493325754@qq.com"
	password := "kfpjhmkeiykmebec" // 邮箱授权码
	smtpServer := "smtp.qq.com:465"
	// 读取图片
	imgPath := "./image.webp"
	imgData, err := ioutil.ReadFile(imgPath)
	if err != nil {
		log.Fatalf("无法读取图片: %v", err)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	// 邮件内容
	body := `
		<h1>这是一级标题</h1>
		<h2>这是二级标题</h2>
		<p>这是 <strong>` + `加粗` + `</strong></p>
		<p>这是 <em>` + `斜体` + `</em></p>
		<p>这是 <u>` + `下划线` + `</u></p>
		<p>这是 <s>` + `删除线` + `</s></p>
		<p>下面是一张图片</p>
		<img src="cid:image001" alt="image" width="180" height="180">
	`

	// 邮件头部
	header := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      "这是一封测试邮件",
		"MIME-Version": "1.0",
		"Content-Type": `multipart/related; boundary="BOUNDARY"`,
	}

	var message bytes.Buffer
	// 添加头部
	for k, v := range header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")

	// 添加正文部分
	message.WriteString("--BOUNDARY\r\n")
	message.WriteString(`Content-Type: text/html; charset="UTF-8"` + "\r\n\r\n")
	message.WriteString(body + "\r\n")

	// 添加图片部分
	message.WriteString("--BOUNDARY\r\n")
	message.WriteString("Content-Type: image/webp\r\n")
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString("Content-ID: <image001>\r\n\r\n")
	message.WriteString(imgBase64 + "\r\n")

	// MIME 结束
	message.WriteString("--BOUNDARY--")

	// 设置 PlainAuth
	// 第一个 "" 可以看作一个可选参数，多数情况下不需要设置，传空即可。
	// 它的存在是为了满足 SMTP 标准协议中的扩展需求，但实际应用中很少需要自定义。
	auth := smtp.PlainAuth("", from, password, "smtp.qq.com")

	// 创建 tls 配置
	// InsecureSkipVerify: true：表示跳过对服务器证书的验证。这在生产环境中是不安全的，通常只在开发或测试环境中使用。
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.qq.com",
	}

	// 连接到 SMTP 服务器
	conn, err := tls.Dial("tcp", smtpServer, tlsconfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, "smtp.qq.com")
	if err != nil {
		return fmt.Errorf("SMTP 客户端创建失败: %v", err)
	}
	defer client.Quit()

	// 使用 auth 进行认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("认证失败: %v", err)
	}

	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("发件人设置失败: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("收件人设置失败: %v", err)
	}

	// 写入邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("数据写入失败: %v", err)
	}
	defer wc.Close()

	// 发送邮件
	_, err = wc.Write(message.Bytes())
	if err != nil {
		return fmt.Errorf("消息发送失败: %v", err)
	}
	return nil
}
```

运行结果：

![5](/img/blog/GoSendEmail/5.png)

#### 解释一下：

1. `cid`是什么？
    
    `cid`是`Content-ID`的缩写，用于标识 MIME 消息(Multipurpose Internet Mail Extensions，多用途互联网邮件扩展 是一种互联网标准，最初设计用于扩展电子邮件的功能，使其支持不仅仅是纯文本内容，还可以包含多种格式的内容)中的资源。
    
    它是一个唯一的标识符，通常通过 HTML 中的 `<img>` 或其他标签引用。例如：表示邮件正文中的图片资源，其 `Content-ID` 为 `image001`。
    
2. 为什么需要`cid`？
    
    邮件客户端默认会阻止外部图片（`<img src="https://...">`）的加载，除非用户明确允许。而使用 `cid` 将图片嵌入到邮件中，可以避免外部图片的加载限制，确保图片能够直接显示。
    

### 怎么添加附件？

```go
func sendEmailByQQEmailAndAppendix(to string) error {
	from := "2493325754@qq.com"
	password := "kfpjhmkeiykmebec" // 邮箱授权码
	smtpServer := "smtp.qq.com:465"
	code := fmt.Sprintf("%06d", rand.Intn(900000)+100000) // 生成6位随机验证码
	attachmentPath := "test.txt"                          // 附件路径
	subject := "Verification Code"
	body := `
		<h1>Verification Code</h1>
		<p>Your verification code is: <strong>` + code + `</strong></p>
		<p>This verification code is valid for 15 minutes</p>
		<p>If you are not doing it yourself, please ignore it !</p>
	`

	// 创建 MIME 消息
	var msg bytes.Buffer
	writer := multipart.NewWriter(&msg)

	// 设置邮件头
	msg.WriteString("From: Sender Name <" + from + ">\r\n")
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: multipart/mixed; boundary=" + writer.Boundary() + "\r\n")
	msg.WriteString("\r\n")

	// 添加邮件正文
	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	if err != nil {
		return fmt.Errorf("创建邮件正文失败: %v", err)
	}
	part.Write([]byte(body))

	// 添加附件
	attachment, err := os.Open(attachmentPath)
	if err != nil {
		return fmt.Errorf("打开附件失败: %v", err)
	}
	defer attachment.Close()

	part, err = writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":        {"application/octet-stream"},
		"Content-Disposition": {"attachment; filename=\"" + filepath.Base(attachmentPath) + "\""},
	})
	if err != nil {
		return fmt.Errorf("创建附件部分失败: %v", err)
	}
	if _, err = io.Copy(part, attachment); err != nil {
		return fmt.Errorf("复制附件内容失败: %v", err)
	}

	// 设置 PlainAuth
	auth := smtp.PlainAuth("", from, password, "smtp.qq.com")

	// 创建 tls 配置
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.qq.com",
	}

	// 连接到 SMTP 服务器
	conn, err := tls.Dial("tcp", smtpServer, tlsconfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, "smtp.qq.com")
	if err != nil {
		return fmt.Errorf("SMTP 客户端创建失败: %v", err)
	}
	defer client.Quit()

	// 使用 auth 进行认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("认证失败: %v", err)
	}

	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("发件人设置失败: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("收件人设置失败: %v", err)
	}

	// 写入邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("数据写入失败: %v", err)
	}
	defer wc.Close()

	writer.Close()

	// 发送邮件
	_, err = wc.Write(msg.Bytes())
	if err != nil {
		return fmt.Errorf("消息发送失败: %v", err)
	}
	return nil
}
```

运行结果：

![](/img/blog/GoSendEmail/6.png)

---

## 有没有其他方法

主播主播，你的方法确实强，但太吃操作了，有没有更加简单又强势的方法推荐一下？有的兄弟有的！这么强的方法当然是不止一个，一共有九位，都是当前版本T0.5的强势方法。掌握一到两个方法，当个小皇帝都没有问题……

- `gomail`包
- `email`包