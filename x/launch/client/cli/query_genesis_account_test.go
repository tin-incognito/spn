package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/spn/x/launch/client/cli"
	"github.com/tendermint/spn/x/launch/types"
)

func (suite *QueryTestSuite) TestShowGenesisAccount() {
	ctx := suite.Network.Validators[0].ClientCtx
	accs := suite.LaunchState.GenesisAccountList

	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc       string
		idLaunchID string
		idAddress  string

		args []string
		err  error
		obj  types.GenesisAccount
	}{
		{
			desc:       "should show an existing genesis account",
			idLaunchID: strconv.Itoa(int(accs[0].LaunchID)),
			idAddress:  accs[0].Address,

			args: common,
			obj:  accs[0],
		},
		{
			desc:       "should send error for a non existing genesis account",
			idLaunchID: strconv.Itoa(100000),
			idAddress:  strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	} {
		suite.T().Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idLaunchID,
				tc.idAddress,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowGenesisAccount(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetGenesisAccountResponse
				require.NoError(t, suite.Network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.GenesisAccount)
				require.Equal(t, tc.obj, resp.GenesisAccount)
			}
		})
	}
}

func (suite *QueryTestSuite) TestListGenesisAccount() {
	ctx := suite.Network.Validators[0].ClientCtx
	accs := suite.LaunchState.GenesisAccountList

	chainID := accs[0].LaunchID
	request := func(chainID uint64, next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			strconv.Itoa(int(chainID)),
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	suite.T().Run("should allow listing genesis accounts by offset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(accs); i += step {
			args := request(chainID, nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListGenesisAccount(), args)
			require.NoError(t, err)
			var resp types.QueryAllGenesisAccountResponse
			require.NoError(t, suite.Network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.GenesisAccount), step)
			require.Subset(t, accs, resp.GenesisAccount)
		}
	})
	suite.T().Run("should allow listing genesis accounts by key", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(accs); i += step {
			args := request(chainID, next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListGenesisAccount(), args)
			require.NoError(t, err)
			var resp types.QueryAllGenesisAccountResponse
			require.NoError(t, suite.Network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.GenesisAccount), step)
			require.Subset(t, accs, resp.GenesisAccount)
			next = resp.Pagination.NextKey
		}
	})
	suite.T().Run("should allow listing all genesis accounts", func(t *testing.T) {
		args := request(chainID, nil, 0, uint64(len(accs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListGenesisAccount(), args)
		require.NoError(t, err)
		var resp types.QueryAllGenesisAccountResponse
		require.NoError(t, suite.Network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(accs), int(resp.Pagination.Total))
		require.ElementsMatch(t, accs, resp.GenesisAccount)
	})
}
