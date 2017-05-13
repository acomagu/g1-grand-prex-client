package main

import (
	"fmt"
	"net/http"
	"io"
	"bufio"
	"os"
	"os/exec"
)

const (
	incomingWebhookURL = ""
	outgoingWebhookURL = ""
)

func main() {
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		io.Copy(os.Stdout, stderr)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(stdin, r.Body); err != nil {
			fmt.Println(err)
			return
		}

		reader := bufio.NewReader(stdout)
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			return
		}

		if _, err := w.Write(line); err != nil {
			fmt.Println(err)
			return
		}
	})
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(http.ListenAndServe(":" + os.Args[1], nil))
}
