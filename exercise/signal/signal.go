package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	sigrecv1 := make(chan os.Signal, 1)
	sigsli1 := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	signal.Notify(sigrecv1, sigsli1...)

	sigrecv2 := make(chan os.Signal, 1)
	sigsli2 := []os.Signal{syscall.SIGQUIT}
	signal.Notify(sigrecv2, sigsli2...)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for sig := range sigrecv1 {
			fmt.Println("received siganl from sigrecv1:", sig)
		}
		fmt.Println("sigrecv1 end...")
		wg.Done()
	}()

	go func() {
		for sig := range sigrecv2 {
			fmt.Println("received siganl from sigrecv2:", sig)
		}
		fmt.Println("sigrecv2 end...")
		wg.Done()
	}()

	fmt.Println("wait for 20 seconds...")
	time.Sleep(20 * time.Second)

	fmt.Println("stop notify sigrecv1")
	signal.Stop(sigrecv1)

	fmt.Println("close sigrecv1")
	close(sigrecv1)

	wg.Wait()

	fmt.Println("main over")

	return
}
