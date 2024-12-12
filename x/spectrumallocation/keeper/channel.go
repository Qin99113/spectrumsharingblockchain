package keeper

import (
	"fmt"

	"spectrumSharingBlockchain/x/spectrumallocation/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetChannel saves a channel object into the store
func (k Keeper) SetChannel(ctx sdk.Context, channel types.Channel) {

	// Open the KVStore for the current context
	store := k.storeService.OpenKVStore(ctx)

	k.CleanInvalidChannels(ctx)

	// Create a key for the channel using its ID
	key := types.GetChannelKey(channel.Id)

	// Marshal the Channel object into binary data for storage
	bz := k.cdc.MustMarshal(&channel)

	// Save the marshaled channel object into the KVStore using the key
	store.Set(key, bz)
}

// GetChannel retrieves a channel object by its ID
func (k Keeper) GetChannel(ctx sdk.Context, id int32) (types.Channel, bool) {
	// Open the KVStore for the current context
	store := k.storeService.OpenKVStore(ctx)

	k.CleanInvalidChannels(ctx)

	// Create the key using the channel ID
	key := types.GetChannelKey(id)

	// Retrieve the binary data stored under the key
	bz, err := store.Get(key)

	// Check if an error occurred or if the data doesn't exist
	if err != nil || bz == nil {
		// Return an empty Channel object and a "not found" flag
		return types.Channel{}, false
	}

	// Initialize a Channel object to hold the unmarshaled data
	var channel types.Channel

	// Unmarshal the binary data into the Channel object
	k.cdc.MustUnmarshal(bz, &channel)

	// Return the Channel object and a "found" flag
	return channel, true
}

// GetAllChannels retrieves all channels from the store
func (k Keeper) GetAllChannels(ctx sdk.Context) []types.Channel {
	// Open the KVStore for the current context
	store := k.storeService.OpenKVStore(ctx)

	k.CleanInvalidChannels(ctx)

	// Retrieve an iterator for all keys and values in the store
	iterator, err := store.Iterator(nil, nil)
	if err != nil {
		// Handle iterator initialization failure, panic as this is critical
		panic(err)
	}
	defer iterator.Close() // Ensure iterator resources are released after use

	// Initialize an empty slice to store the retrieved channels
	channels := []types.Channel{}

	// Iterate over all key-value pairs in the store
	for ; iterator.Valid(); iterator.Next() {
		// Initialize a Channel object to hold unmarshaled data
		var channel types.Channel

		// Unmarshal the binary value into the Channel object
		k.cdc.MustUnmarshal(iterator.Value(), &channel)

		// Append the unmarshaled Channel object to the result slice
		channels = append(channels, channel)
	}

	// Return the list of all retrieved Channel objects
	return channels
}

// ReleaseChannels releases all channels allocated to a specific allocation.
func (k Keeper) ReleaseChannels(ctx sdk.Context, allocation types.SpectrumAllocation) error {
	for _, channel := range allocation.Channels {
		// Retrieve the channel from the store
		channelObj, found := k.GetChannel(ctx, channel.Id)
		if !found {
			return fmt.Errorf("channel with ID %d not found", channel.Id)
		}

		// 创建独立的 channel 副本，避免修改原始引用
		updatedChannel := types.Channel{
			Id:            channelObj.Id,
			Frequency:     channelObj.Frequency,
			Bandwidth:     channelObj.Bandwidth,
			ChannelStatus: "Available",
		}

		// Save the updated channel back to the store
		k.SetChannel(ctx, updatedChannel)
	}
	return nil
}

// CleanInvalidChannels removes invalid channels from the store.
func (k Keeper) CleanInvalidChannels(ctx sdk.Context) {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var channel types.Channel
		k.cdc.MustUnmarshal(iterator.Value(), &channel)

		// Check if the channel is invalid
		if channel.Frequency == 0 || channel.Bandwidth == 0 || channel.ChannelStatus == "" {
			k.Logger().Warn(fmt.Sprintf("Deleting invalid channel: %+v", channel))
			err := store.Delete(iterator.Key())
			if err != nil {
				k.Logger().Error(fmt.Sprintf("Failed to delete invalid channel: %+v", channel))
			}
		}
	}
}
