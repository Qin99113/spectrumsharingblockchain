package spectrumsharingblockchain_test

import (
	"testing"

	keepertest "spectrumSharingBlockchain/testutil/keeper"
	"spectrumSharingBlockchain/testutil/nullify"
	spectrumsharingblockchain "spectrumSharingBlockchain/x/spectrumsharingblockchain/module"
	"spectrumSharingBlockchain/x/spectrumsharingblockchain/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SpectrumsharingblockchainKeeper(t)
	spectrumsharingblockchain.InitGenesis(ctx, k, genesisState)
	got := spectrumsharingblockchain.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
