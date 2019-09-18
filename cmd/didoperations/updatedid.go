package didoperations

import (
	"context"
	"encoding/hex"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/iotexproject/iotex-DID/util"
)

var UpdateDIDCmd = &cobra.Command{
	Use:   "update DID",
	Short: "Update DID for io device",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return updateDID()
	},
}

func updateDID() error {
	conn, err := iotex.NewDefaultGRPCConn(IOEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to set up grpc connection")
	}
	defer conn.Close()

	c, err := getAuthedClient(conn)
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
	docHash := util.MustFetchNonEmptyParam("DOCUMENT_HASH")
	h1, err := c.Contract(caddr, iotexDIDABI).Execute("updateHash", didString, stringToBytes32(docHash)).
		SetGasPrice(GasPrice).SetGasLimit(GasLimit).Call(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to execute updateHash function")
	}

	docURI := util.MustFetchNonEmptyParam("DOCUMENT_URI")
	h2, err := c.Contract(caddr, iotexDIDABI).Execute("updateURI", didString, docURI).
		SetGasPrice(GasPrice).SetGasLimit(GasLimit).Call(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to execute updateURI function")
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

	resp, err = c.API().GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{
		ActionHash: hex.EncodeToString(h2[:]),
	})
	if err != nil {
		return err
	}
	if resp.ReceiptInfo.Receipt.Status != 1 {
		return errors.Errorf("updating uri failed: %x", h2)
	}
	return nil
}
