package v0_16

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	app "github.com/percosis-labs/fury/app"
	v015jinx "github.com/percosis-labs/fury/x/jinx/legacy/v0_15"
	v016jinx "github.com/percosis-labs/fury/x/jinx/types"
)

type migrateTestSuite struct {
	suite.Suite

	addresses []sdk.AccAddress
	cdc       codec.Codec
	legacyCdc *codec.LegacyAmino
}

func (s *migrateTestSuite) SetupTest() {
	app.SetSDKConfig()
	config := app.MakeEncodingConfig()
	s.cdc = config.Marshaler

	legacyCodec := codec.NewLegacyAmino()
	s.legacyCdc = legacyCodec

	_, accAddresses := app.GeneratePrivKeyAddressPairs(10)
	s.addresses = accAddresses
}

func (s *migrateTestSuite) TestMigrate_JSON() {
	file := filepath.Join("testdata", "v15-jinx.json")
	data, err := ioutil.ReadFile(file)
	s.Require().NoError(err)
	var v15genstate v015jinx.GenesisState
	err = s.legacyCdc.UnmarshalJSON(data, &v15genstate)
	s.Require().NoError(err)
	genstate := Migrate(v15genstate)
	actual := s.cdc.MustMarshalJSON(genstate)

	file = filepath.Join("testdata", "v16-jinx.json")
	expected, err := ioutil.ReadFile(file)
	s.Require().NoError(err)
	s.Require().JSONEq(string(expected), string(actual))
}

