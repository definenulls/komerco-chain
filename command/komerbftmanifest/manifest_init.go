package komerbftmanifest

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"
	"strings"

	"github.com/definenulls/komerco-chain/command"
	"github.com/definenulls/komerco-chain/command/genesis"
	"github.com/definenulls/komerco-chain/consensus/komerbft"
	"github.com/definenulls/komerco-chain/types"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
)

const (
	manifestPathFlag      = "path"
	premineValidatorsFlag = "premine-validators"
	stakeFlag             = "stake"
	validatorsFlag        = "validators"
	validatorsPathFlag    = "validators-path"
	validatorsPrefixFlag  = "validators-prefix"
	chainIDFlag           = "chain-id"

	defaultValidatorPrefixPath = "test-chain-"
	defaultManifestPath        = "./manifest.json"

	ecdsaAddressLength = 40
	blsKeyLength       = 256
	blsSignatureLength = 128
)

var (
	params = &manifestInitParams{}
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "manifest",
		Short:   "Initializes manifest file. It is applicable only to komerbft consensus protocol.",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(cmd)

	return cmd
}

func runPreRun(_ *cobra.Command, _ []string) error {
	return params.validateFlags()
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.manifestPath,
		manifestPathFlag,
		defaultManifestPath,
		"the file path where manifest file is going to be stored",
	)

	cmd.Flags().StringVar(
		&params.validatorsPath,
		validatorsPathFlag,
		"./",
		"root path containing komerbft validator keys",
	)

	cmd.Flags().StringVar(
		&params.validatorsPrefixPath,
		validatorsPrefixFlag,
		defaultValidatorPrefixPath,
		"folder prefix names for komerbft validator keys",
	)

	cmd.Flags().StringArrayVar(
		&params.validators,
		validatorsFlag,
		[]string{},
		"validators defined by user (format: <P2P multi address>:<ECDSA address>:<public BLS key>:<BLS signature>)",
	)

	cmd.Flags().StringArrayVar(
		&params.premineValidators,
		premineValidatorsFlag,
		[]string{},
		fmt.Sprintf(
			"the premined validators and balances (format: <address>[:<balance>]). Default premined balance: %d",
			command.DefaultPremineBalance,
		),
	)

	cmd.Flags().Int64Var(
		&params.chainID,
		chainIDFlag,
		command.DefaultChainID,
		"the ID of the chain",
	)

	cmd.Flags().StringArrayVar(
		&params.stakes,
		stakeFlag,
		[]string{},
		fmt.Sprintf(
			"validators staked amount (format: <address>[:<amount>]). Default stake amount: %d",
			command.DefaultStake,
		),
	)

	cmd.MarkFlagsMutuallyExclusive(validatorsFlag, validatorsPathFlag)
	cmd.MarkFlagsMutuallyExclusive(validatorsFlag, validatorsPrefixFlag)
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	validators, err := params.getValidatorAccounts()
	if err != nil {
		outputter.SetError(fmt.Errorf("failed to get validator accounts: %w", err))

		return
	}

	manifest := &komerbft.Manifest{GenesisValidators: validators, ChainID: params.chainID}
	if err = manifest.Save(params.manifestPath); err != nil {
		outputter.SetError(fmt.Errorf("failed to save manifest file '%s': %w", params.manifestPath, err))

		return
	}

	outputter.SetCommandResult(params.getResult())
}

type manifestInitParams struct {
	manifestPath         string
	validatorsPath       string
	validatorsPrefixPath string
	premineValidators    []string
	stakes               []string
	validators           []string
	chainID              int64
}

func (p *manifestInitParams) validateFlags() error {
	if _, err := os.Stat(p.validatorsPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("provided validators path '%s' doesn't exist", p.validatorsPath)
	}

	return nil
}

// getValidatorAccounts gathers validator accounts info either from CLI or from provided local storage
func (p *manifestInitParams) getValidatorAccounts() ([]*komerbft.Validator, error) {
	// populate validators premine info
	premineMap := make(map[types.Address]*genesis.PremineInfo, len(p.premineValidators))
	stakeMap := make(map[types.Address]*genesis.PremineInfo, len(p.stakes))

	for _, premine := range p.premineValidators {
		premineInfo, err := genesis.ParsePremineInfo(premine)
		if err != nil {
			return nil, err
		}

		premineMap[premineInfo.Address] = premineInfo
	}

	for _, stake := range p.stakes {
		stakeInfo, err := genesis.ParsePremineInfo(stake)
		if err != nil {
			return nil, fmt.Errorf("invalid stake amount provided '%s' : %w", stake, err)
		}

		stakeMap[stakeInfo.Address] = stakeInfo
	}

	if len(p.validators) > 0 {
		validators := make([]*komerbft.Validator, len(p.validators))
		for i, validator := range p.validators {
			parts := strings.Split(validator, ":")
			if len(parts) != 4 {
				return nil, fmt.Errorf("expected 4 parts provided in the following format "+
					"<P2P multi address:ECDSA address:public BLS key:BLS signature>, but got %d part(s)",
					len(parts))
			}

			if _, err := multiaddr.NewMultiaddr(parts[0]); err != nil {
				return nil, fmt.Errorf("invalid P2P multi address '%s' provided: %w ", parts[0], err)
			}

			trimmedAddress := strings.TrimPrefix(parts[1], "0x")
			if len(trimmedAddress) != ecdsaAddressLength {
				return nil, fmt.Errorf("invalid ECDSA address: %s", parts[1])
			}

			trimmedBLSKey := strings.TrimPrefix(parts[2], "0x")
			if len(trimmedBLSKey) != blsKeyLength {
				return nil, fmt.Errorf("invalid BLS key: %s", parts[2])
			}

			if len(parts[3]) != blsSignatureLength {
				return nil, fmt.Errorf("invalid BLS signature: %s", parts[3])
			}

			addr := types.StringToAddress(trimmedAddress)
			validators[i] = &komerbft.Validator{
				MultiAddr:    parts[0],
				Address:      addr,
				BlsKey:       trimmedBLSKey,
				BlsSignature: parts[3],
				Balance:      getPremineAmount(addr, premineMap, command.DefaultPremineBalance),
				Stake:        getPremineAmount(addr, stakeMap, command.DefaultStake),
			}
		}

		return validators, nil
	}

	validatorsPath := p.validatorsPath
	if validatorsPath == "" {
		validatorsPath = path.Dir(p.manifestPath)
	}

	validators, err := genesis.ReadValidatorsByPrefix(validatorsPath, p.validatorsPrefixPath)
	if err != nil {
		return nil, err
	}

	for _, v := range validators {
		v.Balance = getPremineAmount(v.Address, premineMap, command.DefaultPremineBalance)
		v.Stake = getPremineAmount(v.Address, stakeMap, command.DefaultStake)
	}

	return validators, nil
}

// getPremineAmount retrieves amount from the premine map or if not provided, returns default amount
func getPremineAmount(addr types.Address, premineMap map[types.Address]*genesis.PremineInfo,
	defaultAmount *big.Int) *big.Int {
	if premine, exists := premineMap[addr]; exists {
		return premine.Amount
	}

	return defaultAmount
}

func (p *manifestInitParams) getResult() command.CommandResult {
	return &result{
		message: fmt.Sprintf("Manifest file written to %s\n", p.manifestPath),
	}
}

type result struct {
	message string
}

func (r *result) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[MANIFEST INITIALIZATION SUCCESS]\n")
	buffer.WriteString(r.message)

	return buffer.String()
}
