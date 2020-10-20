# go 条件语句

- [go 条件语句](#go-条件语句)
  - [if](#if)
    - [if 语句](#if-语句)
    - [if...else 语句](#ifelse-语句)
    - [if 嵌套语句](#if-嵌套语句)
  - [switch](#switch)
    - [表达式 switch](#表达式-switch)
      - [switch 可包含简短的声明语句](#switch-可包含简短的声明语句)
      - [没有条件的 switch 语句](#没有条件的-switch-语句)
      - [case 中的 fallthrough/break 语句](#case-中的-fallthroughbreak-语句)
    - [类型 switch 语句](#类型-switch-语句)
  - [select](#select)

## if

### if 语句

- 类似于 for 循环，但是不需要小括号，需要大括号

  ```go
  if x < 0 {
      return 0
  }
  ```

  - if 语句可以在执行条件语句之前有一个简短的语句，在这个语句声明的变量的作用范围在 if 语句末尾结束

    ```go
    func pow(x, n, lim float64) float64 {
        if v := math.Pow(x, n); v < lim {
            return v
        }
        return lim
    }
    ```

### if...else 语句

- if 语句中声明的变量对于匹配的 else 代码块也是可见的

  ```go
  func pow(x, n, lim float64) float64 {
      if v := math.Pow(x, n); v < lim {
          return v
      } else {
          fmt.Printf("%g >= %g\n", v, lim)
      }
      return lim
  }
  ```

### if 嵌套语句

```go
if {
    //...
    if {
        //...
    }
}
```

## switch

### 表达式 switch

- 每个 case 分支都是唯一的，从上到下测试直到匹配，**只执行匹配项**，匹配项后面不用加 break(每个测试项后面自动加上 break)，和 C++ 不同
- 另外一个区别是每个 case 不需要是常数，值也不必是整数
- default case 最多只能有一个，且不是必须作为最后一个 case 出现。如果没有找到匹配的 case，存在 default case 时会执行 default case

  ```go
  switch var1 { //case 的值必须是相同类型
      case var1[, var11, var111]://可同时测试多个条件，用逗号分隔
          //...
      case var2:
          //...
      default:
          //...
  }
  ```

#### switch 可包含简短的声明语句

- switch 也有简短的声明语句，声明变量只对 switch 范围可见

  ```go
  package main

  import (
      "fmt"
      "runtime"
  )

  func main() {
      switch os := runtime.GOOS; os {
      case "darwin":
          fmt.Println("OS X")
      case "linux":
          fmt.Println("Linux")
      default:
          fmt.Printf("%s\n", os)
      }
  }
  ```

#### 没有条件的 switch 语句

- 没有条件语句的 switch 和 `switch true` 相同。这个可以用于实现比较长的 if-then-else 链

  ```go
  package main

  import (
      "fmt"
      "time"
  )

  func main() {
      t := time.Now()
      switch {
      case t.Hour() < 12:
          fmt.Println("Good morning")
      case t.Hour() < 17:
          fmt.Println("Good afternoon")
      default:
          fmt.Println("Good evening")
      }
  }
  ```

#### case 中的 fallthrough/break 语句

- case 语句列表最后一句可以是 fallthrough 语句。这会将流程控制权转移到下一条 case 语句。
  - fallthrough 只能是 case 语句列表的最后一句
  - fallthrough 不能在最后一个 case
- case 语句列表也可以出现 break 语句。break 语句包含一个关键字 break 和一个可选的标记

### 类型 switch 语句

- 类型 switch 语句对；类型进行判定，而不是值。其他方面与表达式 switch 相同。这个 switch 表达式类似于类型断言的表达式

```go
switch v.(type) {
case int, uint:
    fmt.Printf("type is integer %d\n", v)
case string:
    fmt.Printf("type is string %s\n", v.(string))
default:
    fmt.Printf("I don't know about type %T!\n", v)
}
```

- 这里的变量 v 必须是接口类型。否则，已知 v 的类型就不需要类型选择。另外，case 中的类型字面量必须是 v 的类型的实现类型
- 通用的方案是将 v 的类型设置为 `interface{}`，这样所有类型都是它的实现类型
- case 表达式中的类型字面量可以是 nil
- case 表达式中不能出现 fallthrough
- 直接使用 `v.(type)` 优化流程：不再需要显式类型转换，相当于每个 case 语句列表开始的时候自动调用 `i := v.(actual_type)`

```go
switch i := v.(type) {
case int, uint:
    fmt.Printf("type is integer %d\n", i)
case string:
    fmt.Printf("type is string %s\n", i)
default:
    fmt.Printf("I don't know about type %T!\n", i)
}
```

## select

- `select` 语句使一个 goroutine 可以等待多个通信操作
- 类似于 `switch` 语句，会随机执行一个可允许的 `case`，如果没有 `case` 可以允许则阻塞到有 `case` 可以运行

  ```go
  select {
      case condition1:
          //...
      case condition2:
          //...
      default:
          //...
  }
  ```

  - `case` 必须是一个通信，所有 `channel` 表达式会被求值
  - **case 右边的发送和接收语句中的通道表达式和元素表达式都会被先求值。求值顺序是自上而下、从左到右的**
  - 任意某个 `channel` 可以进行，就会执行，其他的被忽略
  - 如果多个 `case` 可以执行，会随机公平选择一个执行，忽略其他
  - 没有可以执行的 `case` 语句
    - 如果有 `default`，则执行 `default`
    - 否则阻塞至某个通信可以运行，go 不会重新对 `channel` 或值进行求值
    - 为了在尝试发送或者接收时不发生阻塞，可使用 `default` 分支

```go
package main

import "fmt"

func fibonacci(c, quit chan int) {
    x, y := 0, 1
    for {
        select {
        case c <- x:
            x, y = y, x+y
        case <-quit:
            fmt.Println("quit")
            return
        }
    }
}

func main() {
    c := make(chan int)
    quit := make(chan int)
    go func() {
        for i := 0; i < 10; i++ {
            fmt.Println(<-c)
        }
        quit <- 0
    }()
    fibonacci(c, quit)
}
```
