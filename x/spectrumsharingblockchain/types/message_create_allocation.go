package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateAllocation{}

func NewMsgCreateAllocation(creator string, requestId uint64, startTime int32, endTime int32, allocationType string) *MsgCreateAllocation {
	return &MsgCreateAllocation{
		Creator:        creator,
		RequestId:      requestId,
		StartTime:      startTime,
		EndTime:        endTime,
		AllocationType: allocationType,
	}
}

func (msg *MsgCreateAllocation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
