---
type: Post
title: Go语言实战数据库
tags: Go
category: 开发
category_bar: true
abbrlink: 34030
date: 2024-03-23 23:23:56
---

## 连接数据库

### 下载依赖

```go
go get -u github.com/go-sql-driver/mysql
```

### 导入数据驱动

```go
import (
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)
```

`_ "github.com/go-sql-driver/mysql"` 的作用是导入mysql驱动包，并执行该包的初始化代码，以便在后续的数据库操作中可以使用该驱动。但由于我们可能并不直接在代码中使用该包内的标识符（例如函数或类型），因此使用下划线来表示不需要直接访问该包内的内容。如果你在代码中确实需要使用包内的标识符，那么就不应该使用下划线，而是直接导入该包并在代码中使用它。

### 链接数据库

```go
func main() {
    //sql.Open(驱动名,数据源) (*DB,err)
    //数据源："用户名:密码@[连接方式](主机名:端口号)/数据库名"
    db, err1 := sql.Open("mysql", "root:040906@(localhost:3306)/library")
    if err1 != nil {
        fmt.Printf("db error is %s", err1)
        return
    }
    defer db.Close()

    err2 := db.Ping()
    if err2 != nil {
        fmt.Printf("数据库连接失败 %s", err2)
        return
    }
}
```

`defer` 关键字用于延迟执行函数调用，如果不使用会导致报错：`数据库连接失败 sql: database is closed`，如果不想使用 defer，依然可以将 db.Close( ) 放在结尾。

另外，你还可以使用 `DB.SetMaxIdleConns` 和 `DB.SetMaxOpenConns` 方法来设置最大空闲连接数和最大打开连接数，以防止数据库连接超过最大限制而被关闭。

sql.Open( )并不会与数据库正真的连接，而是在实际需要执行查询或操作时才会建立连接。它的返回值是一个 `*sql.DB` 类型的对象，代表了数据库连接池。db.Ping( ): 这个方法用于检查当前连接是否有效。在数据库操作之前，有时需要确保连接是可用的，以防止执行操作时出现意外错误。db.Ping() 方法会尝试与数据库建立连接，如果连接成功则表示连接有效，否则会返回错误信息。

## 执行sql语句

Go将数据库操作分为两类：Query与Exec

- Query表示查询，它会从数据库获取查询结果（一系列行，可能为空）。
- Exec表示执行语句，它不会返回行。

常见数据库操作模式：

- QueryRow只返回一行的查询，作为Query的一个常见特例。
- Prepare准备一个需要多次使用的语句，供后续执行用。

### 查询

查询分为两种，一种是单行查询`QueryRow()`，一种是多行查询`Query()` 。QueryRow( )总是返回非nil的值，直到返回值的Scan方法被调用时，才会返回被延迟的错误，如：未找到结果。多行查询Query( )执行一次查询，返回多行结果，一般用于执行select命令。

```go
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

type Book struct {
    Id     int    `json:"id"`
    Name   string `json:"name"`
    Author string `json:"author"`
}

var book Book

func main() {
    db, err1 := sql.Open("mysql", "root:040906@(localhost:3306)/library")
    if err1 != nil {
        fmt.Printf("open database failed , error is %s", err1)
        return
    }
    db.SetMaxIdleConns(10)     //设置数据库最大连接数量
    db.SetConnMaxLifetime(100) //设置数据库空闲时最大连接数量

    err2 := db.Ping()
    if err2 != nil {
        fmt.Printf("connect database failed , error is %s", err2)
        return
    }

    rows, err3 := db.Query("select * from books where id in (1,2,3)")
    if err3 != nil {
        fmt.Printf("fetching book error is %s\n", err3)
        return
    }
    for rows.Next() {
        e := rows.Scan(&book.Id, &book.Name, &book.Author)
        if e == nil {
            fmt.Printf("Book: %+v\n", book)
        }
    }
    err4 := rows.Close()
    if err4 != nil {
        fmt.Printf("failed to close rows %s", err4)
        return
    }

    //单行查询
    err5 := db.QueryRow("select * from books where id=4").Scan(&book.Id, &book.Name, &book.Author)
    if err5 != nil {
        fmt.Printf("fetching book error is %s\n", err5)
        return
    }
    fmt.Printf("Book: %+v\n", book)

    //单行查询
    err6 := db.QueryRow("select * from books where id=5").Scan(&book.Id, &book.Name, &book.Author)
    if err6 != nil {
        fmt.Printf("fetching book error is %s\n", err6)
        return
    }
    fmt.Printf("Book: %+v\n", book)
}
```

