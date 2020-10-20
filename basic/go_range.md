# go 范围

- range 关键字右侧的表达式只在迭代开始前计算一次；左边的表达式每次迭代都会计算
- range 关键字用于 for 循环中迭代数组 array、切片 slice、通道 channel、字典 map 或字符串 string 的元素
  - 数组、切片返回两个值：元素的索引、索引对应的值的**拷贝**
  - 通道返回一个值：元素的值
  - 字典返回两个值：key、value (迭代之前删掉的元素不会显示；迭代过程中添加的元素不确定是否可以显示)
  - 字符串返回两个值：元素的索引、索引对应的元素值(类型为 rune，因此索引可能是不连续的)
- 如果 range 表达式的结果类型是数组或某个执行数组值的指针类型，且只获取第一个迭代值，那么这个 range 表达式只会被部分求值

  ```go
  package main

  import (
      "fmt"
  )

  func main() {
      var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}
      for i, v := range pow {
          fmt.Printf("2**%d = %d\n", i, v)
      }
  }
  ```

- 上述返回两个值的情况，如果不显示获取可以只获取第一个值
- 可以赋值给 `_` 跳过索引或值
  - `for i, _ := range pow`。如果只想要索引，可以忽视第二个参数 `for i := range pow`
  - `for _, val := range pow`

  ```go
  package main

  import (
      "fmt"
  )

  func main() {
      pow := make([]int, 10)

      for i := range pow {
          pow[i] = 1 << uint(i) //2**i
      }

      for _, val := range pow {
          fmt.Printf("%d\n", val)
      }
  }
  ```
