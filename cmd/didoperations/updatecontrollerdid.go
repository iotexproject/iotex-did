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

var UpdateControllerDIDCmd = &cobra.Command{
	Use:   "update-controller",
	Short: "Update DID for controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return updateControllerDID()
	},
}

func updateControllerDID() error {
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
		return errors.Wrap(err, "failed to parse IoTeX DID ABI")
	}

	ioCommonAddr, err := ioAddrToEvmAddr(c.Account().Address().String())
	if err != nil {
		return errors.Wrap(err, "failed to convert iotex address to eth common address")
	}
	didString := ControllerDIDPrefix + ioCommonAddr.String()
	if _newControllerHash != "" {
		h1, err := c.Contract(caddr, iotexDIDABI).Execute("updateHash", didString, stringToBytes32(_newControllerHash)).
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
		fmt.Println("Updated Hash for controller DID:", ControllerDIDPrefix+strings.ToLower(ioCommonAddr.String()))
	}

	if _newControllerURI != "" {
		h2, err := c.Contract(caddr, iotexDIDABI).Execute("updateURI", didString, _newControllerURI).
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
		fmt.Println("Updated URI for controller DID:", ControllerDIDPrefix+strings.ToLower(ioCommonAddr.String()))
	}
	return nil
}

var _newControllerHash string
var _newControllerURI string

func init() {
	UpdateControllerDIDCmd.Flags().StringVarP(&_password, "password", "p", "", "password for keystore file")
	if err := UpdateControllerDIDCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err.Error())
	}
	UpdateControllerDIDCmd.Flags().StringVar(&_newControllerHash, "hash", "", "updated document hash")
	UpdateControllerDIDCmd.Flags().StringVarP(&_newControllerURI, "uri", "u", "", "updated document uri")
}
