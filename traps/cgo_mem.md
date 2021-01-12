# cgo 内存管理

```go
package main

//#cgo pkg-config: libavutil libavformat
//#include <libavutil/samplefmt.h>
//#include <libavformat/avformat.h>
import "C"
import (
  "fmt"

  "github.com/xueqing/goav/libavformat"
  "github.com/xueqing/goav/libavutil"
)

var (
  pFmtCtx *libavformat.AvFormatContext
)

func openInputFile(fileName string) (err error) {
  if ret := libavformat.AvformatOpenInput(&pFmtCtx, fileName, nil, nil); ret < 0 {
    err = fmt.Errorf("AvformatOpenInput error(%v)", libavutil.ErrorFromCode(ret))
  }
  return
}

func openInputFile2(fileName string) (err error) {
  var pFmtCtx2 *libavformat.AvFormatContext
  defer pFmtCtx2.AvformatCloseInput()
  if ret := libavformat.AvformatOpenInput(&pFmtCtx2, fileName, nil, nil); ret < 0 {
    err = fmt.Errorf("AvformatOpenInput error(%v)", libavutil.ErrorFromCode(ret))
  }
  return
}

// refer doc/examples/filtering_audio.c
func main() {
  fmt.Printf("Start filtering audio example\n")
  libavutil.AvLogSetLevel(48)
  filename := "/home/kiki/github/goav/resource/20s.flv"

  err := openInputFile2(filename)
  fmt.Printf("open input file2 error(%v)\n", err)
  // open input file2 error(<nil>)

  defer pFmtCtx.AvformatCloseInput()
  err = openInputFile(filename)
  fmt.Printf("open input file error(%v)\n", err)
  // panic: runtime error: cgo argument has Go pointer to Go pointer
  return
}
```

一个指针是 Go 指针还是 C 指针是由内存分配方式动态确定的，和指针类型无关。在 cgo 调用的 C 语言函数返回前，cgo 保证传入的 go 语言内存在此期间不会发生移动。

Go 的垃圾回收器在运行时会对内存进行一些操作(栈收缩、栈扩容)，会导致存储数据的地方发生改变。Go 调度器需要控制内存声明周期，将 Go 变量传给 cgo 后，调度器无法知道该变量的状态(何时使用、何时不再使用)。因此 Go 调度器会检测传递给 cgo 的指针，如果将 Go 指针传递给 cgo 会报错：`panic: runtime error: cgo argument has Go pointer to Go pointer`。
