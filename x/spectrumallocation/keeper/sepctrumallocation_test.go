package keeper_test

import (
	"fmt"
	keepertest "spectrumSharingBlockchain/testutil/keeper"
	"spectrumSharingBlockchain/x/spectrumallocation/keeper"
	"spectrumSharingBlockchain/x/spectrumallocation/types"
	requesttypes "spectrumSharingBlockchain/x/spectrumrequest/types"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupAllocationKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	// 假设 SpectrumallocationKeeper 是用于初始化 Keeper 和 Context 的帮助函数
	k, ctx := keepertest.SpectrumallocationKeeper(t)
	return k, ctx
}

func TestInitializeChannels(t *testing.T) {
	k, ctx := setupAllocationKeeper(t)

	// 初始化频道
	k.InitializeChannels(ctx)

	// 获取所有初始化后的频道
	channels := k.GetAllChannels(ctx)

	// 打印频道数量
	t.Logf("Total channels: %d", len(channels))

	// 打印每个频道的详细信息
	for i, channel := range channels {
		t.Logf("Channel %d: ID=%d, Frequency=%d MHz, Bandwidth=%d MHz, Status=%s",
			i+1, channel.Id, channel.Frequency, channel.Bandwidth, channel.ChannelStatus)
	}

	// 确认频道数量是否正确 (6 GHz 范围: 5925 MHz - 7125 MHz，间隔 20 MHz)
	expectedChannelCount := (7125 - 5925) / 20
	require.Equal(t, expectedChannelCount, len(channels), "Channel count mismatch")

	// 验证每个频道的属性
	for _, channel := range channels {
		// 检查频率是否正确
		require.True(t, channel.Frequency >= 5925 && channel.Frequency < 7125, "Channel frequency out of range")

		// 检查频道带宽是否正确
		require.Equal(t, int32(20), channel.Bandwidth, "Channel bandwidth mismatch")

		// 验证频道状态是否符合逻辑
		switch {
		case (channel.Frequency >= 5925 && channel.Frequency < 6425) || (channel.Frequency >= 6525 && channel.Frequency < 6875):
			require.Equal(t, "Available", channel.ChannelStatus, "Channel status mismatch for AFC range")
			require.ElementsMatch(t, []string{"SP", "LPI", "VLP"}, channel.AllowedUsers, "Allowed users mismatch for AFC range")
		case (channel.Frequency >= 6425 && channel.Frequency < 6525) || (channel.Frequency >= 6875 && channel.Frequency < 7125):
			require.Equal(t, "Low Power Indoor Only", channel.ChannelStatus, "Channel status mismatch for LPI range")
			require.ElementsMatch(t, []string{"LPI"}, channel.AllowedUsers, "Allowed users mismatch for LPI range")
		default:
			require.Equal(t, "Protected", channel.ChannelStatus, "Channel status mismatch for Protected range")
			require.Empty(t, channel.AllowedUsers, "Protected channels should have no allowed users")
		}
	}
}

func TestGetChannel(t *testing.T) {
	k, ctx := setupAllocationKeeper(t)

	// 模拟创建一个 Channel
	channel := types.Channel{
		Id:            1,
		Frequency:     6000,
		Bandwidth:     20,
		ChannelStatus: "Allocated",
	}

	// 保存频道到存储
	k.SetChannel(ctx, channel)

	// 调用 GetChannel 方法
	retrievedChannel, found := k.GetChannel(ctx, channel.Id)
	require.True(t, found, "Channel should be found")
	require.Equal(t, channel.Id, retrievedChannel.Id, "Channel ID should match")
	require.Equal(t, channel.Frequency, retrievedChannel.Frequency, "Channel frequency should match")
	require.Equal(t, channel.Bandwidth, retrievedChannel.Bandwidth, "Channel bandwidth should match")
	require.Equal(t, channel.ChannelStatus, retrievedChannel.ChannelStatus, "Channel status should match")
}
func TestReleaseChannels(t *testing.T) {
	k, ctx := setupAllocationKeeper(t)

	// // 模拟创建一个已分配的 SpectrumAllocation
	// allocation := types.SpectrumAllocation{
	// 	AllocationId: 1,
	// 	Channels: []*types.Channel{
	// 		{Id: 1, Frequency: 6000, Bandwidth: 20, ChannelStatus: "Allocated"},
	// 		{Id: 2, Frequency: 6020, Bandwidth: 20, ChannelStatus: "Allocated"},
	// 	},
	// }
	// 模拟创建一个已分配的 SpectrumAllocation
	allocation := types.SpectrumAllocation{
		AllocationId:   1,
		RequestId:      1,
		Creator:        "cosmos1rkrctacqxv4mcxz4ymvpv30yl26lsuxgzt8kc4",
		Organization:   "Test Organization",
		UserType:       "SP",
		Channels:       []*types.Channel{{Id: 1, Frequency: 6000, Bandwidth: int32(20), ChannelStatus: "Allocated"}},
		Bandwidth:      int32(20),
		StartTime:      ctx.BlockHeader().Time.Unix() - 3600,
		EndTime:        ctx.BlockHeader().Time.Unix() + 1800,
		Priority:       10,
		Status:         "Active",
		AllocationType: "Auto",
	}

	// 将频道保存到存储
	for _, channel := range allocation.Channels {
		k.SetChannel(ctx, *channel)
	}

	// 调用 ReleaseChannels 方法
	err := k.ReleaseChannels(ctx, allocation)
	require.NoError(t, err, "ReleaseChannels should not return an error")

	// 验证每个频道是否更新为 "Available"
	for _, channel := range allocation.Channels {
		retrievedChannel, found := k.GetChannel(ctx, channel.Id)
		require.True(t, found, "Channel should be found after release")
		require.Equal(t, "Available", retrievedChannel.ChannelStatus, "Channel status should be 'Available'")
	}
}

