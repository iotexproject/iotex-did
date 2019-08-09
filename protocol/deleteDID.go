package protocol

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/iotexproject/iotex-antenna-go/v2/utils/wait"
)

// DeleteDID is processing public key to DID
func DeleteDID(did string) error {

	// we got our DID in d variable
	// Create grpc connection
	didContract, err := connectTodidContract()
	if err != nil {
		return err
	}

	resultURI, err := didContract.Read("getURI", did).Call(context.Background())
	if err != nil {
		return err
	}

	var uri string
	if err := resultURI.Unmarshal(&uri); err != nil {
		return err
	}

	result := didContract.Execute("deleteDID", did).SetGasLimit(4000000)
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
