package protocol

import (
	"context"
	"database/sql"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

//AddPbKey is add key function
func AddPbKey(did string, pbKey string) error {
	// newPbkey := &PbKey{
	// 	PbKeyHex: pbKey,
	// }
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

	db, err := sql.Open("mysql", uri+"?autocommit=true")
	if err != nil {
		return err
	}
	defer db.Close()

	pbSQL, err := db.Prepare("SELECT public_key FROM didDocumentation WHERE DID = ?")
	if err != nil {
		return err
	}
	var pb []byte
	err = pbSQL.QueryRow(did).Scan(&pb)
	if err != nil {
		return err
	}
	var pbArray []PbKey
	err = json.Unmarshal(pb, &pbArray)
	if err != nil {
		return err
	}
	pbArray = append(pbArray, PbKey{
		PbKeyHex: pbKey,
	})
	publicKeyData, err := json.Marshal(pbArray)
	if err != nil {
		return err
	}
	resultSQL, err := db.Prepare("UPDATE didDocumentation SET public_key = ? WHERE DID = ?")
	resultSQL.Exec(publicKeyData, did)
	if err != nil {
		return err
	}

	return nil
}

// RemovePbKey is remove public key from our database
func RemovePbKey(did string, pbKey string) error {
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

	db, err := sql.Open("mysql", uri+"?autocommit=true")
	if err != nil {
		return err
	}
	defer db.Close()

	pbSQL, err := db.Prepare("SELECT public_key FROM didDocumentation WHERE DID = ?")
	if err != nil {
		return err
	}
	var pb []byte
	err = pbSQL.QueryRow(did).Scan(&pb)
	if err != nil {
		return err
	}
	var pbArray []PbKey
	err = json.Unmarshal(pb, &pbArray)
	if err != nil {
		return err
	}
	for index, value := range pbArray {
		if value.PbKeyHex == pbKey {
			pbArray = append(pbArray[:index], pbArray[index+1:]...)
			break
		}
	}
	publicKeyData, err := json.Marshal(pbArray)
	if err != nil {
		return err
	}
	resultSQL, err := db.Prepare("UPDATE didDocumentation SET public_key = ? WHERE DID = ?")
	resultSQL.Exec(publicKeyData, did)
	if err != nil {
		return err
	}

	return nil
}

//ModifyContext is aim to update the context content
func ModifyContext(did string, contextContent string) error {
	return nil
}
