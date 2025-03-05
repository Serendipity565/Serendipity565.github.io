---
type: Post
title: viper的“陷阱”
slug: viper的“陷阱”
tags: Go
category: 开发
category_bar: true
abbrlink: 10850
date: 2025-02-20 19:57:43
---

## 第一个问题：SetConfigType()真的有用吗？

### 问题引入

当你使用如下方式读取配置时，viper会从`./conf`目录下查找任何以`config`为文件名的配置文件，如果同时存在`./conf/config.json`和`./conf/config.yaml`两个配置文件的话，`viper`会从哪个配置文件加载配置呢？

```go
viper.SetConfigName("config")
viper.AddConfigPath("./conf")
```

#### 复现

下面的 demo 代码模拟这种情况

```go
package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Name string `mapstructure:"name"`
	Age  int    `mapstructure:"age"`
}

// NewConfig 创建并返回配置对象
func NewConfig() *Config {
	viper.AddConfigPath("./conf/") // 配置文件所在的路径
	viper.SetConfigName("config")  // 配置文件名称（不带扩展名）
	// viper.SetConfigType("yaml")    // 如果配置文件的名称中没有扩展名，则需要配置此项

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果读取配置文件失败，返回 nil
		fmt.Println("Error reading config file:", err)
		return nil
	}

	// 解组配置文件到结构体
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		// 解组失败
		fmt.Println("Error unmarshalling config:", err)
		return nil
	}

	// 返回解析后的配置
	return &conf
}

func main() {
	conf := NewConfig()

	fmt.Printf("Name: %s, Age: %d\n", conf.Name, conf.Age)
}
```

```yaml
name: "Serendipity_yaml"
age: 20
```

```json
{
  "name": "Serendipity_json",
  "age": 20
}
```

发现输出的是 `Name: Serendipity_json, Age: 20`

之后我们尝试加上 `viper.SetConfigType("yaml")` 这一行代码，运行发现结果并没有因为加上这一行而改变。

运行结果

![](/img/blog/viper/1.png)

### 分析

首先来看一下为什么在没有指明文件类型的时候会读取到 `json` 文件。

- `viper` 会按照文件系统的顺序查找文件，在你设置的路径下依次尝试加载 `config.json`、`config.yaml`、`config.toml` 等文件格式。
- **默认情况下，`.json` 文件会被优先加载**，如果同时存在 `config.json` 和 `config.yaml`，`viper` 会加载 `config.json` 文件。

明白了这一点后我们再来看为什么加上 `SetConfigType("yaml")` 后结果依旧不变。

我们来看一下 Viper 中 `viper.ReadInConfig()` 的源码

![](/img/blog/viper/2.png)

可以看到，只有在 `stringInSlice()` 中用到过 `v.getConfigType()`，也就是获取文件的种类，在读取文件的时候并没有做文件名称和种类的拼接，导致这个 `SetConfigType()` 并没有起到实质性确定文件类型的作用。

### 解决方案

```go
	viper.AddConfigPath("./conf/")
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
```

## 第二个问题：热更新配置一致性问题

### 问题引入

假设我们当前服务中有一条流水线操作，需要分别调用三个接口A、B、C才能完成对应的功能，其中配置文件存储了调用的接口名称。

现在我们需要将流水线要执行的流程从ABC换成DEF，在替换过程中就可能出现热更新冲突的问题。

```yaml
cfg:
  callee_1: A
  callee_2: B
  callee_3: C
```

#### **复现**

1. 业务协程开始从配置文件读取接口A的配置，读取完成后调用接口A。
2. 在接口A调用尚未返回时，WatchConfig监听到配置文件变化，触发热更新OnConfigChange，配置文件变化如下：
    
    ```yaml
    cfg:
      callee_1: D
      callee_2: E
      callee_3: F
    ```
    
3. 此时协程继续按照流水线流程读取配置，读取到下一个要执行的接口是E，这里就破坏了流程的完整性，与我们理想状态下的ABC或者DEF的执行流程不一致，可能导致无法预估的错误出现。

以下demo程序模拟了这种情况

