---
type: Post
title: MySql使用教程
tags: MySql
category: 开发
category_bar: true
abbrlink: 28400
date: 2024-03-20 21:23:03
---

**SQL** (Structured Query Language:结构化查询语言) 是用于管理关系数据库管理系统（RDBMS）。 SQL 的范围包括数据插入、查询、更新和删除，数据库模式创建和修改，以及数据访问控制。

注意：**MySQL 在Windows和MacOS系统下不区分大小写**，但在Linux 系统下默认区分大小写。但是为了方便使用，我们一般会将关键字全部大写。

## 创建

创建数据库：`CREATE DATABASE 数据库名称;`

创建数据库下的表格：`CREATE TABLE 数据表名称(列名1 数据类型,列名2 数据类型,...,列名n,数据类型);`

创建数据库：`CREATE DATABASE 数据库名称;` 创建数据库下的shu`CREATE TABLE 数据表名称;`

```sql
CREATE  DATABASE library;

USE library;

CREATE TABLE books(
    id INT AUTO_INCREMENT PRIMARY KEY,
    book_name VARCHAR(20) NOT NULL,
    author VARCHAR(20) NULL,
    d DATE NULL
);
```

USE library; 这一命令表示使用library这一数据库，适用于有多个数据库的情况。

![](/img/blog/SQL/1.png)

下面是一些常见的数据类型：

| 数据类型 | 说明 |
| --- | --- |
| varchar(最长255) | 可变长度的字符串，varchar(10)，10表示最大可分配空间，会根据传过来的数据动态分配。节省空间，但是需要动态分配空间，速度慢。 |
| char(最长255) | 定长字符串，可能会导致空间的浪费 |
| int(最长11) | 数字中的长整型 |
| bigint | 数字中的长整型 |
| float | 单精度浮点型数据 |
| double | 双精度浮点型数据 |
| date | 短日期类型 |
| datetime | 长日期类型 |

对与创建列名的时候，可以设置一些相应的默认格式。比如：

| 格式 | 说明 |
| --- | --- |
| NOT NULL | 这一列不能为空 |
| NULL | 这一列可以为空 |
| PRIMARY KEY | 主键 |
| AUTO_INCREMENT | 自动增加数字 |
| UNIQUE | 不允许重复 |

## 插入

`INSERT INTO 数据库名.表格名 (列名1,列名2,…,列名n) VALUES (数值1,数值2,…,数值n)`

```sql
INSERT INTO library.books(id,book_name,author,d) VALUES (1,'红楼梦','曹雪芹',null);
INSERT INTO library.books(id,book_name,author,d) VALUES (2,'水浒传','施耐庵','2024-03-20');
INSERT INTO library.books(id,book_name,author,d) VALUES (3,'三国演义','曹雪芹','2024-03-20');
INSERT INTO library.books(id,book_name,author,d) VALUES (DEFAULT,'百年孤独','加西亚·马尔克斯','2020-05-06');
INSERT INTO library.books(id,book_name,author,d) VALUES (5,'月亮与六便士',NULL,'2020-05-06');
INSERT INTO library.books(id,book_name,author,d) VALUES (DEFAULT,'局外人',NULL,'2020-05-06');
INSERT INTO library.books(id,book_name,author,d) VALUES (DEFAULT,'史记',NULL,'2020-05-06');
INSERT INTO library.books(id,book_name,author,d) VALUES (DEFAULT,'中国通史',NULL,'2020-05-06');
```

这样我们就插入了一些数据

![](/img/blog/SQL/2.png)

其中，DEFAULT表示默认格式。

## 更新

我们也可以再增加一个列名：`ALTER TABLE 数据库名.表格名 ADD 列名 数据类型 默认条件`

更新具体数据：`UPDATE 数据库名.表格名 SET 值 WHERE 条件`

```sql
ALTER TABLE library.books ADD sold FLOAT NULL;
UPDATE library.books SET sold=20.2 WHERE id=1;
```

![](/img/blog/SQL/3.png)

## 删除

删除某条数据：`DELETE FROM 数据库名.表格名 WHERE 条件`

删除表格：`DROP TABLE 数据库名.表格名`

删除数据库：`DROP DATABASE 数据库名`

```sql
DELETE FROM library.books WHERE id=2;
```

![](/img/blog/SQL/4.png)

## 查找

查找全部内容：`SELECT * FROM 数据库名.表格名`

查找某一列内容：`SELECT 列名1，列名2 FROM 数据库名.表格名`

查找不同的内容（即去除重复内容）：`SELECT DISTINCT 列名1 FROM 数据库名.表格名`

查看时排序（默认ASC，即ascending，升序；DESC，descending，降序）：
`SELECT * FROM 数据库名.表格名 ORDER BY 列名 ASC/DESC`

要过滤掉某些信息：`SELECT * FROM 数据库名.表格名 WHERE 条件 ORDER BY 列名 ASC/DESC`

| 运算符 | 说明 |
| --- | --- |
| = | 等于 |
| ! = 或 <> | 不等于 |
| > | 大于 |
| < | 小于 |
| > = | 大于等于 |
| < = | 小于等于 |
| BETWEEN | 介于两者之间 |
| IN | 在一组值内 |
| LIKE | 相似匹配 |
| AND | 与 |
| OR | 或 |
| NOT 或 ! | 非 |

比如想知道这些书分别在哪几天入库：

![](/img/blog/SQL/5.png)

去重后：

![](/img/blog/SQL/6.png)

```sql
SELECT * FROM library.books WHERE author IN ('曹雪芹','施耐庵');
```

![](/img/blog/SQL/7.png)

在 SQL 中，通配符与 SQL LIKE 操作符一起使用，下面具体说说LIKE的用法：

| 通配符 | 描述 |
| --- | --- |
| % | 替代 0 个或多个字符 |
| _ | 替代一个字符 |
| [*charlist*] | 字符列中的任何单一字符 |
| [^*charlist*]或[!*charlist*] | 不在字符列中的任何单一字符 |

```sql
SELECT * FROM library.books WHERE author LIKE 'B%';   --查找名字以B开头的作者
SELECT * FROM library.books WHERE author LIKE '%b';   --查找名字以b开头的作者
SELECT * FROM library.books WHERE author LIKE '__b%';   --查找名字第三个字符是b的作者
SELECT * FROM library.books WHERE author LIKE '^[AB]';   --查找名字以A或B开头的作者
SELECT * FROM library.books WHERE author LIKE '^[A-H]';   --查找名字以A到H开头的作者
```
