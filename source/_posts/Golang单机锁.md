---
type: Post
title: Golang单机锁
tags: Go
category: 开发
category_bar: true
abbrlink: 52635
date: 2024-09-26 23:30:34
---

## **sync.Mutex**

### 核心机制

- 通过 Mutex 内一个状态值标识锁的状态，例如，取 0 表示未加锁，1 表示已加锁；
  - 上锁：把 0 改为 1；
  - 解锁：把 1 置为 0.
  - 上锁时，假若已经是 1，则上锁失败，需要等他人解锁，将状态改为 0.

下面首先拎清两个概念：

- 饥饿模式：当 Mutex 阻塞队列中存在处于饥饿态的 goroutine 时，会进入模式，将抢锁流程由非公平机制转为公平机制.这是 sync.Mutex 为拯救陷入饥荒的老 goroutine 而启用的特殊机制，饥饿模式下，锁的所有权按照阻塞队列的顺序进行依次传递. 新 goroutine 进行流程时不得抢锁，而是进入队列尾部排队.
- 正常模式/非饥饿模式：这是 sync.Mutex 默认采用的模式. 当有 goroutine 从阻塞队列被唤醒时，会和此时先进入抢锁流程的 goroutine 进行锁资源的争夺，假如抢锁失败，会重新回到阻塞队列头部.

### 数据结构

![1.png](/img/blog/golanddjs/1.png)

```go
type Mutex struct {
    state int32
    sema  uint32
}
```

- state：锁中最核心的状态字段，不同 bit 位分别存储了 mutexLocked(是否上锁)、mutexWoken（是否有 goroutine 从阻塞队列中被唤醒）、mutexStarving（是否处于饥饿模式）的信息
- sema：用于阻塞和唤醒 goroutine 的信号量。

```go
const (
    mutexLocked = 1 << iota // mutex is locked
    mutexWoken
    mutexStarving
    mutexWaiterShift = iota
​
    starvationThresholdNs = 1e6
)
```

- mutexLocked = 1：state 最右侧的一个 bit 位标志是否上锁，0-未上锁，1-已上锁；
- mutexWoken = 2：state 右数第二个 bit 位标志是否有 goroutine 从阻塞中被唤醒，0-没有，1-有；
- mutexStarving = 4：state 右数第三个 bit 位标志 Mutex 是否处于饥饿模式，0-非饥饿，1-饥饿；
- mutexWaiterShift = 3：右侧存在 3 个 bit 位标识特殊信息，分别为上述的 mutexLocked、mutexWoken、mutexStarving；
- starvationThresholdNs = 1 ms：sync.Mutex 进入饥饿模式的等待时间阈值。

### **Mutex.Lock()**

![2.png](/img/blog/golanddjs/2.png)

```go
func (m *Mutex) Lock() {
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        return
    }
    // Slow path (outlined so that the fast path can be inlined)
    m.lockSlow()
}
```

- 首先进行一轮 CAS 操作，假如当前未上锁且锁内不存在阻塞协程，则直接 CAS 抢锁成功返回；
- 第一轮初探失败，则进入 lockSlow 流程，下面细谈。

#### 补充：什么是CAS策略

CAS是Compare And Swap（比较并交换）的缩写，是一种非阻塞式并发控制技术，用于保证多个线程在修改同一个共享资源时不会出现竞争条件。更准确的是采用乐观锁技术，实现线程安全的问题。CAS有三个操作数———内存对象（V）、预期原值（A）、新值（B）。

CAS原理就是对V对象进行赋值时，先判断原来的值是否为A，如果为A，就把新值B赋值到V对象上面，如果原来的值不是A（代表V的值发生了变化），就不赋新值。

```go
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
```

- `addr`：指向要更新的整数的指针。
- `old`：期望的旧值。
- `new`：要设置的新值。
- 返回值 `swapped`：如果 `addr` 的当前值等于 `old`，则返回 `true`，并将 `addr` 的值替换为 `new`；否则返回 `false`，且不做任何更改。

### **Mutex.lockSlow()**

```go
func (m *Mutex) lockSlow() {
    var waitStartTime int64
    starving := false
    awoke := false
    iter := 0
    old := m.state
    // ...
}
```

