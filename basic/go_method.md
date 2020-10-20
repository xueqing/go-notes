# go 方法

- [go 方法](#go-%e6%96%b9%e6%b3%95)
  - [方法 method](#%e6%96%b9%e6%b3%95-method)
    - [方法的接收者是 struct](#%e6%96%b9%e6%b3%95%e7%9a%84%e6%8e%a5%e6%94%b6%e8%80%85%e6%98%af-struct)
    - [方法的接收者是非结构体](#%e6%96%b9%e6%b3%95%e7%9a%84%e6%8e%a5%e6%94%b6%e8%80%85%e6%98%af%e9%9d%9e%e7%bb%93%e6%9e%84%e4%bd%93)
  - [指针接收者 vs 值接收者](#%e6%8c%87%e9%92%88%e6%8e%a5%e6%94%b6%e8%80%85-vs-%e5%80%bc%e6%8e%a5%e6%94%b6%e8%80%85)
    - [使用指针接收者](#%e4%bd%bf%e7%94%a8%e6%8c%87%e9%92%88%e6%8e%a5%e6%94%b6%e8%80%85)
    - [使用值接收者](#%e4%bd%bf%e7%94%a8%e5%80%bc%e6%8e%a5%e6%94%b6%e8%80%85)
  - [方法的匿名域](#%e6%96%b9%e6%b3%95%e7%9a%84%e5%8c%bf%e5%90%8d%e5%9f%9f)

## 方法 method

### 方法的接收者是 struct

- go 没有类，但是可以为 struct 类型定义方法
- 方法是一类带特殊的“接收者”参数的函数
  - 方法接收者在参数列表内，位于 `func` 关键字和方法名之间
  - 接收者可以是命名类型或者结构体类型的个值或一个指针
  - 接收者的基本类型不能是接口类型，也不能是指针类型
  - 所有给定类型的方法属于该类型的方法集
- 语法格式

  ```go
  func (v_name v_type) func_name() [return_type] {
      //func body
  }
  ```

- 示例

  ```go
  package main

  import "fmt"

  type Circle struct {
      radius float64
  }

  func (c Circle) getArea() float64 {
      return 3.14 * c.radius * c.radius
  }

  func main() {
      var c1 Circle
      c1.radius = 10.00
      fmt.Println(c1.getArea())
  }
  ```

### 方法的接收者是非结构体

- 也可以为非 struct 类型声明方法，定义作用于一个类型的方法
  - 接收者的类型定义和方法声明必须在同一个包内，不能为内建类型声明方法
  - 可为内置类型起一个别名，然后基于别名作为接收者定义方法

```go
package main

import (
    "fmt"
    "math"
)

type MyFloat float64

func (f MyFloat) Abs() float64 {
    if f < 0 {
        return float64(-f)
    }
    return float64(f)
}

func main() {
    f := MyFloat(-math.Sqrt2)
    fmt.Println(f.Abs())
}
```

## 指针接收者 vs 值接收者

### 使用值接收者

- 使用值接收者，类似于形参，方法内部的修改不影响调用者。方法会得到调用的值的副本。
  - 当接收者类型是引用类型(指针、切片、map、channel)的别名时，那么使用值作为接收者的方法内部的修改对外部也时可见的
- 接收一个值作为参数的函数必须接受一个指定类型的值，而以值作为接收者的方法被调用时，接收者可以是值或者指针

### 使用指针接收者

- 使用指针接收者，在方法内部修改会影响调用者
  - 场景 1：希望方法内部修改影响调用者
  - 场景 2：拷贝数据结构的代价比较大

  ```go
  package main

  import (
      "fmt"
      "math"
  )

  type Vertex struct {
      X, Y float64
  }

  func (v Vertex) Abs() float64 {
      return math.Sqrt(v.X*v.X + v.Y*v.Y)
  }

  func (v *Vertex) Scale(f float64) {
      v.X = v.X * f
      v.Y = v.Y * f
  }

  func main() {
      v := Vertex{3, 4}
      v.Scale(10)
      fmt.Println(v.Abs())
  }
  ```

- 带指针参数的函数必须接受一个指针，而以指针为接收者的方法被调用时，接收者可以是值或者指针，go 会根据接收者类型自动调整
  - 当使用值调用方法时，会自动对值取地址，使用其指针调用对应的方法
- 即某个类型的值只能调用值作为接收者的方法；但是该类型的指针可以调用值作为接收者的方法，也可以调用作为指针作为接收者的方法

  ```go
  package main

  import (
      "fmt"
  )

  type vertex struct {
      X, Y float64
  }

  func (v *vertex) Scale(f float64) {
      v.X = v.X * f
      v.Y = v.Y * f
  }

  func scaleFunc(v *vertex, f float64) {
      v.X = v.X * f
      v.Y = v.Y * f
  }

  func main() {
      v := vertex{3, 4}
      v.Scale(2)
      scaleFunc(&v, 10)

      p := &vertex{4, 3}
      p.Scale(3)
      scaleFunc(p, 8)

      fmt.Println(v, p)
  }
  ```

## 通过接口类型的值调用方法

与直接通过值或者指针调用方法不同，如果通过接口类型的值调用方法，那么：使用指针作为接收者声明者的方法，只能在接口类型的值是一个指针的时候被调用；使用值作为接收者声明的方法，在接口类型的值为值或者指针时，都可以被调用。

```go
package main

import (
  "fmt"
)

type myInterface interface {
  SayHello()
}

type implePointer struct{}

func (i *implePointer) SayHello() {
  fmt.Printf("I am a pointer\n")
}

type impleValue struct{}

func (i impleValue) SayHello() {
  fmt.Printf("I am a value\n")
}

func main() {
  var (
    ip    implePointer
    ipptr *implePointer
    iv    impleValue
    ivptr *impleValue

    im myInterface
  )

  ip.SayHello()
  // im = ip //cannot use ip (variable of type implePointer) as myInterface value in assignment: missing method SayHello

  ipptr.SayHello()
  im = ipptr
  im.SayHello()

  ipptr2 := &ip
  ipptr2.SayHello()
  im = ipptr2
  im.SayHello()

  iv.SayHello()
  im = iv
  im.SayHello()

  // ivptr.SayHello() //panic: runtime error: invalid memory address or nil pointer dereferenc
  im = ivptr
  // im.SayHello() //panic: value method main.impleValue.SayHello called using nil *impleValue pointer

  ivptr2 := &iv
  ivptr2.SayHello()
  im = ivptr2
  im.SayHello()

  return
}
```

## 方法的匿名域

- 方法的接收者是一个结构体的匿名域（结构体中的结构体），可直接调用不指定匿名域

```go
type address struct {
    city string
    state string
}

func (a address) fullAddress() {
    fmt.Println("Full address: %s, %s", a.city. a.state)
}

type person struct {
    firstName string
    lastName string
    address
}

func printPersonInfo(p person) {
    fmt.Println("name: %s %s", p.firstName, p.secondName)
    p.fullAddress()//p.address.fullAddress()
}
```
