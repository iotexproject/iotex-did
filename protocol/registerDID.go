package protocol

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/go-sql-driver/mysql"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-antenna-go/v2/utils/wait"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

const (
	host = "api.testnet.iotex.one:443"
)

const databaseName = "testDID"
const tableName = "didDocumentation"
const contractAddress = "io1vwzr8lh44fx0t0ac29jxvf8d7y0ft2gpa9a087"
const testURIContract = "io1l2gl0p5d2yxk8a6fnhjurcjnehp2zdxsts8k68"

const constDID = `[
	{
		"constant": true,
		"inputs": [
			{
				"name": "value",
				"type": "bytes32"
			}
		],
		"name": "getHexString",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "pure",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "did",
				"type": "string"
			}
		],
		"name": "deleteDID",
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
		"constant": false,
		"inputs": [
			{
				"name": "inputHash",
				"type": "string"
			},
			{
				"name": "userType",
				"type": "uint16"
			},
			{
				"name": "meta",
				"type": "bytes"
			}
		],
		"name": "createDID",
		"outputs": [
			{
				"name": "resultDID",
				"type": "string"
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
				"name": "didString",
				"type": "string"
			}
		],
		"name": "getURI",
		"outputs": [
			{
				"name": "uri",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "userType",
				"type": "uint16"
			},
			{
				"name": "addr",
				"type": "address"
			}
		],
		"name": "addType",
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
		"constant": false,
		"inputs": [
			{
				"name": "inputHash",
				"type": "string"
			}
		],
		"name": "getDID",
		"outputs": [
			{
				"name": "did",
				"type": "string"
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

type PbKey struct {
	PbKeyHex string
}

// IoAddrToEvmAddr converts IoTeX address into evm address
func IoAddrToEvmAddr(ioAddr string) (common.Address, error) {
	address, err := address.FromString(ioAddr)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(address.Bytes()), nil
}

func connectTodidContract() (iotex.Contract, error) {
	conn, err := iotex.NewDefaultGRPCConn(host)
	if err != nil {
		log.Fatal(err)
	}

	// Add account by private key
	acc, err := account.HexStringToAccount("16668401a3f686717b457bcde6c182cb46dffce81b3bff9ea23f6bd41ac4d54a")
	if err != nil {
		log.Fatal(err)
	}
	// create client
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	didABI, err := abi.JSON(strings.NewReader(constDID))
	if err != nil {
		return nil, err
	}
	didContractAddress, err := address.FromString(contractAddress)
	if err != nil {
		return nil, err
	}
	didContract := c.Contract(didContractAddress, didABI)
	return didContract, nil
}

// CreateDIDByPbkey is processing public key to DID
func CreateDIDByPbkey(pbKey string, contextContent string) (string, error) {
	// hash the public key

	// we got our DID in d variable
	// Create grpc connection

	didContract, err := connectTodidContract()
	if err != nil {
		return "", err
	}

	result := didContract.Execute("createDID", pbKey, uint16(0), []byte("")).SetGasLimit(4000000)
	if err := wait.Wait(context.Background(), result); err != nil {
		return "", err
	}

	uriAddress, err := IoAddrToEvmAddr(testURIContract)
	if err != nil {
		return "", err
	}

	addType := didContract.Execute("addType", uint16(0), uriAddress).SetGasLimit(4000000)
	if err := wait.Wait(context.Background(), addType); err != nil {
		return "", err
	}

	resultRead, err := didContract.Read("getDID", pbKey).Call(context.Background())
	if err != nil {

		return "", err
	}
	var did string
	if err := resultRead.Unmarshal(&did); err != nil {

		return "", err
	}
	resultURI, err := didContract.Read("getURI", did).Call(context.Background())
	if err != nil {
		return "", err
	}

	var uri string
	if err := resultURI.Unmarshal(&uri); err != nil {
		return "", err
	}

	var pbArray []*PbKey
	pbArray = append(pbArray, &PbKey{PbKeyHex: pbKey})
	publicKeyData, err := json.Marshal(pbArray)
	if err != nil {
		return "", err
	}

	db, err := sql.Open("mysql", uri+"?autocommit=true")
	if err != nil {
		return "", err
	}
	defer db.Close()

	resultSQL, err := db.Prepare("INSERT INTO didDocumentation (DID,context,public_key) VALUES (?,?,?)")
	resultSQL.Exec(did, contextContent, publicKeyData)
	if err != nil {
		return "", err
	}

	return did, nil
}
