package keeper

import (
	"fmt"
	"strconv"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		// 如果计数器不存在，从 1 开始
		store.Set(types.KeyPrefix(types.RequestIDKey), sdk.Uint64ToBigEndian(1))
		return 1
	}

	// 更新计数器值
	id := sdk.BigEndianToUint64(bz) + 1
	store.Set(types.KeyPrefix(types.RequestIDKey), sdk.Uint64ToBigEndian(id))
	return id
}

func (k Keeper) SetSpectrumRequest(ctx sdk.Context, request types.SpectrumRequest) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.KeyPrefix(types.SpectrumRequestKey + strconv.FormatUint(request.Id, 10))
	bz := k.cdc.MustMarshal(&request)
	err := store.Set(key, bz)
	if err != nil {
		panic(fmt.Sprintf("failed to set SpectrumRequest: %v", err))
	}
}