- waitStartTime：标识当前 goroutine 在抢锁过程中的等待时长，单位：ns；
- starving：标识当前是否处于饥饿模式；
- awoke：标识当前是否已有协程在等锁；
- iter：标识当前 goroutine 参与自旋的次数；
- old：临时存储锁的 state 值.

#### **自旋空转**

![3.png](/img/blog/golanddjs/3.png)

```go
func (m *Mutex) lockSlow() {
    // ...
    for {
        // 进入该 if 分支，说明抢锁失败，处于饥饿模式，但仍满足自旋条件
        if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
            // 进入该 if 分支，说明当前锁阻塞队列有协程，但还未被唤醒，因此需要将
            // mutexWoken 标识置为 1，避免再有其他协程被唤醒和自己抢锁
            if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
                atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
                awoke = true
            }
            runtime_doSpin()
            iter++
            old = m.state
            continue
        }

        // ...
    }
}
```

- 走进 for 循环；
- 假如满足三个条件：I 锁已被占用、 II 锁为正常模式、III 满足自旋条件（runtime_canSpin 方法），则进入自旋后处理环节；
- 在自旋后处理中，假如当前锁有尚未唤醒的阻塞协程，则通过 CAS 操作将 state 的 mutexWoken 标识置为 1，将局部变量 awoke 置为 true；
- 调用 runtime_doSpin 告知调度器 P 当前处于自旋模式；
- 更新自旋次数 iter 和锁状态值 old；
- 通过 continue 语句进入下一轮尝试.

#### state 新值构造

![4.png](/img/blog/golanddjs/4.png)

```go
func (m *Mutex) lockSlow() {
    // ...
    for {
        // 自旋抢锁失败后处理 ...

        new := old
        if old&mutexStarving == 0 {
            new |= mutexLocked
        }
        if old&(mutexLocked|mutexStarving) != 0 {
            new += 1 << mutexWaiterShift
        }
        if starving && old&mutexLocked != 0 {
            new |= mutexStarving
        }
        if awoke {
            new &^= mutexWoken
        }

        // ...
    }
}
```

- 从自旋中走出来后，会存在两种分支，要么加锁成功，要么陷入自锁，不论是何种情形，都会先对 sync.Mutex 的状态新值 new 进行更新；
- 倘若当前是非饥饿模式，则在新值 new 中置为已加锁，即尝试抢锁；
- 倘若旧值为已加锁或者处于饥饿模式，则当前 goroutine 在这一轮注定无法抢锁成功，可以直接令新值的阻塞协程数加1；
- 倘若当前进入饥饿模式且旧值已加锁，则将新值置为饥饿模式；
- 倘若局部变量标识是已有唤醒协程抢锁，说明 Mutex.state 中的 mutexWoken 是被当前 gouroutine 置为 1 的，但由于当前 goroutine 接下来要么抢锁成功，要么被阻塞挂起，因此需要在新值中将该 mutexWoken 标识更新置 0.

##### 补充：&^ 是什么

a &^ b ：bit clear，清零a中，ab都是1的位。

#### **state 新旧值替换**

![5.png](/img/blog/golanddjs/5.png)

```go
func (m *Mutex) lockSlow() {
    // ...
    for {
        // 自旋抢锁失败后处理 ...

        // new old 状态值更新 ...

        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            // case1 加锁成功
            // case2 将当前协程挂起

            // ...
        }else {
            old = m.state
        }
        // ...
    }
}
```

- 通过 CAS 操作，用构造的新值替换旧值；
- 倘若失败（即旧值被其他协程介入提前修改导致不符合预期），则将旧值更新为此刻的 Mutex.State，并开启一轮新的循环；
- 倘若 CAS 替换成功，则进入最后一轮的二择一局面：I 倘若当前 goroutine 加锁成功，则返回；II 倘若失败，则将 goroutine 挂起添加到阻塞队列.

#### 上锁成功分支

![6.png](/img/blog/golanddjs/6.png)

```go
func (m *Mutex) lockSlow() {
    // ...
    for {
        // 自旋抢锁失败后处理 ...

        // new old 状态值更新 ...

        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            if old&(mutexLocked|mutexStarving) == 0 {
                break
            }

            // ...
        }
        // ...
    }
}
```

- 延续 1.2.2.4 的思路，此时已经成功将 Mutex.state 由旧值替换为新值；
- 接下来进行判断，倘若旧值是未加锁状态且为正常模式，则意味着加锁标识位正是由当前 goroutine 完成的更新，说明加锁成功，返回即可；
- 倘若旧值中锁未释放或者处于饥饿模式，则当前 goroutine 需要进入阻塞队列挂起.

