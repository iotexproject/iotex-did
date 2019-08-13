package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/iotexproject/iotex-DID/protocol"
)

func main() {
	pbKey := "029a4774d53094deaf34663e67272e2f33b6d91b0b7995fade0fab23"
	contextContent := "www.iotex.io"

	newPbKey := "0268ccc80007f82d49c2f2ee25a9dae856559330611f0a62356e59ec8cdb566e69"

	did, err := protocol.CreateDIDByPbkey(pbKey, contextContent)
	if err != nil {
		fmt.Println(err)
	}
	displayDoc()

	err = protocol.AddPbKey(did, newPbKey)
	if err != nil {
		return
	}
	displayDoc()

	err = protocol.RemovePbKey(did, newPbKey)
	if err != nil {
		return
	}
	displayDoc()

	err = protocol.DeleteDID(did)
	if err != nil {
		return
	}

}
func displayDoc() {
	url := "http://localhost:3003/029a4774d53094deaf34663e67272e2f33b6d91b0b7995fade0fab23"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("\n")
		return
	}

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", html)
}
