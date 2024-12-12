package types

const (
	// ModuleName defines the module name
	ModuleName = "spectrumallocation"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_spectrumallocation"

	AllocationIDKey       = "AllocationIDKey"
	SpectrumAllocationKey = "SpectrumAllocation-"
)

var (
	ParamsKey = []byte("p_spectrumallocation")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
