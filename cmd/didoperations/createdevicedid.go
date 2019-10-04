package didoperations

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var CreateDeviceDIDCmd = &cobra.Command{
	Use:   "create-device",
	Short: "Create DID for device",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return createDeviceDID()
	},
}

func createDeviceDID() error {
	conn, err := iotex.NewDefaultGRPCConn(IOEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to set up grpc connection")
	}
	defer conn.Close()

	c, err := getAuthedClient(conn, _password)
	if err != nil {
		return errors.Wrap(err, "failed to get authed client")
	}

	caddr, err := address.FromString(MockDeviceContractAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get contract address")
	}

	ctx := context.Background()
	mockDeviceABI, err := abi.JSON(strings.NewReader(MockDeviceDIDABI))
	if err != nil {
		return errors.Wrap(err, "failed to parse mock device ABI")
	}

	proof, err := hex.DecodeString(_signature)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex signature to bytes")
	}
	h, err := c.Contract(caddr, mockDeviceABI).Execute("createDID", _uuid, proof, stringToBytes32(_deviceHash), _deviceURI).
		SetGasPrice(GasPrice).SetGasLimit(GasLimit).Call(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to execute createDID function")
	}

	time.Sleep(30 * time.Second)

	resp, err := c.API().GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{
		ActionHash: hex.EncodeToString(h[:]),
	})
	if err != nil {
		return err
	}
	if resp.ReceiptInfo.Receipt.Status != 1 {
		return errors.Errorf("creating device DID failed: %x", h)
	}

	fmt.Println("Created device DID:", MockDeviceDIDPrefix+_uuid)
	return nil
}

var _deviceHash string
var _deviceURI string

func init() {
	CreateDeviceDIDCmd.Flags().StringVarP(&_password, "password", "p", "", "password for keystore file")
	if err := CreateDeviceDIDCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err.Error())
	}
	CreateDeviceDIDCmd.Flags().StringVarP(&_uuid, "uuid", "i", "", "device uuid")
	if err := CreateDeviceDIDCmd.MarkFlagRequired("uuid"); err != nil {
		log.Fatal(err.Error())
	}
	CreateDeviceDIDCmd.Flags().StringVarP(&_signature, "signature", "s", "", "device binding proof")
	if err := CreateDeviceDIDCmd.MarkFlagRequired("signature"); err != nil {
		log.Fatal(err.Error())
	}
	CreateDeviceDIDCmd.Flags().StringVar(&_deviceHash, "hash", "", "device document hash")
	if err := CreateDeviceDIDCmd.MarkFlagRequired("hash"); err != nil {
		log.Fatal(err.Error())
	}
	CreateDeviceDIDCmd.Flags().StringVarP(&_deviceURI, "uri", "u", "", "device document uri")
	if err := CreateDeviceDIDCmd.MarkFlagRequired("uri"); err != nil {
		log.Fatal(err.Error())
	}
}

