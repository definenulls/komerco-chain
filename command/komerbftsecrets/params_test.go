package komerbftsecrets

import (
	"encoding/hex"
	"os"
	"path"
	"testing"

	bls "github.com/definenulls/komerco-chain/consensus/komerbft/signer"
	"github.com/definenulls/komerco-chain/secrets/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo/wallet"
)

// Test initKeys
func Test_initKeys(t *testing.T) {
	t.Parallel()

	// Creates test directory
	dir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)

	defer os.RemoveAll(dir)

	sm, err := helper.SetupLocalSecretsManager(dir)
	require.NoError(t, err)

	ip := &initParams{
		generatesAccount: false,
		generatesNetwork: false,
		chainID:          1,
	}

	_, err = ip.initKeys(sm)
	require.NoError(t, err)

	assert.False(t, fileExists(path.Join(dir, "consensus/validator.key")))
	assert.False(t, fileExists(path.Join(dir, "consensus/validator-bls.key")))
	assert.False(t, fileExists(path.Join(dir, "libp2p/libp2p.key")))
	assert.False(t, fileExists(path.Join(dir, "consensus/validator.sig")))

	ip.generatesAccount = true
	res, err := ip.initKeys(sm)
	require.NoError(t, err)
	assert.Len(t, res, 3)

	assert.True(t, fileExists(path.Join(dir, "consensus/validator.key")))
	assert.True(t, fileExists(path.Join(dir, "consensus/validator-bls.key")))
	assert.True(t, fileExists(path.Join(dir, "consensus/validator.sig")))
	assert.False(t, fileExists(path.Join(dir, "libp2p/libp2p.key")))

	ip.generatesNetwork = true
	res, err = ip.initKeys(sm)
	require.NoError(t, err)
	assert.Len(t, res, 1)

	assert.True(t, fileExists(path.Join(dir, "libp2p/libp2p.key")))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func Test_getResult(t *testing.T) {
	t.Parallel()

	// Creates test directory
	dir, err := os.MkdirTemp("", "test")
	require.NoError(t, err)

	defer os.RemoveAll(dir)

	sm, err := helper.SetupLocalSecretsManager(dir)
	require.NoError(t, err)

	ip := &initParams{
		generatesAccount: true,
		generatesNetwork: true,
		printPrivateKey:  true,
		chainID:          1,
	}

	_, err = ip.initKeys(sm)
	require.NoError(t, err)

	res, err := ip.getResult(sm, []string{})
	require.NoError(t, err)

	// Test BLS signature
	sir := res.(*SecretsInitResult) //nolint:forcetypeassert
	ds, err := hex.DecodeString(sir.BLSSignature)
	require.NoError(t, err)

	_, err = bls.UnmarshalSignature(ds)
	require.NoError(t, err)

	// Test public key serialization
	privKey, err := hex.DecodeString(sir.PrivateKey)
	require.NoError(t, err)
	k, err := wallet.NewWalletFromPrivKey(privKey)
	require.NoError(t, err)

	pubKey := k.Address().String()
	assert.Equal(t, sir.Address.String(), pubKey)

	// Test BLS public key serialization
	blsPrivKeyRaw, err := hex.DecodeString(sir.BLSPrivateKey)
	require.NoError(t, err)
	blsPrivKey, err := bls.UnmarshalPrivateKey(blsPrivKeyRaw)
	require.NoError(t, err)

	blsPubKey := hex.EncodeToString(blsPrivKey.PublicKey().Marshal())
	assert.Equal(t, sir.BLSPubkey, blsPubKey)
}
