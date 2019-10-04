package didoperations

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var GetControllerDIDCmd = &cobra.Command{
	Use:   "get-controller",
	Short: "Get DID document hash and URI for controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return getControllerDID()
	},
}

func getControllerDID() error {
	conn, err := iotex.NewDefaultGRPCConn(IOEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to set up grpc connection")
	}
	defer conn.Close()

	c, err := getAuthedClient(conn, _password)
	if err != nil {
		return errors.Wrap(err, "failed to get authed client")
	}

	caddr, err := address.FromString(ControllerContractAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get contract address")
	}

	ctx := context.Background()
	iotexDIDABI, err := abi.JSON(strings.NewReader(IoTeXDIDABI))
	if err != nil {
		return errors.Wrap(err, "failed to parse DID ABI")
	}

	data, err := c.Contract(caddr, iotexDIDABI).Read("getHash", _controllerDID).Call(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get DID document hash")
	}
	var hash [32]byte
	if err := data.Unmarshal(&hash); err != nil {
		return errors.Wrap(err, "failed to unmarshal hash")
	}

	data, err = c.Contract(caddr, iotexDIDABI).Read("getURI", _controllerDID).Call(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get DID document URI")
	}
	var uri string
	if err := data.Unmarshal(&uri); err != nil {
		return errors.Wrap(err, "failed to unmarshal uri")
	}

	fmt.Println("Document Hash:", string(hash[:]))
	fmt.Println("Document URI:", uri)
	return nil
}

var _controllerDID string

func init() {
	GetControllerDIDCmd.Flags().StringVarP(&_password, "password", "p", "", "password for keystore file")
	if err := GetControllerDIDCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err.Error())
	}
	GetControllerDIDCmd.Flags().StringVarP(&_controllerDID, "did-string", "d", "", "DID string")
	if err := GetControllerDIDCmd.MarkFlagRequired("did-string"); err != nil {
		log.Fatal(err.Error())
	}
}
