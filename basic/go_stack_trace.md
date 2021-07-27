# go 堆栈跟踪分析

- [go 堆栈跟踪分析](#go-堆栈跟踪分析)
  - [程序输出堆栈跟踪](#程序输出堆栈跟踪)
    - [程序 panic](#程序-panic)
    - [debug.PrintStack](#debugprintstack)
    - [pprof.Lookup("goroutine").WriteTo](#pproflookupgoroutinewriteto)
    - [runtime.Stack](#runtimestack)
    - [http/pprof](#httppprof)
  - [示例1](#示例1)
    - [源码](#源码)
    - [运行程序](#运行程序)
      - [完整的堆栈跟踪](#完整的堆栈跟踪)
      - [panic 的堆栈跟踪](#panic-的堆栈跟踪)
  - [示例2](#示例2)
  - [示例3](#示例3)
  - [参考](#参考)

## 程序输出堆栈跟踪

### 程序 panic

程序中有不能恢复的 panic，或者出现运行时异常，程序默认会打印当前 goroutine 的堆栈跟踪。

```go
// main.go
package main
import (
  "time"
)
func main() {
  go a()
  m1()
}
func m1() {
  m2()
}
func m2() {
  m3()
}
func m3() {
  panic("panic from m3")
}
func a() {
  time.Sleep(time.Hour)
}
```

```sh
kiki@ubuntu:~/gopro/test$ go run main.go
panic: panic from m3

goroutine 1 [running]:
main.m3(...)
  /home/kiki/gopro/test/main.go:18
main.m2(...)
  /home/kiki/gopro/test/main.go:15
main.m1(...)
  /home/kiki/gopro/test/main.go:12
main.main()
  /home/kiki/gopro/test/main.go:9 +0x54
exit status 2
# 查看所有 goroutine 的堆栈跟踪
kiki@ubuntu:~/gopro/test$ GOTRACEBACK="all" go run main.go
panic: panic from m3

goroutine 1 [running]:
main.m3(...)
  /home/kiki/gopro/test/main.go:18
main.m2(...)
  /home/kiki/gopro/test/main.go:15
main.m1(...)
  /home/kiki/gopro/test/main.go:12
main.main()
  /home/kiki/gopro/test/main.go:9 +0x54

goroutine 5 [runnable]:
main.a()
  /home/kiki/gopro/test/main.go:20
created by main.main
  /home/kiki/gopro/test/main.go:8 +0x35
exit status 2
```

### debug.PrintStack

通过 `debug.PrintStack()` 方法可以将当前所在的 goroutine 的堆栈跟踪打印出来。

```go
// main.go
package main
import (
  "runtime/debug"
  "time"
)
func main() {
  go a()
  m1()
}
func m1() {
  m2()
}
func m2() {
  m3()
}
func m3() {
  debug.PrintStack()
  time.Sleep(time.Hour)
}
func a() {
  time.Sleep(time.Hour)
}
```

```sh
kiki@ubuntu:~/gopro/test$ go run main.go
goroutine 1 [running]:
runtime/debug.Stack(0x1, 0x8, 0x43091e)
  /usr/local/go/src/runtime/debug/stack.go:24 +0x9d
runtime/debug.PrintStack()
  /usr/local/go/src/runtime/debug/stack.go:16 +0x22
main.m3()
  /home/kiki/gopro/test/main.go:19 +0x22
main.m2(...)
  /home/kiki/gopro/test/main.go:16
main.m1(...)
  /home/kiki/gopro/test/main.go:13
main.main()
  /home/kiki/gopro/test/main.go:10 +0x3c
```

### pprof.Lookup("goroutine").WriteTo

可以通过 `pprof.Lookup("goroutine").WriteTo` 将所有的 goroutine 的堆栈跟踪都打印出来。

```go
// main.go
package main
import (
  "os"
  "runtime/pprof"
  "time"
)
func main() {
  go a()
  m1()
}
func m1() {
  m2()
}
func m2() {
  m3()
}
func m3() {
  pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
  time.Sleep(time.Hour)
}
func a() {
  time.Sleep(time.Hour)
}
```

```sh
kiki@ubuntu:~/gopro/test$ go run main.go
goroutine profile: total 2
1 @ 0x4b2455 0x4b2270 0x4aee9a 0x4b8f75 0x4b8efc 0x4b8ef7 0x4b8ef6 0x42d18e 0x455be1
#  0x4b2454  runtime/pprof.writeRuntimeProfile+0x94  /usr/local/go/src/runtime/pprof/pprof.go:708
#  0x4b226f  runtime/pprof.writeGoroutine+0x9f  /usr/local/go/src/runtime/pprof/pprof.go:670
#  0x4aee99  runtime/pprof.(*Profile).WriteTo+0x3d9  /usr/local/go/src/runtime/pprof/pprof.go:329
#  0x4b8f74  main.m3+0x64        /home/kiki/gopro/test/main.go:20
#  0x4b8efb  main.m2+0x3b        /home/kiki/gopro/test/main.go:17
#  0x4b8ef6  main.m1+0x36        /home/kiki/gopro/test/main.go:14
#  0x4b8ef5  main.main+0x35        /home/kiki/gopro/test/main.go:11
#  0x42d18d  runtime.main+0x21d      /usr/local/go/src/runtime/proc.go:203

1 @ 0x4b8fa0 0x455be1
#  0x4b8fa0  main.a+0x0  /home/kiki/gopro/test/main.go:23
```

### runtime.Stack

```go
// main.go
package main
import (
  "fmt"
  "os"
  "os/signal"
  "runtime"
  "syscall"
  "time"
)
func main() {
  setupSigusr1Trap()
  go a()
  m1()
}
func m1() {
  m2()
}
func m2() {
  m3()
}
func m3() {
  time.Sleep(time.Hour)
}
func a() {
  time.Sleep(time.Hour)
}
func setupSigusr1Trap() {
  c := make(chan os.Signal, 1)
  signal.Notify(c, syscall.SIGUSR1)
  go func() {
    for range c {
      DumpStacks()
    }
  }()
}
func DumpStacks() {
  buf := make([]byte, 16384)
  buf = buf[:runtime.Stack(buf, true)]
  fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===", buf)
}
```

### http/pprof

导入 `net/http/pprof` 包为 `debug/pprof` 注册一个 HTPP handler：

```go
import _ "net/http/pprof"
import _ "net/http"
```

如果媒体运行的 HTTP 监听者开启一个：

```go
go func() {
  log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

然后在浏览器访问 `http://localhost:6060/debug/pprof` 查看菜单，或者访问 `http://localhost:8888/debug/pprof/goroutine?debug=2` 访问所有的 goroutine 的堆栈。

## 示例1

### 源码

```go
// main.go
package main

import (
  "fmt"
  "sync"
  "time"
)

func main() {
  slice := make([]string, 2, 4)
  wg := sync.WaitGroup{}
  stop := make(chan interface{})
  for i := 0; i < 3; i++ {
    wg.Add(1)
    go func(slice []string, i int, stop chan interface{}) {
      defer wg.Done()
      example(slice, "hello", i, stop)
    }(slice, i, stop)
  }
  wg.Wait()
  close(stop)
}

func example(slice []string, str string, i int, stop chan interface{}) {
  fmt.Println("example begin", str, i)
  if i == 0 {
    time.Sleep(2 * time.Second)
    fmt.Println(1 / i)
  }
  <-stop
  fmt.Println("example end", str, i)
}
```

### 运行程序

#### 完整的堆栈跟踪

```sh
# GOTRACEBACK=1 显示所有 goroutine 的堆栈跟踪
kiki@ubuntu:~/gopro/test$ GOTRACEBACK=1 go run main.go
example begin hello 2
example begin hello 1
example begin hello 0
panic: runtime error: integer divide by zero

goroutine 6 [running]:
main.example(0xc0000260c0, 0x2, 0x4, 0x4c22f3, 0x5, 0x0, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:28 +0x211
main.main.func1(0xc000014110, 0xc0000260c0, 0x2, 0x4, 0x0, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:17 +0xae
created by main.main
  /home/kiki/gopro/test/main.go:15 +0xfc

goroutine 1 [semacquire]:
sync.runtime_Semacquire(0xc000014118)
  /usr/local/go/src/runtime/sema.go:56 +0x42
sync.(*WaitGroup).Wait(0xc000014110)
  /usr/local/go/src/sync/waitgroup.go:130 +0x64
main.main()
  /home/kiki/gopro/test/main.go:20 +0x122

goroutine 7 [chan receive]:
main.example(0xc0000260c0, 0x2, 0x4, 0x4c22f3, 0x5, 0x1, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:30 +0x12f
main.main.func1(0xc000014110, 0xc0000260c0, 0x2, 0x4, 0x1, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:17 +0xae
created by main.main
  /home/kiki/gopro/test/main.go:15 +0xfc

goroutine 8 [chan receive]:
main.example(0xc0000260c0, 0x2, 0x4, 0x4c22f3, 0x5, 0x2, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:30 +0x12f
main.main.func1(0xc000014110, 0xc0000260c0, 0x2, 0x4, 0x2, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:17 +0xae
created by main.main
  /home/kiki/gopro/test/main.go:15 +0xfc
exit status 2
```

上述堆栈信息显示了 panic 时所有的 goroutine 的状态，发生 panic 的 goroutine 会显示在最上面。

#### panic 的堆栈跟踪

```sh
# 最先 panic 的是 goroutine 6
goroutine 6 [running]:
# panic 位于 main.example 中，且定位到 main.go 第 28 行
# 参数传递
## 第一个参数是 []string，slice 是引用类型，前三个参数分别是 slice 的指针、长度和容量
## 第二个参数是 string，string 是引用类型，分别是 string 的指针和长度
## 第三个参数是 int，表示传入的 int 值
main.example(0xc0000260c0, 0x2, 0x4, 0x4c22f3, 0x5, 0x0, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:28 +0x211
## 第一个参数是 wg 的地址
main.main.func1(0xc000014110, 0xc0000260c0, 0x2, 0x4, 0x0, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:17 +0xae
created by main.main
  /home/kiki/gopro/test/main.go:15 +0xfc
```

## 示例2

```go
// main.go
package main

import (
  "fmt"
  "sync"
  "time"
)

type trace struct{}

func main() {
  slice := make([]string, 2, 4)
  var t trace
  wg := sync.WaitGroup{}
  stop := make(chan interface{})
  for i := 0; i < 3; i++ {
    wg.Add(1)
    go func(slice []string, i int, stop chan interface{}) {
      defer wg.Done()
      t.example(slice, "hello", i, stop)
    }(slice, i, stop)
  }
  wg.Wait()
  close(stop)
}

func (t *trace) example(slice []string, str string, i int, stop chan interface{}) {
  fmt.Println("example begin", str, i)
  if i == 0 {
    time.Sleep(2 * time.Second)
    fmt.Println(1 / i)
  }
  <-stop
  fmt.Println("example end", str, i)
}
```

```sh
kiki@ubuntu:~/gopro/test$ go run main.go
example begin hello 2
example begin hello 0
example begin hello 1
panic: runtime error: integer divide by zero

goroutine 6 [running]:
# 调用指针接收者
# 参数传递
## 第一个参数是结构体地址
main.(*trace).example(0x57ab60, 0xc0000260c0, 0x2, 0x4, 0x4c22f3, 0x5, 0x0, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:31 +0x211
main.main.func1(0xc000014100, 0x57ab60, 0xc0000260c0, 0x2, 0x4, 0x0, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:20 +0xbb
created by main.main
  /home/kiki/gopro/test/main.go:18 +0x10b
exit status 2
```

## 示例3

```go
// main.go
package main

import (
  "fmt"
  "sync"
  "time"
)

type trace struct{}

func main() {
  slice := make([]string, 2, 4)
  var t trace
  wg := sync.WaitGroup{}
  stop := make(chan interface{})
  for i := 0; i < 3; i++ {
    wg.Add(1)
    go func(slice []string, i int, stop chan interface{}) {
      defer wg.Done()
      t.example(true, false, true, uint8(i), stop)
    }(slice, i, stop)
  }
  wg.Wait()
  close(stop)
}

func (t *trace) example(b1, b2, b3 bool, i uint8, stop chan interface{}) {
  fmt.Println("example begin", i)
  if i == 0 {
    time.Sleep(2 * time.Second)
    fmt.Println(1 / i)
  }
  <-stop
  fmt.Println("example end", i)
}
```

```sh
kiki@ubuntu:~/gopro/test$ GOTRACEBACK=all go run main.go
example begin 2
example begin 0
example begin 1
panic: runtime error: integer divide by zero

goroutine 6 [running]:
# 如果多个参数可以填充到一个字，则这些参数值会合并传递
## 3 个布尔型和 1 个无符号 8 位整数可以合并到一个字 0x00 01 00 01
main.(*trace).example(0x57ab60, 0x10001, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:31 +0x176
main.main.func1(0xc000014110, 0x57ab60, 0xc0000260c0, 0x2, 0x4, 0x0, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:20 +0x84
created by main.main
  /home/kiki/gopro/test/main.go:18 +0x10b

goroutine 1 [semacquire]:
sync.runtime_Semacquire(0xc000014118)
  /usr/local/go/src/runtime/sema.go:56 +0x42
sync.(*WaitGroup).Wait(0xc000014110)
  /usr/local/go/src/sync/waitgroup.go:130 +0x64
main.main()
  /home/kiki/gopro/test/main.go:23 +0x131

goroutine 7 [chan receive]:
## 0x01 01 00 01
main.(*trace).example(0x57ab60, 0x1010001, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:33 +0xe1
main.main.func1(0xc000014110, 0x57ab60, 0xc0000260c0, 0x2, 0x4, 0x1, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:20 +0x84
created by main.main
  /home/kiki/gopro/test/main.go:18 +0x10b

goroutine 8 [chan receive]:
## 0x02 01 00 01
main.(*trace).example(0x57ab60, 0x2010001, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:33 +0xe1
main.main.func1(0xc000014110, 0x57ab60, 0xc0000260c0, 0x2, 0x4, 0x2, 0xc0000200c0)
  /home/kiki/gopro/test/main.go:20 +0x84
created by main.main
  /home/kiki/gopro/test/main.go:18 +0x10b
exit status 2
```

## 参考

- [Stack Traces In Go](https://www.ardanlabs.com/blog/2015/01/stack-traces-in-go.html)
- [调试利器：dump goroutine 的 stacktrace](https://colobu.com/2016/12/21/how-to-dump-goroutine-stack-traces/)
- [A whirlwind tour of Go’s runtime environment variables](https://dave.cheney.net/2015/11/29/a-whirlwind-tour-of-gos-runtime-environment-variables)
