package didoperations

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var GetDIDCmd = &cobra.Command{
	Use:   "get DID",
	Short: "Get DID document hash and URI for io device",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return getHash()
	},
}

func getHash() error {
	conn, err := iotex.NewDefaultGRPCConn(IOEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to set up grpc connection")
	}
	defer conn.Close()

	c, err := getAuthedClient(conn)
	if err != nil {
		return errors.Wrap(err, "failed to get authed client")
	}

	caddr, err := address.FromString(OperatorContractAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get contract address")
	}

	ctx := context.Background()
	didABI, err := abi.JSON(strings.NewReader(DecentralizedIdentifierABI))
	if err != nil {
		return errors.Wrap(err, "failed to parse DID ABI")
	}

	ioCommonAddr, err := ioAddrToEvmAddr(c.Account().Address().String())
	if err != nil {
		return errors.Wrap(err, "failed to convert iotex address to eth common address")
	}
	didString := DIDPrefix + ioCommonAddr.String()
	data, err := c.Contract(caddr, didABI).Read("getHash", stringToBytes32(Namespace), didString).Call(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get DID document hash")
	}
	var hash [32]byte
	if err := data.Unmarshal(&hash); err != nil {
		return errors.Wrap(err, "failed to unmarshal hash")
	}

	data, err = c.Contract(caddr, didABI).Read("getURI", stringToBytes32(Namespace), didString).Call(ctx)
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
