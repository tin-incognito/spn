package campaign_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	spntypes "github.com/tendermint/spn/pkg/types"
	testkeeper "github.com/tendermint/spn/testutil/keeper"
	"github.com/tendermint/spn/testutil/sample"
	"github.com/tendermint/spn/x/campaign"
)

/*
// We use a genesis template from sample package, therefore this placeholder is not used
// this line is used by starport scaffolding # genesis/test/state
*/

func TestGenesis(t *testing.T) {
	ctx, tk, _ := testkeeper.NewTestSetup(t)
	r := sample.Rand()

	genesisState := sample.CampaignGenesisStateWithAccounts(r)
	campaign.InitGenesis(ctx, *tk.CampaignKeeper, genesisState)
	got := *campaign.ExportGenesis(ctx, *tk.CampaignKeeper)

	require.ElementsMatch(t, genesisState.CampaignChainsList, got.CampaignChainsList)

	require.ElementsMatch(t, genesisState.CampaignList, got.CampaignList)
	require.Equal(t, genesisState.CampaignCounter, got.CampaignCounter)

	require.ElementsMatch(t, genesisState.MainnetAccountList, got.MainnetAccountList)

	require.Equal(t, genesisState.Params, got.Params)

	maxShares := tk.CampaignKeeper.GetTotalShares(ctx)
	require.Equal(t, uint64(spntypes.TotalShareNumber), maxShares)

	// this line is used by starport scaffolding # genesis/test/assert
}
