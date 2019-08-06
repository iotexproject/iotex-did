package protocol

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/iotexproject/iotex-address/address"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multihash"
	"golang.org/x/crypto/sha3"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-antenna-go/v2/utils/wait"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

const (
	host = "api.testnet.iotex.one:443"
)

const contractAddress = "io1hfgtmdc27uzd7g47ky42q00wt9z37st87ghqjp"

const constDID = `[
	{
		"constant": false,
		"inputs": [
			{
				"name": "idStr",
				"type": "string"
			}
		],
		"name": "createDID",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": true,
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "idStr",
				"type": "string"
			}
		],
		"name": "isExist",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	}
]`

// ProcessPbkey is processing public key to DID
func ProcessPbkey() error {
	// hash the public key
	pbKey := "029a4774d543094deaf342663ae672728e12f03b3b6d9816b0b79995fade0fab23"
	pbHash := sha3.Sum256([]byte(pbKey))
	idString := pbHash[:]
	idString = pbHash[len(idString)-20:]
	// prepend the multihash label for the hash algo, skip the varint length of the multihash, since that is fixed to 20
	idString = append([]byte{multihash.SHA3_256}, idString...)
	// base58 encode the above value
	id := base58.Encode(idString)
	d := "did:iotex:" + id
	// we got our DID in d variable
	// Create grpc connection
	conn, err := iotex.NewDefaultGRPCConn(host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// Add account by private key
	acc, err := account.HexStringToAccount("16668401a3f686717b457bcde6c182cb46dffce81b3bff9ea23f6bd41ac4d54a")
	if err != nil {
		log.Fatal(err)
	}
	// create client
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	didABI, err := abi.JSON(strings.NewReader(constDID))
	if err != nil {
		return err
	}
	didContractAddress, err := address.FromString(contractAddress)
	if err != nil {
		return err
	}
	didContract := c.Contract(didContractAddress, didABI)
	result := didContract.Execute("createDID", d).SetGasLimit(5000000)
	if err := wait.Wait(context.Background(), result); err != nil {
		return err
	}
	result1, err := didContract.Read("isExist", d).Call(context.Background())
	var inte bool
	if err := result1.Unmarshal(&inte); err != nil {
		return err
	}

	fmt.Println(inte)
	return nil
}
