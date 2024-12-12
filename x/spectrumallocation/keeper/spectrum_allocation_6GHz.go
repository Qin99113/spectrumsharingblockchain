package keeper

import (
	"fmt"

	"spectrumSharingBlockchain/x/spectrumallocation/types"
	requests "spectrumSharingBlockchain/x/spectrumrequest/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AllocateChannels(ctx sdk.Context, request requests.SpectrumRequest) ([]types.Channel, error) {
	channels := k.GetAllChannels(ctx)      // Retrieve all available channels
	requiredBandwidth := request.Bandwidth // Total bandwidth required
	allocatedChannels := []types.Channel{} // List of allocated channels
	currentBandwidth := int32(0)           // Track cumulative allocated bandwidth

	// Filter channels based on user type
	filteredChannels := k.filterChannelsByUserType(channels, request.UserType)
	if len(filteredChannels) == 0 {
		return nil, fmt.Errorf("no channels available for user type: %s", request.UserType)
	}

	// Iterate through filtered channels to allocate bandwidth
	for _, channel := range filteredChannels {
		// Check if channel is continuous with the last allocated channel
		if !CheckChannelConfict(allocatedChannels, channel) {
			// Reset allocation if channels are non-continuous
			allocatedChannels = []types.Channel{}
			currentBandwidth = 0
			continue
		}

		// Add channel to allocation
		allocatedChannels = append(allocatedChannels, channel)
		currentBandwidth += channel.Bandwidth

		// Mark channel as Allocated in the store
		channel.ChannelStatus = "Allocated"
		k.SetChannel(ctx, channel)

		// Stop allocation if required bandwidth is satisfied
		if currentBandwidth >= requiredBandwidth {
			return allocatedChannels, nil
		}
	}

	// If bandwidth requirement cannot be fulfilled, return an error
	return nil, fmt.Errorf("insufficient bandwidth to meet request for user type: %s", request.UserType)
}

// Helper function to filter channels based on user type
func (k Keeper) filterChannelsByUserType(channels []types.Channel, userType string) []types.Channel {
	filtered := []types.Channel{}
	for _, channel := range channels {
		switch userType {
		case "SP", "VLP":
			if channel.ChannelStatus == "Available" {
				filtered = append(filtered, channel)
			}
		case "LPI":
			if channel.ChannelStatus == "Available" || channel.ChannelStatus == "Low Power Indoor Only" {
				filtered = append(filtered, channel)
			}
		default:
			continue
		}
	}
	return filtered
}

// isChannelContinuous checks if a new channel is continuous with the last allocated channel.
func CheckChannelConfict(allocatedChannels []types.Channel, newChannel types.Channel) bool {
	// If no channels are allocated yet, the new channel can be added
	if len(allocatedChannels) == 0 {
		return true
	}

	// Check for continuity: last allocated channel's end frequency should match new channel's start
	lastChannel := allocatedChannels[len(allocatedChannels)-1]
	return lastChannel.Frequency+lastChannel.Bandwidth == newChannel.Frequency
}

// CalculatePriority computes the priority for a spectrum request.
func (k Keeper) CalculatePriority(request requests.SpectrumRequest) int32 {
	basePriority := int32(0)

	// User type weight
	switch request.UserType {
	case "SP":
		basePriority += 100 // Standard power users with AFC
	case "LPI":
		basePriority += 90 // Low power indoor users
	case "VLP":
		basePriority += 70 // Very low power users
	default:
		basePriority += 50
	}

	// 2. Organization weight (optional)
	organizationPriority := getOrganizationPriority(request.Organization)
	basePriority += organizationPriority

	// Bandwidth penalty: Higher bandwidth reduces priority
	// Dynamic adjustment could be based on sub-band congestion or policies.
	basePriority -= (request.Bandwidth / 10) * 5

	// Duration penalty: Longer duration reduces priority
	// Capped to prevent excessive reductions for very long requests.
	durationPenalty := (request.Duration / 3600) * 10
	if durationPenalty > 100 { // Cap duration penalty to 100
		durationPenalty = 100
	}
	basePriority -= durationPenalty

	// Bid amount normalization: Higher bids increase priority
	// Cap bid impact to prevent dominance
	bidPriority := int32(request.BidAmount.Amount.Int64() / 1000)
	if bidPriority > 200 { // Cap bid impact to 200 points
		bidPriority = 200
	}
	basePriority += bidPriority

	// Ensure priority is non-negative
	if basePriority < 0 {
		basePriority = 0
	}

	return basePriority
}

// Helper function to get organization priority
func getOrganizationPriority(organization string) int32 {
	switch organization {
	case "Government":
		return 50 // Higher priority for government requests
	case "EmergencyServices":
		return 100 // Top priority for emergency services
	case "Commercial":
		return 30 // Moderate priority for commercial organizations
	case "NonProfit":
		return 20 // Lower priority for non-profits
	default:
		return 10 // Default for unknown organizations
	}
}
