# sync.WaitGroup 不是引用类型

WaitGroup 对象不是引用类型，在通过函数传递时需要使用其指针，否则会进入死锁状态。

```go
func work(id int, wg *sync.WaitGroup) {
  fmt.Println("work: id=", id)
  wg.Done()
}

func main() {
  var wg sync.WaitGroup
  wg.Add(10)
  for i:=0; i<10; i++ {
    go work(i, &wg)
  }
  wg.Wait()
}
```
