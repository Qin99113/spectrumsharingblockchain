package keeper

import (
	"fmt"
	"sort"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"spectrumSharingBlockchain/x/spectrumallocation/types"
	requests "spectrumSharingBlockchain/x/spectrumrequest/types"
)

type (
	Keeper struct {
		cdc                   codec.BinaryCodec
		storeService          store.KVStoreService
		logger                log.Logger
		SpectrumRequestKeeper types.SpectrumrequestKeeper
		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	SpectrumRequestKeeper types.SpectrumrequestKeeper,
	authority string,

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:                   cdc,
		storeService:          storeService,
		authority:             authority,
		logger:                logger,
		SpectrumRequestKeeper: SpectrumRequestKeeper,
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

func (k Keeper) SetSpectrumAllocation(ctx sdk.Context, allocation types.SpectrumAllocation) error {
	store := k.storeService.OpenKVStore(ctx)

	// 清理无效的分配记录
	removedCount := k.CleanInvalidAllocations(ctx)
	if removedCount > 0 {
		k.Logger().Info(fmt.Sprintf("Cleaned up %d invalid allocations before setting new allocation", removedCount))
	}
	key := types.GetSpectrumAllocationsKey(allocation.AllocationId)
	bz := k.cdc.MustMarshal(&allocation)

	err := store.Set(key, bz)
	if err != nil {
		panic(fmt.Sprintf("Failed to set SpectrumAllocation: %v", err))
	}
	return nil
}

func (k Keeper) GetNextAllocationID(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.KeyPrefix(types.AllocationIDKey))
	if err != nil || bz == nil {
		store.Set(types.KeyPrefix(types.AllocationIDKey), sdk.Uint64ToBigEndian(1))
		return 1
	}

	id := sdk.BigEndianToUint64(bz) + 1
	store.Set(types.KeyPrefix(types.AllocationIDKey), sdk.Uint64ToBigEndian(id))
	return id
}

func (k Keeper) GetSpectrumAllocation(ctx sdk.Context, id uint64) types.SpectrumAllocation {
	// 打开 KVStore
	store := k.storeService.OpenKVStore(ctx)

	// 生成存储 Key
	key := types.GetSpectrumAllocationsKey(id)

	// 从 KVStore 中获取数据
	bz, err := store.Get(key)
	if err != nil {
		// 错误处理，例如记录日志或返回默认值
		k.logger.Error(fmt.Sprintf("failed to get SpectrumRequest with ID %d: %v", id, err))
		return types.SpectrumAllocation{}
	}

	if bz == nil {
		// 数据为空，表示没有找到对应的请求
		return types.SpectrumAllocation{}
	}

	var request types.SpectrumAllocation
	if err := k.cdc.Unmarshal(bz, &request); err != nil {
		k.logger.Error(fmt.Sprintf("failed to unmarshal SpectrumRequest with ID %d: %v", id, err))
		return types.SpectrumAllocation{}
	}

	// 成功返回
	return request
}

// GetAllSpectrumAllocations retrieves all spectrum allocation records from the store.
func (k Keeper) GetAllSpectrumAllocations(ctx sdk.Context) []types.SpectrumAllocation {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(types.KeyPrefix(types.SpectrumAllocationsKey), nil) // Create an iterator to scan all keys
	if err != nil {
		panic(fmt.Sprintf("failed to create store iterator: %s", err))
	}
	defer iterator.Close()

	allocations := []types.SpectrumAllocation{}
	for ; iterator.Valid(); iterator.Next() {
		var allocation types.SpectrumAllocation
		k.cdc.MustUnmarshal(iterator.Value(), &allocation) // Unmarshal the allocation from the store
		allocations = append(allocations, allocation)
	}
	return allocations
}

func (k Keeper) CleanInvalidAllocations(ctx sdk.Context) int {
	store := k.storeService.OpenKVStore(ctx)
	prefix := types.KeyPrefix(types.SpectrumAllocationsKey)

	iterator, err := store.Iterator(prefix, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create store iterator: %s", err))
	}
	defer iterator.Close()

	removedCount := 0

	for ; iterator.Valid(); iterator.Next() {
		var allocation types.SpectrumAllocation
		bz := iterator.Value()

		// 尝试解码
		err := k.cdc.Unmarshal(bz, &allocation)
		if err != nil {
			// 如果解码失败，删除该记录
			store.Delete(iterator.Key())
			removedCount++
			k.Logger().Warn(fmt.Sprintf("Removed invalid allocation due to unmarshaling error: Key=%x", iterator.Key()))
			continue
		}

		// 检查记录是否无效
		if allocation.AllocationId == 0 || allocation.RequestId == 0 || allocation.Bandwidth == 0 || allocation.Status == "" {
			store.Delete(iterator.Key())
			removedCount++
			k.Logger().Warn(fmt.Sprintf("Removed invalid allocation: %+v", allocation))
		}
	}

	k.Logger().Info(fmt.Sprintf("Cleaned up %d invalid allocations", removedCount))
	return removedCount
}

func (k Keeper) InitializeChannels(ctx sdk.Context) {
	channels := []types.Channel{}
	// 6 GHz band from 5925 MHz to 7125 MHz
	bandwidth := 20   // Bandwidth in MHz
	startFreq := 5925 // Starting frequency in MHz
	endFreq := 7125   // Ending frequency in MHz

	// Initialize channels in the 6 GHz band
	channelId := 0
	for freq := startFreq; freq < endFreq; freq += bandwidth {

		var status string
		var allowedUsers []string

		// Ensure the channel frequency is valid
		if freq+bandwidth > endFreq {
			break
		}
		switch {
		// U-NII-5 (5.925–6.425 GHz) and U-NII-7 (6.525–6.875 GHz)
		// These require AFC for SP users, and VLP users have no restrictions
		case (freq >= 5925 && freq < 6425) || (freq >= 6525 && freq < 6875):
			status = "Available"
			allowedUsers = []string{"SP", "LPI", "VLP"}
		// U-NII-6 (6.425–6.525 GHz) and U-NII-8 (6.875–7.125 GHz)
		// These are restricted to low-power indoor users only
		case (freq >= 6425 && freq < 6525) || (freq >= 6875 && freq < 7125):
			status = "Low Power Indoor Only" // LPI users only
			allowedUsers = []string{"LPI"}   // Only LPI users
		// Any undefined or protected frequency
		default:
			status = "Protected"      // Mark as protected or unavailable
			allowedUsers = []string{} // No users allowed
		}

		// Add channel to the list
		channels = append(channels, types.Channel{
			Id:            int32(channelId), // 确保字段名与生成的 Go 类型一致
			Frequency:     int32(freq),      // 中心频率
			Bandwidth:     int32(bandwidth), // 频道带宽
			ChannelStatus: status,           // 状态
			AllowedUsers:  allowedUsers,     // 允许的用户类型
		})
		channelId++

	}

	// 保存初始化的频道
	for _, channel := range channels {

		k.SetChannel(ctx, channel)

	}
	k.CleanInvalidChannels(ctx)

	k.Logger().Info(fmt.Sprintf("Initialized %d channels for 6 GHz band", len(channels)))
}

// AutoAllocateRequests scans all pending requests and attempts to allocate resources for them.
func (k Keeper) AutoAllocateRequests(ctx sdk.Context) {
	// Retrieve all pending spectrum requests and sort them by priority
	pendingRequests := k.SpectrumRequestKeeper.GetPendingRequests(ctx)
	sort.SliceStable(pendingRequests, func(i, j int) bool {
		pi, pj := k.CalculatePriority(pendingRequests[i]), k.CalculatePriority(pendingRequests[j])
		if pi == pj {
			return pendingRequests[i].RequestTime < pendingRequests[j].RequestTime // 二级排序按提交时间
		}
		return pi > pj
	})

	for _, request := range pendingRequests {
		allocationConflicts := k.CheckConflictingAllocations(ctx, request)
		// Release lower-priority allocations if necessary
		// Handle conflicts by releasing low-priority allocations
		err := k.ReleaseLowPriorityAllocations(ctx, request, allocationConflicts)
		if err != nil {
			k.Logger().Error(fmt.Sprintf("Request ID: %d skipped due to insufficient bandwidth", request.Id))
			continue
		}
		// Allocate channels for the current request
		allocatedChannels, err := k.AllocateChannels(ctx, request)
		if err != nil {
			allocationID := k.GetNextAllocationID(ctx)
			allocation := types.SpectrumAllocation{
				AllocationId:   allocationID,
				RequestId:      request.Id,
				Organization:   request.Organization,
				Creator:        request.Creator,
				UserType:       request.UserType,
				Channels:       nil, // No channels allocated
				Bandwidth:      request.Bandwidth,
				StartTime:      ctx.BlockHeader().Time.Unix(),
				EndTime:        ctx.BlockHeader().Time.Unix() + int64(request.Duration),
				Priority:       0,        // Priority is irrelevant when allocation fails
				Status:         "Failed", // Status indicates allocation failure
				AllocationType: "Auto",
			}
			request.Status = "Failed"
			k.SetSpectrumAllocation(ctx, allocation)
			// Emit event for failed allocation
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"SpectrumAllocation",
					sdk.NewAttribute("RequestID", fmt.Sprintf("%d", request.Id)),
					sdk.NewAttribute("Creator", request.Creator),
					sdk.NewAttribute("Status", "Failed"),
				),
			)
			// Remove the processed request from pending
			// k.SpectrumRequestKeeper.RemovePendingRequest(ctx, allocation.RequestId)

			k.Logger().Error(fmt.Sprintf("Failed to allocate for Request ID: %d, Reason: %s", request.Id, err.Error()))
			continue
		}

		// Generate unique allocation ID
		allocationID := k.GetNextAllocationID(ctx)

		// Calculate priority
		priority := k.CalculatePriority(request)

		// Convert allocated channels to pointers
		allocatedChannelPointers := make([]*types.Channel, len(allocatedChannels))
		for i := range allocatedChannels {
			allocatedChannelPointers[i] = &allocatedChannels[i]
		}

		// Create allocation record
		allocation := types.SpectrumAllocation{
			AllocationId:   allocationID,
			RequestId:      request.Id,
			Organization:   request.Organization,
			Creator:        request.Creator,
			UserType:       request.UserType,
			Channels:       allocatedChannelPointers,
			Bandwidth:      request.Bandwidth,
			StartTime:      ctx.BlockHeader().Time.Unix(),
			EndTime:        ctx.BlockHeader().Time.Unix() + int64(request.Duration),
			Priority:       priority,
			Status:         "Active",
			AllocationType: "Auto",
		}

		// Save allocation record
		k.SetSpectrumAllocation(ctx, allocation)

		// Update request status to "Allocated"
		request.Status = "Allocated"
		// k.spectrumRequestKeeper.SetSpectrumRequest(ctx, request)
		err = k.SpectrumRequestKeeper.SetSpectrumRequest(ctx, request)
		if err != nil {
			k.Logger().Error(fmt.Sprintf("Failed to update SpectrumRequest ID: %d, Reason: %v", request.Id, err))
		}
		// Emit event for successful allocation
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"SpectrumAllocation",
				sdk.NewAttribute("RequestID", fmt.Sprintf("%d", request.Id)),
				sdk.NewAttribute("Creator", request.Creator),
				sdk.NewAttribute("Status", "Allocated"),
				sdk.NewAttribute("AllocationID", fmt.Sprintf("%d", allocationID)),
			),
		)

		// k.SpectrumRequestKeeper.RemovePendingRequest(ctx, allocation.RequestId)
		k.Logger().Info(fmt.Sprintf("Successfully allocated Request ID: %d with Allocation ID: %d", request.Id, allocationID))
	}
}

