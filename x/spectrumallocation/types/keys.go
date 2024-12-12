package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "spectrumallocation"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_spectrumallocation"

	AllocationIDKey        = "AllocationIDKey"
	SpectrumAllocationsKey = "SpectrumAllocation-"
	ChannelKey             = "Channel-"
)

var (
	ParamsKey = []byte("p_spectrumallocation")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetSpectrumRequestKey generates the key for a SpectrumRequest with the given ID
func GetSpectrumAllocationsKey(id uint64) []byte {
	return append([]byte(SpectrumAllocationsKey), sdk.Uint64ToBigEndian(id)...)
}

// GetChannelKey generates a unique key for a channel using its ID.
func GetChannelKey(channelID int32) []byte {
	return []byte(fmt.Sprintf("%s%d", ChannelKey, channelID))
}
