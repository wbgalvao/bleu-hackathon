# Bleutrade Golang + Blockchain Hackathon - Challenge 1
This project contains the solution of the first problem addressed at the Bleutrade Golang + Blockchain Hackathon event. The challenge consisted in creating a HTTP client using the Go programming language. This client would communicate with the [Bleutrade API](https://classic.bleutrade.com/help/API) and make operations for a given user.

The HTTP client is presented in the `client` package. We (me + [@ramonhpr](https://github.com/ramonhpr)) also built a Telegram Chatbot which uses the HTTP client and enables the enduser to use the [Bleutrade Platform](https://malta.bleutrade.com/) via Telegram.

## Chatobot instructions
To access the Telegram chatbot, use this  [link](t.me/donamariabot).

### Commands
* `/setup [APIKey] [APISecret]`: Enables chatbot to operate in the Bleutrade's platform. This command should be executed prior to other operations, otherwise the chatobot won't haver permission to access user information. The `APIKey` and `APISecret` must be retrieved from the cryptoexchange website.
* `/saldo [Currency]`: Retrieves the total amount of a given currency for the setup user wallet.
* `/saque [Currency] [Quantity] [DestinationWallet]`: Makes a transaction of a given currency, in a given quantity to a third party cryptocurrency wallet (also given as an input).
* `/selllimit [Market] [Quantity]`: Sells a given quantity of a given market (crypto currency).
* `/buylimit [Market] [Quantity]`: Buys a given quantity of a given market (crypto currency).
* `/wallet`: Retrieves the address of the setup user cryto currency wallet.

## build instructions with Docker

build the image:
`$ docker build . -t donamaria`

run container:
`$ docker container run -d donamaria`