// func DeepCopyAllocation(allocation types.SpectrumAllocation) types.SpectrumAllocation {
// 	// 深拷贝 Channels
// 	newChannels := make([]*types.Channel, len(allocation.Channels))
// 	for i, ch := range allocation.Channels {
// 		copy := *ch
// 		newChannels[i] = &copy
// 	}

// 	// 返回深拷贝的对象
// 	return types.SpectrumAllocation{
// 		AllocationId:   allocation.AllocationId,
// 		RequestId:      allocation.RequestId,
// 		Creator:        allocation.Creator,
// 		Organization:   allocation.Organization,
// 		UserType:       allocation.UserType,
// 		Channels:       newChannels,
// 		Bandwidth:      allocation.Bandwidth,
// 		StartTime:      allocation.StartTime,
// 		EndTime:        allocation.EndTime,
// 		Priority:       allocation.Priority,
// 		Status:         allocation.Status,
// 		AllocationType: allocation.AllocationType,
// 	}
// }

func (k Keeper) ReleaseAllocation(ctx sdk.Context, allocation types.SpectrumAllocation) {
	// Update allocation directly and save
	allocation = k.GetSpectrumAllocation(ctx, allocation.AllocationId)
	allocation.Status = "Released"
	// 返回深拷贝的对象
	allocation = types.SpectrumAllocation{
		AllocationId:   allocation.AllocationId,
		RequestId:      allocation.RequestId,
		Creator:        allocation.Creator,
		Organization:   allocation.Organization,
		UserType:       allocation.UserType,
		Channels:       allocation.Channels,
		Bandwidth:      allocation.Bandwidth,
		StartTime:      allocation.StartTime,
		EndTime:        allocation.EndTime,
		Priority:       allocation.Priority,
		Status:         "Released",
		AllocationType: allocation.AllocationType,
	}

	k.SetSpectrumAllocation(ctx, allocation)
	// // Release channels
	// Channels := make([]*types.Channel, len(allocation.Channels))
	// for i, channel := range allocation.Channels {
	// 	Channel := &types.Channel{
	// 		Id:            channel.Id,
	// 		Frequency:     channel.Frequency,
	// 		Bandwidth:     channel.Bandwidth,
	// 		ChannelStatus: channel.ChannelStatus,
	// 	}
	// 	Channels[i] = Channel
	// }
	// allocation.Channels = Channels

	// err := k.ReleaseChannels(ctx, allocation)
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to release channels: %v", err))
	// }

}