说明：

1. 使用db.Query()来发送查询到数据库，获取结果集Rows，并检查错误。
2. 使用rows.Next()作为循环条件，迭代读取结果集。
3. 使用rows.Scan从结果集中获取一行结果。
4. 使用rows.Close()关闭结果集，释放连接。

![](/img/blog/GoSql/1.png)

输出：

```go
Book: {Id:1 Name:红楼梦 Author:曹雪芹}
Book: {Id:2 Name:水浒传 Author:施耐庵}
Book: {Id:3 Name:三国演义 Author:曹雪芹}
Book: {Id:4 Name:百年孤独 Author:加西亚·马尔克斯}
fetching book error is sql: Scan error on column index 2, name "author": converting NULL to string is unsupported
```

下面以函数的形式给出：

```go
// 单行查询
func my_queryrow(a int) {
    sqlStr := "select id, book_name, author from books where id=?"
    // 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
    err := db.QueryRow(sqlStr, a).Scan(&book.Id, &book.Name, &book.Author)
    if err != nil {
        fmt.Printf("scan failed, err:%v\n", err)
        return
    }
    fmt.Printf("id:%d name:%s age:%s\n", book.Id, book.Name, book.Author)
}

// 多行查询
func my_query(a int) {
    sqlStr := "select id, book_name, author from books where id> ?"
    rows, err := db.Query(sqlStr, a)
    if err != nil {
        fmt.Printf("query failed, err:%v\n", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&book.Id, &book.Name, &book.Author)
        if err != nil {
            fmt.Printf("scan failed, err:%v\n", err)
            return
        }
        fmt.Printf("id:%d name:%s age:%s\n", book.Id, book.Name, book.Author)
    }
}
```

增删改用的都是Exec方法，Exec执行一次命令（包括查询、删除、更新、插入等），返回的Result是对已执行的SQL命令的总结。

### 插入

```go
// 插入数据
func my_insert(i int, n string, a string) {
    sqlStr := "insert into books (id,book_name, author) values (?,?,?)"
    ret, err := db.Exec(sqlStr, i, n, a)
    if err != nil {
        fmt.Printf("insert failed, err:%v\n", err)
        return
    }

    // 获取新插入数据的id，可以不写
    last_id, err := ret.LastInsertId()
    if err != nil {
        fmt.Printf("get lastinsert ID failed, err:%v\n", err)
        return
    }
    fmt.Printf("insert success, the id is %d.\n", last_id)
}
```

### 更新

```go
// 更新数据
func my_update_author(a string, i int) {
    sqlStr := "update books set author=? where id = ?"
    ret, err := db.Exec(sqlStr, a, i)
    if err != nil {
        fmt.Printf("update failed, err:%v\n", err)
        return
    }

    // 操作影响的行数，可以不写
    newupdate, err := ret.RowsAffected()
    if err != nil {
        fmt.Printf("get RowsAffected failed, err:%v\n", err)
        return
    }
    fmt.Printf("update success, affected rows:%d\n", newupdate)
}
```

### 删除

