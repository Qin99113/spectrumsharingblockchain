package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "spectrumSharingBlockchain/testutil/keeper"
	"spectrumSharingBlockchain/x/spectrumrequest/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.SpectrumrequestKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
