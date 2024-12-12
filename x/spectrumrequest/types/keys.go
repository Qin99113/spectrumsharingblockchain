package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "spectrumrequest"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_spectrumrequest"
	// RequestIDKey defines the requestID key
	RequestIDKey = "RequestIDKey"

	// SpectrumRequestKeyPrefix defines the prefix for SpectrumRequest keys
	SpectrumRequestKeyPrefix = "SpectrumRequest-"
)

var (
	ParamsKey = []byte("p_spectrumrequest")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetSpectrumRequestKey generates the key for a SpectrumRequest with the given ID
func GetSpectrumRequestKey(id uint64) []byte {
	return append([]byte(SpectrumRequestKeyPrefix), sdk.Uint64ToBigEndian(id)...)
}
