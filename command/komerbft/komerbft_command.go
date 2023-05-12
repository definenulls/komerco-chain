package komerbft

import (
	"github.com/definenulls/komerco-chain/command/sidechain/registration"
	"github.com/definenulls/komerco-chain/command/sidechain/staking"
	"github.com/definenulls/komerco-chain/command/sidechain/unstaking"
	"github.com/definenulls/komerco-chain/command/sidechain/validators"

	"github.com/definenulls/komerco-chain/command/sidechain/whitelist"
	"github.com/definenulls/komerco-chain/command/sidechain/withdraw"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	komerbftCmd := &cobra.Command{
		Use:   "komerbft",
		Short: "Komerbft command",
	}

	komerbftCmd.AddCommand(
		staking.GetCommand(),
		unstaking.GetCommand(),
		withdraw.GetCommand(),
		validators.GetCommand(),
		whitelist.GetCommand(),
		registration.GetCommand(),
	)

	return komerbftCmd
}
