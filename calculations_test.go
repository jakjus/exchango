package exchango

import (
	"testing"
)

func TestBuyNewCurrency(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"usd": 10, "eur": 5}
	ProvideLiquidity(&newBankBalance, newPair)

        myBalance := map[string]float64{"eur": 200, "jjt": 0, "gc":100}
	myToffer := TradeOffer{
		WantAmount:   5,
		WantCurrency: "usd",
		GiveCurrency: "eur",
	}

	tl, err := InitTrade(TradeDetails{myToffer, &newBankBalance, &myBalance})
	if err != nil {
		t.Fatalf(`Should not error, but got %v.`, err)
	}
	err = ExecuteTrade(tl)
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}

	if myBalance["usd"] != 5 {
		t.Fatalf(`myBalance["usd"] should be 5, but got %v.`, myBalance["usd"])
	}
}

func TestStagedExchangeEquilibrium(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"usd": 500, "eur": 600}
	ProvideLiquidity(&newBankBalance, newPair)

	Tax = 0
        myBalance := map[string]float64{"eur": 100, "usd": 100}

	myToffer := TradeOffer{
		WantAmount:   5,
		WantCurrency: "usd",
		GiveCurrency: "eur",
	}
	tl, err := InitTrade(TradeDetails{myToffer, &newBankBalance, &myBalance})
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}
	err = ExecuteTrade(tl)
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}

	myToffer = TradeOffer{
		WantAmount:   5,
		WantCurrency: "usd",
		GiveCurrency: "eur",
	}
	tl, err = InitTrade(TradeDetails{myToffer, &newBankBalance, &myBalance})
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}
	err = ExecuteTrade(tl)
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}


	WantAmountHere := 100 - myBalance["eur"]

	myToffer = TradeOffer{
		WantAmount:   WantAmountHere,
		WantCurrency: "eur",
		GiveCurrency: "usd",
	}
	tl, err = InitTrade(TradeDetails{myToffer, &newBankBalance, &myBalance})
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}
	err = ExecuteTrade(tl)
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}

	if (myBalance["eur"] != 100 || myBalance["usd"] != 100) {
		t.Fatalf(`myBalance should be map["eur":100, "usd":100], but got %v.`, myBalance)
	}
}

func TestReverseExchangeWithTax(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"usd": 500, "eur": 600}
	ProvideLiquidity(&newBankBalance, newPair)

	Tax = 0.1
        myBalance := map[string]float64{"eur": 100, "usd": 100}

	myToffer := TradeOffer{
		WantAmount:   10,
		WantCurrency: "eur",
		GiveCurrency: "usd",
	}
	tl, err := InitTrade(TradeDetails{myToffer, &newBankBalance, &myBalance})
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}
	err = ExecuteTrade(tl)
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}


	myToffer = TradeOffer{
		WantAmount:   99.9-myBalance["usd"],
		WantCurrency: "usd",
		GiveCurrency: "eur",
	}
	tl, err = InitTrade(TradeDetails{myToffer, &newBankBalance, &myBalance})
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}
	err = ExecuteTrade(tl)

	if (myBalance["eur"] >= 100 || myBalance["usd"] >= 100) {
		t.Fatalf(`Values of myBalance should be less than 100, but got %v.`, myBalance)
	}
	Tax = 0
}