```go
package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

// 模拟业务流程调用
func process(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Config file read start time:%v\n", time.Now().Format("2006-01-02 15:04:05"))
	// 第一步：获取配置接口A
	InterfaceFirst := viper.GetString("cfg.callee_1")
	fmt.Printf("Step 1: Interface types: %v\n", InterfaceFirst)

	// 模拟接口执行过程
	time.Sleep(5 * time.Second)

	// 第二步：获取配置接口B
	InterfaceSecond := viper.GetString("cfg.callee_2")
	fmt.Printf("Step 2: Interface types: %v\n", InterfaceSecond)
	// 模拟接口执行过程
	time.Sleep(5 * time.Second)

	// 第三步：获取配置接口C
	InterfaceThird := viper.GetString("cfg.callee_3")
	fmt.Printf("Step 3: Interface types: %v\n", InterfaceThird)
	fmt.Printf("Config file read end time:%v\n", time.Now().Format("2006-01-02 15:04:05"))
}

func InitConfig() {
	// 初始化 Viper 配置
	viper.SetConfigName("node")    // 配置文件名 (不带后缀)
	viper.SetConfigType("yaml")    // 配置文件类型
	viper.AddConfigPath("./conf/") // 配置文件路径

	// 加载初始配置
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 开启热更新监听
	viper.WatchConfig()

	// 在配置变更时触发的回调
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed:%v,time:%v\n", e.Name, time.Now().Format("2006-01-02 15:04:05"))
	})
}

func main() {
	InitConfig()

	// 使用 WaitGroup 模拟多个协程并发读取
	var wg sync.WaitGroup
	wg.Add(1)

	go process(&wg)

	// 在等待5s内手动更新配置文件config.yaml

	wg.Wait()
}
```

运行结果：

![](/img/blog/viper/3.png)

### 分析

我们先来看一下 `WatchConfig()` 的源码

![](/img/blog/viper/4.png)

简单来说就是通过 `WatchConfig()` 先初始化一个 `viper` 对象，完成后开始进行文件变更事件的监听，`OnConfigChange()` 负责将用户定义的回调逻辑赋值给 `viper` 对象的 `onConfigChange()`，此时文件如果发生对应的变更，则会触发对应的回调逻辑。

因此我们可以得出产生热更新冲突点原因：

1. 并发读写未同步 `viper` 默认未对内存中的配置数据加锁，当多个 `goroutine` 同时读写配置时，会引发竞争。
2. 配置存储非原子性配置文件写入中途被读取（如文件未完全写入），会导致读取到损坏或不完整的数据。

### 解决方案

#### 方案一：加读写锁

**实现思路**

1. **全局配置对象的建立**：为了便于管理多个系统的共享配置资源，我们将所有系统的相关配置集中存储在一个全局的配置对象中。通过这种设计，可以避免因同一配置对象被不同部分重复读取导致的操作异常。
2. **线程安全机制的实现**：在对全局配置进行更新时，直接修改原对象可能会引发多线程竞争和不一致性问题。为此，我们采用“读写锁“的方式，确保对配置对象的所有操作均需通过锁进行互斥处理。这种机制能够在任何时间点确保只有一个线程对配置对象进行修改。
3. **防止更新阻塞的技术**：为了避免因配置更新导致的线程阻塞问题以及确保数据一致性，我们在全局配置发生更新后采取以下措施：首先，在相关组件中触发回调机制，以通知其获取最新的配置信息；其次，启动一个协程来执行更新操作，这种方式可以有效避免因单一操作引发的资源阻塞，并确保其他协程线程能够正常读取和处理数据。

源码

