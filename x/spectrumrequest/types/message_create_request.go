package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateRequest{}

const (
	StatusPending  = "Pending"
	StatusApproved = "Approved"
	StatusRejected = "Rejected"
)

func NewMsgCreateRequest(id uint64, creator string, frequency int32, bandwidth int32, duration int32, bidAmount *sdk.Coin) *MsgCreateRequest {
	return &MsgCreateRequest{
		Id:          id,
		Creator:     creator,
		Frequency:   frequency,
		Bandwidth:   bandwidth,
		Duration:    duration,
		BidAmount:   bidAmount,
		Status:      StatusPending, // default status is "Pending"
		RequestTime: time.Now().Unix(),
	}
}

func (msg *MsgCreateRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// 验证 Frequency、Bandwidth、Duration 为正值
	if msg.Frequency <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "frequency must be positive")
	}
	if msg.Bandwidth <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bandwidth must be positive")
	}
	if msg.Duration <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duration must be positive")
	}

	// 验证 BidAmount 是否为有效的 Coin 类型且不为零
	if msg.BidAmount == nil || !msg.BidAmount.IsValid() || msg.BidAmount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid bid amount")
	}
	return nil
}
