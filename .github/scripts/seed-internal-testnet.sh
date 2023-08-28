#!/bin/bash
set -ex

# configure fury binary to talk to the desired chain endpoint
fury config node "${CHAIN_API_URL}"
fury config chain-id "${CHAIN_ID}"

# use the test keyring to allow scriptng key generation
fury config keyring-backend test

# wait for transactions to be committed per CLI command
fury config broadcast-mode block

# setup dev wallet
echo "${DEV_WALLET_MNEMONIC}" | fury keys add --recover dev-wallet
DEV_TEST_WALLET_ADDRESS="0x17D321E9F64A4FC3640EA9B71AFD16FACEAE9878"
WEBAPP_E2E_WHALE_ADDRESS="0x5FEA9D19024A52B034F310F03B1956E8F8EB5E60"

# setup fury ethereum compatible account for deploying
# erc20 contracts to the fury chain
echo "uphold elbow off winter pill grunt reward drive joke common slab purpose sheriff custom spend conduct inch kidney multiply food spray bunker eager oven" | fury keys add --recover --eth dev-erc20-deployer-wallet

# fund evm-contract-deployer account (using issuance)
fury tx issuance issue 200000000ufury fury1van3znl6597xgwwh46jgquutnqkwvwszjg04fz --from dev-wallet --gas-prices 0.5ufury -y

