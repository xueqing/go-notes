# go 散列表

- [go 散列表](#go-散列表)
  - [map 定义和初始化](#map-定义和初始化)
  - [修改 map](#修改-map)
  - [map 的键类型](#map-的键类型)
  - [map 是引用类型](#map-是引用类型)
  - [map 不是并发安全的类型](#map-不是并发安全的类型)

## map 定义和初始化

- map 是一种无序的键值对的集合，可以通过 key 快速检索数据，使用 hash 表实现
- 定义集合
  - `var map_name[key_type]val_type`
  - 使用 make 函数`map_name := make(map[key_type]val_type)`
- 不初始化 map，得到的是一个 nil map，不能用于存放键值对
- map 文法：类似 struct，但是需要键名
  - 当顶级类型是一个 type 的名字时，可以忽视

```go
package main

import "fmt"

type vertex struct {
    Lat, Long float64
}

var m = map[string]vertex{
    "Bell lab": vertex{
        40.68433, -74.39967,
    },
    "Google": vertex{
        37.42202, -122.08408,
    },
}

var m1 = map[string]vertex{
    "Bell lab": { 40.68433, -74.39967},
    "Google": { 37.42202, -122.08408},
}

func main() {
    fmt.Println(m)
    fmt.Println(m1)
}
```

## 修改 map

- 使用 `map_name[key]` 查看元素在集合中是否存在
  - 如果元素存在，返回的第一个元素是对应的键，第二个元素是 true
  - 元素不存在，返回值类型对应的零值，返回第二个元素是 false
- 插入或更新元素 `map_name[key] = val`
- delete 函数用于删除集合的元素`delete(map_name, key)`，指定元素名和对应的键

```go
package main

import (
    "fmt"
)

func main() {
    m := make(map[string]int)

    m["answer"] = 42
    fmt.Println("The value:", m["answer"])

    m["answer"] = 48
    fmt.Println("The value:", m["answer"])

    delete(m, "answer")
    fmt.Println("The value:", m["answer"])

    val, ok := m["answer"]
    fmt.Println("The value:", val, "Present?", ok)
}
```

## map 的键类型

**map 的键类型不能是函数、字典或切片**。因为键的类型必须是可比较的，即键的值可以作为比较操作符 `==` 和 `!=` 的操作数。

如果字典类型的键是接口类型，要求在程序运行期间，该类型的字典值的每个键值的动态类型都是可比较的，否则会发生运行时恐慌。

## map 是引用类型

map 持有一个指向某个底层数据结构值的引用，当传递 map 给一个函数时，函数内部对 map 的修改对外部是可见的。

## map 不是并发安全的类型

Go 语言官方认为，大多数使用字典值的场景不需要多线程场景下的安全访问控制，为了少数的并发使用场景二强制要求所有字典值满足互斥操作将降低大多数程序的速度，这是得不偿失的。
