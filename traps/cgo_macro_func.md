# cgo 调用 C 中的宏函数

## 结论

cgo 目前只支持可以规约为常量或变量的宏，不支持宏函数。需要添加封装函数才能在 cgo 中使用。比如下面的代码

```go
package main

/*
#define SUM(a,b) (a)+(b)
int sum(int a, int b) {
  return SUM(a,b);
}
#define sum2(a,b) sum(a,b)
*/
import "C"

func main() {
  // println(C.SUM(1, 2))  // error: could not determine kind of name for C.SUM
  println(C.sum(1, 2)) // ok
  // println(C.sum2(1, 2)) // error: could not determine kind of name for C.sum2
}
```

## 参考

- [Calling cgo macro function](https://stackoverflow.com/questions/51904947/calling-cgo-macro-function)
- [github issue](https://github.com/golang/go/commit/03876af91c50c6e0227218a856f037dd20a45729)
