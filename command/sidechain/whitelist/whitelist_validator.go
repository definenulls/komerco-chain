package whitelist

import (
	"fmt"
	"time"

	"github.com/definenulls/komerco-chain/command"
	"github.com/definenulls/komerco-chain/command/helper"
	"github.com/definenulls/komerco-chain/command/komerbftsecrets"
	sidechainHelper "github.com/definenulls/komerco-chain/command/sidechain"
	"github.com/definenulls/komerco-chain/consensus/komerbft/contractsapi"
	"github.com/definenulls/komerco-chain/contracts"
	"github.com/definenulls/komerco-chain/txrelayer"
	"github.com/definenulls/komerco-chain/types"
	"github.com/spf13/cobra"
	"github.com/umbracle/ethgo"
)

var params whitelistParams

func GetCommand() *cobra.Command {
	registerCmd := &cobra.Command{
		Use:     "whitelist-validator",
		Short:   "whitelist a new validator",
		PreRunE: runPreRun,
		RunE:    runCommand,
	}

	setFlags(registerCmd)

	return registerCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.accountDir,
		komerbftsecrets.AccountDirFlag,
		"",
		komerbftsecrets.AccountDirFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.accountConfig,
		komerbftsecrets.AccountConfigFlag,
		"",
		komerbftsecrets.AccountConfigFlagDesc,
	)

	cmd.Flags().StringVar(
		&params.newValidatorAddress,
		newValidatorAddressFlag,
		"",
		"account address of a possible validator",
	)

	cmd.MarkFlagsMutuallyExclusive(komerbftsecrets.AccountDirFlag, komerbftsecrets.AccountConfigFlag)
	helper.RegisterJSONRPCFlag(cmd)
}

func runPreRun(cmd *cobra.Command, _ []string) error {
	params.jsonRPC = helper.GetJSONRPCAddress(cmd)

	return params.validateFlags()
}

func runCommand(cmd *cobra.Command, _ []string) error {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	ownerAccount, err := sidechainHelper.GetAccount(params.accountDir, params.accountConfig)
	if err != nil {
		return fmt.Errorf("enlist validator failed: %w", err)
	}

	txRelayer, err := txrelayer.NewTxRelayer(txrelayer.WithIPAddress(params.jsonRPC),
		txrelayer.WithReceiptTimeout(150*time.Millisecond))
	if err != nil {
		return fmt.Errorf("enlist validator failed: %w", err)
	}

	whitelistFn := &contractsapi.AddToWhitelistChildValidatorSetFn{
		WhitelistAddreses: []ethgo.Address{ethgo.Address(types.StringToAddress(params.newValidatorAddress))},
	}

	encoded, err := whitelistFn.EncodeAbi()
	if err != nil {
		return fmt.Errorf("enlist validator failed: %w", err)
	}

	txn := &ethgo.Transaction{
		From:     ownerAccount.Ecdsa.Address(),
		Input:    encoded,
		To:       (*ethgo.Address)(&contracts.ValidatorSetContract),
		GasPrice: sidechainHelper.DefaultGasPrice,
	}

	receipt, err := txRelayer.SendTransaction(txn, ownerAccount.Ecdsa)
	if err != nil {
		return fmt.Errorf("enlist validator failed %w", err)
	}

	if receipt.Status == uint64(types.ReceiptFailed) {
		return fmt.Errorf("enlist validator transaction failed on block %d", receipt.BlockNumber)
	}

	var (
		whitelistEvent contractsapi.AddedToWhitelistEvent
		result         = &enlistResult{}
		foundLog       = false
	)

	for _, log := range receipt.Logs {
		doesMatch, err := whitelistEvent.ParseLog(log)
		if !doesMatch {
			continue
		}

		if err != nil {
			return err
		}

		result.newValidatorAddress = whitelistEvent.Validator.String()
		foundLog = true

		break
	}

	if !foundLog {
		return fmt.Errorf("could not find an appropriate log in receipt that enlistment happened")
	}

	outputter.WriteCommandResult(result)

	return nil
}
