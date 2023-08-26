package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/percosis-labs/fury/x/incentive/types"
)

const (
	flagOwner    = "owner"
	flagType     = "type"
	flagUnsynced = "unsynced"
	flagDenom    = "denom"

	typeDelegator   = "delegator"
	typeJinx        = "jinx"
	typeUSDFMinting = "usdf-minting"
	typeSwap        = "swap"
	typeSavings     = "savings"
	typeEarn        = "earn"
)

var rewardTypes = []string{typeDelegator, typeJinx, typeUSDFMinting, typeSwap, typeEarn}

// GetQueryCmd returns the cli query commands for the incentive module
func GetQueryCmd() *cobra.Command {
	incentiveQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the incentive module",
	}

	cmds := []*cobra.Command{
		queryParamsCmd(),
		queryRewardsCmd(),
		queryRewardFactorsCmd(),
	}

	for _, cmd := range cmds {
		flags.AddQueryFlagsToCmd(cmd)
	}

	incentiveQueryCmd.AddCommand(cmds...)

	return incentiveQueryCmd
}

func queryRewardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards",
		Short: "query claimable rewards",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query rewards with optional flags for owner and type

			Example:
			$ %[1]s query %[2]s rewards
			$ %[1]s query %[2]s rewards --owner fury15qdefkmwswysgg4qxgqpqr35k3m49pkx2jdfnw
			$ %[1]s query %[2]s rewards --type jinx
			$ %[1]s query %[2]s rewards --type usdf-minting
			$ %[1]s query %[2]s rewards --type delegator
			$ %[1]s query %[2]s rewards --type swap
			$ %[1]s query %[2]s rewards --type savings
			$ %[1]s query %[2]s rewards --type earn
			$ %[1]s query %[2]s rewards --type jinx --owner fury15qdefkmwswysgg4qxgqpqr35k3m49pkx2jdfnw
			$ %[1]s query %[2]s rewards --type jinx --unsynced
			`,
				version.AppName, types.ModuleName)),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			page, _ := cmd.Flags().GetInt(flags.FlagPage)
			limit, _ := cmd.Flags().GetInt(flags.FlagLimit)
			strOwner, _ := cmd.Flags().GetString(flagOwner)
			strType, _ := cmd.Flags().GetString(flagType)
			boolUnsynced, _ := cmd.Flags().GetBool(flagUnsynced)

			// Prepare params for querier
			var owner sdk.AccAddress
			if strOwner != "" {
				if owner, err = sdk.AccAddressFromBech32(strOwner); err != nil {
					return err
				}
			}

			switch strings.ToLower(strType) {
			case typeJinx:
				params := types.NewQueryRewardsParams(page, limit, owner, boolUnsynced)
				claims, err := executeJinxRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				return cliCtx.PrintObjectLegacy(claims)
			case typeUSDFMinting:
				params := types.NewQueryRewardsParams(page, limit, owner, boolUnsynced)
				claims, err := executeUSDFMintingRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				return cliCtx.PrintObjectLegacy(claims)
			case typeDelegator:
				params := types.NewQueryRewardsParams(page, limit, owner, boolUnsynced)
				claims, err := executeDelegatorRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				return cliCtx.PrintObjectLegacy(claims)
			case typeSwap:
				params := types.NewQueryRewardsParams(page, limit, owner, boolUnsynced)
				claims, err := executeSwapRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				return cliCtx.PrintObjectLegacy(claims)
			case typeSavings:
				params := types.NewQueryRewardsParams(page, limit, owner, boolUnsynced)
				claims, err := executeSavingsRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				return cliCtx.PrintObjectLegacy(claims)
			case typeEarn:
				params := types.NewQueryRewardsParams(page, limit, owner, boolUnsynced)
				claims, err := executeEarnRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				return cliCtx.PrintObjectLegacy(claims)
			default:
				params := types.NewQueryRewardsParams(page, limit, owner, boolUnsynced)

				jinxClaims, err := executeJinxRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				usdfMintingClaims, err := executeUSDFMintingRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				delegatorClaims, err := executeDelegatorRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				swapClaims, err := executeSwapRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				savingsClaims, err := executeSavingsRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}
				earnClaims, err := executeEarnRewardsQuery(cliCtx, params)
				if err != nil {
					return err
				}

				if len(jinxClaims) > 0 {
					if err := cliCtx.PrintObjectLegacy(jinxClaims); err != nil {
						return err
					}
				}
				if len(usdfMintingClaims) > 0 {
					if err := cliCtx.PrintObjectLegacy(usdfMintingClaims); err != nil {
						return err
					}
				}
				if len(delegatorClaims) > 0 {
					if err := cliCtx.PrintObjectLegacy(delegatorClaims); err != nil {
						return err
					}
				}
				if len(swapClaims) > 0 {
					if err := cliCtx.PrintObjectLegacy(swapClaims); err != nil {
						return err
					}
				}
				if len(savingsClaims) > 0 {
					if err := cliCtx.PrintObjectLegacy(savingsClaims); err != nil {
						return err
					}
				}
				if len(earnClaims) > 0 {
					if err := cliCtx.PrintObjectLegacy(earnClaims); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().String(flagOwner, "", "(optional) filter by owner address")
	cmd.Flags().String(flagType, "", fmt.Sprintf("(optional) filter by a reward type: %s", strings.Join(rewardTypes, "|")))
	cmd.Flags().Bool(flagUnsynced, false, "(optional) get unsynced claims")
	cmd.Flags().Int(flags.FlagPage, 1, "pagination page rewards of to to query for")
	cmd.Flags().Int(flags.FlagLimit, 100, "pagination limit of rewards to query for")
	return cmd
}

func queryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "get the incentive module parameters",
		Long:  "Get the current global incentive module parameters.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Query
			route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryGetParams)
			res, height, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}
			cliCtx = cliCtx.WithHeight(height)

			// Decode and print results
			var params types.Params
			if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &params); err != nil {
				return fmt.Errorf("failed to unmarshal params: %w", err)
			}
			return cliCtx.PrintObjectLegacy(params)
		},
	}
}

func queryRewardFactorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-factors",
		Short: "get current global reward factors",
		Long:  `Get current global reward factors for all reward types.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Execute query
			route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryGetRewardFactors)
			res, height, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}
			cliCtx = cliCtx.WithHeight(height)

			// Decode and print results
			var response types.QueryGetRewardFactorsResponse
			if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &response); err != nil {
				return fmt.Errorf("failed to unmarshal reward factors: %w", err)
			}
			return cliCtx.PrintObjectLegacy(response)
		},
	}
	cmd.Flags().String(flagDenom, "", "(optional) filter reward factors by denom")
	return cmd
}

