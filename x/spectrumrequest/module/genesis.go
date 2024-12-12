package spectrumrequest

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"spectrumSharingBlockchain/x/spectrumrequest/keeper"
	"spectrumSharingBlockchain/x/spectrumrequest/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// if err := k.SetParams(ctx, genState.Params); err != nil {
	// 	panic(err)
	// }
	// Initialize pending requests
	for _, request := range genState.PendingRequests {
		k.SetSpectrumRequest(ctx, *request)
	}

	// Set module parameters
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	k.Logger().Info(fmt.Sprintf("Initialized %d pending requests from genesis state.", len(genState.PendingRequests)))
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	// genesis := types.DefaultGenesis()
	// genesis.Params = k.GetParams(ctx)

	// // this line is used by starport scaffolding # genesis/module/export

	// return genesis
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export pending requests
	pendingRequests := k.GetPendingRequests(ctx)
	genesis.PendingRequests = make([]*types.SpectrumRequest, len(pendingRequests))
	for i, request := range pendingRequests {
		genesis.PendingRequests[i] = &request
	}

	k.Logger().Info(fmt.Sprintf("Exported %d pending requests to genesis state.", len(genesis.PendingRequests)))
	return genesis
}
