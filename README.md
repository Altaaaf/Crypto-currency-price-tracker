# Crypto-currency-price-tracker
Simple Golang tool for tracking the prices of multiple cryptocurrency tokens and receiving alerts when the prices go above a given threshold. The tool uses the Coinbase API to fetch the latest prices and send a notification when the price of a token goes above the threshold. The tool is built using multithreading for efficient processing of multiple tokens.

## Installation
To get started, follow these steps:

1. Clone the repository to your local machine.
``` {.sourceCode}
git clone https://github.com/Altaaaf/Crypto-currency-price-tracker.git
```
2. Run command to start Go application
``` {.sourceCode}
go run main.go
```

## Usage

You first need to create a configuration file named config.json. The configuration file should contain a list of tokens that you want to track, along with their symbols, threshold values. Here's an example configuration file:
```
{
    "tokens": [
        {
            "name": "Bitcoin",
            "symbol": "BTC",
            "threshold": 30000
        },
        {
            "name": "Ethereum",
            "symbol": "ETH",
            "threshold": 1200
        },
        {
            "name": "Litecoin",
            "symbol": "LTC",
            "threshold": 110
        }
    ]
}

```
