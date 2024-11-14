package keeper

import (
	"spectrumSharingBlockchain/x/spectrumrequest/types"
)

var _ types.QueryServer = Keeper{}
