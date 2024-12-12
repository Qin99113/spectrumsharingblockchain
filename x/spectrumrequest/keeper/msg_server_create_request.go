package keeper

import (
	"context"
	"fmt"

	"spectrumSharingBlockchain/x/spectrumrequest/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateRequest(goCtx context.Context, msg *types.MsgCreateRequest) (*types.MsgCreateRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	// _ = ctx

	requestID := k.GetNextRequestID(ctx)

	request := types.SpectrumRequest{
		Id:           requestID,
		Creator:      msg.Creator,
		Organization: msg.Organization,
		UserType:     msg.UserType,
		Bandwidth:    msg.Bandwidth,
		Duration:     msg.Duration,
		BidAmount:    msg.BidAmount,
		Status:       types.StatusPending,
		RequestTime:  msg.RequestTime,
	}

	// Attempt to save the SpectrumRequest
	err := k.SetSpectrumRequest(ctx, request)
	if err != nil {
		// Log the error and return a failure response
		k.Logger().Error(fmt.Sprintf("Failed to create SpectrumRequest with ID: %d, error: %v", requestID, err))
		return &types.MsgCreateRequestResponse{
			Status:  "failure",
			Message: fmt.Sprintf("Failed to create request: %v", err),
		}, err
	}

	// Log success and return a success response
	k.Logger().Info(fmt.Sprintf("Request created successfully with ID: %d", requestID))
	return &types.MsgCreateRequestResponse{
		Status:  "success",
		Message: fmt.Sprintf("Request created successfully with ID: %d", requestID),
	}, nil
}
