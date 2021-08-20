# defer 语句需要注意的

- [defer 语句需要注意的](#defer-语句需要注意的)
  - [概述](#概述)
  - [在循环中使用 defer](#在循环中使用-defer)
  - [在代码块中使用 defer](#在代码块中使用-defer)
  - [调用 os.Exit 时不会执行 defer](#调用-osexit-时不会执行-defer)
  - [在匿名返回值和命名返回值中的结果不同](#在匿名返回值和命名返回值中的结果不同)
  - [defer 方法的接收者是指针还是值](#defer-方法的接收者是指针还是值)
  - [defer 函数为 nil](#defer-函数为-nil)
  - [defer 调用含有闭包的函数](#defer-调用含有闭包的函数)
  - [参考](#参考)

## 概述

defer 的执行时机：

- 包裹 defer 的函数返回时
- 包裹 defer 的函数执行到末尾时
- 所在的 goroutine 发生 panic 时

defer 的执行顺序：

- 先进后出，先遇到的 defer 语句入栈，最后再出栈执行

问题：调用 defer 增加了参数拷贝、入栈出栈的操作，导致代码执行耗时增加。

## 在循环中使用 defer

defer 按照先进后出的顺序执行，**在循环退出时才会执行栈内所有操作**，所以某些情况下不能及时释放循环内的资源。

另外，和直接调用相比，执行 defer 有额外开销，比如 defer 会对需要的参数进行内存拷贝，还要对 defer 结构进行入栈出栈操作。因此在循环中使用 defer 可能导致大量的资源开销。

```go
var body io.ReadCloser
for {
  body, err = openURL(sURL)
  if err != nil {
    logv.Errorf("get reader err(%v)", err)
    return err
  }
  defer body.Close()
  // handle body
}
```

修改方法：

1. 不使用 defer，处理完成之后手动释放资源
2. 封装函数(命名函数或匿名函数)，处理单次循环，在封装函数内可使用 defer

## 在代码块中使用 defer

代码块中的 defer 只会在代码块所属的函数执行结束后才会执行，比如 for/switch 代码块。匿名函数块不包含在内。

```go
func testCodeBlock() {
  fmt.Println("testCodeBlock: enter")
  {
    defer func() {
      fmt.Println("testCodeBlock: exit code block")
    }()
    fmt.Println("testCodeBlock: enter code block")
  }
  fmt.Println("testCodeBlock: exit")
}

func main() {
  testCodeBlock()
}

/*
testCodeBlock: enter
testCodeBlock: enter code block
testCodeBlock: exit
testCodeBlock: exit code block
*/
```

## 调用 os.Exit 时不会执行 defer

```go
func testExit() {
  fmt.Println("testExit: enter")
  defer fmt.Println("testExit: exit") // 不会打印
  os.Exit(-1)
}
```

## 在匿名返回值和命名返回值中的结果不同

对于匿名返回值，defer 内的修改不会影响实际返回结果；对于命名返回值，defer 内的修改会影响返回结果。

```go
func testUnNamedReturn() int {
  ret := 1
  fmt.Println("testUnNamedReturn: enter")
  defer func() {
    fmt.Printf("testUnNamedReturn: defer ret(%d)\n", ret) // 2
    ret++
    fmt.Printf("testUnNamedReturn: defer ret(%d)\n", ret) // 3
  }()
  ret++
  return ret
}

func testNamedReturn() (ret int) {
  ret = 1
  fmt.Println("testNamedReturn: enter")
  defer func() {
    fmt.Printf("testNamedReturn: defer ret(%d)\n", ret) // 2
    ret++
    fmt.Printf("testNamedReturn: defer ret(%d)\n", ret) // 3
  }()
  ret++
  return ret
}

func main() {
  fmt.Printf("testUnNamedReturn(%d)\n", testUnNamedReturn()) // 2
  fmt.Printf("testNamedReturn(%d)\n", testNamedReturn())     // 3
}
```

## defer 方法的接收者是指针还是值

defer 方法的接收者是值时，调用 defer 入栈时会拷贝值，后续对其操作对方法不可见；defer 方法的接收者是指针时，入栈的是指向同一对象的指针，后续对该对象的操作在 defer 方法真正执行时是可见的。

```go
type people struct {
  name string
  age  int
}

func (p people) getName() {
  fmt.Println("name: ", p.name)
}

func (p *people) getAge() {
  fmt.Println("age: ", p.age)
}

func testMethodReceiver() {
  p := people{
    name: "kiki",
    age:  8,
  }
  defer p.getName() // kiki
  defer p.getAge()  // 10

  p.name = "kiki2"
  p.age = 10
}

func main() {
  testMethodReceiver()
}
```

## defer 函数为 nil

将 defer 函数赋值为 nil，在调用 defer 入栈的时候并没有执行函数，不会崩溃。在包裹 defer 的函数执行结束时真正调用 defer 入栈的函数时才会崩溃。

```go

func testNilDefer() {
  fmt.Println("testNilDefer: enter")
  var run func() = nil
  defer run() // no panic here

  fmt.Println("runs")
  fmt.Println("testNilDefer: exit")
}

/*
testNilDefer: enter
runs
testNilDefer: exit
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x499180]

goroutine 1 [running]:
main.testNilDefer()
        /home/kiki/test/main.go:56 +0x160
main.main()
        /home/kiki/test/main.go:59 +0x25
exit status 2
*/
```

## defer 调用含有闭包的函数

```go
type database struct{}

func (db *database) connect() (disconnect func()) {
  fmt.Println("connect")

  return func() {
    fmt.Println("disconnect")
  }
}

func testClosureFunc() {
  db := &database{}
  // defer db.connect() // will not print "disconnect"
  disconnect := db.connect()
  defer disconnect()

  fmt.Println("query db...")
}

func main() {
  testClosureFunc()
}
```

## 参考

- [Go语言中defer的一些坑](https://www.jianshu.com/p/79c029c0bd58)
- [新手必看：Go 中 defer 的 5 个坑 - 第一部分](https://zhuanlan.zhihu.com/p/134121503)
