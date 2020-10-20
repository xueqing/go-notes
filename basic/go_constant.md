# go 常量

- [go 常量](#go-常量)
  - [常量](#常量)
  - [iota](#iota)
  - [数值型常量](#数值型常量)

## 常量

- 分类：
  - 布尔常量
  - 数值型常量：rune 常量(字符常量)、整数常量、浮点数常量、复数常量
  - 字符串常量
- 定义类似于变量声明，但是需要 `const` 关键字
  - `const vname [vtype] = value`
    - 显式类型定义`const vname vtype = value`
    - 隐式类型定义`const vname = value`
- **不能使用 `:=` 声明**
- 相同类型声明 `const vname1, vname2, vname3 = value1, value2, value3`
- 用于枚举

  ```go
  const {
      Unknown = 0
      Famale = 1
      Male = 2
  }
  ```

- 常量可使用 `len()`, `cap()`, `unsafe.Sizeof()` 函数计算表达式的值，函数必须是内置函数，否则编译错误

  ```go
  a = "abc"
  unsafe.Sizeof(a) //16，字符串类型在 go 中是个结构，包括指向数组的指针和长度，每部分都是 8 字节，所以是 16 个字节
  ```

- 常量必须在编译时确定(在编译时创建，包含定义在函数内部的局部变量)
- 常量可是无类型的，也可以是有类型的
  - 由字面量表示的常量以及由仅以无类型的常量作为其操作数的常量表达式的结果只都属于无类型常量
- 如果有一个未被显式赋值的常量，那么与它同一行的常量(如果有的话)的赋值也必须被省略
- 在未包含显式赋值的哪一行常量声明中的常量标识符的数量必须与它最上面、最近且包含显式赋值的哪一行常量声明中的常量标识符的数量相等

  ```go
  const (
    utc1, utc2 = 6.3, false
    utc3, utc4 // = 6.3, false
    utc5       = "C"
    utc6       //"C"
    utc7       //"C"
  )
  ```

## iota

- iota: 特殊常量，一个可被编译器修改的常量，代表了连续的、无类型的整数常量
  - 在 `const` 关键字出现时被重置为 0（const 内部的第一行之前）
  - `const` 中每新增一行常量声明，`iota` 计数一次

    ```go
    const (
      a = iota
      b // 1
      c // 2
    )

    const (
      a = 1 << iota
      b // 2
      c // 4
    )

    const (
      a, i = iota, 1 << iota
      b, j // 1, 2
      _, _
      c, k // 3, 8
    )
    ```

## 数值型常量

- 数值常量是高精度值，一个没有类型的常量根据上下文确定自身的类型

  ```go
  package main

  import "fmt"

  const (
    Big = 1 << 100
    Small = Big >> 99
  )

  func needInt(x int) int { return x*10 + 1 }
  func needFloat(x float64) float64 {
    return x * 0.1
  }

  func main() {
    fmt.Println(needInt(Small))
    fmt.Println(needInt(Big))     //error: constant 1267650600228229401496703205376 overflows int
    fmt.Println(needFloat(Small))
    fmt.Println(needFloat(Big))
  }
  ```
