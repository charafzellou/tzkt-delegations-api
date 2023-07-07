# TZKT DELEGATIONS API

## Original assignement :

Build a service that :
- gathers new [delegations](https://opentezos.com/baking/delegating/) made on the Tezos protocol,
- exposes them through a public API. 

The original assignement document can be viewed [here](./docs/ASSIGNEMENT.md).

## Usage :
### Using `Docker-compose` :
Using `Docker-compose`, you can execute the following commands :

```bash
docker-compose up -d
```


### Using Go locally :
- Using a local install of Go, you can execute the following commands :

```bash
cp .env.dist .env
```

- Set up your Environement Variables, then :

```bash
cd app/indexer
go build . -o indexer
./indexer
```
```bash
cd app/api
go build . -o api
./api
```

## Checklist :

- [X] The service will poll the new delegations from this Tzkt API endpoint: https://api.tzkt.io/#operation/Operations_GetDelegations
- [X] For each delegation, save the following information: sender's address, timestamp, amount, and block.
- [X] Expose the collected data through a public API at the endpoint `/xtz/delegations`.
    - [X] The expected response format is:
    
        ```json
        {
        "data": [ 
            {
                "timestamp": "2022-05-05T06:29:14Z",
                "amount": "125896",
                "delegator": "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
                "block": "2338084"
            },
            {
                "timestamp": "2021-05-07T14:48:07Z",
                "amount": "9856354",
                "delegator": "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
                "block": "1461334"
            }
        ],
        }
        ```
    
    - [X] The sender’s address is the delegator.
    - [X] The delegations must be listed most recent first.
    - [X] The endpoint takes one optional query parameter `year` , which is specified in the format `YYYY` and will result in the data being filtered for that year only.
- [ ] Ensure the service is production-ready, considering factors like performance, scalability, error handling, and reliability.
