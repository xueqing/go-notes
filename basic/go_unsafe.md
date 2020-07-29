# 关于 unsafe 包

- [关于 unsafe 包](#关于-unsafe-包)
  - [概述](#概述)
  - [Go 指针](#go-指针)
  - [Sizeof 函数](#sizeof-函数)
  - [Alignof 函数](#alignof-函数)
  - [参考](#参考)

## 概述

unsafe 包可以想 C 一样操作计算机内存，绕过了 Go 的内存安全原则，不利于程序的扩展和维护，不建议使用。

## Go 指针

Go 指针包括：

- `*`: 普通指针，用于传递对象地址，不能进行指针运算
- `unsafe.Pointer`: 通用类型指针，用于转换不同类型的指针，不能进行指针运算
- `uintptr`: 用于指针运算，GC 不把 uintptr 当做指针，会回收 uintptr 的目标对象；uintptr 不能持有对象

其中，普通指针和 unsafe.Pointer 可以互相转换，unsafe.Pointer 和 uintptr 可以互相转换。因此，**unsafe.Pointer 实现任意类型的指针互相转换，也可以将任意指针转换为 uintptr 进行指针运算**。

Pointer 和 uinptr 在 Go 中和 int 字节长度一致。

```go
func main() {
  var a int
  p := unsafe.Pointer(&a)
  // 8 8
  fmt.Println(unsafe.Sizeof(p), unsafe.Sizeof(uintptr(p)))
}
```

## Sizeof 函数

```go
func main() {
  var a int
  p := unsafe.Pointer(&a)
  // uintptr
  fmt.Println(reflect.TypeOf(unsafe.Sizeof(p)))

type st1 struct {
    a byte
    b int32
    c int64
  }
  var t1 st1
  // 8 16 1 4 8
  fmt.Println(unsafe.Sizeof(&t1), unsafe.Sizeof(t1), unsafe.Sizeof(t1.a), unsafe.Sizeof(t1.b), unsafe.Sizeof(t1.c))

  type st2 struct {
    a byte
    c int64
    b int32
  }
  var t2 st2
  // 8 24 1 4 8
  fmt.Println(unsafe.Sizeof(&t2), unsafe.Sizeof(t2), unsafe.Sizeof(t2.a), unsafe.Sizeof(t2.b), unsafe.Sizeof(t2.c))
}
```

## Alignof 函数

```go
func main() {
  var a int
  p := unsafe.Pointer(&a)
  // 8 8
  fmt.Println(unsafe.Sizeof(p), unsafe.Sizeof(uintptr(p)))

  // uintptr
  fmt.Println(reflect.TypeOf(unsafe.Sizeof(p)))

  type st1 struct {
    a byte
    b int32
    c int64
  }
  var t1 st1
  // 8 8 1 4 8
  fmt.Println(unsafe.Alignof(&t1), unsafe.Alignof(t1), unsafe.Alignof(t1.a), unsafe.Alignof(t1.b), unsafe.Alignof(t1.c))

  type st2 struct {
    a byte
    c int64
    b int32
  }
  var t2 st2
  // 8 8 1 4 8
  fmt.Println(unsafe.Alignof(&t2), unsafe.Alignof(t2), unsafe.Alignof(t2.a), unsafe.Alignof(t2.b), unsafe.Alignof(t2.c))
}
```

## 参考

- [Golang unsafe](https://golang.org/pkg/unsafe/)
- [Go 语言之 unsafe 包介绍及使用](https://www.jianshu.com/p/c85fc3e31249)
- [go unsafe包](https://blog.csdn.net/cc7756789w/article/details/51241382)