#### 阻塞挂起

![7.png](/img/blog/golanddjs/7.png)

```go
func (m *Mutex) lockSlow() {
    // ...
    for {
        // 自旋抢锁失败后处理 ...

        // new old 状态值更新 ...

        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            // 加锁成功后返回的逻辑分支 ...

            queueLifo := waitStartTime != 0
            if waitStartTime == 0 {
                waitStartTime = runtime_nanotime()
            }
            runtime_SemacquireMutex(&m.sema, queueLifo, 1)
            // ...
        }
        // ...
    }
}
```

承接上节，走到此处的情形有两种：要么是抢锁失败，要么是锁已处于饥饿模式，而当前 goroutine 不是从阻塞队列被唤起的协程. 不论处于哪种情形，当前 goroutine 都面临被阻塞挂起的命运.

- 基于 queueLifo 标识当前 goroutine 是从阻塞队列被唤起的老客还是新进流程的新客；
- 倘若等待的起始时间为零，则为新客；倘若非零，则为老客；
- 倘若是新客，则对等待的起始时间进行更新，置为当前时刻的 ns 时间戳；
- 将当前协程添加到阻塞队列中，倘若是老客则挂入队头；倘若是新客，则挂入队尾；
- 挂起当前协程.

#### **从阻塞态被唤醒**

![8.png](/img/blog/golanddjs/8.png)

```go
func (m *Mutex) lockSlow() {
    // ...
    for {
        // 自旋抢锁失败后处理...

        // new old 状态值更新 ...

        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            // 加锁成功后返回的逻辑分支 ...

            // 挂起前处理 ...
            runtime_SemacquireMutex(&m.sema, queueLifo, 1)
            // 从阻塞队列被唤醒了
            starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
            old = m.state
            if old&mutexStarving != 0 {
                delta := int32(mutexLocked - 1<<mutexWaiterShift)
                if !starving || old>>mutexWaiterShift == 1 {
                    delta -= mutexStarving
                }
                atomic.AddInt32(&m.state, delta)
                break
            }
            awoke = true
            iter = 0
        }
        // ...
    }
}
```

- 走入此处，说明当前 goroutine 是从 Mutex 的阻塞队列中被唤起的；
- 判断一下，此刻需要进入阻塞态，倘若当前 goroutine 进入阻塞队列时间长达 1 ms，则说明需要；此时会更新 starving 局部变量，并在下一轮循环中完成对 Mutex.state 中 starving 标识位的更新；
- 获取此时锁的状态，通过 old 存储；
- 倘若此时锁是饥饿模式，则当前 goroutine 无需竞争可以直接获得锁；
- 饥饿模式下，goroutine 获取锁前需要更新锁的状态，包含 mutexLocked、锁阻塞队列等待协程数以及 mutexStarving 三个信息；均通过 delta 变量记录差值，最终通过原子操作添加到 Mutex.state 中；
- mutexStarving 的更新要作前置判断，倘若当前局部变量 starving 为 false，或者当前 goroutine 就是 Mutex 阻塞队列的最后一个 goroutine，则将 Mutex.state 置为正常模式.

### **Unlock**

![9.png](/img/blog/golanddjs/9.png)

```go
func (m *Mutex) Unlock() {
    new := atomic.AddInt32(&m.state, -mutexLocked)
    if new != 0 {
        m.unlockSlow(new)
    }
}
```

- 通过原子操作解锁；
- 倘若解锁时发现，目前参与竞争的仅有自身一个 goroutine，则直接返回即可；
- 倘若发现锁中还有阻塞协程，则走入 unlockSlow 分支.

#### 补充：什么是原子操作

原子操作（atomic operation）指的是由多步操作组成的一个操作。如果该操作不能原子地执行，则要么执行完所有步骤，要么一步也不执行，不可能只执行所有步骤的一个子集。

原子操作是进行过程中不能被中断的操作，针对某个值的原子操作在被进行的过程中，CPU绝不会再去进行其他的针对该值的操作。为了实现这样的严谨性，原子操作仅会由一个独立的CPU指令代表和完成。原子操作是无锁的，常常直接通过CPU指令直接实现。