```go
// 更新数据
func my_update_author(a string, i int) {
    sqlStr := "update books set author=? where id = ?"
    ret, err := db.Exec(sqlStr, a, i)
    if err != nil {
        fmt.Printf("update failed, err:%v\n", err)
        return
    }
    newupdate, err := ret.RowsAffected() // 操作影响的行数
    if err != nil {
        fmt.Printf("get RowsAffected failed, err:%v\n", err)
        return
    }
    fmt.Printf("update success, affected rows:%d\n", newupdate)
}
```

### 完整代码

```go
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

type Book struct {
    Id     int    `json:"id"`
    Name   string `json:"name"`
    Author string `json:"author"`
}

var book Book
var err1 error
var db *sql.DB

// 单行查询
func my_queryrow(a int) {
    sqlStr := "select id, book_name, author from books where id=?"
    // 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
    err := db.QueryRow(sqlStr, a).Scan(&book.Id, &book.Name, &book.Author)
    if err != nil {
        fmt.Printf("scan failed, err:%v\n", err)
        return
    }
    fmt.Printf("id:%d name:%s author:%s\n", book.Id, book.Name, book.Author)
}

// 多行查询
func my_query(a int) {
    sqlStr := "select id, book_name, author from books where id> ?"
    rows, err := db.Query(sqlStr, a)
    if err != nil {
        fmt.Printf("query failed, err:%v\n", err)
        return
    }
    defer rows.Close()

    // 循环读取结果集中的数据
    for rows.Next() {
        err := rows.Scan(&book.Id, &book.Name, &book.Author)
        if err != nil {
            fmt.Printf("scan failed, err:%v\n", err)
            return
        }
        fmt.Printf("id:%d name:%s author:%s\n", book.Id, book.Name, book.Author)
    }
}

// 插入数据
func my_insert(i int, n string, a string) {
    sqlStr := "insert into books (id,book_name, author) values (?,?,?)"
    ret, err := db.Exec(sqlStr, i, n, a)
    if err != nil {
        fmt.Printf("insert failed, err:%v\n", err)
        return
    }
    last_id, err := ret.LastInsertId() // 新插入数据的id
    if err != nil {
        fmt.Printf("get lastinsert ID failed, err:%v\n", err)
        return
    }
    fmt.Printf("insert success, the id is %d.\n", last_id)
}

// 更新数据
func my_update_author(a string, i int) {
    sqlStr := "update books set author=? where id = ?"
    ret, err := db.Exec(sqlStr, a, i)
    if err != nil {
        fmt.Printf("update failed, err:%v\n", err)
        return
    }
    newupdate, err := ret.RowsAffected() // 操作影响的行数
    if err != nil {
        fmt.Printf("get RowsAffected failed, err:%v\n", err)
        return
    }
    fmt.Printf("update success, affected rows:%d\n", newupdate)
}

// 删除数据
func my_delete(i int) {
    sqlStr := "delete from books where id = ?"
    ret, err := db.Exec(sqlStr, i)
    if err != nil {
        fmt.Printf("delete failed, err:%v\n", err)
        return
    }
    // 操作影响的行数，可以不写
    newdelete, err := ret.RowsAffected()
    if err != nil {
        fmt.Printf("get RowsAffected failed, err:%v\n", err)
        return
    }
    fmt.Printf("delete success, affected rows:%d\n", newdelete)
}
func main() {
    db, err1 = sql.Open("mysql", "root:040906@(localhost:3306)/library")
    if err1 != nil {
        fmt.Printf("open database failed , err:%v", err1)
        return
    }
    db.SetMaxIdleConns(10)     //设置数据库最大连接数量
    db.SetConnMaxLifetime(100) //设置数据库空闲时最大连接数量

    err2 := db.Ping()
    if err2 != nil {
        fmt.Printf("connect database failed , err: %v", err2)
        return
    }

    for i := 1; i <= 3; i++ {
        my_queryrow(i)
    }
    my_insert(9, "世界通史", "unkown")
    my_query(8)
    my_update_author("L·S·斯塔夫里阿诺斯", 9)
    my_query(8)
    my_delete(9)
}
```

