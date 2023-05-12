package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/definenulls/komerco-chain/command/backup"
	"github.com/definenulls/komerco-chain/command/bridge"
	"github.com/definenulls/komerco-chain/command/genesis"
	"github.com/definenulls/komerco-chain/command/helper"
	"github.com/definenulls/komerco-chain/command/ibft"
	"github.com/definenulls/komerco-chain/command/license"
	"github.com/definenulls/komerco-chain/command/monitor"
	"github.com/definenulls/komerco-chain/command/peers"
	"github.com/definenulls/komerco-chain/command/komerbft"
	"github.com/definenulls/komerco-chain/command/komerbftmanifest"
	"github.com/definenulls/komerco-chain/command/komerbftsecrets"
	"github.com/definenulls/komerco-chain/command/regenesis"
	"github.com/definenulls/komerco-chain/command/rootchain"
	"github.com/definenulls/komerco-chain/command/secrets"
	"github.com/definenulls/komerco-chain/command/server"
	"github.com/definenulls/komerco-chain/command/status"
	"github.com/definenulls/komerco-chain/command/txpool"
	"github.com/definenulls/komerco-chain/command/version"
	"github.com/definenulls/komerco-chain/command/whitelist"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "Komerco Edge is a framework for building Ethereum-compatible Blockchain networks",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		version.GetCommand(),
		txpool.GetCommand(),
		status.GetCommand(),
		secrets.GetCommand(),
		peers.GetCommand(),
		rootchain.GetCommand(),
		monitor.GetCommand(),
		ibft.GetCommand(),
		backup.GetCommand(),
		genesis.GetCommand(),
		server.GetCommand(),
		whitelist.GetCommand(),
		license.GetCommand(),
		komerbftsecrets.GetCommand(),
		komerbft.GetCommand(),
		komerbftmanifest.GetCommand(),
		bridge.GetCommand(),
		regenesis.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
