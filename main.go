package main

import (
	"fmt"

	"github.com/iotexproject/iotex-DID/protocol"
)

func main() {

	err := protocol.AddPbKey("did:iotex:627523e8023f485b7676e74cdfa77de2f098e0abfe36f03b295db5c3b01d34ef", "ffffff")
	if err != nil {
		fmt.Println(err)
	}
}
