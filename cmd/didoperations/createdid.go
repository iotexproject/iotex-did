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

var CreateDIDCmd = &cobra.Command{
	Use:   "create DID",
	Short: "Create DID for io device",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return createDID()
	},
}

func createDID() error {
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
	docHash := util.MustFetchNonEmptyParam("DOCUMENT_HASH")
	docURI := util.MustFetchNonEmptyParam("DOCUMENT_URI")
	h, err := c.Contract(caddr, iotexDIDABI).Execute("createDID", ioCommonAddr.String(), stringToBytes32(docHash), docURI).
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
		return errors.Errorf("creating IoTeX DID failed: %x", h)
	}
	return nil
}
