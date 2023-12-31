package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proposaltypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/percosis-labs/fury/app"
	"github.com/percosis-labs/fury/x/committee"
	committeekeeper "github.com/percosis-labs/fury/x/committee/keeper"
	committeetypes "github.com/percosis-labs/fury/x/committee/types"
	furydisttypes "github.com/percosis-labs/fury/x/furydist/types"
	"github.com/percosis-labs/fury/x/incentive/keeper"
	"github.com/percosis-labs/fury/x/incentive/testutil"
	"github.com/percosis-labs/fury/x/incentive/types"
	"github.com/percosis-labs/fury/x/jinx"
	jinxkeeper "github.com/percosis-labs/fury/x/jinx/keeper"
)

type BorrowIntegrationTests struct {
	testutil.IntegrationTester

	genesisTime time.Time
	addrs       []sdk.AccAddress
}

func TestBorrowIntegration(t *testing.T) {
	suite.Run(t, new(BorrowIntegrationTests))
}

// SetupTest is run automatically before each suite test
func (suite *BorrowIntegrationTests) SetupTest() {
	_, suite.addrs = app.GeneratePrivKeyAddressPairs(5)

	suite.genesisTime = time.Date(2020, 12, 15, 14, 0, 0, 0, time.UTC)
}

func (suite *BorrowIntegrationTests) TestSingleUserAccumulatesRewardsAfterSyncing() {
	userA := suite.addrs[0]

	authBulder := app.NewAuthBankGenesisBuilder().
		WithSimpleModuleAccount(furydisttypes.ModuleName, cs(c("jinx", 1e18))). // Fill furydist with enough coins to pay out any reward
		WithSimpleAccount(userA, cs(c("bnb", 1e12)))                            // give the user some coins

	incentBuilder := testutil.NewIncentiveGenesisBuilder().
		WithGenesisTime(suite.genesisTime).
		WithMultipliers(types.MultipliersPerDenoms{{
			Denom:       "jinx",
			Multipliers: types.Multipliers{types.NewMultiplier("large", 12, d("1.0"))}, // keep payout at 1.0 to make maths easier
		}}).
		WithSimpleBorrowRewardPeriod("bnb", cs(c("jinx", 1e6))) // only borrow rewards

	suite.SetApp()
	suite.WithGenesisTime(suite.genesisTime)
	suite.StartChain(
		NewPricefeedGenStateMultiFromTime(suite.App.AppCodec(), suite.genesisTime),
		NewJinxGenStateMulti(suite.genesisTime).BuildMarshalled(suite.App.AppCodec()),
		authBulder.BuildMarshalled(suite.App.AppCodec()),
		incentBuilder.BuildMarshalled(suite.App.AppCodec()),
	)

	// Create a borrow (need to first deposit to allow it)
	suite.NoError(suite.DeliverJinxMsgDeposit(userA, cs(c("bnb", 1e11))))
	suite.NoError(suite.DeliverJinxMsgBorrow(userA, cs(c("bnb", 1e10))))

	// Let time pass to accumulate interest on the borrow
	// Use one long block instead of many to reduce any rounding errors, and speed up tests.
	suite.NextBlockAfter(1e6 * time.Second) // about 12 days

	// User borrows and repays just to sync their borrow.
	suite.NoError(suite.DeliverJinxMsgRepay(userA, cs(c("bnb", 1))))
	suite.NoError(suite.DeliverJinxMsgBorrow(userA, cs(c("bnb", 1))))

	// Accumulate more rewards.
	// The user still has the same percentage of all borrows (100%) so their rewards should be the same as in the previous block.
	suite.NextBlockAfter(1e6 * time.Second) // about 12 days

	msg := types.NewMsgClaimJinxReward(userA.String(), types.Selections{
		types.NewSelection("jinx", "large"),
	})

	// User claims all their rewards
	suite.NoError(suite.DeliverIncentiveMsg(&msg))

	// The users has always had 100% of borrows, so they should receive all rewards for the previous two blocks.
	// Total rewards for each block is block duration * rewards per second
	accuracy := 1e-10 // using a very high accuracy to flag future small calculation changes
	suite.BalanceInEpsilon(userA, cs(c("bnb", 1e12-1e11+1e10), c("jinx", 2*1e6*1e6)), accuracy)
}

// Test suite used for all keeper tests
type BorrowRewardsTestSuite struct {
	suite.Suite

	keeper          keeper.Keeper
	jinxKeeper      jinxkeeper.Keeper
	committeeKeeper committeekeeper.Keeper

	app app.TestApp
	ctx sdk.Context

	genesisTime time.Time
	addrs       []sdk.AccAddress
}

// SetupTest is run automatically before each suite test
func (suite *BorrowRewardsTestSuite) SetupTest() {
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixes(config)

	_, suite.addrs = app.GeneratePrivKeyAddressPairs(5)

	suite.genesisTime = time.Date(2020, 12, 15, 14, 0, 0, 0, time.UTC)
}

func (suite *BorrowRewardsTestSuite) SetupApp() {
	suite.app = app.NewTestApp()

	suite.keeper = suite.app.GetIncentiveKeeper()
	suite.jinxKeeper = suite.app.GetJinxKeeper()
	suite.committeeKeeper = suite.app.GetCommitteeKeeper()

	suite.ctx = suite.app.NewContext(true, tmproto.Header{Height: 1, Time: suite.genesisTime})
}

func (suite *BorrowRewardsTestSuite) SetupWithGenState(authBuilder *app.AuthBankGenesisBuilder, incentBuilder testutil.IncentiveGenesisBuilder, jinxBuilder testutil.JinxGenesisBuilder) {
	suite.SetupApp()

	suite.app.InitializeFromGenesisStatesWithTime(
		suite.genesisTime,
		authBuilder.BuildMarshalled(suite.app.AppCodec()),
		NewPricefeedGenStateMultiFromTime(suite.app.AppCodec(), suite.genesisTime),
		jinxBuilder.BuildMarshalled(suite.app.AppCodec()),
		NewCommitteeGenesisState(suite.app.AppCodec(), 1, suite.addrs[:2]...),
		incentBuilder.BuildMarshalled(suite.app.AppCodec()),
	)
}

