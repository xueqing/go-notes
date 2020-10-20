# go runtime 环境变量

- [go runtime 环境变量](#go-runtime-环境变量)
  - [介绍](#介绍)
  - [GOGC](#gogc)
  - [GOTRACEBACK](#gotraceback)
    - [GOTRACEBACK=0](#gotraceback0)
    - [GOTRACEBACK=1](#gotraceback1)
    - [GOTRACEBACK=2](#gotraceback2)
    - [GOTRACEBACK=crash](#gotracebackcrash)
    - [Go1.6 将引入的关于 GOTRACEBACK 的变化](#go16-将引入的关于-gotraceback-的变化)
  - [GOMAXPROCS](#gomaxprocs)
  - [GODEBUG](#godebug)
    - [gctrace](#gctrace)
    - [堆清除程序(scavenger)](#堆清除程序scavenger)
    - [schedtrace](#schedtrace)
  - [推荐阅读](#推荐阅读)
  - [参考](#参考)

## 介绍

go runtime，除了提供垃圾收集、goroutine 调度、定时器、网络 polling 等日常服务，也包含工具使能额外输出调试，甚至是改变 runtime 本身的行为。

这些工具由传递给 Go 程序的环境变量控制。这篇文章描述 runtime 支持的一些主要的环境变量的功能。

## GOGC

`GOGC` 是 go runtime 支持的最早的环境变量之一。它可能比 `GOROOT` 还早，但远未众所周知。

`GOGC` 控制垃圾收集器的侵占性(aggressiveness)。这个值默认设为 100，表示直到堆自上次分配后增长 100% 才会触发垃圾收集器。实际上，`GOGC=100`(默认)意味着每次活动堆翻倍时，垃圾收集器就会运行。

把这个值设置的更高，比如 `GOGC=200`，会延迟垃圾收集器周期的开始，直到活动堆增加到先前大小的 200%。把这个值设置的更低，比如 `GOGC=20`，会导致更频繁地触发垃圾收集器，因为在触发收集之前，可以在堆上分配的新数据较少。

设置 `GOGC=OFF` 会完全禁用垃圾收集器。

随着 Go1.5 引入低延迟收集器，类似"触发垃圾收集器周期"的短语变得更加流畅，但是底层信息不变，即 `GOGC` 大于 100 意味着垃圾收集器运行更不频繁，`GOGC` 小于 100 意味着垃圾收集器运行更频繁。

## GOTRACEBACK

`GOTRACEBACK` 控制当你的程序出现一个 panic 时的细节等级。在 Go1.5 `GOTRACEBACK` 有四个有效值：

- `GOTRACEBACK=0` 会抑制所有跟踪(traceback)，你只能获得 panic 信息
- `GOTRACEBACK=1` 是默认行为，会显示所有 goroutine 的堆栈跟踪，但是 runtime 相关的堆栈帧被抑制
- `GOTRACEBACK=2` 和上一个值相同，但是 runtime 相关的堆栈帧也会显示，这会显示 runtime 本身开启的 goroutine
- `GOTRACEBACK=crash` 和上一个值相同，但是 runtime 会导致程序段错误(segfault)，且操作系统允许时会触发一个核心转储(core dump)，而不是调用 `os.Exit`

可用一个简单的程序看 `GOTRACEBACK` 的效果。

```go
package main

func main() {
  panic("hello")
}
```

### GOTRACEBACK=0

编译并使用 `GOTRACEBACK=0` 运行这个程序，显示所有 goroutine 的堆栈跟踪被抑制。

```sh
kiki@ubuntu:~/gopro/test$ env GOTRACEBACK=0 ./test
panic: hello
kiki@ubuntu:~/gopro/test$ echo $?
2
```

### GOTRACEBACK=1

```sh
kiki@ubuntu:~/gopro/test$ env GOTRACEBACK=1 ./test
panic: hello

goroutine 1 [running]:
main.main()
  /home/kiki/gopro/test/main.go:4 +0x39
kiki@ubuntu:~/gopro/test$ echo $?
2
```

### GOTRACEBACK=2

```sh
kiki@ubuntu:~/gopro/test$ env GOTRACEBACK=2 ./test
panic: hello

goroutine 1 [running]:
panic(0x45e680, 0x481900)
  /usr/local/go/src/runtime/panic.go:722 +0x2c2 fp=0xc000048740 sp=0xc0000486b0 pc=0x424762
main.main()
  /home/kiki/gopro/test/main.go:4 +0x39 fp=0xc000048760 sp=0xc000048740 pc=0x452369
runtime.main()
  /usr/local/go/src/runtime/proc.go:203 +0x21e fp=0xc0000487e0 sp=0xc000048760 pc=0x42646e
runtime.goexit()
  /usr/local/go/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc0000487e8 sp=0xc0000487e0 pc=0x44bc91

goroutine 2 [force gc (idle)]:
runtime.gopark(0x477f60, 0x4c4ca0, 0x1411, 0x1)
  /usr/local/go/src/runtime/proc.go:304 +0xe0 fp=0xc000048fb0 sp=0xc000048f90 pc=0x426830
runtime.goparkunlock(...)
  /usr/local/go/src/runtime/proc.go:310
runtime.forcegchelper()
  /usr/local/go/src/runtime/proc.go:253 +0xb7 fp=0xc000048fe0 sp=0xc000048fb0 pc=0x4266e7
runtime.goexit()
  /usr/local/go/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc000048fe8 sp=0xc000048fe0 pc=0x44bc91
created by runtime.init.5
  /usr/local/go/src/runtime/proc.go:242 +0x35

goroutine 3 [GC sweep wait]:
runtime.gopark(0x477f60, 0x4c4d60, 0x140c, 0x1)
  /usr/local/go/src/runtime/proc.go:304 +0xe0 fp=0xc0000497a8 sp=0xc000049788 pc=0x426830
runtime.goparkunlock(...)
  /usr/local/go/src/runtime/proc.go:310
runtime.bgsweep(0xc000024070)
  /usr/local/go/src/runtime/mgcsweep.go:70 +0x9c fp=0xc0000497d8 sp=0xc0000497a8 pc=0x41b39c
runtime.goexit()
  /usr/local/go/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc0000497e0 sp=0xc0000497d8 pc=0x44bc91
created by runtime.gcenable
  /usr/local/go/src/runtime/mgc.go:210 +0x5c

goroutine 4 [GC scavenge wait]:
runtime.gopark(0x477f60, 0x4c4dc0, 0x140d, 0x1)
  /usr/local/go/src/runtime/proc.go:304 +0xe0 fp=0xc000049f38 sp=0xc000049f18 pc=0x426830
runtime.goparkunlock(...)
  /usr/local/go/src/runtime/proc.go:310
runtime.bgscavenge(0xc000024070)
  /usr/local/go/src/runtime/mgcscavenge.go:257 +0xe1 fp=0xc000049fd8 sp=0xc000049f38 pc=0x41a9f1
runtime.goexit()
  /usr/local/go/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc000049fe0 sp=0xc000049fd8 pc=0x44bc91
created by runtime.gcenable
  /usr/local/go/src/runtime/mgc.go:211 +0x7e
kiki@ubuntu:~/gopro/test$ echo $?
2
```

### GOTRACEBACK=crash

```sh
# 输出的结果为 0，说明默认是关闭 core dump 的
kiki@ubuntu:~/gopro/test$ ulimit -c
0
# 开启 core dump
kiki@ubuntu:~/gopro/test$ ulimit -c unlimited
# 修改 core 文件保存的路径，文件名格式为“core-命令名-pid-时间戳”
kiki@ubuntu:~/gopro/test$ sudo su
root@ubuntu:/home/kiki/gopro/test# echo "/opt/coredump/core-%e-%p-%t" > /proc/sys/kernel/core_pattern
root@ubuntu:/home/kiki/gopro/test# exit
exit
kiki@ubuntu:~/gopro/test$ env GOTRACEBACK=crash ./test
panic: hello

goroutine 1 [running]:
panic(0x45e680, 0x481900)
  /usr/local/go/src/runtime/panic.go:722 +0x2c2 fp=0xc000048740 sp=0xc0000486b0 pc=0x424762
main.main()
  /home/kiki/gopro/test/main.go:4 +0x39 fp=0xc000048760 sp=0xc000048740 pc=0x452369
runtime.main()
  /usr/local/go/src/runtime/proc.go:203 +0x21e fp=0xc0000487e0 sp=0xc000048760 pc=0x42646e
runtime.goexit()
  /usr/local/go/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc0000487e8 sp=0xc0000487e0 pc=0x44bc91

goroutine 2 [force gc (idle)]:
runtime.gopark(0x477f60, 0x4c4ca0, 0x1411, 0x1)
  /usr/local/go/src/runtime/proc.go:304 +0xe0 fp=0xc000048fb0 sp=0xc000048f90 pc=0x426830
runtime.goparkunlock(...)
  /usr/local/go/src/runtime/proc.go:310
runtime.forcegchelper()
  /usr/local/go/src/runtime/proc.go:253 +0xb7 fp=0xc000048fe0 sp=0xc000048fb0 pc=0x4266e7
runtime.goexit()
  /usr/local/go/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc000048fe8 sp=0xc000048fe0 pc=0x44bc91
created by runtime.init.5
  /usr/local/go/src/runtime/proc.go:242 +0x35

goroutine 3 [GC sweep wait]:
runtime.gopark(0x477f60, 0x4c4d60, 0x140c, 0x1)
  /usr/local/go/src/runtime/proc.go:304 +0xe0 fp=0xc0000497a8 sp=0xc000049788 pc=0x426830
runtime.goparkunlock(...)
  /usr/local/go/src/runtime/proc.go:310
runtime.bgsweep(0xc000024070)
  /usr/local/go/src/runtime/mgcsweep.go:70 +0x9c fp=0xc0000497d8 sp=0xc0000497a8 pc=0x41b39c
runtime.goexit()
  /usr/local/go/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc0000497e0 sp=0xc0000497d8 pc=0x44bc91
created by runtime.gcenable
  /usr/local/go/src/runtime/mgc.go:210 +0x5c

goroutine 4 [GC scavenge wait]:
runtime.gopark(0x477f60, 0x4c4dc0, 0x140d, 0x1)
  /usr/local/go/src/runtime/proc.go:304 +0xe0 fp=0xc000049f38 sp=0xc000049f18 pc=0x426830
runtime.goparkunlock(...)
  /usr/local/go/src/runtime/proc.go:310
runtime.bgscavenge(0xc000024070)
  /usr/local/go/src/runtime/mgcscavenge.go:257 +0xe1 fp=0xc000049fd8 sp=0xc000049f38 pc=0x41a9f1
runtime.goexit()
  /usr/local/go/src/runtime/asm_amd64.s:1357 +0x1 fp=0xc000049fe0 sp=0xc000049fd8 pc=0x44bc91
created by runtime.gcenable
  /usr/local/go/src/runtime/mgc.go:211 +0x7e
Aborted (core dumped)
kiki@ubuntu:~/gopro/test$ ls /opt/coredump/
core-test-60094-1597993246
kiki@ubuntu:~/gopro/test$ echo $?
0
```

### Go1.6 将引入的关于 GOTRACEBACK 的变化

对于 Go1.6，对 `GOTRACEBACK` 的解释正在改变。`GOTRACEBACK` 的新值会是：

- `GOTRACEBACK=none` 会抑制所有跟踪，你只能获得 panic 信息
- `GOTRACEBACK=single` 是新的默认行为，即只打印认为已经导致 panic 的 goroutine
- `GOTRACEBACK=all` 导致显示所有 goroutine 的堆栈跟踪，但是 runtime 相关的堆栈帧被抑制
- `GOTRACEBACK=system` 和上一个值相同，但是 runtime 相关的堆栈帧也会显示，这会显示 runtime 本身开启的 goroutine
- `GOTRACEBACK=crash` 跟 Go1.5 没有变化

为了兼容 Go1.5，值 0 映射到 `none`，1 映射到 `all`，2 映射到 `system`。

这个修改主要是在 Go1.6 中，默认的 panic 信息只会打印错误的 goroutine 的堆栈跟踪。可以参考 [issue 12366](https://github.com/golang/go/issues/12366) 和 [CL 16512](https://go-review.googlesource.com/#/c/16512/) 查看更多细节。

## GOMAXPROCS

`GOMAXPROCS` 是众所周知的值(并且通过 runtime.GOMAXPROCS 对应物进行训练(cargo culted))，该值控制程序中分配给 goroutine 的操作系统线程数。

从 Go1.5 开始，`GOMAXPROCS` 的默认值是启动时对程序可见的 CPU 的数目(无论你的操作系统是否是一个 CPU)。

注意：一个 Go 程序使用的操作系统线程数包括为 cgo 调用服务的线程、在操作系统调用中阻塞的线程，而且可能大于 `GOMAXPROCS` 的值。

## GODEBUG

保存最好的是 `GODEBUG`。`GODEBUG` 的内容被解释为以逗号分隔的 "名称=值" 对的列表，其中每个名称是 runtime 调试工具。有一个实例是在启用垃圾收集器和调度跟踪的情况下调用 godoc：

```sh
kiki@ubuntu:~/gopro/test$ env GODEBUG=gctrace=1,schedtrace=1000 godoc -http=:8080
```

这篇文章的剩余部分会讨论 `GODEBUG` 调试工具，这些工具对于诊断 Go 程序很有用。

### gctrace

在所有 `GODEBUG` 工具中，`gctrace` 是我发现最有用的。这里是 `godoc -http` 服务启用 `gctrace` 调试的前几毫秒的输出：

```sh
kiki@ubuntu:~/gopro/test$ env GODEBUG=gctrace=1 godoc -http=:8080 -index
gc #1 @0.042s 4%: 0.051+1.1+0.026+16+0.43 ms clock, 0.10+1.1+0+2.0/6.7/0+0.86 ms cpu, 4->32->10 MB, 4 MB goal, 4 P
gc #2 @0.062s 5%: 0.044+1.0+0.017+2.3+0.23 ms clock, 0.044+1.0+0+0.46/2.0/0+0.23 ms cpu, 4->12->3 MB, 8 MB goal, 4 P
gc #3 @0.067s 6%: 0.041+1.1+0.078+4.0+0.31 ms clock, 0.082+1.1+0+0/2.8/0+0.62 ms cpu, 4->6->4 MB, 8 MB goal, 4 P
gc #4 @0.073s 7%: 0.044+1.3+0.018+3.1+0.27 ms clock, 0.089+1.3+0+0/2.9/0+0.54 ms cpu, 4->7->4 MB, 6 MB goal, 4 P
```

这个输出的格式随着每个 Go 版本变化，但是你总会找到共性，比如不同 gc 短语的时间总数 `0.051+1.1+0.026+16+0.43 ms clock`，以及在垃圾收集周期期间的堆大小 `4->6->4 MB`。这个跟踪也包含 gc 周期相对于程序的开始时间结束的时间戳，但是旧版本的 Go 忽视这个信息。

单个的输出行可能对于[分析](https://dave.cheney.net/2014/07/11/visualising-the-go-garbage-collector)很有用，但是我发现以聚合方式查看它们更加有用。比如，如果你使能 gc 跟踪，且输出是连续的，明显表明程序受分配限制。同样的，如果报告的堆大小随着时间持续增长，明显表明有内存泄漏，再次情况下，按照预期应该被释放的引用正保存在某处全局结构。

使能 `gctrace` 的开销对于生产部署实际是 0，因为这些统计信息总是被收集，但通常是抑制的。我建议你至少将应用的生产部署的某些代表性示例使能。

注意：设置 `gctrace` 值大于 1 导致每个垃圾收集周期运行两倍。这将执行某些需要两个垃圾收集周期才能完成的终结工作。你不应将此作为修改程序终结性能的机制，因为你不应编写正确性取决于终结的程序。

### 堆清除程序(scavenger)

到目前为止，通过 `gctrace=1` 启用的最有用的输出片段是堆清除程序的输出：

```txt
scvg143: inuse: 8, idle: 104, sys: 113, released: 104, consumed: 8 (MB)
```

清除程序的工作是周期性地清理堆，以查看未使用的操作系统分页。然后，清除程序通过通知操作系统堆中这些内存页未使用而释放它们。没有工具强制操作系统收回这些页，且许多操作系统选择忽视这个建议，或至少推迟所有操作到机器缺少空闲内存。

清除程序的输出是我知道的最好方式，来辨别你的 Go 程序正在使用的虚拟地址空间大小。预计这些值将与类似 free(1) 和 top(1) 工具的报告出入很大。你应该相信清除程序报告的值。

### schedtrace

因为 Go 的 runtime 管理着大量 goroutine 分配给少量操作系统线程，因此从外部观察程序可能不会提供足够细节来了解其性能。你可能需要直接查看 runtime 调度器的操作。这个输出由 `schedtrace` 值控制：

```sh
kiki@ubuntu:~/gopro/test$ env GODEBUG=schedtrace=1000 godoc -http=:8080 -index
SCHED 0ms: gomaxprocs=4 idleprocs=2 threads=4 spinningthreads=1 idlethreads=0 runqueue=0 [0 0 0 0]
SCHED 1001ms: gomaxprocs=4 idleprocs=0 threads=8 spinningthreads=0 idlethreads=2 runqueue=0 [189 197 231 142]
SCHED 2004ms: gomaxprocs=4 idleprocs=0 threads=9 spinningthreads=0 idlethreads=1 runqueue=0 [54 45 38 86]
SCHED 3011ms: gomaxprocs=4 idleprocs=0 threads=9 spinningthreads=0 idlethreads=2 runqueue=2 [85 0 67 111]
SCHED 4018ms: gomaxprocs=4 idleprocs=3 threads=9 spinningthreads=0 idlethreads=4 runqueue=0 [0 0 0 0]
```

`schedtrace` 输出更详细的讨论可参考 Dmitry Vyukov 的 [excellent blog post from the Intel DeveloperZone](https://software.intel.com/content/www/us/en/develop/blogs/debugging-performance-issues-in-go-programs.html)。

追加 `scheddetail=1` 会导致 runtime 除了输出摘要，还输出每一个 goroutine 的状态，产生非常冗长的输出。

```sh
kiki@ubuntu:~/gopro/test$ env GODEBUG=scheddetail=1,schedtrace=1000 godoc -http=:8080 -index
SCHED 0ms: gomaxprocs=4 idleprocs=3 threads=3 spinningthreads=0 idlethreads=0 runqueue=0 gcwaiting=0 nmidlelocked=0 stopwait=0 sysmonwait=0
  P0: status=1 schedtick=0 syscalltick=0 m=0 runqsize=0 gfreecnt=0
  P1: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0
  P2: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0
  P3: status=0 schedtick=0 syscalltick=0 m=-1 runqsize=0 gfreecnt=0
  M2: p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=1 dying=0 helpgc=0 spinning=false blocked=false lockedg=-1
  M1: p=-1 curg=17 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 helpgc=0 spinning=false blocked=false lockedg=17
  M0: p=0 curg=1 mallocing=0 throwing=0 preemptoff= locks=2 dying=0 helpgc=0 spinning=false blocked=false lockedg=1
  G1: status=2(stack growth) m=0 lockedm=0
  G17: status=3() m=1 lockedm=1
  G2: status=1() m=-1 lockedm=-1
```

这个输出可能对于调试泄漏的 goroutine 很有用，但是其他工具比如 [net/http/pprof](https://golang.org/pkg/net/http/pprof/) 可能更有用。

## 推荐阅读

你的 Go 版本可用的所有环节变量在 [runtime 包的 godoc](https://godoc.org/runtime#hdr-Environment_Variables) 中。

## 参考

- [A whirlwind tour of Go’s runtime environment variables](https://dave.cheney.net/2015/11/29/a-whirlwind-tour-of-gos-runtime-environment-variables)
