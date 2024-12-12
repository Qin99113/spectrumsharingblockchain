package types

import (
	"context"

	"spectrumSharingBlockchain/x/spectrumrequest/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// SpectrumRequestKeeper defines the expected interface for spectrumrequest Keeper
type SpectrumrequestKeeper interface {
	//
	GetSpectrumRequest(ctx sdk.Context, id uint64) (types.SpectrumRequest, bool)
	//
	SetSpectrumRequest(ctx sdk.Context, request types.SpectrumRequest) error
	// 提供筛选 Pending 状态的请求方法
	GetPendingRequests(ctx sdk.Context) []types.SpectrumRequest
	//
	RemovePendingRequest(ctx sdk.Context, requestID uint64)
}
