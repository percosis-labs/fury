# Testnet-4000

Resources and examples for running and interacting with fury-testnet-4000

## Rest server requests

### Setup

Before making a request, query account information for the signing account. Note the 'accountnumber' and 'sequence' fields, we'll need them later in order to send our request:

```bash
    fucli q auth account $(fucli keys show accB -a)
```

If testing locally, start the Fury rest server:

```bash
    fucli rest-server
```

Now we'll create an unsigned request, sign it, and broadcast it to the Fury blockchain via the rest server. Note that if you're using the mainnet or testnet, the host IP address will need to be updated to point at an active rest server instead of http://127.0.0.1.

### Create CDP example request

Format the base request in create-cdp.json. You'll need to update the 'from', 'chain-id', 'account_number', 'sequence', and 'gas' as appropriate. Then, populate the CDP creation request's params 'owner', 'collateral', and 'principal'. An example formatted request can be found in `example-create-cdp.json`.

```bash
    # Create an unsigned request
    curl -H "Content-Type: application/json" -X PUT -d @./contrib/testnet-4000/rest_examples/cdp/create-cdp.json http://127.0.0.1:1317/cdp | jq > ./contrib/testnet-4000/rest_examples/cdp/create-cdp-unsigned.json

    # Sign the request
    fucli tx sign ./contrib/testnet-4000/rest_examples/cdp/create-cdp-unsigned.json --from accB --offline --chain-id testing --sequence 1 --account-number 2 | jq  > ./contrib/testnet-4000/rest_examples/cdp/broadcast-create-cdp.json

    # Broadcast the request
    fucli tx broadcast ./contrib/testnet-4000/rest_examples/cdp/broadcast-create-cdp.json
```

Congratulations, you've just created a CDP on Fury using the rest server!

### Post market price example request

Note that only market oracles can post prices, other senders will have their transactions rejected by Fury.

Format the base request in post-price.json. You'll need to update the 'from', 'chain-id', 'account_number', 'sequence', and 'gas' as appropriate. Then, populate the post price request's params 'from', 'market_id', 'price', and 'expiry'. An example formatted request can be found in `example-post-price.json`.

```bash
    # Create an unsigned request
	curl -H "Content-Type: application/json" -X PUT -d @./contrib/testnet-4000/rest_examples/pricefeed/post-price.json http://127.0.0.1:1317/pricefeed/postprice | jq > ./contrib/testnet-4000/rest_examples/pricefeed/post-price-unsigned.json


    # Sign the request
    fucli tx sign ./contrib/testnet-4000/rest_examples/pricefeed/post-price-unsigned.json --from validator --offline --chain-id testing --sequence 96 --account-number 0 | jq > ./contrib/testnet-4000/rest_examples/pricefeed/broadcast-post-price.json

    # Broadcast the request
    fucli tx broadcast ./contrib/testnet-4000/rest_examples/pricefeed/broadcast-post-price.json
```

Congratulations, you've just posted a current market price on Fury using the rest server!

## Governance proposals

Example governance proposals are located in `/proposal_examples`.

## Dredd automated REST API testing

To run the dredd tests to hit the endpoints and test the api run the following command from the `fury` folder:

```bash
    make test_dredd
```