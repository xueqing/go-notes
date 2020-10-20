package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

func main() {
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "code")

	stdout1, err := cmd1.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := cmd1.Start(); err != nil {
		fmt.Println(err)
		return
	}

	stdin2, err := cmd2.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	bufr := bufio.NewReader(stdout1)
	bufr.WriteTo(stdin2)
	bufw := bytes.Buffer{}
	cmd2.Stdout = &bufw

	if err := cmd2.Start(); err != nil {
		fmt.Println(err)
		return
	}

	// important: use close to complete data transfer
	if err := stdin2.Close(); err != nil {
		fmt.Println(err)
		return
	}

	if err := cmd2.Wait(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bufw.String())
	return
}
