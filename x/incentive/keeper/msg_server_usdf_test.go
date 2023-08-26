package keeper_test

import (
	"time"

	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/percosis-labs/fury/x/incentive/types"
)

func (suite *HandlerTestSuite) TestPayoutUSDFClaim() {
	userAddr, receiverAddr := suite.addrs[0], suite.addrs[1]

	authBulder := suite.authBuilder().
		WithSimpleAccount(userAddr, cs(c("bnb", 1e12))).
		WithSimpleAccount(receiverAddr, nil)

	incentBuilder := suite.incentiveBuilder().
		WithSimpleUSDFRewardPeriod("bnb-a", c(types.USDFMintingRewardDenom, 1e6))

	suite.SetupWithGenState(authBulder, incentBuilder)

	// mint some usdf
	err := suite.DeliverMsgCreateCDP(userAddr, c("bnb", 1e9), c("usdf", 1e7), "bnb-a")
	suite.NoError(err)
	// accumulate some rewards
	suite.NextBlockAfter(7 * time.Second)

	preClaimBal := suite.GetBalance(userAddr)

	msg := types.NewMsgClaimUSDFMintingReward(userAddr.String(), "large")

	// Claim a single denom
	err = suite.DeliverIncentiveMsg(&msg)
	suite.NoError(err)

	// Check rewards were paid out
	expectedRewards := cs(c(types.USDFMintingRewardDenom, 7*1e6))
	suite.BalanceEquals(userAddr, preClaimBal.Add(expectedRewards...))

	suite.VestingPeriodsEqual(userAddr, []vestingtypes.Period{
		{Length: (17+31+28+31+30+31+30+31+31+30+31+30+31)*secondsPerDay - 7, Amount: expectedRewards},
	})
	// Check that claimed coins have been removed from a claim's reward
	suite.USDFRewardEquals(userAddr, c(types.USDFMintingRewardDenom, 0))
}
