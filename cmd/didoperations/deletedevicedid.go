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

var DeleteDeviceDIDCmd = &cobra.Command{
	Use:   "delete-device",
	Short: "Delete DID for device",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return deleteDeviceDID()
	},
}

func deleteDeviceDID() error {
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

	h, err := c.Contract(caddr, iotexDIDABI).Execute("deleteDID", _uuid, proof).
		SetGasPrice(GasPrice).SetGasLimit(GasLimit).Call(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to execute deleteDID function")
	}

	time.Sleep(30 * time.Second)

	resp, err := c.API().GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{
		ActionHash: hex.EncodeToString(h[:]),
	})
	if err != nil {
		return err
	}
	if resp.ReceiptInfo.Receipt.Status != 1 {
		return errors.Errorf("deleting device DID failed: %x", h)
	}
	fmt.Println("Deleted device DID:", MockDeviceDIDPrefix+_uuid)
	return nil
}

func init() {
	DeleteDeviceDIDCmd.Flags().StringVarP(&_password, "password", "p", "", "password for keystore file")
	if err := DeleteDeviceDIDCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err.Error())
	}
	DeleteDeviceDIDCmd.Flags().StringVarP(&_uuid, "uuid", "i", "", "device uuid")
	if err := DeleteDeviceDIDCmd.MarkFlagRequired("uuid"); err != nil {
		log.Fatal(err.Error())
	}
	DeleteDeviceDIDCmd.Flags().StringVarP(&_signature, "signature", "s", "", "device binding proof")
	if err := DeleteDeviceDIDCmd.MarkFlagRequired("signature"); err != nil {
		log.Fatal(err.Error())
	}
}

