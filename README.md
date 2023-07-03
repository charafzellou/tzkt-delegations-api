# TZKT DELEGATIONS API

## Original assignement :

Build a service that :
- gathers new [delegations](https://opentezos.com/baking/delegating/) made on the Tezos protocol,
- exposes them through a public API. 

The original assignement document can be viewed [here](./docs/ASSIGNEMENT.md).

## Usage :
### Using `Make` :
Using `Make`, you can execute the following commands :

```bash
make install
make start
```


### Using Go locally :
- Using a local install of Go, you can execute the following commands :

```bash
source .env
cd src
go build . -o app
./app
```

## Checklist :

<!-- - [ ] Mercury -->
<!-- - [x] Venus -->
- [ ] The service will poll the new delegations from this Tzkt API endpoint: https://api.tzkt.io/#operation/Operations_GetDelegations
- [ ] For each delegation, save the following information: sender's address, timestamp, amount, and block.
- [ ] Expose the collected data through a public API at the endpoint `/xtz/delegations`.
    - [ ] The expected response format is:
    
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
    
    - [ ] The sender’s address is the delegator.
    - [ ] The delegations must be listed most recent first.
    - [ ] The endpoint takes one optional query parameter `year` , which is specified in the format `YYYY` and will result in the data being filtered for that year only.
- [ ] Ensure the service is production-ready, considering factors like performance, scalability, error handling, and reliability.