func (suite *BorrowRewardsTestSuite) TestAccumulateJinxBorrowRewards() {
	type args struct {
		borrow                sdk.Coin
		rewardsPerSecond      sdk.Coins
		timeElapsed           int
		expectedRewardIndexes types.RewardIndexes
	}
	type test struct {
		name string
		args args
	}
	testCases := []test{
		{
			"single reward denom: 7 seconds",
			args{
				borrow:                c("bnb", 1000000000000),
				rewardsPerSecond:      cs(c("jinx", 122354)),
				timeElapsed:           7,
				expectedRewardIndexes: types.RewardIndexes{types.NewRewardIndex("jinx", d("0.000000856478000001"))},
			},
		},
		{
			"single reward denom: 1 day",
			args{
				borrow:                c("bnb", 1000000000000),
				rewardsPerSecond:      cs(c("jinx", 122354)),
				timeElapsed:           86400,
				expectedRewardIndexes: types.RewardIndexes{types.NewRewardIndex("jinx", d("0.010571385600010177"))},
			},
		},
		{
			"single reward denom: 0 seconds",
			args{
				borrow:                c("bnb", 1000000000000),
				rewardsPerSecond:      cs(c("jinx", 122354)),
				timeElapsed:           0,
				expectedRewardIndexes: types.RewardIndexes{types.NewRewardIndex("jinx", d("0.0"))},
			},
		},
		{
			"multiple reward denoms: 7 seconds",
			args{
				borrow:           c("bnb", 1000000000000),
				rewardsPerSecond: cs(c("jinx", 122354), c("ufury", 122354)),
				timeElapsed:      7,
				expectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("0.000000856478000001")),
					types.NewRewardIndex("ufury", d("0.000000856478000001")),
				},
			},
		},
		{
			"multiple reward denoms: 1 day",
			args{
				borrow:           c("bnb", 1000000000000),
				rewardsPerSecond: cs(c("jinx", 122354), c("ufury", 122354)),
				timeElapsed:      86400,
				expectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("0.010571385600010177")),
					types.NewRewardIndex("ufury", d("0.010571385600010177")),
				},
			},
		},
		{
			"multiple reward denoms: 0 seconds",
			args{
				borrow:           c("bnb", 1000000000000),
				rewardsPerSecond: cs(c("jinx", 122354), c("ufury", 122354)),
				timeElapsed:      0,
				expectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("0.0")),
					types.NewRewardIndex("ufury", d("0.0")),
				},
			},
		},
		{
			"multiple reward denoms with different rewards per second: 1 day",
			args{
				borrow:           c("bnb", 1000000000000),
				rewardsPerSecond: cs(c("jinx", 122354), c("ufury", 555555)),
				timeElapsed:      86400,
				expectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("0.010571385600010177")),
					types.NewRewardIndex("ufury", d("0.047999952000046210")),
				},
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userAddr := suite.addrs[3]
			authBuilder := app.NewAuthBankGenesisBuilder().WithSimpleAccount(
				userAddr,
				cs(c("bnb", 1e15), c("ufury", 1e15), c("btcb", 1e15), c("xrp", 1e15), c("zzz", 1e15)),
			)

			incentBuilder := testutil.NewIncentiveGenesisBuilder().
				WithGenesisTime(suite.genesisTime).
				WithSimpleBorrowRewardPeriod(tc.args.borrow.Denom, tc.args.rewardsPerSecond)

			suite.SetupWithGenState(authBuilder, incentBuilder, NewJinxGenStateMulti(suite.genesisTime))

			// User deposits and borrows to increase total borrowed amount
			err := suite.jinxKeeper.Deposit(suite.ctx, userAddr, sdk.NewCoins(sdk.NewCoin(tc.args.borrow.Denom, tc.args.borrow.Amount.Mul(sdkmath.NewInt(2)))))
			suite.Require().NoError(err)
			err = suite.jinxKeeper.Borrow(suite.ctx, userAddr, sdk.NewCoins(tc.args.borrow))
			suite.Require().NoError(err)

			// Set up chain context at future time
			runAtTime := suite.ctx.BlockTime().Add(time.Duration(int(time.Second) * tc.args.timeElapsed))
			runCtx := suite.ctx.WithBlockTime(runAtTime)

			// Run Jinx begin blocker in order to update the denom's index factor
			jinx.BeginBlocker(runCtx, suite.jinxKeeper)

			// Accumulate jinx borrow rewards for the deposit denom
			multiRewardPeriod, found := suite.keeper.GetJinxBorrowRewardPeriods(runCtx, tc.args.borrow.Denom)
			suite.Require().True(found)
			suite.keeper.AccumulateJinxBorrowRewards(runCtx, multiRewardPeriod)

			// Check that each expected reward index matches the current stored reward index for the denom
			globalRewardIndexes, found := suite.keeper.GetJinxBorrowRewardIndexes(runCtx, tc.args.borrow.Denom)
			suite.Require().True(found)
			for _, expectedRewardIndex := range tc.args.expectedRewardIndexes {
				globalRewardIndex, found := globalRewardIndexes.GetRewardIndex(expectedRewardIndex.CollateralType)
				suite.Require().True(found)
				suite.Require().Equal(expectedRewardIndex, globalRewardIndex)
			}
		})
	}
}

