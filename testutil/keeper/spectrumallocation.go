package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"spectrumSharingBlockchain/x/spectrumallocation/keeper"
	"spectrumSharingBlockchain/x/spectrumallocation/types"
	spectrumrequestmodulekeeper "spectrumSharingBlockchain/x/spectrumrequest/keeper"
)

func SpectrumallocationKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	// storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	// db := dbm.NewMemDB()
	// stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	// stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	// require.NoError(t, stateStore.LoadLatestVersion())

	// registry := codectypes.NewInterfaceRegistry()
	// cdc := codec.NewProtoCodec(registry)
	// authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	// SpectrumRequestKeeper := spectrumrequestmodulekeeper.NewKeeper(
	// 	cdc,
	// 	runtime.NewKVStoreService(storetypes.NewKVStoreKey("spectrumrequest")),
	// 	log.NewNopLogger(),
	// 	authority.String(),
	// )

	// k := keeper.NewKeeper(
	// 	cdc,
	// 	runtime.NewKVStoreService(storeKey),
	// 	log.NewNopLogger(),
	// 	SpectrumRequestKeeper,
	// 	authority.String(),
	// )

	// ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// // Initialize params
	// if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
	// 	panic(err)
	// }

	// return k, ctx
	// 定义存储键
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	spectrumRequestStoreKey := storetypes.NewKVStoreKey("spectrumrequest") // 添加 spectrumrequest 存储键

	// 创建内存数据库
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

	// 挂载存储
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(spectrumRequestStoreKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	// 初始化编码器
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	// 设置权限
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// 初始化 SpectrumRequestKeeper
	spectrumRequestKeeper := spectrumrequestmodulekeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(spectrumRequestStoreKey), // 使用挂载的存储
		log.NewNopLogger(),
		authority,
	)

	// 初始化 SpectrumAllocationKeeper
	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		spectrumRequestKeeper,
		authority,
	)

	// 创建上下文
	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// 初始化模块参数
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, ctx
}
