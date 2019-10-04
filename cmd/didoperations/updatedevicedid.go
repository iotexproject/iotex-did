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

var UpdateDeviceDIDCmd = &cobra.Command{
	Use:   "update-device",
	Short: "Update DID for device",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return updateDeviceDID()
	},
}

func updateDeviceDID() error {
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
	iotexDIDABI, err := abi.JSON(strings.NewReader(MockDeviceDIDABI))
	if err != nil {
		return errors.Wrap(err, "failed to parse IoTeX DID ABI")
	}

	proof, err := hex.DecodeString(_signature)
	if err != nil {
		return errors.Wrap(err, "failed to decode hex signature to bytes")
	}

	if _newDeviceHash != "" {
		h1, err := c.Contract(caddr, iotexDIDABI).Execute("updateHash", _uuid, proof, stringToBytes32(_newDeviceHash)).
			SetGasPrice(GasPrice).SetGasLimit(GasLimit).Call(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to execute updateHash function")
		}
		time.Sleep(30 * time.Second)
		resp, err := c.API().GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{
			ActionHash: hex.EncodeToString(h1[:]),
		})
		if err != nil {
			return err
		}
		if resp.ReceiptInfo.Receipt.Status != 1 {
			return errors.Errorf("updating hash failed: %x", h1)
		}
		fmt.Println("Updated Hash for device DID:", MockDeviceDIDPrefix+_uuid)
	}

	if _newDeviceURI != "" {
		h2, err := c.Contract(caddr, iotexDIDABI).Execute("updateURI", _uuid, proof, _newDeviceURI).
			SetGasPrice(GasPrice).SetGasLimit(GasLimit).Call(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to execute updateURI function")
		}
		time.Sleep(30 * time.Second)
		resp, err := c.API().GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{
			ActionHash: hex.EncodeToString(h2[:]),
		})
		if err != nil {
			return err
		}
		if resp.ReceiptInfo.Receipt.Status != 1 {
			return errors.Errorf("updating uri failed: %x", h2)
		}
		fmt.Println("Updated URI for device DID:", MockDeviceDIDPrefix+_uuid)
	}
	return nil
}

var _newDeviceHash string
var _newDeviceURI string

func init() {
	UpdateDeviceDIDCmd.Flags().StringVarP(&_password, "password", "p", "", "password for keystore file")
	if err := UpdateDeviceDIDCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err.Error())
	}
	UpdateDeviceDIDCmd.Flags().StringVarP(&_uuid, "uuid", "i", "", "device uuid")
	if err := UpdateDeviceDIDCmd.MarkFlagRequired("uuid"); err != nil {
		log.Fatal(err.Error())
	}
	UpdateDeviceDIDCmd.Flags().StringVarP(&_signature, "signature", "s", "", "device binding proof")
	if err := UpdateDeviceDIDCmd.MarkFlagRequired("signature"); err != nil {
		log.Fatal(err.Error())
	}
	UpdateDeviceDIDCmd.Flags().StringVar(&_newDeviceHash, "hash", "", "updated document hash")
	UpdateDeviceDIDCmd.Flags().StringVarP(&_newDeviceURI, "uri", "u", "", "updated document uri")
}