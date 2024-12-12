package spectrumallocation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"spectrumSharingBlockchain/testutil/sample"
	spectrumallocationsimulation "spectrumSharingBlockchain/x/spectrumallocation/simulation"
	"spectrumSharingBlockchain/x/spectrumallocation/types"
)

// avoid unused import issue
var (
	_ = spectrumallocationsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateAllocation = "op_weight_msg_create_allocation"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateAllocation int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	spectrumallocationGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&spectrumallocationGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateAllocation int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateAllocation, &weightMsgCreateAllocation, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAllocation = defaultWeightMsgCreateAllocation
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAllocation,
		spectrumallocationsimulation.SimulateMsgCreateAllocation(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateAllocation,
			defaultWeightMsgCreateAllocation,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				spectrumallocationsimulation.SimulateMsgCreateAllocation(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