func (suite *BorrowRewardsTestSuite) TestInitializeJinxBorrowRewards() {
	type args struct {
		moneyMarketRewardDenoms          map[string]sdk.Coins
		deposit                          sdk.Coins
		borrow                           sdk.Coins
		expectedClaimBorrowRewardIndexes types.MultiRewardIndexes
	}
	type test struct {
		name string
		args args
	}

	standardMoneyMarketRewardDenoms := map[string]sdk.Coins{
		"bnb":  cs(c("jinx", 1)),
		"btcb": cs(c("jinx", 1), c("ufury", 1)),
	}

	testCases := []test{
		{
			"single deposit denom, single reward denom",
			args{
				moneyMarketRewardDenoms: standardMoneyMarketRewardDenoms,
				deposit:                 cs(c("bnb", 1000000000000)),
				borrow:                  cs(c("bnb", 100000000000)),
				expectedClaimBorrowRewardIndexes: types.MultiRewardIndexes{
					types.NewMultiRewardIndex(
						"bnb",
						types.RewardIndexes{
							types.NewRewardIndex("jinx", d("0.0")),
						},
					),
				},
			},
		},
		{
			"single deposit denom, multiple reward denoms",
			args{
				moneyMarketRewardDenoms: standardMoneyMarketRewardDenoms,
				deposit:                 cs(c("btcb", 1000000000000)),
				borrow:                  cs(c("btcb", 100000000000)),
				expectedClaimBorrowRewardIndexes: types.MultiRewardIndexes{
					types.NewMultiRewardIndex(
						"btcb",
						types.RewardIndexes{
							types.NewRewardIndex("jinx", d("0.0")),
							types.NewRewardIndex("ufury", d("0.0")),
						},
					),
				},
			},
		},
		{
			"single deposit denom, no reward denoms",
			args{
				moneyMarketRewardDenoms: standardMoneyMarketRewardDenoms,
				deposit:                 cs(c("xrp", 1000000000000)),
				borrow:                  cs(c("xrp", 100000000000)),
				expectedClaimBorrowRewardIndexes: types.MultiRewardIndexes{
					types.NewMultiRewardIndex(
						"xrp",
						nil,
					),
				},
			},
		},
		{
			"multiple deposit denoms, multiple overlapping reward denoms",
			args{
				moneyMarketRewardDenoms: standardMoneyMarketRewardDenoms,
				deposit:                 cs(c("bnb", 1000000000000), c("btcb", 1000000000000)),
				borrow:                  cs(c("bnb", 100000000000), c("btcb", 100000000000)),
				expectedClaimBorrowRewardIndexes: types.MultiRewardIndexes{
					types.NewMultiRewardIndex(
						"bnb",
						types.RewardIndexes{
							types.NewRewardIndex("jinx", d("0.0")),
						},
					),
					types.NewMultiRewardIndex(
						"btcb",
						types.RewardIndexes{
							types.NewRewardIndex("jinx", d("0.0")),
							types.NewRewardIndex("ufury", d("0.0")),
						},
					),
				},
			},
		},
		{
			"multiple deposit denoms, correct discrete reward denoms",
			args{
				moneyMarketRewardDenoms: standardMoneyMarketRewardDenoms,
				deposit:                 cs(c("bnb", 1000000000000), c("xrp", 1000000000000)),
				borrow:                  cs(c("bnb", 100000000000), c("xrp", 100000000000)),
				expectedClaimBorrowRewardIndexes: types.MultiRewardIndexes{
					types.NewMultiRewardIndex(
						"bnb",
						types.RewardIndexes{
							types.NewRewardIndex("jinx", d("0.0")),
						},
					),
					types.NewMultiRewardIndex(
						"xrp",
						nil,
					),
				},
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userAddr := suite.addrs[3]
			authBuilder := app.NewAuthBankGenesisBuilder().WithSimpleAccount(
				userAddr,
				cs(c("bnb", 1e15), c("ufury", 1e15), c("btcb", 1e15), c("xrp", 1e15), c("zzz", 1e15)),
			)

			incentBuilder := testutil.NewIncentiveGenesisBuilder().WithGenesisTime(suite.genesisTime)
			for moneyMarketDenom, rewardsPerSecond := range tc.args.moneyMarketRewardDenoms {
				incentBuilder = incentBuilder.WithSimpleBorrowRewardPeriod(moneyMarketDenom, rewardsPerSecond)
			}

			suite.SetupWithGenState(authBuilder, incentBuilder, NewJinxGenStateMulti(suite.genesisTime))

			// User deposits
			err := suite.jinxKeeper.Deposit(suite.ctx, userAddr, tc.args.deposit)
			suite.Require().NoError(err)
			// User borrows
			err = suite.jinxKeeper.Borrow(suite.ctx, userAddr, tc.args.borrow)
			suite.Require().NoError(err)

			claim, foundClaim := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, userAddr)
			suite.Require().True(foundClaim)
			suite.Require().Equal(tc.args.expectedClaimBorrowRewardIndexes, claim.BorrowRewardIndexes)
		})
	}
}

