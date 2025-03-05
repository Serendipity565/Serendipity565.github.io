---
type: Post
title: Go语言操作Redis
tags: Go
category: 开发
category_bar: true
abbrlink: 8470
date: 2024-04-17 17:23:53
---

## Redis介绍

Redis是一个开源的内存数据库，Redis提供了多种不同类型的数据结构，很多业务场景下的问题都可以很自然地映射到这些数据结构上。除此之外，通过复制、持久化和客户端分片等特性，我们可以很方便地将Redis扩展成一个能够包含数百GB数据、每秒处理上百万次请求的系统。

### Redis支持的数据结构

Redis支持诸如字符串（string）、哈希（hashe）、列表（list）、集合（set）、带范围查询的排序集合（sorted set）、bitmap、hyperloglog、带半径查询的地理空间索引（geospatial index）和流（stream）等数据结构。

### Redis应用场景

- 缓存系统，减轻主数据库（MySQL）的压力。
- 计数场景，比如微博、抖音中的关注数和粉丝数。
- 热门排行榜，需要排序的场景特别适合使用ZSET。
- 利用 LIST 可以实现队列的功能。
- 利用 HyperLogLog 统计UV、PV等数据。
- 使用 geospatial index 进行地理位置相关查询。

## go-redis库

### 安装