# deploy and fund USDC ERC20 contract
MULTICHAIN_USDC_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "USD Coin" USDC 6)
MULTICHAIN_USDC_CONTRACT_ADDRESS=${MULTICHAIN_USDC_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$MULTICHAIN_USDC_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 1000000000000

# # deploy and fund wFury ERC20 contract
wFURY_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "Wrapped Fury" wFury 6)
wFURY_CONTRACT_ADDRESS=${wFURY_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$wFURY_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 1000000000000

# deploy and fund BNB contract
BNB_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "Binance" BNB 8)
BNB_CONTRACT_ADDRESS=${BNB_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$BNB_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 1000000000000

# deploy and fund USDT contract
MULTICHAIN_USDT_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "Multichain USDT" USDT 6)
MULTICHAIN_USDT_CONTRACT_ADDRESS=${MULTICHAIN_USDT_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$MULTICHAIN_USDT_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 1000000000000

# deploy and fund DAI contract
DAI_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "DAI" DAI 18)
DAI_CONTRACT_ADDRESS=${DAI_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$DAI_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 1000000000000

# deploy and fund axlwBTC ERC20 contract
AXL_WBTC_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "Wrapped BTC" BTC 8)
AXL_WBTC_CONTRACT_ADDRESS=${AXL_WBTC_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$AXL_WBTC_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 100000000000000

# deploy and fund wETH ERC20 contract
wETH_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "Wrapped wETH" ETH 18)
wETH_CONTRACT_ADDRESS=${wETH_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$wETH_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 1000000000000000000000

# deploy and fund axlUSDC ERC20 contract
AXL_USDC_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "USD Coin" USDC 6)
AXL_USDC_CONTRACT_ADDRESS=${AXL_USDC_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$AXL_USDC_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 1000000000000

# deploy and fund Multichain wBTC ERC20 contract
MULTICHAIN_wBTC_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "Wrapped BTC" BTC 8)
MULTICHAIN_wBTC_CONTRACT_ADDRESS=${MULTICHAIN_wBTC_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$MULTICHAIN_wBTC_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 100000000000000

# deploy and fund Tether USDT ERC20 contract
TETHER_USDT_CONTRACT_DEPLOY=$(npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" deploy-erc20 "USD₮" USD₮ 6)
TETHER_USDT_CONTRACT_ADDRESS=${TETHER_USDT_CONTRACT_DEPLOY: -42}
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$TETHER_USDT_CONTRACT_ADDRESS" 0x7B5FE22B5446F7C62EA27B8BD71CEF94E03F3DF2 1000000000000

# seed some evm wallets
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$AXL_WBTC_CONTRACT_ADDRESS" "$DEV_TEST_WALLET_ADDRESS" 10000000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$MULTICHAIN_wBTC_CONTRACT_ADDRESS" "$DEV_TEST_WALLET_ADDRESS" 10000000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$MULTICHAIN_USDC_CONTRACT_ADDRESS" "$DEV_TEST_WALLET_ADDRESS" 100000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$wETH_CONTRACT_ADDRESS" "$DEV_TEST_WALLET_ADDRESS" 1000000000000000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$AXL_USDC_CONTRACT_ADDRESS" "$DEV_TEST_WALLET_ADDRESS" 100000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$MULTICHAIN_USDT_CONTRACT_ADDRESS" "$DEV_TEST_WALLET_ADDRESS" 100000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$TETHER_USDT_CONTRACT_ADDRESS" "$DEV_TEST_WALLET_ADDRESS" 1000000000000
# seed webapp E2E whale account
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$AXL_WBTC_CONTRACT_ADDRESS" "$WEBAPP_E2E_WHALE_ADDRESS" 100000000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$MULTICHAIN_wBTC_CONTRACT_ADDRESS" "$WEBAPP_E2E_WHALE_ADDRESS" 10000000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$MULTICHAIN_USDC_CONTRACT_ADDRESS" "$WEBAPP_E2E_WHALE_ADDRESS" 1000000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$wETH_CONTRACT_ADDRESS" "$WEBAPP_E2E_WHALE_ADDRESS" 10000000000000000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$AXL_USDC_CONTRACT_ADDRESS" "$WEBAPP_E2E_WHALE_ADDRESS" 10000000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$MULTICHAIN_USDT_CONTRACT_ADDRESS" "$WEBAPP_E2E_WHALE_ADDRESS" 1000000000000
npx hardhat --network "${ERC20_DEPLOYER_NETWORK_NAME}" mint-erc20 "$TETHER_USDT_CONTRACT_ADDRESS" "$WEBAPP_E2E_WHALE_ADDRESS" 1000000000000

# give dev-wallet enough delegation power to pass proposals by itself

# issue fury to dev wallet for delegating to each validator
fury tx issuance issue 6000000000ufury fury1vlpsrmdyuywvaqrv7rx6xga224sqfwz3fyfhwq \
  --from dev-wallet --gas-prices 0.5ufury -y

# parse space seperated list of validators
# into bash array
read -r -a GENESIS_VALIDATOR_ADDRESS_ARRAY <<<"$GENESIS_VALIDATOR_ADDRESSES"

# delegate 300FURY to each validator
for validator in "${GENESIS_VALIDATOR_ADDRESS_ARRAY[@]}"; do
  fury tx staking delegate "${validator}" 300000000ufury --from dev-wallet --gas-prices 0.5ufury -y
done

# create a text proposal
fury tx gov submit-legacy-proposal --deposit 1000000000ufury --type "Text" --title "Example Proposal" --description "This is an example proposal" --gas auto --gas-adjustment 1.2 --from dev-wallet --gas-prices 0.01ufury -y

# setup god's wallet
echo "${FURY_TESTNET_GOD_MNEMONIC}" | fury keys add --recover god

# create template string for the proposal we want to enact
# https://percosis-labs.atlassian.net/wiki/spaces/ENG/pages/1228537857/Submitting+Governance+Proposals+WIP
PARAM_CHANGE_PROP_TEMPLATE=$(
  cat <<'END_HEREDOC'
{
    "@type": "/cosmos.params.v1beta1.ParameterChangeProposal",
    "title": "Set Initial ERC-20 Contracts",
    "description": "Set Initial ERC-20 Contracts",
    "changes": [
        {
            "subspace": "evmutil",
            "key": "EnabledConversionPairs",
            "value": "[{\"fury_erc20_address\":\"MULTICHAIN_USDC_CONTRACT_ADDRESS\",\"denom\":\"erc20/multichain/usdc\"},{\"fury_erc20_address\":\"MULTICHAIN_USDT_CONTRACT_ADDRESS\",\"denom\":\"erc20/multichain/usdt\"},{\"fury_erc20_address\":\"MULTICHAIN_wBTC_CONTRACT_ADDRESS\",\"denom\":\"erc20/multichain/wbtc\"},{\"fury_erc20_address\":\"AXL_USDC_CONTRACT_ADDRESS\",\"denom\":\"erc20/axelar/usdc\"},{\"fury_erc20_address\":\"AXL_WBTC_CONTRACT_ADDRESS\",\"denom\":\"erc20/axelar/wbtc\"},{\"fury_erc20_address\":\"wETH_CONTRACT_ADDRESS\",\"denom\":\"erc20/axelar/eth\"},{\"fury_erc20_address\":\"TETHER_USDT_CONTRACT_ADDRESS\",\"denom\":\"erc20/tether/usdt\"}]"
        }
    ]
}
END_HEREDOC
)

# substitute freshly deployed contract addresses
finalProposal=$PARAM_CHANGE_PROP_TEMPLATE

finalProposal="${finalProposal/MULTICHAIN_USDC_CONTRACT_ADDRESS/$MULTICHAIN_USDC_CONTRACT_ADDRESS}"
finalProposal="${finalProposal/MULTICHAIN_USDT_CONTRACT_ADDRESS/$MULTICHAIN_USDT_CONTRACT_ADDRESS}"
finalProposal="${finalProposal/MULTICHAIN_wBTC_CONTRACT_ADDRESS/$MULTICHAIN_wBTC_CONTRACT_ADDRESS}"
finalProposal="${finalProposal/AXL_USDC_CONTRACT_ADDRESS/$AXL_USDC_CONTRACT_ADDRESS}"
finalProposal="${finalProposal/AXL_WBTC_CONTRACT_ADDRESS/$AXL_WBTC_CONTRACT_ADDRESS}"
finalProposal="${finalProposal/wETH_CONTRACT_ADDRESS/$wETH_CONTRACT_ADDRESS}"
finalProposal="${finalProposal/TETHER_USDT_CONTRACT_ADDRESS/$TETHER_USDT_CONTRACT_ADDRESS}"

# create unique proposal filename
proposalFileName="$(date +%s)-proposal.json"
touch $proposalFileName

# save proposal as file to disk
echo "$finalProposal" >$proposalFileName

# snapshot original module params
originalEvmUtilParams=$(curl https://api.app.internal.testnet.us-east.production.fury.io/fury/evmutil/v1beta1/params)
printf "original evm util module params\n %s" , "$originalEvmUtilParams"

# change the params of the chain like a god - make it so 🖖🏽
# make sure to update god committee member permissions for the module
# and params being updated (see below for example)
# https://github.com/Percosis-Labs/fury/pull/1556/files#diff-0bd6043650c708661f37bbe6fa5b29b52149e0ec0069103c3954168fc9f12612R900-R903
# committee 1 is the stability committee. on internal testnet, this has only one member.
fury tx committee submit-proposal 1 "$proposalFileName" --gas 2000000 --gas-prices 0.01ufury --from god -y

# vote on the proposal. this assumes no other committee proposal has ever been submitted (id=1)
fury tx committee vote 1 yes --gas 2000000 --gas-prices 0.01ufury --from god -y

# fetch current module params
updatedEvmUtilParams=$(curl https://api.app.internal.testnet.us-east.production.fury.io/fury/evmutil/v1beta1/params)
printf "updated evm util module params\n %s" , "$updatedEvmUtilParams"

# if adding more cosmos coins -> er20s, ensure that the deployment order below remains the same.
# convert 1 JINX to an erc20. doing this ensures the contract is deployed.
fury tx evmutil convert-cosmos-coin-to-erc20 \
  "$DEV_TEST_WALLET_ADDRESS" \
  1000000jinx \
  --from dev-wallet --gas 2000000 --gas-prices 0.001ufury -y