func (suite *BorrowRewardsTestSuite) TestSynchronizeJinxBorrowReward() {
	type args struct {
		incentiveBorrowRewardDenom   string
		borrow                       sdk.Coin
		rewardsPerSecond             sdk.Coins
		blockTimes                   []int
		expectedRewardIndexes        types.RewardIndexes
		expectedRewards              sdk.Coins
		updateRewardsViaCommmittee   bool
		updatedBaseDenom             string
		updatedRewardsPerSecond      sdk.Coins
		updatedExpectedRewardIndexes types.RewardIndexes
		updatedExpectedRewards       sdk.Coins
		updatedTimeDuration          int
	}
	type test struct {
		name string
		args args
	}

	testCases := []test{
		{
			"single reward denom: 10 blocks",
			args{
				incentiveBorrowRewardDenom: "bnb",
				borrow:                     c("bnb", 10000000000),
				rewardsPerSecond:           cs(c("jinx", 122354)),
				blockTimes:                 []int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
				expectedRewardIndexes:      types.RewardIndexes{types.NewRewardIndex("jinx", d("0.001223540000173228"))},
				expectedRewards:            cs(c("jinx", 12235400)),
				updateRewardsViaCommmittee: false,
			},
		},
		{
			"single reward denom: 10 blocks - long block time",
			args{
				incentiveBorrowRewardDenom: "bnb",
				borrow:                     c("bnb", 10000000000),
				rewardsPerSecond:           cs(c("jinx", 122354)),
				blockTimes:                 []int{86400, 86400, 86400, 86400, 86400, 86400, 86400, 86400, 86400, 86400},
				expectedRewardIndexes:      types.RewardIndexes{types.NewRewardIndex("jinx", d("10.571385603126235340"))},
				expectedRewards:            cs(c("jinx", 105713856031)),
			},
		},
		{
			"single reward denom: user reward index updated when reward is zero",
			args{
				incentiveBorrowRewardDenom: "ufury",
				borrow:                     c("ufury", 1), // borrow a tiny amount so that rewards round to zero
				rewardsPerSecond:           cs(c("jinx", 122354)),
				blockTimes:                 []int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
				expectedRewardIndexes:      types.RewardIndexes{types.NewRewardIndex("jinx", d("0.122354003908172328"))},
				expectedRewards:            cs(),
				updateRewardsViaCommmittee: false,
			},
		},
		{
			"multiple reward denoms: 10 blocks",
			args{
				incentiveBorrowRewardDenom: "bnb",
				borrow:                     c("bnb", 10000000000),
				rewardsPerSecond:           cs(c("jinx", 122354), c("ufury", 122354)),
				blockTimes:                 []int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
				expectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("0.001223540000173228")),
					types.NewRewardIndex("ufury", d("0.001223540000173228")),
				},
				expectedRewards: cs(c("jinx", 12235400), c("ufury", 12235400)),
			},
		},
		{
			"multiple reward denoms: 10 blocks - long block time",
			args{
				incentiveBorrowRewardDenom: "bnb",
				borrow:                     c("bnb", 10000000000),
				rewardsPerSecond:           cs(c("jinx", 122354), c("ufury", 122354)),
				blockTimes:                 []int{86400, 86400, 86400, 86400, 86400, 86400, 86400, 86400, 86400, 86400},
				expectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("10.571385603126235340")),
					types.NewRewardIndex("ufury", d("10.571385603126235340")),
				},
				expectedRewards: cs(c("jinx", 105713856031), c("ufury", 105713856031)),
			},
		},
		{
			"multiple reward denoms with different rewards per second: 10 blocks",
			args{
				incentiveBorrowRewardDenom: "bnb",
				borrow:                     c("bnb", 10000000000),
				rewardsPerSecond:           cs(c("jinx", 122354), c("ufury", 555555)),
				blockTimes:                 []int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
				expectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("0.001223540000173228")),
					types.NewRewardIndex("ufury", d("0.005555550000786558")),
				},
				expectedRewards: cs(c("jinx", 12235400), c("ufury", 55555500)),
			},
		},
		{
			"denom is in incentive's jinx borrow reward params and has rewards; add new reward type",
			args{
				incentiveBorrowRewardDenom: "bnb",
				borrow:                     c("bnb", 10000000000),
				rewardsPerSecond:           cs(c("jinx", 122354)),
				blockTimes:                 []int{86400},
				expectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("1.057138560060101160")),
				},
				expectedRewards:            cs(c("jinx", 10571385601)),
				updateRewardsViaCommmittee: true,
				updatedBaseDenom:           "bnb",
				updatedRewardsPerSecond:    cs(c("jinx", 122354), c("ufury", 100000)),
				updatedExpectedRewards:     cs(c("jinx", 21142771202), c("ufury", 8640000000)),
				updatedExpectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("2.114277120120202320")),
					types.NewRewardIndex("ufury", d("0.864000000049120715")),
				},
				updatedTimeDuration: 86400,
			},
		},
		{
			"denom is in jinx's money market params but not in incentive's jinx supply reward params; add reward",
			args{
				incentiveBorrowRewardDenom: "bnb",
				borrow:                     c("zzz", 10000000000),
				rewardsPerSecond:           nil,
				blockTimes:                 []int{100},
				expectedRewardIndexes:      types.RewardIndexes{},
				expectedRewards:            sdk.Coins{},
				updateRewardsViaCommmittee: true,
				updatedBaseDenom:           "zzz",
				updatedRewardsPerSecond:    cs(c("jinx", 100000)),
				updatedExpectedRewards:     cs(c("jinx", 8640000000)),
				updatedExpectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("0.864000000049803065")),
				},
				updatedTimeDuration: 86400,
			},
		},
		{
			"denom is in jinx's money market params but not in incentive's jinx supply reward params; add multiple reward types",
			args{
				incentiveBorrowRewardDenom: "bnb",
				borrow:                     c("zzz", 10000000000),
				rewardsPerSecond:           nil,
				blockTimes:                 []int{100},
				expectedRewardIndexes:      types.RewardIndexes{},
				expectedRewards:            sdk.Coins{},
				updateRewardsViaCommmittee: true,
				updatedBaseDenom:           "zzz",
				updatedRewardsPerSecond:    cs(c("jinx", 100000), c("ufury", 100500), c("swap", 500)),
				updatedExpectedRewards:     cs(c("jinx", 8640000000), c("ufury", 8683200001), c("swap", 43200000)),
				updatedExpectedRewardIndexes: types.RewardIndexes{
					types.NewRewardIndex("jinx", d("0.864000000049803065")),
					types.NewRewardIndex("ufury", d("0.868320000050052081")),
					types.NewRewardIndex("swap", d("0.004320000000249015")),
				},
				updatedTimeDuration: 86400,
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userAddr := suite.addrs[3]
			authBuilder := app.NewAuthBankGenesisBuilder().
				WithSimpleAccount(suite.addrs[2], cs(c("ufury", 1e9))).
				WithSimpleAccount(userAddr, cs(c("bnb", 1e15), c("ufury", 1e15), c("btcb", 1e15), c("xrp", 1e15), c("zzz", 1e15)))

			incentBuilder := testutil.NewIncentiveGenesisBuilder().WithGenesisTime(suite.genesisTime)
			if tc.args.rewardsPerSecond != nil {
				incentBuilder = incentBuilder.WithSimpleBorrowRewardPeriod(tc.args.incentiveBorrowRewardDenom, tc.args.rewardsPerSecond)
			}
			// Set the minimum borrow to 0 to allow testing small borrows
			jinxBuilder := NewJinxGenStateMulti(suite.genesisTime).WithMinBorrow(sdk.ZeroDec())

			suite.SetupWithGenState(authBuilder, incentBuilder, jinxBuilder)

			// Borrow a fixed amount from another user to dilute primary user's rewards per second.
			suite.Require().NoError(
				suite.jinxKeeper.Deposit(suite.ctx, suite.addrs[2], cs(c("ufury", 200_000_000))),
			)
			suite.Require().NoError(
				suite.jinxKeeper.Borrow(suite.ctx, suite.addrs[2], cs(c("ufury", 100_000_000))),
			)

			// User deposits and borrows to increase total borrowed amount
			err := suite.jinxKeeper.Deposit(suite.ctx, userAddr, sdk.NewCoins(sdk.NewCoin(tc.args.borrow.Denom, tc.args.borrow.Amount.Mul(sdkmath.NewInt(2)))))
			suite.Require().NoError(err)
			err = suite.jinxKeeper.Borrow(suite.ctx, userAddr, sdk.NewCoins(tc.args.borrow))
			suite.Require().NoError(err)

			// Check that Jinx hooks initialized a JinxLiquidityProviderClaim
			claim, found := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, userAddr)
			suite.Require().True(found)
			multiRewardIndex, _ := claim.BorrowRewardIndexes.GetRewardIndex(tc.args.borrow.Denom)
			for _, expectedRewardIndex := range tc.args.expectedRewardIndexes {
				currRewardIndex, found := multiRewardIndex.RewardIndexes.GetRewardIndex(expectedRewardIndex.CollateralType)
				suite.Require().True(found)
				suite.Require().Equal(sdk.ZeroDec(), currRewardIndex.RewardFactor)
			}

			// Run accumulator at several intervals
			var timeElapsed int
			previousBlockTime := suite.ctx.BlockTime()
			for _, t := range tc.args.blockTimes {
				timeElapsed += t
				updatedBlockTime := previousBlockTime.Add(time.Duration(int(time.Second) * t))
				previousBlockTime = updatedBlockTime
				blockCtx := suite.ctx.WithBlockTime(updatedBlockTime)

				// Run Jinx begin blocker for each block ctx to update denom's interest factor
				jinx.BeginBlocker(blockCtx, suite.jinxKeeper)

				// Accumulate jinx borrow-side rewards
				multiRewardPeriod, found := suite.keeper.GetJinxBorrowRewardPeriods(blockCtx, tc.args.borrow.Denom)
				if found {
					suite.keeper.AccumulateJinxBorrowRewards(blockCtx, multiRewardPeriod)
				}
			}
			updatedBlockTime := suite.ctx.BlockTime().Add(time.Duration(int(time.Second) * timeElapsed))
			suite.ctx = suite.ctx.WithBlockTime(updatedBlockTime)

			// After we've accumulated, run synchronize
			borrow, found := suite.jinxKeeper.GetBorrow(suite.ctx, userAddr)
			suite.Require().True(found)
			suite.Require().NotPanics(func() {
				suite.keeper.SynchronizeJinxBorrowReward(suite.ctx, borrow)
			})

			// Check that the global reward index's reward factor and user's claim have been updated as expected
			claim, found = suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, userAddr)
			suite.Require().True(found)
			globalRewardIndexes, foundGlobalRewardIndexes := suite.keeper.GetJinxBorrowRewardIndexes(suite.ctx, tc.args.borrow.Denom)
			if len(tc.args.rewardsPerSecond) > 0 {
				suite.Require().True(foundGlobalRewardIndexes)
				for _, expectedRewardIndex := range tc.args.expectedRewardIndexes {
					// Check that global reward index has been updated as expected
					globalRewardIndex, found := globalRewardIndexes.GetRewardIndex(expectedRewardIndex.CollateralType)
					suite.Require().True(found)
					suite.Require().Equal(expectedRewardIndex, globalRewardIndex)

					// Check that the user's claim's reward index matches the corresponding global reward index
					multiRewardIndex, found := claim.BorrowRewardIndexes.GetRewardIndex(tc.args.borrow.Denom)
					suite.Require().True(found)
					rewardIndex, found := multiRewardIndex.RewardIndexes.GetRewardIndex(expectedRewardIndex.CollateralType)
					suite.Require().True(found)
					suite.Require().Equal(expectedRewardIndex, rewardIndex)

					// Check that the user's claim holds the expected amount of reward coins
					suite.Require().Equal(
						tc.args.expectedRewards.AmountOf(expectedRewardIndex.CollateralType),
						claim.Reward.AmountOf(expectedRewardIndex.CollateralType),
					)
				}
			}

			// Only test cases with reward param updates continue past this point
			if !tc.args.updateRewardsViaCommmittee {
				return
			}

			// If are no initial rewards per second, add new rewards through a committee param change
			// 1. Construct incentive's new JinxBorrowRewardPeriods param
			currIncentiveJinxBorrowRewardPeriods := suite.keeper.GetParams(suite.ctx).JinxBorrowRewardPeriods
			multiRewardPeriod, found := currIncentiveJinxBorrowRewardPeriods.GetMultiRewardPeriod(tc.args.borrow.Denom)
			if found {
				// Borrow denom's reward period exists, but it doesn't have any rewards per second
				index, found := currIncentiveJinxBorrowRewardPeriods.GetMultiRewardPeriodIndex(tc.args.borrow.Denom)
				suite.Require().True(found)
				multiRewardPeriod.RewardsPerSecond = tc.args.updatedRewardsPerSecond
				currIncentiveJinxBorrowRewardPeriods[index] = multiRewardPeriod
			} else {
				// Borrow denom's reward period does not exist
				_, found := currIncentiveJinxBorrowRewardPeriods.GetMultiRewardPeriodIndex(tc.args.borrow.Denom)
				suite.Require().False(found)
				newMultiRewardPeriod := types.NewMultiRewardPeriod(true, tc.args.borrow.Denom, suite.genesisTime, suite.genesisTime.Add(time.Hour*24*365*4), tc.args.updatedRewardsPerSecond)
				currIncentiveJinxBorrowRewardPeriods = append(currIncentiveJinxBorrowRewardPeriods, newMultiRewardPeriod)
			}

			// 2. Construct the parameter change proposal to update JinxBorrowRewardPeriods param
			pubProposal := proposaltypes.NewParameterChangeProposal(
				"Update jinx borrow rewards", "Adds a new reward coin to the incentive module's jinx borrow rewards.",
				[]proposaltypes.ParamChange{
					{
						Subspace: types.ModuleName,                         // target incentive module
						Key:      string(types.KeyJinxBorrowRewardPeriods), // target jinx borrow rewards key
						Value:    string(suite.app.LegacyAmino().MustMarshalJSON(currIncentiveJinxBorrowRewardPeriods)),
					},
				},
			)

			// 3. Ensure proposal is properly formed
			err = suite.committeeKeeper.ValidatePubProposal(suite.ctx, pubProposal)
			suite.Require().NoError(err)

			// 4. Committee creates proposal
			committeeMemberOne := suite.addrs[0]
			committeeMemberTwo := suite.addrs[1]
			proposalID, err := suite.committeeKeeper.SubmitProposal(suite.ctx, committeeMemberOne, 1, pubProposal)
			suite.Require().NoError(err)

			// 5. Committee votes and passes proposal
			err = suite.committeeKeeper.AddVote(suite.ctx, proposalID, committeeMemberOne, committeetypes.VOTE_TYPE_YES)
			suite.Require().NoError(err)
			err = suite.committeeKeeper.AddVote(suite.ctx, proposalID, committeeMemberTwo, committeetypes.VOTE_TYPE_YES)
			suite.Require().NoError(err)

			// 6. Check proposal passed
			com, found := suite.committeeKeeper.GetCommittee(suite.ctx, 1)
			suite.Require().True(found)
			proposalPasses := suite.committeeKeeper.GetProposalResult(suite.ctx, proposalID, com)
			suite.Require().NoError(err)
			suite.Require().True(proposalPasses)

			// 7. Run committee module's begin blocker to enact proposal
			suite.NotPanics(func() {
				committee.BeginBlocker(suite.ctx, abci.RequestBeginBlock{}, suite.committeeKeeper)
			})

			// We need to accumulate jinx supply-side rewards again
			multiRewardPeriod, found = suite.keeper.GetJinxBorrowRewardPeriods(suite.ctx, tc.args.borrow.Denom)
			suite.Require().True(found)

			// But new borrow denoms don't have their PreviousJinxBorrowRewardAccrualTime set yet,
			// so we need to call the accumulation method once to set the initial reward accrual time
			if tc.args.borrow.Denom != tc.args.incentiveBorrowRewardDenom {
				suite.keeper.AccumulateJinxBorrowRewards(suite.ctx, multiRewardPeriod)
			}

			// Now we can jump forward in time and accumulate rewards
			updatedBlockTime = previousBlockTime.Add(time.Duration(int(time.Second) * tc.args.updatedTimeDuration))
			suite.ctx = suite.ctx.WithBlockTime(updatedBlockTime)
			suite.keeper.AccumulateJinxBorrowRewards(suite.ctx, multiRewardPeriod)

			// After we've accumulated, run synchronize
			borrow, found = suite.jinxKeeper.GetBorrow(suite.ctx, userAddr)
			suite.Require().True(found)
			suite.Require().NotPanics(func() {
				suite.keeper.SynchronizeJinxBorrowReward(suite.ctx, borrow)
			})

			// Check that the global reward index's reward factor and user's claim have been updated as expected
			globalRewardIndexes, found = suite.keeper.GetJinxBorrowRewardIndexes(suite.ctx, tc.args.borrow.Denom)
			suite.Require().True(found)
			claim, found = suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, userAddr)
			suite.Require().True(found)

			for _, expectedRewardIndex := range tc.args.updatedExpectedRewardIndexes {
				// Check that global reward index has been updated as expected
				globalRewardIndex, found := globalRewardIndexes.GetRewardIndex(expectedRewardIndex.CollateralType)
				suite.Require().True(found)
				suite.Require().Equal(expectedRewardIndex, globalRewardIndex)
				// Check that the user's claim's reward index matches the corresponding global reward index
				multiRewardIndex, found := claim.BorrowRewardIndexes.GetRewardIndex(tc.args.borrow.Denom)
				suite.Require().True(found)
				rewardIndex, found := multiRewardIndex.RewardIndexes.GetRewardIndex(expectedRewardIndex.CollateralType)
				suite.Require().True(found)
				suite.Require().Equal(expectedRewardIndex, rewardIndex)

				// Check that the user's claim holds the expected amount of reward coins
				suite.Require().Equal(
					tc.args.updatedExpectedRewards.AmountOf(expectedRewardIndex.CollateralType),
					claim.Reward.AmountOf(expectedRewardIndex.CollateralType),
				)
			}
		})
	}
}

