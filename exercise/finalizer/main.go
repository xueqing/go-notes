package main

import (
	"fmt"
	"runtime"
	"time"
)

type B struct {
	n int
}

func NewB(n int) *B {
	b := &B{
		n: n,
	}
	runtime.SetFinalizer(b, func(pB *B) {
		fmt.Printf("release B(%v)\n", pB.n)
	})
	return b
}

type A struct {
	ptr *B
}

func NewA(b *B) *A {
	a := &A{
		ptr: b,
	}
	runtime.SetFinalizer(a, func(pA *A) {
		fmt.Printf("release A(%v)\n", pA.ptr)
	})
	return a
}

func test() {
	b := NewB(555)
	fmt.Printf("NewB: %v\n", b)
	a := NewA(b)
	fmt.Printf("NewA: %v\n", a)
}

func main() {
	test()
	// if main exit here, it may not run A and B's Finalizer
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Printf("Sleep %vs\n", i+1)
		runtime.GC()
	}
	// output:
	//	NewB: &{555}
	//	NewA: &{0xc0000140e0}
	//	Sleep 1s
	//	release A(&{555})
	//	Sleep 2s
	//	release B(555)
	//	Sleep 3s
	//	Sleep 4s
	//	Sleep 5s
}