func (k Keeper) ReleaseExpiredAllocations(ctx sdk.Context) {
	allocations := k.GetAllSpectrumAllocations(ctx) // Retrieve all allocations
	currentTime := ctx.BlockHeader().Time.Unix()
	for _, allocation := range allocations {
		// Debug log
		k.Logger().Info(fmt.Sprintf("Checking Allocation ID: %d, EndTime: %d, CurrentTime: %d", allocation.AllocationId, allocation.EndTime, currentTime))

		// Check if the allocation has expired
		if allocation.EndTime <= currentTime && allocation.Status == "Active" {
			k.ReleaseAllocation(ctx, allocation)
			k.Logger().Info(fmt.Sprintf("Released Allocation ID: %d", allocation.AllocationId))
		} else {
			k.Logger().Info(fmt.Sprintf("Allocation ID: %d is not expired or inactive", allocation.AllocationId))
		}
	}
}

// checkConflictingAllocations identifies existing allocations conflicting with the requested bandwidth
func (k Keeper) CheckConflictingAllocations(ctx sdk.Context, request requests.SpectrumRequest) []types.SpectrumAllocation {
	setSpectrumAllocations := k.GetAllSpectrumAllocations(ctx)
	conflictingAllocations := []types.SpectrumAllocation{}
	requiredBandwidth := request.Bandwidth

	for _, allocation := range setSpectrumAllocations {
		// Skip non-active allocations
		if allocation.Status != "Active" {
			continue
		}

		// Check if the allocation's frequency overlaps with the requested bandwidth
		if allocation.Bandwidth >= requiredBandwidth {
			conflictingAllocations = append(conflictingAllocations, allocation)
		}
	}
	return conflictingAllocations
}