func (suite *BorrowRewardsTestSuite) TestUpdateJinxBorrowIndexDenoms() {
	type withdrawModification struct {
		coins sdk.Coins
		repay bool
	}

	type args struct {
		initialDeposit            sdk.Coins
		firstBorrow               sdk.Coins
		modification              withdrawModification
		rewardsPerSecond          sdk.Coins
		expectedBorrowIndexDenoms []string
	}
	type test struct {
		name string
		args args
	}

	testCases := []test{
		{
			"single reward denom: update adds one borrow reward index",
			args{
				initialDeposit:            cs(c("bnb", 10000000000)),
				firstBorrow:               cs(c("bnb", 50000000)),
				modification:              withdrawModification{coins: cs(c("ufury", 500000000))},
				rewardsPerSecond:          cs(c("jinx", 122354)),
				expectedBorrowIndexDenoms: []string{"bnb", "ufury"},
			},
		},
		{
			"single reward denom: update adds multiple borrow supply reward indexes",
			args{
				initialDeposit:            cs(c("btcb", 10000000000)),
				firstBorrow:               cs(c("btcb", 50000000)),
				modification:              withdrawModification{coins: cs(c("ufury", 500000000), c("bnb", 50000000000), c("xrp", 50000000000))},
				rewardsPerSecond:          cs(c("jinx", 122354)),
				expectedBorrowIndexDenoms: []string{"btcb", "ufury", "bnb", "xrp"},
			},
		},
		{
			"single reward denom: update doesn't add duplicate borrow reward index for same denom",
			args{
				initialDeposit:            cs(c("bnb", 100000000000)),
				firstBorrow:               cs(c("bnb", 50000000)),
				modification:              withdrawModification{coins: cs(c("bnb", 50000000000))},
				rewardsPerSecond:          cs(c("jinx", 122354)),
				expectedBorrowIndexDenoms: []string{"bnb"},
			},
		},
		{
			"multiple reward denoms: update adds one borrow reward index",
			args{
				initialDeposit:            cs(c("bnb", 10000000000)),
				firstBorrow:               cs(c("bnb", 50000000)),
				modification:              withdrawModification{coins: cs(c("ufury", 500000000))},
				rewardsPerSecond:          cs(c("jinx", 122354), c("ufury", 122354)),
				expectedBorrowIndexDenoms: []string{"bnb", "ufury"},
			},
		},
		{
			"multiple reward denoms: update adds multiple borrow supply reward indexes",
			args{
				initialDeposit:            cs(c("btcb", 10000000000)),
				firstBorrow:               cs(c("btcb", 50000000)),
				modification:              withdrawModification{coins: cs(c("ufury", 500000000), c("bnb", 50000000000), c("xrp", 50000000000))},
				rewardsPerSecond:          cs(c("jinx", 122354), c("ufury", 122354)),
				expectedBorrowIndexDenoms: []string{"btcb", "ufury", "bnb", "xrp"},
			},
		},
		{
			"multiple reward denoms: update doesn't add duplicate borrow reward index for same denom",
			args{
				initialDeposit:            cs(c("bnb", 100000000000)),
				firstBorrow:               cs(c("bnb", 50000000)),
				modification:              withdrawModification{coins: cs(c("bnb", 50000000000))},
				rewardsPerSecond:          cs(c("jinx", 122354), c("ufury", 122354)),
				expectedBorrowIndexDenoms: []string{"bnb"},
			},
		},
		{
			"single reward denom: fully repaying a denom deletes the denom's supply reward index",
			args{
				initialDeposit:            cs(c("bnb", 1000000000)),
				firstBorrow:               cs(c("bnb", 100000000)),
				modification:              withdrawModification{coins: cs(c("bnb", 1100000000)), repay: true},
				rewardsPerSecond:          cs(c("jinx", 122354)),
				expectedBorrowIndexDenoms: []string{},
			},
		},
		{
			"single reward denom: fully repaying a denom deletes only the denom's supply reward index",
			args{
				initialDeposit:            cs(c("bnb", 1000000000)),
				firstBorrow:               cs(c("bnb", 100000000), c("ufury", 10000000)),
				modification:              withdrawModification{coins: cs(c("bnb", 1100000000)), repay: true},
				rewardsPerSecond:          cs(c("jinx", 122354)),
				expectedBorrowIndexDenoms: []string{"ufury"},
			},
		},
		{
			"multiple reward denoms: fully repaying a denom deletes the denom's supply reward index",
			args{
				initialDeposit:            cs(c("bnb", 1000000000)),
				firstBorrow:               cs(c("bnb", 100000000), c("ufury", 10000000)),
				modification:              withdrawModification{coins: cs(c("bnb", 1100000000)), repay: true},
				rewardsPerSecond:          cs(c("jinx", 122354), c("ufury", 122354)),
				expectedBorrowIndexDenoms: []string{"ufury"},
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userAddr := suite.addrs[3]
			authBuilder := app.NewAuthBankGenesisBuilder().
				WithSimpleAccount(
					userAddr,
					cs(c("bnb", 1e15), c("ufury", 1e15), c("btcb", 1e15), c("xrp", 1e15), c("zzz", 1e15)),
				).
				WithSimpleAccount(
					suite.addrs[0],
					cs(c("bnb", 1e15), c("ufury", 1e15), c("btcb", 1e15), c("xrp", 1e15), c("zzz", 1e15)),
				)

			incentBuilder := testutil.NewIncentiveGenesisBuilder().
				WithGenesisTime(suite.genesisTime).
				WithSimpleBorrowRewardPeriod("bnb", tc.args.rewardsPerSecond).
				WithSimpleBorrowRewardPeriod("ufury", tc.args.rewardsPerSecond).
				WithSimpleBorrowRewardPeriod("btcb", tc.args.rewardsPerSecond).
				WithSimpleBorrowRewardPeriod("xrp", tc.args.rewardsPerSecond)

			suite.SetupWithGenState(authBuilder, incentBuilder, NewJinxGenStateMulti(suite.genesisTime))

			// Fill the jinx supply to allow user to borrow
			err := suite.jinxKeeper.Deposit(suite.ctx, suite.addrs[0], tc.args.firstBorrow.Add(tc.args.modification.coins...))
			suite.Require().NoError(err)

			// User deposits initial funds (so that user can borrow)
			err = suite.jinxKeeper.Deposit(suite.ctx, userAddr, tc.args.initialDeposit)
			suite.Require().NoError(err)

			// Confirm that claim exists but no borrow reward indexes have been added
			claimAfterDeposit, found := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, userAddr)
			suite.Require().True(found)
			suite.Require().Equal(0, len(claimAfterDeposit.BorrowRewardIndexes))

			// User borrows (first time)
			err = suite.jinxKeeper.Borrow(suite.ctx, userAddr, tc.args.firstBorrow)
			suite.Require().NoError(err)

			// Confirm that claim's borrow reward indexes have been updated
			claimAfterFirstBorrow, found := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, userAddr)
			suite.Require().True(found)
			for _, coin := range tc.args.firstBorrow {
				_, hasIndex := claimAfterFirstBorrow.HasBorrowRewardIndex(coin.Denom)
				suite.Require().True(hasIndex)
			}
			suite.Require().True(len(claimAfterFirstBorrow.BorrowRewardIndexes) == len(tc.args.firstBorrow))

			// User modifies their Borrow by either repaying or borrowing more
			if tc.args.modification.repay {
				err = suite.jinxKeeper.Repay(suite.ctx, userAddr, userAddr, tc.args.modification.coins)
			} else {
				err = suite.jinxKeeper.Borrow(suite.ctx, userAddr, tc.args.modification.coins)
			}
			suite.Require().NoError(err)

			// Confirm that claim's borrow reward indexes contain expected values
			claimAfterModification, found := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, userAddr)
			suite.Require().True(found)
			for _, coin := range tc.args.modification.coins {
				_, hasIndex := claimAfterModification.HasBorrowRewardIndex(coin.Denom)
				if tc.args.modification.repay {
					// Only false if denom is repaid in full
					if tc.args.modification.coins.AmountOf(coin.Denom).GTE(tc.args.firstBorrow.AmountOf(coin.Denom)) {
						suite.Require().False(hasIndex)
					}
				} else {
					suite.Require().True(hasIndex)
				}
			}
			suite.Require().True(len(claimAfterModification.BorrowRewardIndexes) == len(tc.args.expectedBorrowIndexDenoms))
		})
	}
}