```go
package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Cfg struct {
	callee1 string
	callee2 string
	callee3 string
}

type ConfigWrapper struct {
	conf Cfg
	sync.RWMutex
}

func (cw *ConfigWrapper) GetConfig() Cfg {
	cw.RLock()
	defer cw.RUnlock()
	return cw.conf
}

func (cw *ConfigWrapper) UpdateConfig(newConfig Cfg) {
	cw.Lock()
	defer cw.Unlock()
	cw.conf = newConfig
}

var globalConfig = &ConfigWrapper{}

// 模拟业务流程调用
func process1(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Config file read start time:%v\n", time.Now().Format("2006-01-02 15:04:05"))

	// 从全局配置中安全读取配置
	config := globalConfig.GetConfig()

	// 第一步：获取配置接口A
	InterfaceFirst := config.callee1
	fmt.Printf("Step 1: Interface types: %v\n", InterfaceFirst)

	// 模拟接口执行过程
	time.Sleep(5 * time.Second)

	// 第二步：获取配置接口B
	InterfaceSecond := config.callee2
	fmt.Printf("Step 2: Interface types: %v\n", InterfaceSecond)
	// 模拟接口执行过程
	time.Sleep(5 * time.Second)

	// 第三步：获取配置接口C
	InterfaceThird := config.callee3
	fmt.Printf("Step 3: Interface types: %v\n", InterfaceThird)
	fmt.Printf("Config file read end time:%v\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	// 初始化 Viper 配置
	viper.SetConfigName("node")    // 配置文件名 (不带后缀)
	viper.SetConfigType("yaml")    // 配置文件类型
	viper.AddConfigPath("./conf/") // 配置文件路径

	// 加载初始配置
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	globalConfig.UpdateConfig(Cfg{
		callee1: viper.GetString("cfg.callee_1"),
		callee2: viper.GetString("cfg.callee_2"),
		callee3: viper.GetString("cfg.callee_3"),
	})

	// 开启热更新监听
	viper.WatchConfig()

	// 在配置变更时触发的回调
	viper.OnConfigChange(func(e fsnotify.Event) {
		go func() {
			newConf := Cfg{
				callee1: viper.GetString("cfg.callee_1"),
				callee2: viper.GetString("cfg.callee_2"),
				callee3: viper.GetString("cfg.callee_3"),
			}
			globalConfig.UpdateConfig(newConf)
			fmt.Printf("Config updated successfully:%v\n", time.Now().Format("2006-01-02 15:04:05"))
		}()
		fmt.Printf("Config file changed:%v,time:%v\n", e.Name, time.Now().Format("2006-01-02 15:04:05"))
	})

	// 使用 WaitGroup 模拟多个协程并发读取
	var wg sync.WaitGroup
	wg.Add(1)

	go process1(&wg)

	// 在等待5s内手动更新配置文件config.yaml

	wg.Wait()
}

```

运行结果：

![](/img/blog/viper/5.png)

#### 方案二：原子操作

**实现思路**

1. **全局配置对象**：通过全局配置对象的Store`方法，在初始化及更新阶段获取最新的配置信息。
2. **复制读取**：Load方法以读取复制一份当前的状态。这种方式确保了每次读取的数据都是最新版本的副本，避免数据不一致的风险。
3. **原子性操作**：由于这些操作采用的是原子性机制，避免了显式的锁管理，因此在性能上具有显著优势。相比传统的加锁方式，这种做法能够有效减少资源竞争和同步开销，从而提升了系统的整体效率。

源码

```go
package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/viper"
)

// Config 定义一个结构体存储配置
type Conf struct {
	Callee1 string
	Callee2 string
	Callee3 string
}

var config atomic.Value // 使用 atomic.Value 存储配置

// loadConfig 将当前 Viper 配置加载到 Config 结构体
func loadConfig() Conf {
	return Conf{
		Callee1: viper.GetString("cfg.callee_1"),
		Callee2: viper.GetString("cfg.callee_2"),
		Callee3: viper.GetString("cfg.callee_3"),
	}
}

