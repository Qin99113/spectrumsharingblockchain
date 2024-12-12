package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateAllocation{}

func NewMsgCreateAllocation(creator string, allocationType string) *MsgCreateAllocation {
	return &MsgCreateAllocation{
		Creator:        creator,
		AllocationType: allocationType,
	}
}

func (msg *MsgCreateAllocation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// Check if allocation type is valid (if needed)
	validTypes := map[string]bool{"Dynamic": true, "Manual": true}
	if _, ok := validTypes[msg.AllocationType]; !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid allocation type: %s", msg.AllocationType)
	}

	return nil

}
