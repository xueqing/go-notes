# 2 对你的代码建模：通信顺序进程

并发 VS 并行：并发属于代码；并行属于一个运行中的程序。

通信顺序进程 (Communicate Sequence Process, CSP)

Go 语言的一个座右铭：

> 使用通信来共享内存，而不是通过共享内存来通信。

Go 语言的并发性哲学：

> 追求简洁，尽量使用 channel，并且认为 goroutine 的使用是没有成本的。
