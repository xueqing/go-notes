# 3 Go 语言并发组件

M:N 调度器：将 M 个绿色线程映射到 N 个 OS 线程

fork-join 的并发模型：

- fork：在程序的任一节点，可将子节支与父节点同时运行
- join：在将来某个时刻，并发的执行分支将会合并在一起

通过匿名函数创建的 goroutine，在它们所创建的相同地址空间内执行，这些闭包内的变量是原值的引用。

goroutine 非常轻量。

sync 包的并发原语：

- WaitGroup：并发-安全的计数器。当不关心并发操作的结果，或者有其他方法收集并发操作的结果时可以使用 WaitGroup
- 互斥锁和读写锁：Mutex 可以独占访问共享资源
- cond：一个 goroutine 的集合点，等待或发布一个 event
  - cont.Wait 不只是阻塞。它挂起当前 goroutine，允许其他 goroutine 在 OS 线程上运行
  - 进入 Wait 后，会调用 Cond 变量的 Locker 的 Unlock 方法
  - 在退出 Wait 时，会调用 Cond 变量的 Locker 的 Lock 方法
  - 内部维护一个 FIFO，Signal 发现等待最长时间的 goroutine 并通知它；Broadcast 向所有等待的 goroutine 发送信号
- Once：Sync.Once 只计算调用 Do 方法的次数，而不是多少次唯一调用 Do 方法。建议包装在一个小的语法块中
- 池 ：是 Pool 模式的并发安全实现。
  - 用途包括
    - 约束创建昂贵的场景(如数据库连接)，以便只创建固定数量的实例
    - 尽可能快地将预先分配的对象缓存加载启动
  - 不适合的场景：使用 Pool 代码所需要的东西不是大概同质的，则检索需要花费时间。比如需要随机和可变长度的切片
- channel：channel 操作的结果给出了 channel 的状态

  | 操作 | channel 状态 | 结果 |
  | --- | --- | --- |
  | Read | nil | 阻塞 |
  | | 打开且非空 | 输出值 |
  | | 打开且空 | 阻塞 |
  | | 关闭的 | <默认值>, false |
  | | 只写 | 编译错误 |
  | Write | nil | 阻塞 |
  | | 打开的但填满 | 阻塞 |
  | | 打开的且不满 | 写入值|
  | | 关闭的 | panic |
  | | 只读 | 编译错误 |
  | Close | nil | panic |
  | | 打开且非空 | 关闭 channel；读取成功，直到通道耗尽，然后读取产生值的默认值 |
  | | 打开但空 | 关闭 channel；读到生产者的默认值 |
  | | 关闭的 | panic |
  | | 只读 | 编译错误 |

- select 语句：select 块中的 case 语句没有测试顺序，如果没有任何满足条件，执行也不会失败
- GOMAXPROCS 控制：这个函数控制的 OS 线程的数量将承载所谓的“工作队列”