具体的原子操作在不同的操作系统中实现是不同的。比如在Intel的CPU架构机器上，主要是使用总线锁的方式实现的。 就是当一个CPU需要操作一个内存块的时候，向总线发送一个LOCK信号，所有CPU收到这个信号后就不对这个内存块进行操作了。 等待操作的CPU执行完操作后，发送UNLOCK信号，才结束。更具体的，在x86体系中, CPU提供了HLOCK pin引线, 允许CPU在执行某一个指令(仅仅是一个指令)时拉低HLOCK pin引线的电位, 直到这个指令执行完毕才放开，从而锁住了总线,，如此在同一总线的CPU就暂时无法通过总线访问内存了，这样就保证了多核处理器的原子性。 在AMD的CPU架构机器上就是使用MESI一致性协议的方式来保证原子操作。 所以我们在看atomic源码的时候，我们看到它针对不同的操作系统有不同汇编语言文件。

##### Go中原子操作的支持

Go语言的`sync/atomic`提供了对原子操作的支持，用于同步访问整数和指针。

- Go语言提供的原子操作都是非入侵式的
- 原子操作支持的类型包括`int32、int64、uint32、uint64、uintptr、unsafe.Pointer`。

### **unlockSlow**

![10.png](/img/blog/golanddjs/10.png)

#### **未加锁的异常情形**

```go
func (m *Mutex) unlockSlow(new int32) {
    if (new+mutexLocked)&mutexLocked == 0 {
        fatal("sync: unlock of unlocked mutex")
    }
    // ...
}
```

解锁时倘若发现 Mutex 此前未加锁，直接抛出 fatal.

#### **正常模式**

```go
func (m *Mutex) unlockSlow(new int32) {
    // ...
    if new&mutexStarving == 0 {
        old := new
        for {

            if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
                return
            }

            new = (old - 1<<mutexWaiterShift) | mutexWoken
            if atomic.CompareAndSwapInt32(&m.state, old, new) {
                runtime_Semrelease(&m.sema, false, 1)
                return
            }
            old = m.state
        }
    }
    // ...
}
```

- 倘若阻塞队列内无 goroutine 或者 mutexLocked、mutexStarving、mutexWoken 标识位任一不为零，三者均说明此时有其他活跃协程已介入，自身无需关心后续流程；
- 基于 CAS 操作将 Mutex.state 中的阻塞协程数减 1，倘若成功，则唤起阻塞队列头部的 goroutine，并退出；
- 倘若减少阻塞协程数的 CAS 操作失败，则更新此时的 Mutex.state 为新的 old 值，开启下一轮循环.

#### **饥饿模式**

```go
func (m *Mutex) unlockSlow(new int32) {
    // ...
    if new&mutexStarving == 0 {
        // ...
    } else {
        runtime_Semrelease(&m.sema, true, 1)
    }
}
```

饥饿模式下，直接唤醒阻塞队列头部的 goroutine 即可.

## **Sync.RWMutex**

### 核心机制

- 从逻辑上，可以把 RWMutex 理解为一把读锁加一把写锁；
- 写锁具有严格的排他性，当其被占用，其他试图取写锁或者读锁的 goroutine 均阻塞；
- 读锁具有有限的共享性，当其被占用，试图取写锁的 goroutine 会阻塞，试图取读锁的 goroutine 可与当前 goroutine 共享读锁；
- 综上可见，RWMutex 适用于读多写少的场景，最理想化的情况，当所有操作均使用读锁，则可实现去无化；最悲观的情况，倘若所有操作均使用写锁，则 RWMutex 退化为普通的 Mutex.

### **数据结构**

![11.png](/img/blog/golanddjs/11.png)

```go
const rwmutexMaxReaders = 1 << 30

type RWMutex struct {
    w           Mutex  // held if there are pending writers
    writerSem   uint32 // semaphore for writers to wait for completing readers
    readerSem   uint32 // semaphore for readers to wait for completing writers
    readerCount int32  // number of pending readers
    readerWait  int32  // number of departing readers
}
```

- rwmutexMaxReaders：共享读锁的 goroutine 数量上限，值为 2^29；
- w：RWMutex 内置的一把普通互斥锁 sync.Mutex；
- writerSem：关联写锁阻塞队列的信号量；
- readerSem：关联读锁阻塞队列的信号量；
- readerCount：正常情况下等于介入读锁流程的 goroutine 数量；当 goroutine 接入写锁流程时，该值为实际介入读锁流程的 goroutine 数量减 rwmutexMaxReaders.
- readerWait：记录在当前 goroutine 获取写锁前，还需要等待多少个 goroutine 释放读锁.

