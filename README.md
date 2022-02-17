# ExchanGo
Digital currencies exchange in style of "one bank, many users", based on liquidity mechanism. Supports multiple currency pairs.

## Install
```
go get github.com/jakjus/exchango
```
or just import the repository into your Go program.

## Usage
```go
package main

import (
    "log"
    "github.com/jakjus/exchango"
)

func main() {
	// Initialize bank and user
	newBankBalance := make([]map[string]float64, 0)
        myBalance := map[string]float64{"eur": 200, "jjt": 2, "coins":100}

	// Provide liquidity (initial money for the bank)
	newPair := map[string]float64{"usd": 201, "eur": 53}
	exchango.ProvideLiquidity(&newBankBalance, newPair)

	// Initialize user's offer
	myToffer := exchango.TradeOffer{
		WantAmount:   5,
		WantCurrency: "usd",
		GiveCurrency: "eur",
	}

	// Perform trade
	tl, err := exchango.InitTrade(exchango.TradeDetails{myToffer, &newBankBalance, &myBalance})
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	err = exchango.ExecuteTrade(tl)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	log.Println("Finished.")
}
```

## Additional notes
### Providing liquidity
Providing liquidity is a process of providing bank with currencies, that can be further exchanged with users. 

[Note] If you want to allow to provide liquidity for *users*, you should verify that the money was in-fact transferred (physically/digitally) and then run `func ProvideLiquidity`. Providing liquidity by users has no use in the current version, as there is no *protocol* to pay out dividends to liquidity provider (all goes to bank).

### Setting exchange rate
The exchange rate is based on available currency ratio, so if you initially want **1 gold** to cost **2 silver**, initiate liquidity with {"gold": 1000, "silver": 2000}. The more currency in exchange pair, the less volatile exchange rate. 

Liquidity of each pair increases with every transaction through `tax` (default: `0.01`). It ensures money gains for the bank and more stable exchange with time. Therefore, we recommend setting the lowest possible liquidity for exchange pair at the start, so that it is more fluid and adjusts to the market faster.

### Exchange dry-up
Once liquidity is set, exchange pairs cannot dry-up (can run indefinetely). The closer one currency in pair is to 0 (`lim(amt_x) -> 0`), the more valued it is (`lim(price_x) -> inf`).

### InitTrade and ExecuteTrade
Trade is split between two functions for two different needs:
  - `InitTrade` - gets initial exchange rate from bank, based on your previously sent proposition
  - `ExecuteTrade` - confirms trade. Errors, if the exchange rate has changed more than `AllowChange`

In an environment with many users, exchange ratio may change before accepting. This architecture is split to support possible REST API architectures built on top of the project.
