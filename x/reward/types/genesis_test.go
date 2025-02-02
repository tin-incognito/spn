package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/spn/testutil/sample"
	"github.com/tendermint/spn/x/reward/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default genesis is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				RewardPoolList: []types.RewardPool{
					sample.RewardPool(r, 1),
					sample.RewardPool(r, 2),
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
		{
			desc: "duplicated rewardPool",
			genState: &types.GenesisState{
				RewardPoolList: []types.RewardPool{
					sample.RewardPool(r, 1),
					sample.RewardPool(r, 1),
				},
			},
			valid: false,
		},
		{
			desc: "invalid rewardPool",
			genState: &types.GenesisState{
				RewardPoolList: []types.RewardPool{
					sample.RewardPool(r, 1),
					{}, // invalid reward pool
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
