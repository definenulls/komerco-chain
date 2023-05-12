package ibft

import (
	"github.com/definenulls/komerco-chain/command/helper"
	"github.com/definenulls/komerco-chain/command/ibft/candidates"
	"github.com/definenulls/komerco-chain/command/ibft/propose"
	"github.com/definenulls/komerco-chain/command/ibft/quorum"
	"github.com/definenulls/komerco-chain/command/ibft/snapshot"
	"github.com/definenulls/komerco-chain/command/ibft/status"
	_switch "github.com/definenulls/komerco-chain/command/ibft/switch"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	ibftCmd := &cobra.Command{
		Use:   "ibft",
		Short: "Top level IBFT command for interacting with the IBFT consensus. Only accepts subcommands.",
	}

	helper.RegisterGRPCAddressFlag(ibftCmd)

	registerSubcommands(ibftCmd)

	return ibftCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// ibft status
		status.GetCommand(),
		// ibft snapshot
		snapshot.GetCommand(),
		// ibft propose
		propose.GetCommand(),
		// ibft candidates
		candidates.GetCommand(),
		// ibft switch
		_switch.GetCommand(),
		// ibft quorum
		quorum.GetCommand(),
	)
}
