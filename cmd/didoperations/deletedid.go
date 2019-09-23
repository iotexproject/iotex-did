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

var DeleteDIDCmd = &cobra.Command{
	Use:   "delete DID",
	Short: "Delete DID for io device",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return deleteDID()
	},
}

func deleteDID() error {
	conn, err := iotex.NewDefaultGRPCConn(IOEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to set up grpc connection")
	}
	defer conn.Close()

	c, err := getAuthedClient(conn, _password)
	if err != nil {
		return errors.Wrap(err, "failed to get authed client")
	}

	caddr, err := address.FromString(ContractAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get contract address")
	}

	ctx := context.Background()
	iotexDIDABI, err := abi.JSON(strings.NewReader(IoTeXDIDABI))
	if err != nil {
		return errors.Wrap(err, "failed to parse IoTeX DID ABI")
	}

	ioCommonAddr, err := ioAddrToEvmAddr(c.Account().Address().String())
	if err != nil {
		return errors.Wrap(err, "failed to convert iotex address to eth common address")
	}
	didString := DIDPrefix + ioCommonAddr.String()
	h, err := c.Contract(caddr, iotexDIDABI).Execute("deleteDID", didString).
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
		return errors.Errorf("deleting IoTeX DID failed: %x", h)
	}
	fmt.Println("Deleted DID:", DIDPrefix+strings.ToLower(ioCommonAddr.String()))
	return nil
}

func init() {
	DeleteDIDCmd.Flags().StringVarP(&_password, "password", "p", "", "password for keystore file")
	if err := DeleteDIDCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err.Error())
	}
}
