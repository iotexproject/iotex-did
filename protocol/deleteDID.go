package protocol

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	_ "github.com/go-sql-driver/mysql"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-antenna-go/v2/utils/wait"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

// CreateDIDByPbkey is processing public key to DID
func deleteDID(did string, meta []byte) error {

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

	result := didContract.Execute("deleteDID", did, meta).SetGasLimit(4000000)
	if err := wait.Wait(context.Background(), result); err != nil {
		return err
	}

	db, err := sql.Open("mysql", uri+"?autocommit=true")
	if err != nil {
		return err
	}
	defer db.Close()
	resultSQL, err := db.Prepare("DELETE FROM didDocumentation WHERE DID = ?")
	resultSQL.Exec(did)
	if err != nil {
		return err
	}

	return nil
}