### **RLock**

```go
func (rw *RWMutex) RLock() {
    if atomic.AddInt32(&rw.readerCount, 1) < 0 {
        runtime_SemacquireMutex(&rw.readerSem, false, 0)
    }
}
```

- 基于原子操作，将 RWMutex 的 readCount 变量加一，表示占用或等待读锁的 goroutine 数加一；
- 倘若 RWMutex.readCount 的新值仍小于 0，说明有 goroutine 未释放写锁，因此将当前 goroutine 添加到读锁的阻塞队列中并阻塞挂起.

### **RUnlock**

```go
func (rw *RWMutex) RUnlock() {
    if r := atomic.AddInt32(&rw.readerCount, -1); r < 0 {
        rw.rUnlockSlow(r)
    }
}
```

- 基于原子操作，将 RWMutex 的 readCount 变量加一，表示占用或等待读锁的 goroutine 数减一；
- 倘若 RWMutex.readCount 的新值小于 0，说明有 goroutine 在等待获取写锁，则走入 RWMutex.rUnlockSlow 的流程中.

### **rUnlockSlow**

```go
func (rw *RWMutex) rUnlockSlow(r int32) {
    if r+1 == 0 || r+1 == -rwmutexMaxReaders {
        fatal("sync: RUnlock of unlocked RWMutex")
    }
    if atomic.AddInt32(&rw.readerWait, -1) == 0 {
        runtime_Semrelease(&rw.writerSem, false, 1)
    }
}
```

- 对 RWMutex.readerCount 进行校验，倘若发现当前协程此前未抢占过读锁，或者介入读锁流程的 goroutine 数量达到上限，则抛出 fatal；

(倘若 r+1 == -rwmutexMaxReaders，说明此时有 goroutine 介入写锁流程，但当前此前未加过读锁；倘若 r+1==0，则要么此前未加过读锁，要么介入读锁流程的 goroutine 数量达到上限.)

- 基于原子操作，对 RWMutex.readerWait 进行减一操作，倘若其新值为 0，说明当前 goroutine 是最后一个介入读锁流程的协程，因此需要唤醒一个等待写锁的阻塞队列的 goroutine.

### **Lock**

```go
func (rw *RWMutex) Lock() {
    rw.w.Lock()
    r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
    if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
        runtime_SemacquireMutex(&rw.writerSem, false, 0)
    }
}
```

- 对 RWMutex 内置的互斥锁进行加锁操作；
- 基于原子操作，对 RWMutex.readerCount 进行减少 -rwmutexMaxReaders 的操作；
- 倘若此时存在未释放读锁的 gouroutine，则基于原子操作在 RWMutex.readerWait 的基础上加上介入读锁流程的 goroutine 数量，并将当前 goroutine 添加到写锁的阻塞队列中挂起.

### **Unlock**

```go
func (rw *RWMutex) Unlock() {
    r := atomic.AddInt32(&rw.readerCount, rwmutexMaxReaders)
    if r >= rwmutexMaxReaders {
        fatal("sync: Unlock of unlocked RWMutex")
    }
    for i := 0; i < int(r); i++ {
        runtime_Semrelease(&rw.readerSem, false, 0)
    }
    rw.w.Unlock()
}
```

- 基于原子操作，将 RWMutex.readerCount 的值加上 rwmutexMaxReaders；
- 倘若发现 RWMutex.readerCount 的新值大于 rwmutexMaxReaders，则说明要么当前 RWMutex 未上过写锁，要么介入读锁流程的 goroutine 数量已经超限，因此直接抛出 fatal；
- 因此唤醒读锁阻塞队列中的所有 goroutine；(可见，竞争读锁的 goroutine 更具备优势)
- 解开 RWMutex 内置的互斥锁.

---

参考资料：

[Golang 单机锁实现原理 (qq.com)](https://mp.weixin.qq.com/s/5o0pR0RDaasKh4veXTctVg)

[分享大纲_哔哩哔哩_bilibili](https://www.bilibili.com/video/BV1kv4y157wj?p=1&vd_source=be87c4fef68f69704d0998e55b81b6a7)（讲的特别详细）
