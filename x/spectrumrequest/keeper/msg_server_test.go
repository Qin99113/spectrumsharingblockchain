package keeper_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "spectrumSharingBlockchain/testutil/keeper"
	"spectrumSharingBlockchain/x/spectrumrequest/keeper"
	"spectrumSharingBlockchain/x/spectrumrequest/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, context.Context) {
	k, ctx := keepertest.SpectrumrequestKeeper(t)
	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, stdCtx := setupMsgServer(t) // `stdCtx` 是标准的 context.Context
	require.NotNil(t, ms)
	require.NotNil(t, stdCtx)
	require.NotEmpty(t, k)

	// 构建 MsgCreateRequest
	msg := &types.MsgCreateRequest{
		Creator:      "cosmos1rkrctacqxv4mcxz4ymvpv30yl26lsuxgzt8kc4", // alice's address
		Organization: "Test Organization",
		UserType:     "SP",
		Bandwidth:    50,
		Duration:     10,
		BidAmount:    &sdk.Coin{Denom: "token", Amount: math.NewInt(1000)}, // 1000 token bid
		RequestTime:  1234567890,
	}

	// 将标准的 context.Context 转换为 sdk.Context
	ctx := sdk.UnwrapSDKContext(stdCtx)

	// 执行 CreateRequest
	resp, err := ms.CreateRequest(stdCtx, msg) // 使用标准上下文
	require.NoError(t, err)                    // 确保没有错误
	require.Equal(t, "success", resp.Status)
	require.Contains(t, resp.Message, "Request created successfully")

	// 打印响应
	fmt.Printf("Response: %+v\n", resp)

	// 验证存储的请求
	request, found := k.GetSpectrumRequest(ctx, 1) // 使用转换后的 sdk.Context
	require.True(t, found)
	require.Equal(t, msg.Creator, request.Creator)
	require.Equal(t, msg.Organization, request.Organization)
	require.Equal(t, msg.Bandwidth, request.Bandwidth)
	require.Equal(t, msg.Duration, request.Duration)
	require.Equal(t, msg.BidAmount.String(), request.BidAmount.String())
	require.Equal(t, "Pending", request.Status)
	// 打印验证结果
	fmt.Printf("Request: %+v\n", request)
}
