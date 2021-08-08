package main

import (
  "fmt"
)

func test() {
  for i := 0; i < 5; i++ {
    fmt.Printf("now i=%v\n", i)
    if i > 2 {
      defer fmt.Printf("test if %v\n", i)
    }
    defer fmt.Printf("test %v\n", i)
    fmt.Printf("exit i=%v\n", i)
  }
}

func main() {
  test()
}

/*
output
    now i=0
    exit i=0
    now i=1
    exit i=1
    now i=2
    exit i=2
    now i=3
    exit i=3
    now i=4
    exit i=4
    test 4
    test if 4
    test 3
    test if 3
    test 2
    test 1
    test 0
*/
