package main

import (
	"fmt"

	"github.com/qqcc123/orange-ci-agent/backend/tcpserver"
)

func main() {
	fmt.Println("---------------")
	tcpserver.RunSocket()
}
