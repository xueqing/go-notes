# 作用域

- 编译器遇到一个名字的引用时，从最内层的封闭词法块到全局块寻找其声明
  - 没有找到会报 “undeclared name” 错误
  - 内层和外层都存在声明时，内层的先被找到。此时内层声明会覆盖外部声明，外部声明将不可访问
- **短变量声明**依赖一个明确的作用域。**只有在同一个词法块中已经存在变量**的情况下，短声明的行为才和赋值操作一样，外层声明将被忽略。下面的代码容易覆盖全局声明的 `cwd`
  - 因为 cwd 和 err 在函数块内部都尚未声明，所以 `:=` 语句将它们视为局部变量。内存 cwd 声明使得外部声明不可见
  - 解决方法是不使用 `:=`，而是使用 `var` 声明变量

    ```go
    package main

    import (
      "fmt"
      "os"
    )

    var cwd string

    func main() {
      cwd, err := os.Getwd() //compile error: cwd declared and not used
      if err != nil {
        fmt.Printf("err=%s\n", err)
      }
    }
    ```

- 下面的代码，main 函数中的 cwd 是局部变量，会覆盖全局变量

```go
package main

import (
  "log"
  "os"
)

var cwd string = "."

func main() {
  mylog()
  cwd, err := os.Getwd()
  if err != nil {
    log.Fatalf("os.Getwd failed: %v\n", err)
  }
  log.Printf("os.Getwd success: %s\n", cwd)
  mylog()
}

func mylog() {
  log.Printf("global cwd: %s\n", cwd)
}
```