func TestSetandGetAllSpectrumAllocations(t *testing.T) {
	k, ctx := setupAllocationKeeper(t)

	// 创建并存储分配记录
	allocation1 := types.SpectrumAllocation{
		AllocationId: 1,
		RequestId:    101,
		Creator:      "cosmos1...",
		Status:       "Active",
	}
	allocation2 := types.SpectrumAllocation{
		AllocationId: 2,
		RequestId:    102,
		Creator:      "cosmos2...",
		Status:       "Released",
	}

	err := k.SetSpectrumAllocation(ctx, allocation1)
	require.NoError(t, err)
	err = k.SetSpectrumAllocation(ctx, allocation2)
	require.NoError(t, err)

	// 读取所有分配记录
	allocations := k.GetAllSpectrumAllocations(ctx)
	require.Len(t, allocations, 2)
	require.Equal(t, allocations[0].AllocationId, uint64(1))
	require.Equal(t, allocations[1].AllocationId, uint64(2))
}

func TestAutoAllocateRequests(t *testing.T) {
	k, ctx := setupAllocationKeeper(t)

	// 初始化频道
	k.InitializeChannels(ctx)

	// 创建一个 SpectrumRequest
	request := requesttypes.SpectrumRequest{
		Id:           2,
		Creator:      "cosmoslkpeag6fkcmwezgj025ntalynphv9hagwwfngrp",
		Organization: "TestOrg2",
		UserType:     "SP",
		Bandwidth:    40, // 需要 40 MHz
		Duration:     3600,
		BidAmount:    &sdk.Coin{Denom: "token", Amount: math.NewInt(1000)}, // 1000 token bid
		Status:       "Pending",
		RequestTime:  1234567890,
	}

	// 模拟添加到 pending requests
	k.SpectrumRequestKeeper.SetSpectrumRequest(ctx, request)
	// 读取存储中的数据
	retrievedRequest, found := k.SpectrumRequestKeeper.GetSpectrumRequest(ctx, request.Id)
	fmt.Printf("Debug - Retrieved SpectrumRequest: %+v %t\n", retrievedRequest, found)

	// // 自动分配
	// k.AutoAllocateRequests(ctx)

	// // 验证分配结果
	// allocations := k.GetAllSpectrumAllocations(ctx)
	// require.NotEmpty(t, allocations)
	// require.Equal(t, request.Id, allocations[0].RequestId)
	// require.Equal(t, "Active", allocations[0].Status)
}