// GetProtectionWindow returns the protection window duration based on the user type.
func (k Keeper) GetProtectionWindow(userType string) int64 {
	switch userType {
	case "SP": // Automated Frequency Coordination
		return 3600 // 1 hour in seconds
	case "LPI": // Low-Power Indoor
		return 1800 // 30 minutes in seconds
	case "VLP": // Standard Power
		return 2400 // 40 minutes in seconds
	default:
		return 600 // 10 minutes in seconds for low-priority users
	}
}

// ReleaseLowPriorityAllocations releases lower-priority allocations for high-priority requests.
func (k Keeper) ReleaseLowPriorityAllocations(ctx sdk.Context, request requests.SpectrumRequest, conflictingAllocations []types.SpectrumAllocation) error {
	requiredBandwidth := request.Bandwidth
	releasedBandwidth := int32(0)
	currentTime := ctx.BlockHeader().Time.Unix()

	// Calculate priority of the incoming request
	requestPriority := k.CalculatePriority(request)

	// Iterate through conflicting allocations (already sorted by priority)
	for _, conflict_allocation := range conflictingAllocations {
		// If request priority is not greater, stop further checks
		if requestPriority <= conflict_allocation.Priority {
			k.Logger().Info(fmt.Sprintf(
				"Skipped remaining allocations due to higher or equal priority. Request ID: %d, Priority: %d",
				request.Id, requestPriority,
			))
			break
		}

		// Skip if the allocation is within its protection window
		if currentTime < conflict_allocation.StartTime+k.GetProtectionWindow(conflict_allocation.UserType) {
			k.Logger().Info(fmt.Sprintf(
				"Skipped Allocation ID: %d due to protection window", conflict_allocation.AllocationId,
			))
			continue
		}

		// Release the allocation
		k.ReleaseAllocation(ctx, conflict_allocation)
		releasedBandwidth += conflict_allocation.Bandwidth

		k.Logger().Info(fmt.Sprintf(
			"Released Allocation ID: %d (Priority: %d) for Request ID: %d (Priority: %d)",
			conflict_allocation.AllocationId, conflict_allocation.Priority, request.Id, requestPriority,
		))

		// Stop if enough bandwidth has been released
		if releasedBandwidth >= requiredBandwidth {
			return nil
		}
	}

	// If we exit the loop without sufficient bandwidth released, return an error
	return fmt.Errorf("unable to release sufficient bandwidth for Request ID: %d", request.Id)
}
