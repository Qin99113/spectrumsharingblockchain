package spectrumrequest_test

import (
	"testing"

	keepertest "spectrumSharingBlockchain/testutil/keeper"
	"spectrumSharingBlockchain/testutil/nullify"
	spectrumrequest "spectrumSharingBlockchain/x/spectrumrequest/module"
	"spectrumSharingBlockchain/x/spectrumrequest/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SpectrumrequestKeeper(t)
	spectrumrequest.InitGenesis(ctx, k, genesisState)
	got := spectrumrequest.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