// 模拟业务流程调用
func process2(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Config file read start time: %v\n", time.Now().Format("2006-01-02 15:04:05"))

	// 获取当前配置的副本
	cfg := config.Load().(Conf)

	// 第一步：获取配置接口A
	fmt.Printf("Step 1: Interface types: %v\n", cfg.Callee1)
	time.Sleep(5 * time.Second)

	// 第二步：获取配置接口B
	fmt.Printf("Step 2: Interface types: %v\n", cfg.Callee2)
	time.Sleep(5 * time.Second)

	// 第三步：获取配置接口C
	fmt.Printf("Step 3: Interface types: %v\n", cfg.Callee3)
	fmt.Printf("Config file read end time: %v\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	// 初始化 Viper 配置
	viper.SetConfigName("node")    // 配置文件名 (不带后缀)
	viper.SetConfigType("yaml")    // 配置文件类型
	viper.AddConfigPath("./conf/") // 配置文件路径

	// 加载初始配置
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 初始化 atomic.Value 存储初始配置
	config.Store(loadConfig())

	// 开启热更新监听
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %v, time: %v\n", e.Name, time.Now().Format("2006-01-02 15:04:05"))
		go func() {
			// 更新 atomic.Value 存储的新配置
			config.Store(loadConfig())
		}()
	})

	// 使用 WaitGroup 模拟多个协程并发读取
	var wg sync.WaitGroup
	wg.Add(1)
	go process2(&wg)

	// 在等待 5 秒内手动更新配置文件 config.yaml
	wg.Wait()
}
```

运行结果：

![](/img/blog/viper/6.png)

换成多个协程不同时运看看效果

```go
var wg sync.WaitGroupwg.Add(2)
go process3(&wg, 1)
time.Sleep(5 * time.Second)
go process3(&wg, 2)
```

![](/img/blog/viper/7.png)
## 第三个问题：AddConfigPath()路径究竟怎写？

### 问题引入

正常启动项目，main 函数，都使用同位置可以正常读取配置，但是在单元测试调用初始化函数时，出现了无法找到配置文件的问题。

#### 复现

```go
func Init() error {  
    v := viper.New()  
    v.AddConfigPath("./config")  // 设置配置文件目录
    v.SetConfigName("config")  // 设置配置文件名
    v.SetConfigType("yaml")  // 设置配置文件后缀
    v.WatchConfig()  
    err := v.ReadInConfig()  
    if err != nil {  
        panic(fmt.Errorf("read config failed, %v", err))  
    }  
    if err := v.Unmarshal(&Conf); err != nil {  
        panic(fmt.Errorf("unmarshal to Conf failed, %v", err))  
    }  
    return err  
}
```

### 分析

列出当前关键目录

```
@MacBook-Air Muxi-Micro-Layout % tree
.
├── conf
│	├── config.go
│	├── config.yaml
│	├── config_model.go
│	└── config_test.go
├── go.mod
├── go.sum
├── main.go
├── wire.go
└── wire_gen.go
```

尝试在测试测试函数中输出当前路径，发现输出结果为 `Muxi-Micro-Layout/conf/conf` ，这显然不符合我们的预期。

主要原因是对 `./` 有误解，我一直以为 `./` 的意思是当前项目的根目录，**在 go 中 `./` 是基于执行命令的目录的**，也就是说在不同的目录下调用 `Init()`，`./` 所代表的意义不同。

### 解决方案

因为是直接获取的 `config.go` 文件的绝对目录，所以无论在哪里调用配置初始化函数，都不会出现找不到文件的问题了

```go
func Init() error {  
    v := viper.New()
    _, filename, _, _ := runtime.Caller(0) // 获取当前文件（config.go）路径
	confPath := path.Dir(filename)         // 获取当前文件目录
	viper.AddConfigPath(confPath)          // 设置配置文件目录
    v.SetConfigName("config")  
    v.SetConfigType("yaml")  
    v.WatchConfig()  
    err := v.ReadInConfig()  
    if err != nil {  
        panic(fmt.Errorf("read config failed, %v", err))  
    }  
    if err := v.Unmarshal(&Conf); err != nil {  
        panic(fmt.Errorf("unmarshal to Conf failed, %v", err))  
    }  
    return err  
}
```

---

参考链接：
官网：[https://github.com/spf13/viper](https://github.com/spf13/viper)

[https://zhuanlan.zhihu.com/p/23237101950](https://zhuanlan.zhihu.com/p/23237101950)

[https://juejin.cn/post/7259715675475558437](https://juejin.cn/post/7259715675475558437)