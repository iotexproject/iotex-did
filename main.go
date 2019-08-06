package main

import (
	"fmt"

	"github.com/iotexproject/iotex-DID/protocol"
)

func main() {
	err := protocol.ProcessPbkey()
	if err != nil {
		fmt.Println(err)
	}
}
