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

var GetDeviceDIDCmd = &cobra.Command{
	Use:   "get-device",
	Short: "Get DID document hash and URI for device",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return getDeviceDID()
	},
}

func getDeviceDID() error {
	conn, err := iotex.NewDefaultGRPCConn(IOEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to set up grpc connection")
	}
	defer conn.Close()

	c, err := getAuthedClient(conn, _password)
	if err != nil {
		return errors.Wrap(err, "failed to get authed client")
	}

	caddr, err := address.FromString(DeviceOperatorContractAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get contract address")
	}

	ctx := context.Background()
	didABI, err := abi.JSON(strings.NewReader(DeviceDecentralizedIdentifierABI))
	if err != nil {
		return errors.Wrap(err, "failed to parse DID ABI")
	}

	data, err := c.Contract(caddr, didABI).Read("getHash", stringToBytes32(_namespace), _deviceDID).Call(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get DID document hash")
	}
	var hash [32]byte
	if err := data.Unmarshal(&hash); err != nil {
		return errors.Wrap(err, "failed to unmarshal hash")
	}

	data, err = c.Contract(caddr, didABI).Read("getURI", stringToBytes32(_namespace), _deviceDID).Call(ctx)
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

var _namespace string
var _deviceDID string

func init() {
	GetDeviceDIDCmd.Flags().StringVarP(&_password, "password", "p", "", "password for keystore file")
	if err := GetDeviceDIDCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err.Error())
	}
	GetDeviceDIDCmd.Flags().StringVarP(&_namespace, "namespace", "n", "", "device DID namespace")
	if err := GetDeviceDIDCmd.MarkFlagRequired("namespace"); err != nil {
		log.Fatal(err.Error())
	}
	GetDeviceDIDCmd.Flags().StringVarP(&_deviceDID, "did-string", "d", "", "DID string")
	if err := GetDeviceDIDCmd.MarkFlagRequired("did-string"); err != nil {
		log.Fatal(err.Error())
	}
}