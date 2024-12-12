package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	spectrumallocation "spectrumSharingBlockchain/x/spectrumallocation/types"
	"spectrumSharingBlockchain/x/spectrumrequest/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetNextRequestID(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.KeyPrefix(types.RequestIDKey)) // 获取当前存储的计数器值
	if err != nil {
		// 处理错误，例如在日志中记录或抛出 panic
		panic(fmt.Sprintf("failed to get RequestIDKey: %v", err))
	}

	if bz == nil {
		k.logger.Info("Initializing RequestIDKey counter to 1.")
		store.Set(types.KeyPrefix(types.RequestIDKey), sdk.Uint64ToBigEndian(1))
		return 1
	}

	// 更新计数器值
	id := sdk.BigEndianToUint64(bz) + 1
	store.Set(types.KeyPrefix(types.RequestIDKey), sdk.Uint64ToBigEndian(id))
	return id
}

func (k Keeper) SetSpectrumRequest(ctx sdk.Context, request types.SpectrumRequest) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetSpectrumRequestKey(request.Id)
	bz := k.cdc.MustMarshal(&request)

	err := store.Set(key, bz)
	if err != nil {

		return fmt.Errorf("failed to set SpectrumRequest with ID %d: %w", request.Id, err)
	}
	return nil
}

func (k Keeper) GetSpectrumRequest(ctx sdk.Context, id uint64) (types.SpectrumRequest, bool) {
	// 打开 KVStore
	store := k.storeService.OpenKVStore(ctx)

	// 生成存储 Key
	key := types.GetSpectrumRequestKey(id)

	// 从 KVStore 中获取数据
	bz, err := store.Get(key)
	if err != nil {
		// 错误处理，例如记录日志或返回默认值
		k.logger.Error(fmt.Sprintf("failed to get SpectrumRequest with ID %d: %v", id, err))
		return types.SpectrumRequest{}, false
	}

	if bz == nil {
		// 数据为空，表示没有找到对应的请求
		return types.SpectrumRequest{}, false
	}

	// 反序列化为 SpectrumRequest
	var request types.SpectrumRequest
	if err := k.cdc.Unmarshal(bz, &request); err != nil {
		k.logger.Error(fmt.Sprintf("failed to unmarshal SpectrumRequest with ID %d: %v", id, err))
		return types.SpectrumRequest{}, false
	}

	// 成功返回
	return request, true
}
func (k Keeper) GetPendingRequests(ctx sdk.Context) []types.SpectrumRequest {
	store := k.storeService.OpenKVStore(ctx)  // Open the KVStore
	iterator, err := store.Iterator(nil, nil) // Iterate through all keys
	if err != nil {
		panic(err) // Handle iterator errors
	}
	defer iterator.Close()

	pendingRequests := []types.SpectrumRequest{}

	// Iterate over all requests in the store
	for ; iterator.Valid(); iterator.Next() {
		var request types.SpectrumRequest
		k.cdc.MustUnmarshal(iterator.Value(), &request) // Unmarshal the request object

		if request.Status == "Pending" { // Check if the request status is "Pending"
			pendingRequests = append(pendingRequests, request)
		}
	}

	return pendingRequests
}

// RemovePendingRequest removes a pending SpectrumRequest from the store by its ID.
func (k Keeper) RemovePendingRequest(ctx sdk.Context, requestID uint64) {
	store := k.storeService.OpenKVStore(ctx)

	requestKey := types.GetSpectrumRequestKey(requestID)

	if bz, err := store.Get(requestKey); err != nil || bz == nil {
		k.Logger().Info(fmt.Sprintf("Attempted to delete non-existent request with ID: %d", requestID))
		return
	}

	if err := store.Delete(requestKey); err != nil {
		k.logger.Error(fmt.Sprintf("Failed to delete spectrum request with ID: %d, error: %v", requestID, err))
		return
	}

	k.Logger().Info(fmt.Sprintf("Successfully removed pending spectrum request with ID: %d", requestID))
}

var _ spectrumallocation.SpectrumrequestKeeper = Keeper{}