func (suite *BorrowRewardsTestSuite) TestSimulateJinxBorrowRewardSynchronization() {
	type args struct {
		borrow                sdk.Coin
		rewardsPerSecond      sdk.Coins
		blockTimes            []int
		expectedRewardIndexes types.RewardIndexes
		expectedRewards       sdk.Coins
	}
	type test struct {
		name string
		args args
	}

	testCases := []test{
		{
			"10 blocks",
			args{
				borrow:                c("bnb", 10000000000),
				rewardsPerSecond:      cs(c("jinx", 122354)),
				blockTimes:            []int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
				expectedRewardIndexes: types.RewardIndexes{types.NewRewardIndex("jinx", d("0.001223540000173228"))},
				expectedRewards:       cs(c("jinx", 12235400)),
			},
		},
		{
			"10 blocks - long block time",
			args{
				borrow:                c("bnb", 10000000000),
				rewardsPerSecond:      cs(c("jinx", 122354)),
				blockTimes:            []int{86400, 86400, 86400, 86400, 86400, 86400, 86400, 86400, 86400, 86400},
				expectedRewardIndexes: types.RewardIndexes{types.NewRewardIndex("jinx", d("10.571385603126235340"))},
				expectedRewards:       cs(c("jinx", 105713856031)),
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userAddr := suite.addrs[3]
			authBuilder := app.NewAuthBankGenesisBuilder().WithSimpleAccount(userAddr, cs(c("bnb", 1e15), c("ufury", 1e15), c("btcb", 1e15), c("xrp", 1e15), c("zzz", 1e15)))

			incentBuilder := testutil.NewIncentiveGenesisBuilder().
				WithGenesisTime(suite.genesisTime).
				WithSimpleBorrowRewardPeriod(tc.args.borrow.Denom, tc.args.rewardsPerSecond)

			suite.SetupWithGenState(authBuilder, incentBuilder, NewJinxGenStateMulti(suite.genesisTime))

			// User deposits and borrows to increase total borrowed amount
			err := suite.jinxKeeper.Deposit(suite.ctx, userAddr, sdk.NewCoins(sdk.NewCoin(tc.args.borrow.Denom, tc.args.borrow.Amount.Mul(sdkmath.NewInt(2)))))
			suite.Require().NoError(err)
			err = suite.jinxKeeper.Borrow(suite.ctx, userAddr, sdk.NewCoins(tc.args.borrow))
			suite.Require().NoError(err)

			// Run accumulator at several intervals
			var timeElapsed int
			previousBlockTime := suite.ctx.BlockTime()
			for _, t := range tc.args.blockTimes {
				timeElapsed += t
				updatedBlockTime := previousBlockTime.Add(time.Duration(int(time.Second) * t))
				previousBlockTime = updatedBlockTime
				blockCtx := suite.ctx.WithBlockTime(updatedBlockTime)

				// Run Jinx begin blocker for each block ctx to update denom's interest factor
				jinx.BeginBlocker(blockCtx, suite.jinxKeeper)

				// Accumulate jinx borrow-side rewards
				multiRewardPeriod, found := suite.keeper.GetJinxBorrowRewardPeriods(blockCtx, tc.args.borrow.Denom)
				suite.Require().True(found)
				suite.keeper.AccumulateJinxBorrowRewards(blockCtx, multiRewardPeriod)
			}
			updatedBlockTime := suite.ctx.BlockTime().Add(time.Duration(int(time.Second) * timeElapsed))
			suite.ctx = suite.ctx.WithBlockTime(updatedBlockTime)

			// Confirm that the user's claim hasn't been synced
			claimPre, foundPre := suite.keeper.GetJinxLiquidityProviderClaim(suite.ctx, userAddr)
			suite.Require().True(foundPre)
			multiRewardIndexPre, _ := claimPre.BorrowRewardIndexes.GetRewardIndex(tc.args.borrow.Denom)
			for _, expectedRewardIndex := range tc.args.expectedRewardIndexes {
				currRewardIndex, found := multiRewardIndexPre.RewardIndexes.GetRewardIndex(expectedRewardIndex.CollateralType)
				suite.Require().True(found)
				suite.Require().Equal(sdk.ZeroDec(), currRewardIndex.RewardFactor)
			}

			// Check that the synced claim held in memory has properly simulated syncing
			syncedClaim := suite.keeper.SimulateJinxSynchronization(suite.ctx, claimPre)
			for _, expectedRewardIndex := range tc.args.expectedRewardIndexes {
				// Check that the user's claim's reward index matches the expected reward index
				multiRewardIndex, found := syncedClaim.BorrowRewardIndexes.GetRewardIndex(tc.args.borrow.Denom)
				suite.Require().True(found)
				rewardIndex, found := multiRewardIndex.RewardIndexes.GetRewardIndex(expectedRewardIndex.CollateralType)
				suite.Require().True(found)
				suite.Require().Equal(expectedRewardIndex, rewardIndex)

				// Check that the user's claim holds the expected amount of reward coins
				suite.Require().Equal(
					tc.args.expectedRewards.AmountOf(expectedRewardIndex.CollateralType),
					syncedClaim.Reward.AmountOf(expectedRewardIndex.CollateralType),
				)
			}
		})
	}
}

func TestBorrowRewardsTestSuite(t *testing.T) {
	suite.Run(t, new(BorrowRewardsTestSuite))
}
