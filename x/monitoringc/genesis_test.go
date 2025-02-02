package monitoringc_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeeper "github.com/tendermint/spn/testutil/keeper"
	"github.com/tendermint/spn/testutil/nullify"
	"github.com/tendermint/spn/x/monitoringc"
	"github.com/tendermint/spn/x/monitoringc/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		VerifiedClientIDList: []types.VerifiedClientID{
			{
				LaunchID:  0,
				ClientIDs: []string{"0"},
			},
			{
				LaunchID:  1,
				ClientIDs: []string{"0"},
			},
		},
		ProviderClientIDList: []types.ProviderClientID{
			{
				LaunchID: 0,
			},
			{
				LaunchID: 1,
			},
		},
		LaunchIDFromVerifiedClientIDList: []types.LaunchIDFromVerifiedClientID{
			{
				ClientID: "0",
			},
			{
				ClientID: "1",
			},
		},
		LaunchIDFromChannelIDList: []types.LaunchIDFromChannelID{
			{
				ChannelID: "0",
			},
			{
				ChannelID: "1",
			},
		},
		MonitoringHistoryList: []types.MonitoringHistory{
			{
				LaunchID: 0,
			},
			{
				LaunchID: 1,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	ctx, tk, _ := testkeeper.NewTestSetup(t)
	monitoringc.InitGenesis(ctx, *tk.MonitoringConsumerKeeper, genesisState)
	got := monitoringc.ExportGenesis(ctx, *tk.MonitoringConsumerKeeper)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	require.ElementsMatch(t, genesisState.VerifiedClientIDList, got.VerifiedClientIDList)
	require.ElementsMatch(t, genesisState.ProviderClientIDList, got.ProviderClientIDList)
	require.ElementsMatch(t, genesisState.LaunchIDFromVerifiedClientIDList, got.LaunchIDFromVerifiedClientIDList)
	require.ElementsMatch(t, genesisState.LaunchIDFromChannelIDList, got.LaunchIDFromChannelIDList)
	require.ElementsMatch(t, genesisState.MonitoringHistoryList, got.MonitoringHistoryList)
	// this line is used by starport scaffolding # genesis/test/assert
}
