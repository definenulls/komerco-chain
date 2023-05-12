package server

import (
	"github.com/definenulls/komerco-chain/chain"
	"github.com/definenulls/komerco-chain/consensus"
	consensusDev "github.com/definenulls/komerco-chain/consensus/dev"
	consensusDummy "github.com/definenulls/komerco-chain/consensus/dummy"
	consensusIBFT "github.com/definenulls/komerco-chain/consensus/ibft"
	consensusKomerBFT "github.com/definenulls/komerco-chain/consensus/komerbft"
	"github.com/definenulls/komerco-chain/secrets"
	"github.com/definenulls/komerco-chain/secrets/awsssm"
	"github.com/definenulls/komerco-chain/secrets/gcpssm"
	"github.com/definenulls/komerco-chain/secrets/hashicorpvault"
	"github.com/definenulls/komerco-chain/secrets/local"
	"github.com/definenulls/komerco-chain/state"
)

type GenesisFactoryHook func(config *chain.Chain, engineName string) func(*state.Transition) error

type ConsensusType string

const (
	DevConsensus     ConsensusType = "dev"
	IBFTConsensus    ConsensusType = "ibft"
	KomerBFTConsensus ConsensusType = "komerbft"
	DummyConsensus   ConsensusType = "dummy"
)

var consensusBackends = map[ConsensusType]consensus.Factory{
	DevConsensus:     consensusDev.Factory,
	IBFTConsensus:    consensusIBFT.Factory,
	KomerBFTConsensus: consensusKomerBFT.Factory,
	DummyConsensus:   consensusDummy.Factory,
}

// secretsManagerBackends defines the SecretManager factories for different
// secret management solutions
var secretsManagerBackends = map[secrets.SecretsManagerType]secrets.SecretsManagerFactory{
	secrets.Local:          local.SecretsManagerFactory,
	secrets.HashicorpVault: hashicorpvault.SecretsManagerFactory,
	secrets.AWSSSM:         awsssm.SecretsManagerFactory,
	secrets.GCPSSM:         gcpssm.SecretsManagerFactory,
}

var genesisCreationFactory = map[ConsensusType]GenesisFactoryHook{
	KomerBFTConsensus: consensusKomerBFT.GenesisPostHookFactory,
}

func ConsensusSupported(value string) bool {
	_, ok := consensusBackends[ConsensusType(value)]

	return ok
}
