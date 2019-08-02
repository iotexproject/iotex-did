package protocol

import (
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multihash"
	"github.com/ockam-network/did"
	"golang.org/x/crypto/sha3"
)

func processPbkey() {
	// hash the public key
	pbKey := "029a4774d543094deaf342663ae672728e12f03b3b6d9816b0b79995fade0fab23"
	pbHash := sha3.Sum256([]byte(pbKey))
	idString := pbHash[:]
	idString = pbHash[len(idString)-20:]
	// prepend the multihash label for the hash algo, skip the varint length of the multihash, since that is fixed to 20
	idString = append([]byte{multihash.SHA3_256}, idString...)
	// base58 encode the above value
	id := base58.Encode(idString)
	d := &did.DID{Method: "iotex", ID: id}
	fmt.Println(d)
	// we got our DID in d variable
}
