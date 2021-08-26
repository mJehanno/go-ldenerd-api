# Go-ldenerd-api

This project is a simple api made in go with [fiber](https://gofiber.io/).
It use an [eventstore](https://www.eventstore.com/) database to manage incoming transactions and [arangoDb](https://www.arangodb.com/) to store cached value.

Thoses cached values are refreshed at every new incoming trnasaction so datas in arango are always up-to-date.

## Endpoint

The fiber api exposes few endpoint :

Verb | Endpoint | Action
---|---|---
Get| `/api/gold` | Get the current amount of Gold converted in Gold coins
Get | `/api/gold/detail`| Get the current amount of Gold with detail amount for every coins (Copp,er Silver, Electrum, Gold, Platinum)
Get | `/api/transactions/history` | Get the whole history of transactions
Post | `/api/transactions`| Add a new transaction. This route is supposed to be protected.
