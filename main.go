package main

import (
	"fmt"
	"os"
	"rabbitmq-golang/example"
)

func main() {
	if len(os.Args) < 2 {
		panic("specify what you want to run")
	}
	fmt.Println(os.Args[1])
	if os.Args[1] == "worker" {
		example.StartWorker()
	}

	if os.Args[1] == "client" {
		example.StartClient()
	}
}
