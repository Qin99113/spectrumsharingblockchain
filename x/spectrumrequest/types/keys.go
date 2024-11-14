package types

const (
	// ModuleName defines the module name
	ModuleName = "spectrumrequest"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_spectrumrequest"
)

var (
	ParamsKey = []byte("p_spectrumrequest")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}