func (s *migrateTestSuite) TestMigrate_GenState() {
	v15genstate := v015jinx.GenesisState{
		Params: v015jinx.Params{
			MoneyMarkets: v015jinx.MoneyMarkets{
				{
					Denom: "fury",
					BorrowLimit: v015jinx.BorrowLimit{
						HasMaxLimit:  true,
						MaximumLimit: sdk.MustNewDecFromStr("0.1"),
						LoanToValue:  sdk.MustNewDecFromStr("0.2"),
					},
					SpotMarketID:     "spot-market-id",
					ConversionFactor: sdkmath.NewInt(110),
					InterestRateModel: v015jinx.InterestRateModel{
						BaseRateAPY:    sdk.MustNewDecFromStr("0.1"),
						BaseMultiplier: sdk.MustNewDecFromStr("0.2"),
						Kink:           sdk.MustNewDecFromStr("0.3"),
						JumpMultiplier: sdk.MustNewDecFromStr("0.4"),
					},
					ReserveFactor:          sdk.MustNewDecFromStr("0.5"),
					KeeperRewardPercentage: sdk.MustNewDecFromStr("0.6"),
				},
			},
		},
		PreviousAccumulationTimes: v015jinx.GenesisAccumulationTimes{
			{
				CollateralType:           "fury",
				PreviousAccumulationTime: time.Date(1998, time.January, 1, 12, 0, 0, 1, time.UTC),
				SupplyInterestFactor:     sdk.MustNewDecFromStr("0.1"),
				BorrowInterestFactor:     sdk.MustNewDecFromStr("0.2"),
			},
		},
		Deposits: v015jinx.Deposits{
			{
				Depositor: s.addresses[0],
				Amount:    sdk.NewCoins(sdk.NewCoin("fury", sdkmath.NewInt(100))),
				Index: v015jinx.SupplyInterestFactors{
					{
						Denom: "fury",
						Value: sdk.MustNewDecFromStr("1.12"),
					},
				},
			},
		},
		Borrows: v015jinx.Borrows{
			{
				Borrower: s.addresses[1],
				Amount:   sdk.NewCoins(sdk.NewCoin("fury", sdkmath.NewInt(100))),
				Index: v015jinx.BorrowInterestFactors{
					{
						Denom: "fury",
						Value: sdk.MustNewDecFromStr("1.12"),
					},
				},
			},
		},
		TotalSupplied: sdk.NewCoins(sdk.NewCoin("fury", sdkmath.NewInt(100))),
		TotalBorrowed: sdk.NewCoins(sdk.NewCoin("bnb", sdkmath.NewInt(200))),
		TotalReserves: sdk.NewCoins(sdk.NewCoin("xrp", sdkmath.NewInt(300))),
	}
	expected := v016jinx.GenesisState{
		Params: v016jinx.Params{
			MoneyMarkets: v016jinx.MoneyMarkets{
				{
					Denom: "fury",
					BorrowLimit: v016jinx.BorrowLimit{
						HasMaxLimit:  true,
						MaximumLimit: sdk.MustNewDecFromStr("0.1"),
						LoanToValue:  sdk.MustNewDecFromStr("0.2"),
					},
					SpotMarketID:     "spot-market-id",
					ConversionFactor: sdkmath.NewInt(110),
					InterestRateModel: v016jinx.InterestRateModel{
						BaseRateAPY:    sdk.MustNewDecFromStr("0.1"),
						BaseMultiplier: sdk.MustNewDecFromStr("0.2"),
						Kink:           sdk.MustNewDecFromStr("0.3"),
						JumpMultiplier: sdk.MustNewDecFromStr("0.4"),
					},
					ReserveFactor:          sdk.MustNewDecFromStr("0.5"),
					KeeperRewardPercentage: sdk.MustNewDecFromStr("0.6"),
				},
				{
					Denom: UATOM_IBC_DENOM,
					BorrowLimit: v016jinx.BorrowLimit{
						HasMaxLimit:  true,
						MaximumLimit: sdk.NewDec(25000000000),
						LoanToValue:  sdk.MustNewDecFromStr("0.5"),
					},
					SpotMarketID:     "atom:usd:30",
					ConversionFactor: sdkmath.NewInt(1000000),
					InterestRateModel: v016jinx.InterestRateModel{
						BaseRateAPY:    sdk.ZeroDec(),
						BaseMultiplier: sdk.MustNewDecFromStr("0.05"),
						Kink:           sdk.MustNewDecFromStr("0.8"),
						JumpMultiplier: sdk.NewDec(5),
					},
					ReserveFactor:          sdk.MustNewDecFromStr("0.025"),
					KeeperRewardPercentage: sdk.MustNewDecFromStr("0.02"),
				},
			},
		},
		PreviousAccumulationTimes: v016jinx.GenesisAccumulationTimes{
			{
				CollateralType:           "fury",
				PreviousAccumulationTime: time.Date(1998, time.January, 1, 12, 0, 0, 1, time.UTC),
				SupplyInterestFactor:     sdk.MustNewDecFromStr("0.1"),
				BorrowInterestFactor:     sdk.MustNewDecFromStr("0.2"),
			},
		},
		Deposits: v016jinx.Deposits{
			{
				Depositor: s.addresses[0],
				Amount:    sdk.NewCoins(sdk.NewCoin("fury", sdkmath.NewInt(100))),
				Index: v016jinx.SupplyInterestFactors{
					{
						Denom: "fury",
						Value: sdk.MustNewDecFromStr("1.12"),
					},
				},
			},
		},
		Borrows: v016jinx.Borrows{
			{
				Borrower: s.addresses[1],
				Amount:   sdk.NewCoins(sdk.NewCoin("fury", sdkmath.NewInt(100))),
				Index: v016jinx.BorrowInterestFactors{
					{
						Denom: "fury",
						Value: sdk.MustNewDecFromStr("1.12"),
					},
				},
			},
		},
		TotalSupplied: sdk.NewCoins(sdk.NewCoin("fury", sdkmath.NewInt(100))),
		TotalBorrowed: sdk.NewCoins(sdk.NewCoin("bnb", sdkmath.NewInt(200))),
		TotalReserves: sdk.NewCoins(sdk.NewCoin("xrp", sdkmath.NewInt(300))),
	}
	genState := Migrate(v15genstate)
	s.Require().Equal(expected, *genState)
}

func TestJinxMigrateTestSuite(t *testing.T) {
	suite.Run(t, new(migrateTestSuite))
}
