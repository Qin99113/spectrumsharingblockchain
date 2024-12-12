package spectrumallocation_test

import (
	"testing"

	keepertest "spectrumSharingBlockchain/testutil/keeper"
	"spectrumSharingBlockchain/testutil/nullify"
	spectrumallocation "spectrumSharingBlockchain/x/spectrumallocation/module"
	"spectrumSharingBlockchain/x/spectrumallocation/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SpectrumallocationKeeper(t)
	spectrumallocation.InitGenesis(ctx, k, genesisState)
	got := spectrumallocation.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