func executeJinxRewardsQuery(cliCtx client.Context, params types.QueryRewardsParams) (types.JinxLiquidityProviderClaims, error) {
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return types.JinxLiquidityProviderClaims{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryGetJinxRewards)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return types.JinxLiquidityProviderClaims{}, err
	}

	cliCtx = cliCtx.WithHeight(height)

	var claims types.JinxLiquidityProviderClaims
	if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &claims); err != nil {
		return types.JinxLiquidityProviderClaims{}, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return claims, nil
}

func executeUSDFMintingRewardsQuery(cliCtx client.Context, params types.QueryRewardsParams) (types.USDFMintingClaims, error) {
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return types.USDFMintingClaims{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryGetUSDFMintingRewards)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return types.USDFMintingClaims{}, err
	}

	cliCtx = cliCtx.WithHeight(height)

	var claims types.USDFMintingClaims
	if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &claims); err != nil {
		return types.USDFMintingClaims{}, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return claims, nil
}

func executeDelegatorRewardsQuery(cliCtx client.Context, params types.QueryRewardsParams) (types.DelegatorClaims, error) {
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return types.DelegatorClaims{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryGetDelegatorRewards)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return types.DelegatorClaims{}, err
	}

	cliCtx = cliCtx.WithHeight(height)

	var claims types.DelegatorClaims
	if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &claims); err != nil {
		return types.DelegatorClaims{}, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return claims, nil
}

func executeSwapRewardsQuery(cliCtx client.Context, params types.QueryRewardsParams) (types.SwapClaims, error) {
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return types.SwapClaims{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryGetSwapRewards)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return types.SwapClaims{}, err
	}

	cliCtx = cliCtx.WithHeight(height)

	var claims types.SwapClaims
	if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &claims); err != nil {
		return types.SwapClaims{}, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return claims, nil
}

func executeSavingsRewardsQuery(cliCtx client.Context, params types.QueryRewardsParams) (types.SavingsClaims, error) {
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return types.SavingsClaims{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryGetSavingsRewards)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return types.SavingsClaims{}, err
	}

	cliCtx = cliCtx.WithHeight(height)

	var claims types.SavingsClaims
	if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &claims); err != nil {
		return types.SavingsClaims{}, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return claims, nil
}

func executeEarnRewardsQuery(cliCtx client.Context, params types.QueryRewardsParams) (types.EarnClaims, error) {
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return types.EarnClaims{}, err
	}

	route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryGetEarnRewards)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return types.EarnClaims{}, err
	}

	cliCtx = cliCtx.WithHeight(height)

	var claims types.EarnClaims
	if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &claims); err != nil {
		return types.EarnClaims{}, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return claims, nil
}
