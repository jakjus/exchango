package main

import (
	"testing"
)

func TestProvideLiquidity(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"coins": 15000, "jjt": 10}
	provideLiquidity(&newBankBalance, newPair)
	for key := range newPair {
		if newBankBalance[0][key] != newPair[key]{
			t.Fatalf(`BankBalance should be [%v], but got %v.`, newPair, newBankBalance)
		}
	}
}

func TestProvideLiquidityAdd(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"abc": 10, "def": 5}
	provideLiquidity(&newBankBalance, newPair)
	newPair2 := map[string]float64{"abc": 20, "def": 10}
	provideLiquidity(&newBankBalance, newPair2)

	want := map[string]float64{"abc": 30, "def": 15}
	for key := range want {
		if newBankBalance[0][key] != want[key] {
			t.Fatalf(`newBankBalance[0] should be %v, but got %v.`, want, newBankBalance)
		}
	}
}

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
		t.Fatalf(`Should be no error, but got %v.`, err)
	}
	if myBalance["usd"] != 5 {
		t.Fatalf(`myBalance["usd"] should be 5, but got %v.`, myBalance["usd"])
	}
}

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
