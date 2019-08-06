package main

import (
	"fmt"

	"github.com/iotexproject/iotex-DID/protocol"
)

func main() {
	err := protocol.CreateDIDByProcessPbkey()
	if err != nil {
		fmt.Println(err)
	}
}