func TestReleaseAllocation(t *testing.T) {
	k, ctx := setupAllocationKeeper(t)

	// // 模拟创建一个已分配的 SpectrumAllocation
	// allocation := types.SpectrumAllocation{
	// 	AllocationId:   1,
	// 	RequestId:      1,
	// 	Creator:        "cosmos1rkrctacqxv4mcxz4ymvpv30yl26lsuxgzt8kc4",
	// 	Organization:   "Test Organization",
	// 	UserType:       "SP",
	// 	Channels:       []*types.Channel{{Id: 1, Frequency: 6000, Bandwidth: int32(20), ChannelStatus: "Allocated"}},
	// 	Bandwidth:      int32(20),
	// 	StartTime:      ctx.BlockHeader().Time.Unix() - 3600,
	// 	EndTime:        ctx.BlockHeader().Time.Unix() + 1800,
	// 	Priority:       10,
	// 	Status:         "Active",
	// 	AllocationType: "Auto",
	// }
	allocation := types.SpectrumAllocation{
		AllocationId:   1,
		RequestId:      1,
		Creator:        "test_creator",
		Organization:   "test_org",
		UserType:       "SP",
		Channels:       []*types.Channel{{Id: 1, Frequency: 6000, Bandwidth: int32(20), ChannelStatus: "Allocated"}},
		Bandwidth:      int32(20),
		StartTime:      0,
		EndTime:        1000,
		Priority:       10,
		Status:         "Active",
		AllocationType: "Auto",
	}

	// 保存分配
	err := k.SetSpectrumAllocation(ctx, allocation)
	require.NoError(t, err)

	// 调用 ReleaseAllocation
	k.ReleaseAllocation(ctx, allocation)

	// // Verify Allocation Status
	allAllocations := k.GetAllSpectrumAllocations(ctx)
	require.Len(t, allAllocations, 1, "Expected only one allocation")
	require.Equal(t, "Released", allAllocations[0].Status, "Allocation status should be Released")
	fmt.Printf("Stored Allocation: %+v\n", allAllocations)

	// // Verify Channel Status
	// channel, found := k.GetChannel(ctx, allocation.Channels[0].Id)
	// require.True(t, found, "Channel not found")
	// require.Equal(t, "Available", channel.ChannelStatus, "Channel status should be Available")
}

func TestReleaseExpiredAllocations(t *testing.T) {
	k, ctx := setupAllocationKeeper(t)

	// 模拟创建一个已分配的 SpectrumAllocation
	allocation := types.SpectrumAllocation{
		AllocationId:   1,
		RequestId:      1,
		Creator:        "cosmos1rkrctacqxv4mcxz4ymvpv30yl26lsuxgzt8kc4",
		Organization:   "Test Organization",
		UserType:       "SP",
		Channels:       []*types.Channel{{Id: 1, Frequency: 6000, Bandwidth: 20, ChannelStatus: "Allocated"}},
		Bandwidth:      20,
		StartTime:      ctx.BlockHeader().Time.Unix() - 3600,
		EndTime:        ctx.BlockHeader().Time.Unix() - 1800,
		Priority:       10,
		Status:         "Active",
		AllocationType: "Auto",
	}

	k.SetSpectrumAllocation(ctx, allocation)
	allAllocations := k.GetAllSpectrumAllocations(ctx)
	require.Equal(t, "Active", allAllocations[0].Status)

	// 调用释放过期分配的方法
	k.ReleaseExpiredAllocations(ctx)

	// 再次获取分配数据，验证状态是否更新
	allAllocations = k.GetAllSpectrumAllocations(ctx)
	t.Logf("Updated Allocations: %+v", allAllocations)

	// 验证分配是否被释放
	require.Equal(t, "Released", allAllocations[0].Status)
}

func TestReleaseLowPriorityAllocations(t *testing.T) {
	k, ctx := setupAllocationKeeper(t)

	// 创建两个分配，优先级不同
	allocation1 := types.SpectrumAllocation{
		AllocationId: 1,
		RequestId:    1,
		Bandwidth:    20,
		StartTime:    ctx.BlockHeader().Time.Unix(),
		Status:       "Active",
		Priority:     50, // 低优先级
	}

	allocation2 := types.SpectrumAllocation{
		AllocationId: 2,
		RequestId:    2,
		Bandwidth:    40,
		StartTime:    ctx.BlockHeader().Time.Unix(),
		Status:       "Active",
		Priority:     100, // 高优先级
	}

	k.SetSpectrumAllocation(ctx, allocation1)
	k.SetSpectrumAllocation(ctx, allocation2)

	// 模拟一个新请求，高优先级
	request := requesttypes.SpectrumRequest{
		Id:           3,
		Creator:      "cosmos1creator",
		Organization: "TestOrg",
		UserType:     "SP",
		Bandwidth:    60,
	}

	conflicts := k.CheckConflictingAllocations(ctx, request)
	err := k.ReleaseLowPriorityAllocations(ctx, request, conflicts)

	require.NoError(t, err)

	// 验证低优先级分配是否已释放
	allAllocations := k.GetAllSpectrumAllocations(ctx)
	require.Equal(t, "Released", allAllocations[0].Status)
	require.Equal(t, "Active", allAllocations[1].Status)
}
