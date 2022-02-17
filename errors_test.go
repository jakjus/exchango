package main

import (
	"testing"
)

func TestUserNoMoney(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"usd": 10, "eur": 5}
	provideLiquidity(&newBankBalance, newPair)

        myBalance := map[string]float64{"eur": 1, "jjt": 0, "gc":100}
	myToffer := TradeOffer{
		wantAmount:   5,
		wantCurrency: "usd",
		giveCurrency: "eur",
	}
	err := trade(myToffer, &newBankBalance, &myBalance)

	if err == nil {
		t.Fatalf(`Should error, but got "err == nil".`)
	}
}

func TestBankNoMoney(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"usd": 10, "eur": 5}
	provideLiquidity(&newBankBalance, newPair)

        myBalance := map[string]float64{"eur": 100, "jjt": 0, "gc":100}
	myToffer := TradeOffer{
		wantAmount:   15,
		wantCurrency: "usd",
		giveCurrency: "eur",
	}
	err := trade(myToffer, &newBankBalance, &myBalance)

	if err == nil {
		t.Fatalf(`Should error, but got "err == nil".`)
	}
}
