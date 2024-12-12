package spectrumallocation

import (
	"fmt"

	"spectrumSharingBlockchain/x/spectrumallocation/keeper"
	"spectrumSharingBlockchain/x/spectrumallocation/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Validate params before setting them
	if err := genState.Params.Validate(); err != nil {
		panic(fmt.Sprintf("invalid module parameters: %v", err))
	}
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(fmt.Sprintf("failed to set params: %v", err))
	}

	// Initialize spectrum channels if not already initialized
	if len(k.GetAllChannels(ctx)) == 0 {
		k.InitializeChannels(ctx)
		k.Logger().Info("Spectrum channels initialized during genesis.")
	} else {
		k.Logger().Info("Spectrum channels already initialized.")
	}

	// Initialize allocation records
	if len(genState.Allocations) > 0 {
		for _, allocation := range genState.Allocations {
			k.SetSpectrumAllocation(ctx, *allocation)

			// Update channel statuses based on allocations
			for _, channel := range allocation.Channels {
				existingChannel, found := k.GetChannel(ctx, channel.Id)
				if !found {
					panic(fmt.Sprintf("channel with ID %d not found during allocation initialization", channel.Id))
				}
				if existingChannel.ChannelStatus != "Allocated" {
					existingChannel.ChannelStatus = "Allocated"
					k.SetChannel(ctx, existingChannel)
				}
			}
		}
		k.Logger().Info(fmt.Sprintf("Initialized %d allocations from genesis state.", len(genState.Allocations)))
	} else {
		k.Logger().Info("No allocations found in genesis state; initialized with an empty state.")
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export spectrum allocations
	allocations := k.GetAllSpectrumAllocations(ctx)
	for _, allocation := range allocations {
		genesis.Allocations = append(genesis.Allocations, &allocation)
	}

	// Export spectrum channels
	channels := k.GetAllChannels(ctx)
	for _, channel := range channels {
		genesis.Channels = append(genesis.Channels, &channel)
	}

	k.Logger().Info(fmt.Sprintf("Exported %d allocations and %d channels in genesis state.", len(allocations), len(channels)))
	return genesis
}