运行结果：

![](/img/blog/GoSql/2.png)

## 预处理

### 什么是预处理

普通sql执行：

1. 客户端对SQL语句进行占位符替换得到完整的SQL语句。
2. 客户端发送完整SQL语句到MySQL服务端
3. MySQL服务端执行完整的SQL语句并将结果返回给客户端。

预处理：

1. 把SQL语句分成两部分，命令部分与数据部分。
2. 先把命令部分发送给MySQL服务端，MySQL服务端进行SQL预处理。
3. 然后把数据部分发送给MySQL服务端，MySQL服务端对SQL语句进行占位符替换。
4. MySQL服务端执行完整的SQL语句并将结果返回给客户端。

### **为什么要预处理**

1. 优化MySQL服务器重复执行SQL的方法，可以提升服务器性能。
2. 避免SQL注入问题。

### **Go实现MySQL预处理**

`func (db *DB) Prepare(query string) (*Stmt, error)`

预处理和普通执行在代码上就多出一部分，下面仅以查询作展示：

```go
// 预处理查询示例
func prepareQuery(a int) {
    sqlStr := "select id, book_name, author from books where id > ?"
    stmt, err := db.Prepare(sqlStr)
    if err != nil {
        fmt.Printf("prepare failed, err:%v\n", err)
        return
    }
    defer stmt.Close()
    rows, err := stmt.Query(a)
    if err != nil {
        fmt.Printf("query failed, err:%v\n", err)
        return
    }
    defer rows.Close()
    // 循环读取结果集中的数据
    for rows.Next() {
        var book Book
        err := rows.Scan(&book.Id, &book.Name, &book.Author)
        if err != nil {
            fmt.Printf("scan failed, err:%v\n", err)
            return
        }
        fmt.Printf("id:%d book_name:%s author:%s\n", book.Id, book.Name, book.Author)
    }
}
```

## 事务

事务是一个最小的不可再分的工作单元。一个事务中的所有操作，要么全部完成，要么全部不完成，不会结束在中间某个环节。事务在执行过程中发生错误，会被回滚到事务开始前的状态，就像这个事务从来没有执行过一样。事务处理结束后，对数据的修改就是永久的，即便系统故障也不会丢失。

```go
func (db *DB) Begin() (*Tx, error)  //开启事务
func (tx *Tx) Commit() error        //提交事务
func (tx *Tx) Rollback() error      //回滚事务
```

示例：

```go
// 事务操作示例
func transaction() {
    tx, err := db.Begin() // 开启事务
    if err != nil {
        if tx != nil {
            tx.Rollback() // 回滚
        }
        fmt.Printf("begin trans failed, err:%v\n", err)
        return
    }
    sqlStr1 := "Update books set author='serendipity' where id=?"
    ret1, err := tx.Exec(sqlStr1, 6)
    if err != nil {
        tx.Rollback() // 回滚
        fmt.Printf("exec sql1 failed, err:%v\n", err)
        return
    }
    affRow1, err := ret1.RowsAffected()
    if err != nil {
        tx.Rollback() // 回滚
        fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
        return
    }

    sqlStr2 := "Update books set author='unknow' where id=?"
    ret2, err := tx.Exec(sqlStr2, 8)
    if err != nil {
        tx.Rollback() // 回滚
        fmt.Printf("exec sql2 failed, err:%v\n", err)
        return
    }
    affRow2, err := ret2.RowsAffected()
    if err != nil {
        tx.Rollback() // 回滚
        fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
        return
    }

    fmt.Println(affRow1, affRow2)
    if affRow1 == 1 && affRow2 == 1 {
        tx.Commit() // 提交事务
    } else {
        tx.Rollback() //回滚
        fmt.Println("something wrong")
    }

    fmt.Println("exec trans success!")
}
```