Go 社区中目前有很多成熟的 redis client 库，比如[https://github.com/gomodule/redigo](https://github.com/gomodule/redigo) 和[https://github.com/redis/go-redis](https://github.com/redis/go-redis) 等。本文使用 go-redis 这个库来操作 Redis 数据库。

使用以下命令下安装 go-redis 库。

安装`v8`版本：

```go
go get github.com/redis/go-redis/v8
```

安装`v9`版本：

```go
go get github.com/redis/go-redis/v9
```

### 连接

在项目中导入 `go-redis`库（以`v9`版本为例）。

```go
import "github.com/redis/go-redis/v9"
```

#### 普通连接模式

go-redis 库中使用 redis.NewClient 函数连接 Redis 服务器。

```go
rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // 密码
    DB:       0,  // 数据库
    PoolSize: 20, // 连接池大小
})
```

除此之外，还可以使用 redis.ParseURL 函数从表示数据源的字符串中解析得到 Redis 服务器的配置信息。

```go
opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
if err != nil {
    panic(err)
}
rdb := redis.NewClient(opt)
```

#### TLS连接模式

如果使用的是 TLS 连接方式，则需要使用 tls.Config 配置。

```go
rdb := redis.NewClient(&redis.Options{
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
        // Certificates: []tls.Certificate{cert},
    // ServerName: "your.domain.com",
    },
})
```

#### Redis Sentinel模式

使用下面的命令连接到由 Redis Sentinel 管理的 Redis 服务器。

```go
rdb := redis.NewFailoverClient(&redis.FailoverOptions{
    MasterName:    "master-name",
    SentinelAddrs: []string{":9126", ":9127", ":9128"},
})
```

#### Redis Cluster模式

使用下面的命令连接到 Redis Cluster，go-redis 支持按延迟或随机路由命令。

```go
rdb := redis.NewClusterClient(&redis.ClusterOptions{
    Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},

    // 若要根据延迟或随机路由命令，请启用以下命令之一
    // RouteByLatency: true,
    // RouteRandomly: true,
})
```

## 基本使用

### 执行命令

下面的示例代码演示了 go-redis 库的基本使用。

```go
func doCommand() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    // 执行命令获取结果
    val, err := rdb.Get(ctx, "key").Result()
    fmt.Println(val, err)

    // 先获取到命令对象
    cmder := rdb.Get(ctx, "key")
    fmt.Println(cmder.Val()) // 获取值
    fmt.Println(cmder.Err()) // 获取错误

    // 直接执行命令获取错误
    err = rdb.Set(ctx, "key", 10, time.Hour).Err()

    // 直接执行命令获取值
    value := rdb.Get(ctx, "key").Val()
    fmt.Println(value)
}
```

- `context.Background()`：这是创建一个空的上下文（context），作为其他上下文的根。
- `context.WithTimeout()`：这是创建一个带有超时时间的上下文（context）。它接受一个父上下文和一个超时时间作为参数，并返回一个派生的上下文（context），该派生上下文将在超时时间到达时自动取消。
- `500*time.Millisecond`：这是超时时间，表示 500 毫秒。在这个例子中，如果操作没有在 500 毫秒内完成，上下文将会被取消。
- `defer cancel()`：这是在函数结束时调用 `cancel()` 函数，以确保在函数执行完毕后及时取消上下文。这个 `cancel()` 函数用于取消与这个上下文相关联的所有操作，释放资源，避免资源泄漏的发生。

这段代码的作用是创建一个具有 500 毫秒超时的上下文，这样在执行操作时，如果操作在 500 毫秒内没有完成，就会自动取消。

### 执行任意命令

go-redis 还提供了一个执行任意命令或自定义命令的 Do 方法，特别是一些 go-redis 库暂时不支持的命令都可以使用该方法执行。具体使用方法如下。

```go
func doDemo() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    // 直接执行命令获取错误
    err := rdb.Do(ctx, "set", "key", 10, "EX", 3600).Err()
    fmt.Println(err)

    // 执行命令获取结果
    val, err := rdb.Do(ctx, "get", "key").Result()
    fmt.Println(val, err)
}
```

### redis.Nil

go-redis 库提供了一个 redis.Nil 错误来表示 Key 不存在的错误。因此在使用 go-redis 时需要注意对返回错误的判断。在某些场景下我们应该区别处理 redis.Nil 和其他不为 nil 的错误。

```go
func getValueFromRedis(key, defaultValue string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    val, err := rdb.Get(ctx, key).Result()
    if err != nil {
        // 如果返回的错误是key不存在
        if errors.Is(err, redis.Nil) {
            return defaultValue, nil
        }
        // 出其他错了
        return "", err
    }
    return val, nil
}
```

## 其他示例

### zset示例

下面的示例代码演示了如何使用 go-redis 库操作 zset。

```go
func zsetDemo() {
    // key
    zsetKey := "language_rank"
    // value
    // 注意：v8版本使用[]*redis.Z；v9版本使用[]redis.Z
    languages := []redis.Z{
        {Score: 90.0, Member: "Golang"},
        {Score: 98.0, Member: "Java"},
        {Score: 95.0, Member: "Python"},
        {Score: 97.0, Member: "JavaScript"},
        {Score: 99.0, Member: "C/C++"},
    }
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    // ZADD
    err := rdb.ZAdd(ctx, zsetKey, languages...).Err()
    if err != nil {
        fmt.Printf("zadd failed, err:%v\n", err)
        return
    }
    fmt.Println("zadd success")

    // 把Golang的分数加10
    newScore, err := rdb.ZIncrBy(ctx, zsetKey, 10.0, "Golang").Result()
    if err != nil {
        fmt.Printf("zincrby failed, err:%v\n", err)
        return
    }
    fmt.Printf("Golang's score is %f now.\n", newScore)

    // 取分数最高的3个
    ret := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Val()
    for _, z := range ret {
        fmt.Println(z.Member, z.Score)
    }

    // 取95~100分的
    op := &redis.ZRangeBy{
        Min: "95",
        Max: "100",
    }
    ret, err = rdb.ZRangeByScoreWithScores(ctx, zsetKey, op).Result()
    if err != nil {
        fmt.Printf("zrangebyscore failed, err:%v\n", err)
        return
    }
    for _, z := range ret {
        fmt.Println(z.Member, z.Score)
    }
}
```

执行上面的函数将得到如下输出结果。

```go
zadd success
Golang's score is 100.000000 now.
Golang 100
C/C++ 99
Java 98
Python 95
JavaScript 97
Java 98
C/C++ 99
Golang 100
```

### 扫描或遍历所有key

在Redis中可以使用[`KEYS prefix*`](https://redis.io/commands/keys/) 命令按前缀查询所有符合条件的 key，`go-redis`库中提供了`Keys`方法实现类似查询key的功能。

例如使用以下命令查询以`user:`为前缀的所有key（`user:cart:00`、`user:order:2023`等）。

`vals, err **:=** rdb.Keys(ctx, "user:*").Result()`

但是如果需要扫描数百万的 key ，那速度就会比较慢。这种场景下你可以使用`Scan`命令来遍历所有符合要求的 key。

```go
// scanKeysDemo1 按前缀查找所有key示例
func scanKeysDemo1() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    var cursor uint64
    for {
        var keys []string
        var err error
        // 将redis中所有以prefix:为前缀的key都扫描出来
        keys, cursor, err = rdb.Scan(ctx, cursor, "prefix:*", 0).Result()
        if err != nil {
            panic(err)
        }

        for _, key := range keys {
            fmt.Println("key", key)
        }
        
        // 扫描不到前缀为prefix:的key
        if cursor == 0 {
            break
        }
    }
}
```

- **`Scan()`**：这是 Redis 客户端库中用于执行 SCAN 命令的方法。它接受一个上下文（context）对象作为第一个参数，然后是游标（cursor）、匹配模式（pattern）和 COUNT 参数。
- **`cursor`**：这是 SCAN 命令中的游标参数，用于标识当前迭代的位置。在第一次调用时，通常设置为 0，后续调用会返回新的游标值，以便进行下一次迭代。
- **`0`**：这是 SCAN 命令中的 COUNT 参数，用于指定每次迭代返回的最大元素数量。0 表示返回所有匹配的键。

**为什么要写成两个for循环？**

外部的 **`for`** 循环是用来处理整个 SCAN 过程的迭代。在每次迭代中，它执行一次 SCAN 命令，获取一批匹配的键，并处理这批键。然后，它检查游标值是否为 0。如果游标为 0，表示已经扫描完所有的键，就退出循环，结束整个 SCAN 过程。如果游标不为 0，表示还有更多的键需要扫描，就继续下一次迭代，执行下一次 SCAN 命令。

**游标是什么？**

游标（cursor）在 Redis 中用于处理 SCAN 命令，它是一种用于分页扫描大量键的机制。SCAN 命令可以用于迭代遍历 Redis 数据库中的键，而不会阻塞服务器，因此在处理大量键时非常有用。使用游标可以将扫描结果分批返回，避免一次性返回大量数据给客户端，减轻客户端和服务器的压力。

针对这种需要遍历大量key的场景，`go-redis`中提供了一个简化方法——`Iterator`，其使用示例如下。

```go
// scanKeysDemo2 按前缀扫描key示例
func scanKeysDemo2() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    // 按前缀扫描key
    iter := rdb.Scan(ctx, 0, "prefix:*", 0).Iterator()
    for iter.Next(ctx) {
        fmt.Println("keys", iter.Val())
    }
    if err := iter.Err(); err != nil {
        panic(err)
    }
}
```

例如，我们可以写出一个将所有匹配指定模式的 key 删除的示例。

```go
// delKeysByMatch 按match格式扫描所有key并删除
func delKeysByMatch(match string, timeout time.Duration) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    iter := rdb.Scan(ctx, 0, match, 0).Iterator()
    for iter.Next(ctx) {
        err := rdb.Del(ctx, iter.Val()).Err()
        if err != nil {
            panic(err)
        }
    }
    if err := iter.Err(); err != nil {
        panic(err)
    }
}
```

此外，对于 Redis 中的 set、hash、zset 数据类型，`go-redis` 也支持类似的遍历方法。

```go
iter := rdb.SScan(ctx, "set-key", 0, "prefix:*", 0).Iterator()
iter := rdb.HScan(ctx, "hash-key", 0, "prefix:*", 0).Iterator()
iter := rdb.ZScan(ctx, "sorted-hash-key", 0, "prefix:*", 0).Iterator()
```

## Pipeline

Redis Pipeline 允许通过使用单个 client-server-client 往返执行多个命令来提高性能。区别于一个接一个地执行100个命令，你可以将这些命令放入 pipeline 中，然后使用1次读写操作像执行单个命令一样执行它们。这样做的好处是节省了执行命令的网络往返时间（RTT）。

下面的示例代码中演示了使用 pipeline 通过一个 write + read 操作来执行多个命令。

```go
pipe := rdb.Pipeline()

incr := pipe.Incr(ctx, "pipeline_counter")
pipe.Expire(ctx, "pipeline_counter", time.Hour)

cmds, err := pipe.Exec(ctx)
if err != nil {
    panic(err)
}

// 在执行pipe.Exec之后才能获取到结果
fmt.Println(incr.Val())
```

上面的代码相当于将以下两个命令一次发给 Redis Server 端执行，与不使用 Pipeline 相比能减少一次RTT。

```go
INCR pipeline_counter
EXPIRE pipeline_counts 3600
```

或者，你也可以使用`Pipelined` 方法，它会在函数退出时调用 Exec。

```go
var incr *redis.IntCmd

cmds, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
    incr = pipe.Incr(ctx, "pipelined_counter")
    pipe.Expire(ctx, "pipelined_counter", time.Hour)
    return nil
})
if err != nil {
    panic(err)
}

// 在pipeline执行后获取到结果
fmt.Println(incr.Val())
```

我们可以遍历 pipeline 命令的返回值依次获取每个命令的结果。下方的示例代码中使用pipiline一次执行了100个 Get 命令，在pipeline 执行后遍历取出100个命令的执行结果。

```go
cmds, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
    for i := 0; i < 100; i++ {
        pipe.Get(ctx, fmt.Sprintf("key%d", i))
    }
    return nil
})
if err != nil {
    panic(err)
}

for _, cmd := range cmds {
    fmt.Println(cmd.(*redis.StringCmd).Val())
}
```

在那些我们需要一次性执行多个命令的场景下，就可以考虑使用 pipeline 来优化。

## 事务

在 Redis 中，MULTI 是一个事务（transaction）命令，它用于标记一个事务的开始。在 MULTI 命令之后，所有后续的命令都会被添加到一个事务队列中，而不会立即执行。只有在 EXEC 命令被调用时，才会执行所有在 MULTI 和 EXEC 之间添加到事务队列中的命令。

MULTI 命令不接受任何参数，它只是一个简单的标记，表示后续的命令应该被视为一个事务的一部分。在调用 MULTI 命令后，Redis 服务器会进入事务状态，并在接收到 EXEC 命令时执行事务队列中的所有命令。

使用事务可以保证一系列的 Redis 命令在执行过程中不会被其他客户端的命令中断。如果在 MULTI 和 EXEC 之间的某个时间点发生了错误，Redis 会取消事务，并且事务队列中的所有命令都不会执行。

举个例子，以下是一个在 Redis 使用 MULTI 命令创建事务的示例：

```sql
MULTI
SET key1 value1
SET key2 value2
GET key1
EXEC
```

Redis 是单线程执行命令的，因此单个命令始终是原子的，但是来自不同客户端的两个给定命令可以依次执行，例如在它们之间交替执行。但是，`Multi/exec`能够确保在`multi/exec`两个语句之间的命令之间没有其他客户端正在执行命令。

在这种场景我们需要使用 TxPipeline 或 TxPipelined 方法将 pipeline 命令使用 `MULTI` 和`EXEC`包裹起来。

```go
// TxPipeline
pipe := rdb.TxPipeline()
incr := pipe.Incr(ctx, "tx_pipeline_counter")
pipe.Expire(ctx, "tx_pipeline_counter", time.Hour)
_, err := pipe.Exec(ctx)
fmt.Println(incr.Val(), err)

// TxPipelined
var incr2 *redis.IntCmd
_, err = rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
    incr2 = pipe.Incr(ctx, "tx_pipeline_counter")
    pipe.Expire(ctx, "tx_pipeline_counter", time.Hour)
    return nil
})
fmt.Println(incr2.Val(), err)
```

上面代码相当于在一个RTT下执行了下面的redis命令：

```go
MULTI
INCR pipeline_counter
EXPIRE pipeline_counts 3600
EXEC
```

**注意：**

1. **TxPipeline**：

    **`TxPipeline`** 是在事务中使用管道（pipeline）的方式。使用 **`TxPipeline`** 创建的事务对象可以执行多个命令，并将这些命令一次性发送给 Redis 服务器，然后一次性获取所有命令的响应。这样做可以减少网络往返次数，提高性能。但是，**`TxPipeline`** 中的每个命令都是原子执行的，即在执行期间不会中断事务。

2. **TxPipelined**：

    **`TxPipelined`** 是在事务中使用异步管道（pipeline）的方式。使用 **`TxPipelined`** 创建的事务对象允许在事务中的每个命令之间执行其他 Go 代码。这些命令被添加到事务队列中，并在调用 **`TxPipelined`** 对象的 **`Exec`** 方法时执行。这种方式允许在事务执行过程中执行其他任务，而不必等待事务执行完成。

例如：

```go
// 使用 TxPipeline 创建事务并执行命令
func exampleTxPipeline() {
    // 开启事务
    pipe := rdb.TxPipeline()
    
    // 在事务中执行多个命令
    pipe.Set("key1", "value1", 0)
    pipe.Set("key2", "value2", 0)
    pipe.Get("key1")
    
    // 执行事务
    _, err := pipe.Exec()
    if err != nil {
        panic(err)
    }
}

// 使用 TxPipelined 创建事务并执行命令
func exampleTxPipelined() {
    // 开启事务
    pipe := rdb.TxPipelined()
    
    // 在事务中执行多个命令
    pipe.Set("key1", "value1", 0)
    pipe.Set("key2", "value2", 0)
    pipe.Get("key1")
    
    // 在事务中执行其他任务，例如调用其他函数
    go func() {
        // 执行其他任务
        fmt.Println("Do something else while waiting for transaction to finish")
    }()
    
    // 执行事务
    _, err := pipe.Exec()
    if err != nil {
        panic(err)
    }
}
```

### Watch

我们通常搭配 `WATCH`命令来执行事务操作。从使用`WATCH`命令监视某个 key 开始，直到执行`EXEC`命令的这段时间里，如果有其他用户抢先对被监视的 key 进行了替换、更新、删除等操作，那么当用户尝试执行`EXEC`的时候，事务将失败并返回一个错误，用户可以根据这个错误选择重试事务或者放弃事务。

Watch方法接收一个函数和一个或多个key作为参数。

```go
Watch(fn func(*Tx) error, keys ...string) error
```

下面的代码片段演示了 Watch 方法搭配 TxPipelined 的使用示例。

```go
// watchDemo 在key值不变的情况下将其值+1
func watchDemo(ctx context.Context, key string) error {
    return rdb.Watch(ctx, func(tx *redis.Tx) error {
        n, err := tx.Get(ctx, key).Int()
        if err != nil && err != redis.Nil {
            return err
        }
        // 假设操作耗时5秒
        // 5秒内我们通过其他的客户端修改key，当前事务就会失败
        time.Sleep(5 * time.Second)
        _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
            pipe.Set(ctx, key, n+1, time.Hour)
            return nil
        })
        return err
    }, key)
}
```

将上面的函数执行并打印其返回值，如果我们在程序运行后的5秒内修改了被 watch 的 key 的值，那么该事务操作失败，返回`redis: transaction failed`错误。

最后我们来看一个 go-redis 官方文档中使用 `GET` 、`SET`和`WATCH`命令实现一个 INCR 命令的完整示例。

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/redis/go-redis/v9"
    "sync"
    "time"
)

func main() {
    opt, err := redis.ParseURL("redis://@localhost:6379/0")
    if err != nil {
        panic(err)
    }

    // 此处rdb为初始化的redis连接客户端
    rdb := redis.NewClient(opt)

    const routineCount = 100

    // 设置5秒超时
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // increment 是一个自定义对key进行递增（+1）的函数
    // 使用 GET + SET + WATCH 实现，类似 INCR
    increment := func(key string) error {
        txf := func(tx *redis.Tx) error {
            // 获得当前值或零值
            n, err := tx.Get(ctx, key).Int()
            if err != nil && !errors.Is(err, redis.Nil) {
                return err
            }

            // 实际操作（乐观锁定中的本地操作）
            n++ // 仅在监视的Key保持不变的情况下运行
            _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
                // pipe 处理错误情况
                pipe.Set(ctx, key, n, 0)
                return nil
            })
            return err
        }

        // 最多重试100次
        for retries := routineCount; retries > 0; retries-- {
            err := rdb.Watch(ctx, txf, key)
            if !errors.Is(err, redis.TxFailedErr) {
                return err
            }
            // 乐观锁丢失
        }
        return errors.New("increment reached maximum number of retries")
    }

    // 开启100个goroutine并发调用increment
    // 相当于对key执行100次递增
    var wg sync.WaitGroup
    wg.Add(routineCount)
    for i := 0; i < routineCount; i++ {
        go func() {
            defer wg.Done()

            if err := increment("counter3"); err != nil {
                fmt.Println("increment error:", err)
            }
        }()
    }
    wg.Wait()

    n, err := rdb.Get(ctx, "counter3").Int()
    fmt.Println("最终结果：", n, err)
}
```

在这个示例中使用了 `redis.TxFailedErr` 来检查事务是否失败。

更多详情请查阅[官方文档](https://redis.uptrace.dev/zh/)。

---

redis命令查询：[https://www.runoob.com/redis/redis-keys.html](https://www.runoob.com/redis/redis-keys.html)
