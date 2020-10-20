# go 类型转换

- [go 类型转换](#go-类型转换)
  - [类型声明](#类型声明)
  - [类型转换](#类型转换)

## 类型声明

`type` 声明定义一个新的命名类型，它和某个已有类型使用同样的底层类型。命名类型提供了一种方式来区分底层类型的不同或不兼容使用，这样就不会在无意中混用他们。

比如，把不同计量单位的温度值转换为不同的类型。即使使用的底层类型相同，二者是不是相同的类型，不能使用算术表达式进行比较和合并。

别名类型与源类型的内部结构是一致的。两种类型之间的转换不会创建新的值。

```go
package main

import "fmt"

type celsius float64
type fahrenheit float64

const (
  absoluteZeroC celsius = -273.15
  freezingC     celsius = 0
  boilingC      celsius = 100
)

func main() {
  fmt.Printf("cToF(%2f) = (%2f)\n", absoluteZeroC, cToF(absoluteZeroC))
  // error:cannot use boilingC (type celsius) as type fahrenheit in argument to fToC
  fmt.Printf("fToC(%2f) = (%2f)\n", boilingC, fToC(boilingC))
  fmt.Printf("fToC(%2f) = (%2f)\n", boilingC, fToC(fahrenheit(boilingC)))
}

func cToF(c celsius) fahrenheit { return fahrenheit(c*9/5 + 32) }
func fToC(f fahrenheit) celsius { return celsius((f - 32) * 5 / 9) }
```

## 类型转换

对于每个类型 T，都有一个对应的类型转换操作 T(x) 将值转化为类型 T。如果两个类型具有相同的底层类型或二者都是指向相同底层类型变量的未命名指针类型，则二者可以相互转化。类型转化不改变类型值的表达方式，仅改变类型。如果 x 对于类型 T 是可赋值的，类型转化也是允许的，但是通常是不必要的。

- 类型转换用于将一种数据类型的变量转换为另一种类型的变量，go 不支持隐式类型转换。两种不同的类型即使互相兼容，也不能互相赋值。编译器不会对不同类型的值做隐式转换。

```go
package main

import (
  "fmt"
  "time"
)

func main() {
  var dur time.Duration
  var i int64

  i = 1000
  dur = time.Duration(i)
  // dur = i          //cannot use i (variable of type int64) as time.Duration value in assignment
  fmt.Println(dur) //1µs

  // i = dur // cannot use dur (variable of type time.Duration) as int64 value in assignmentcompil
  i = int64(dur)
  fmt.Println(i) //1000

  return
}
```

  ```go
  package main

  import (
    "fmt"
    "math"
  )

  func main() {
    var x, y int = 3, 4
    var f float64 = math.Sqrt(float64(x*x + y*y))
    var z uint = f //error: cannot use f (type float64) as type uint in assignment
    fmt.Println(x, y, z)
  }
  ```

- 格式`type_name(expression)`
  - type_name 是类型
  - expression 是表达式
