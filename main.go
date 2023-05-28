package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Usage: %s [RPC-ENDPOINT]\n", args[0])
		return
	}
	rpc := args[1]

	encodingConfig := params.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	simapp.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	simapp.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	cfg := sdk.GetConfig()
	cfg.Seal()

	cl, err := client.NewClientFromNode(rpc)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir("upgradetracker").
		WithViper("").
		WithNodeURI(rpc).
		WithClient(cl)

	qc := stakingtypes.NewQueryClient(ctx)

	pr := query.PageRequest{
		Offset:     0,
		Limit:      150,
		CountTotal: true,
		Reverse:    false,
	}

	qvr := stakingtypes.QueryValidatorsRequest{
		Status:     stakingtypes.BondStatusBonded,
		Pagination: &pr,
	}

	validators, err := qc.Validators(context.Background(), &qvr)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.OpenFile("validators.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, validator := range validators.Validators {
		name := validator.Description.Moniker
		tokens := validator.Tokens

		_, err = datawriter.WriteString(fmt.Sprintf("%s,%d\n", name, tokens.Int64()))
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
	}

	datawriter.Flush()
	file.Close()
}
