package main

import (
	"testing"
)

func TestBuyNewCurrency(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"usd": 10, "eur": 5}
	provideLiquidity(&newBankBalance, newPair)

        myBalance := map[string]float64{"eur": 200, "jjt": 0, "gc":100}
	myToffer := TradeOffer{
		wantAmount:   5,
		wantCurrency: "usd",
		giveCurrency: "eur",
	}
	err := trade(myToffer, &newBankBalance, &myBalance)

	if err != nil {
		t.Fatalf(`Should not error, but got %v.`, err)
	}
	if myBalance["usd"] != 5 {
		t.Fatalf(`myBalance["usd"] should be 5, but got %v.`, myBalance["usd"])
	}
}

func TestReverseExchange(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"usd": 500, "eur": 600}
	provideLiquidity(&newBankBalance, newPair)

	tax = 0
        myBalance := map[string]float64{"eur": 100, "usd": 100}

	myToffer := TradeOffer{
		wantAmount:   5,
		wantCurrency: "usd",
		giveCurrency: "eur",
	}
	err := trade(myToffer, &newBankBalance, &myBalance)
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}

	myToffer = TradeOffer{
		wantAmount:   5,
		wantCurrency: "usd",
		giveCurrency: "eur",
	}
	err = trade(myToffer, &newBankBalance, &myBalance)
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}

	wantAmountHere := 100 - myBalance["eur"]

	myToffer = TradeOffer{
		wantAmount:   wantAmountHere,
		wantCurrency: "eur",
		giveCurrency: "usd",
	}
	err = trade(myToffer, &newBankBalance, &myBalance)
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
	provideLiquidity(&newBankBalance, newPair)

	tax = 0.1
        myBalance := map[string]float64{"eur": 100, "usd": 100}

	myToffer := TradeOffer{
		wantAmount:   10,
		wantCurrency: "eur",
		giveCurrency: "usd",
	}
	err := trade(myToffer, &newBankBalance, &myBalance)
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}

	myToffer = TradeOffer{
		wantAmount:   99.9-myBalance["usd"],
		wantCurrency: "usd",
		giveCurrency: "eur",
	}
	err = trade(myToffer, &newBankBalance, &myBalance)
	if (err != nil) {
		t.Fatalf(`Should not error, but got %v.`, err)
	}

	if (myBalance["eur"] >= 100 || myBalance["usd"] >= 100) {
		t.Fatalf(`Values of myBalance should be less than 100, but got %v.`, myBalance)
	}
	tax = 0
}
