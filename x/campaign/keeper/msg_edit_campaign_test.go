package keeper_test

import (
	"testing"

	testkeeper "github.com/tendermint/spn/testutil/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/spn/testutil/sample"
	"github.com/tendermint/spn/x/campaign/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
)

func TestMsgUpdateCampaignName(t *testing.T) {
	var (
		coordAddr           = sample.Address(r)
		coordAddrNoCampaign = sample.Address(r)
		campaign            = sample.Campaign(r, 0)

		sdkCtx, tk, ts = testkeeper.NewTestSetup(t)
		ctx            = sdk.WrapSDKContext(sdkCtx)
	)
	res, err := ts.ProfileSrv.CreateCoordinator(ctx, &profiletypes.MsgCreateCoordinator{
		Address:     coordAddr,
		Description: sample.CoordinatorDescription(r),
	})
	require.NoError(t, err)
	campaign.CoordinatorID = res.CoordinatorID
	campaign.CampaignID = tk.CampaignKeeper.AppendCampaign(sdkCtx, campaign)

	res, err = ts.ProfileSrv.CreateCoordinator(ctx, &profiletypes.MsgCreateCoordinator{
		Address:     coordAddrNoCampaign,
		Description: sample.CoordinatorDescription(r),
	})
	require.NoError(t, err)

	for _, tc := range []struct {
		name string
		msg  types.MsgEditCampaign
		err  error
	}{
		{
			name: "invalid campaign id",
			msg: types.MsgEditCampaign{
				Coordinator: coordAddr,
				CampaignID:  100,
				Name:        sample.CampaignName(r),
				Metadata:    sample.Metadata(r, 20),
			},
			err: types.ErrCampaignNotFound,
		},
		{
			name: "invalid coordinator address",
			msg: types.MsgEditCampaign{
				Coordinator: sample.Address(r),
				CampaignID:  campaign.CampaignID,
				Name:        sample.CampaignName(r),
				Metadata:    sample.Metadata(r, 20),
			},
			err: profiletypes.ErrCoordAddressNotFound,
		},
		{
			name: "wrong coordinator id",
			msg: types.MsgEditCampaign{
				Coordinator: coordAddrNoCampaign,
				CampaignID:  campaign.CampaignID,
				Name:        sample.CampaignName(r),
				Metadata:    sample.Metadata(r, 20),
			},
			err: profiletypes.ErrCoordInvalid,
		},
		{
			name: "valid transaction - both modified",
			msg: types.MsgEditCampaign{
				Coordinator: coordAddr,
				CampaignID:  campaign.CampaignID,
				Name:        sample.CampaignName(r),
				Metadata:    sample.Metadata(r, 20),
			},
		},
		{
			name: "valid transaction - unmodified metadata",
			msg: types.MsgEditCampaign{
				Coordinator: coordAddr,
				CampaignID:  campaign.CampaignID,
				Name:        sample.CampaignName(r),
				Metadata:    []byte{},
			},
		},
		{
			name: "valid transaction - unmodified name",
			msg: types.MsgEditCampaign{
				Coordinator: coordAddr,
				CampaignID:  campaign.CampaignID,
				Name:        "",
				Metadata:    sample.Metadata(r, 20),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			previousCampaign, found := tk.CampaignKeeper.GetCampaign(sdkCtx, tc.msg.CampaignID)
			_, err := ts.CampaignSrv.EditCampaign(ctx, &tc.msg)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			campaign, found := tk.CampaignKeeper.GetCampaign(sdkCtx, tc.msg.CampaignID)
			require.True(t, found)

			if len(tc.msg.Name) > 0 {
				require.EqualValues(t, tc.msg.Name, campaign.CampaignName)
			} else {
				require.EqualValues(t, previousCampaign.CampaignName, campaign.CampaignName)
			}

			if len(tc.msg.Metadata) > 0 {
				require.EqualValues(t, tc.msg.Metadata, campaign.Metadata)
			} else {
				require.EqualValues(t, previousCampaign.Metadata, campaign.Metadata)
			}
		})
	}
}
