package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/spn/testutil/sample"
	"github.com/tendermint/spn/x/launch/types"
)

func TestGenesisState_Validate(t *testing.T) {
	var (
		launchID1        = uint64(0)
		launchID2        = uint64(1)
		noExistLaunchID  = uint64(2)
		addr1            = sample.Address(r)
		addr2            = sample.Address(r)
		vestingAddress   = sample.Address(r)
		genesisValidator = sample.GenesisValidator(r, launchID1, addr1)
		genesisChainID   = sample.GenesisChainID(r)

		// Those are samples we can use for each fields when they are not the one to test
		sampleChainList = []types.Chain{
			{
				LaunchID:       launchID1,
				GenesisChainID: genesisChainID,
			},
			{
				LaunchID:       launchID2,
				GenesisChainID: genesisChainID,
			},
		}
		sampleGenesisAccountList = []types.GenesisAccount{
			{
				LaunchID: launchID1,
				Address:  addr1,
			},
			{
				LaunchID: launchID1,
				Address:  addr2,
			},
			{
				LaunchID: launchID2,
				Address:  addr1,
			},
			{
				LaunchID: launchID2,
				Address:  addr2,
			},
		}
		sampleVestingAccountList = []types.VestingAccount{
			{
				LaunchID: launchID1,
				Address:  vestingAddress,
			},
			{
				LaunchID: launchID2,
				Address:  vestingAddress,
			},
		}
		sampleGenesisValidatorList = []types.GenesisValidator{genesisValidator}
		sampleRequestList          = []types.Request{
			{
				LaunchID:  launchID1,
				RequestID: 0,
			},
			{
				LaunchID:  launchID1,
				RequestID: 1,
			},
		}
		sampleRequestCounterList = []types.RequestCounter{
			{
				LaunchID: launchID1,
				Counter:  10,
			},
		}
	)

	for _, tc := range []struct {
		desc          string
		genState      *types.GenesisState
		shouldBeValid bool
	}{
		{
			desc:          "should validate default genesis",
			genState:      types.DefaultGenesis(),
			shouldBeValid: true,
		},
		{
			desc: "should validate valid genesis state",
			genState: &types.GenesisState{
				ChainList:            sampleChainList,
				ChainCounter:         10,
				GenesisAccountList:   sampleGenesisAccountList,
				VestingAccountList:   sampleVestingAccountList,
				GenesisValidatorList: sampleGenesisValidatorList,
				RequestList:          sampleRequestList,
				RequestCounterList:   sampleRequestCounterList,
				Params:               types.DefaultParams(),
				// this line is used by starport scaffolding # types/genesis/validField
			},
			shouldBeValid: true,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
		{
			desc: "should prevent validate genesis with an invalid chain",
			genState: &types.GenesisState{
				ChainList: []types.Chain{
					{
						LaunchID:       launchID1,
						GenesisChainID: "invalid_chain_id",
					},
				},
				ChainCounter: 10,
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with duplicated chains",
			genState: &types.GenesisState{
				ChainList: []types.Chain{
					{
						LaunchID:       launchID1,
						GenesisChainID: genesisChainID,
					},
					{
						LaunchID:       launchID1,
						GenesisChainID: genesisChainID,
					},
				},
				ChainCounter: 10,
				Params:       types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with a chain with a chain id number above the chain counter",
			genState: &types.GenesisState{
				ChainList: []types.Chain{
					{
						LaunchID:       12,
						GenesisChainID: genesisChainID,
					},
				},
				ChainCounter: 10,
				Params:       types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with duplicated accounts",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				GenesisAccountList: []types.GenesisAccount{
					{
						LaunchID: launchID1,
						Address:  addr1,
					},
					{
						LaunchID: launchID1,
						Address:  addr1,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with an account not associated with chain",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				GenesisAccountList: []types.GenesisAccount{
					{
						LaunchID: noExistLaunchID,
						Address:  addr1,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with duplicated vesting accounts",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				VestingAccountList: []types.VestingAccount{
					{
						LaunchID: launchID1,
						Address:  vestingAddress,
					},
					{
						LaunchID: launchID1,
						Address:  vestingAddress,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with a vesting account not associated with chain",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				VestingAccountList: []types.VestingAccount{
					{
						LaunchID: noExistLaunchID,
						Address:  vestingAddress,
					},
				},
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with one address present in a genesis account and vesting account",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				GenesisAccountList: []types.GenesisAccount{
					{
						LaunchID: launchID1,
						Address:  addr1,
					},
				},
				VestingAccountList: []types.VestingAccount{
					{
						LaunchID: launchID1,
						Address:  addr1,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with a genesis validator not associated to a chain",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				GenesisValidatorList: []types.GenesisValidator{
					sample.GenesisValidator(r, noExistLaunchID, addr1),
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with duplicated genesis validator",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				GenesisValidatorList: []types.GenesisValidator{
					sample.GenesisValidator(r, launchID1, addr1),
					sample.GenesisValidator(r, launchID1, addr1),
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with a validator address not associated to a chain",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				GenesisValidatorList: []types.GenesisValidator{
					sample.GenesisValidator(r, noExistLaunchID, addr1),
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with duplicated requests",
			genState: &types.GenesisState{
				ChainList:          sampleChainList,
				ChainCounter:       10,
				RequestCounterList: sampleRequestCounterList,
				RequestList: []types.Request{
					{
						LaunchID:  launchID1,
						RequestID: 0,
					},
					{
						LaunchID:  launchID1,
						RequestID: 0,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with request not associated with chain",
			genState: &types.GenesisState{
				ChainList:          sampleChainList,
				ChainCounter:       10,
				RequestCounterList: sampleRequestCounterList,
				RequestList: []types.Request{
					{
						LaunchID:  noExistLaunchID,
						RequestID: 0,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with request while no request count for the chain",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				RequestCounterList: []types.RequestCounter{
					{
						LaunchID: launchID2,
						Counter:  1,
					},
				},
				RequestList: []types.Request{
					{
						LaunchID:  launchID1,
						RequestID: 0,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with duplicated request counter",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				RequestCounterList: []types.RequestCounter{
					{
						LaunchID: launchID1,
						Counter:  0,
					},
					{
						LaunchID: launchID1,
						Counter:  1,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with a request counter not associated with chain",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				RequestCounterList: []types.RequestCounter{
					{
						LaunchID: noExistLaunchID,
						Counter:  0,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
		{
			desc: "should prevent validate genesis with a request counter below a request id",
			genState: &types.GenesisState{
				ChainList:    sampleChainList,
				ChainCounter: 10,
				RequestCounterList: []types.RequestCounter{
					{
						LaunchID: launchID1,
						Counter:  5,
					},
				},
				RequestList: []types.Request{
					{
						LaunchID:  launchID1,
						RequestID: 10,
					},
				},
				Params: types.DefaultParams(),
			},
			shouldBeValid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if !tc.shouldBeValid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			launchIDMap := make(map[uint64]struct{})
			for _, elem := range tc.genState.ChainList {
				launchIDMap[elem.LaunchID] = struct{}{}
			}

			for _, acc := range tc.genState.RequestList {
				// check if the chain exist for requests
				_, ok := launchIDMap[acc.LaunchID]
				require.True(t, ok)
			}

			for _, acc := range tc.genState.GenesisValidatorList {
				// check if the chain exist for validators
				_, ok := launchIDMap[acc.LaunchID]
				require.True(t, ok)
			}

			for _, acc := range tc.genState.GenesisAccountList {
				// check if the chain exist for genesis accounts
				_, ok := launchIDMap[acc.LaunchID]
				require.True(t, ok)
			}

			for _, acc := range tc.genState.VestingAccountList {
				// check if the chain exist for vesting accounts
				_, ok := launchIDMap[acc.LaunchID]
				require.True(t, ok)
			}
		})
	}
}

func TestGenesisState_ValidateParams(t *testing.T) {
	for _, tc := range []struct {
		desc          string
		genState      types.GenesisState
		shouldBeValid bool
	}{
		{
			desc: "should prevent validate genesis with invalid params",
			genState: types.GenesisState{
				Params: types.NewParams(types.DefaultMinLaunchTime, types.MaxParametrableLaunchTime+1, types.DefaultRevertDelay, types.DefaultChainCreationFee),
			},
			shouldBeValid: false,
		},
		{
			desc: "should validate genesis with valid params",
			genState: types.GenesisState{
				Params: types.NewParams(types.DefaultMinLaunchTime, types.DefaultMaxLaunchTime, types.DefaultRevertDelay, types.DefaultChainCreationFee),
			},
			shouldBeValid: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.shouldBeValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
