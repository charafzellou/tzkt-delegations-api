# Exercise: Tezos Delegation Service

In this exercise, you will build a service that gathers new [delegations](https://opentezos.com/baking/delegating/) made on the Tezos protocol and exposes them through a public API. 

## Requirements:

- The service will poll the new delegations from this Tzkt API endpoint: https://api.tzkt.io/#operation/Operations_GetDelegations
- For each delegation, save the following information: sender's address, timestamp, amount, and block.
- Expose the collected data through a public API at the endpoint `/xtz/delegations`.
    - The expected response format is:
    
    ```jsx
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
    
    - The sender’s address is the delegator.
    - The delegations must be listed most recent first.
    - The endpoint takes one optional query parameter `year` , which is specified in the format YYYY and will result in the data being filtered for that year only.
- Ensure the service is production-ready, considering factors like performance, scalability, error handling, and reliability.