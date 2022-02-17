package exchango

import (
	"testing"
)

func TestUserNoMoney(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"usd": 10, "eur": 5}
	ProvideLiquidity(&newBankBalance, newPair)

        myBalance := map[string]float64{"eur": 1, "jjt": 0, "gc":100}
	myToffer := TradeOffer{
		WantAmount:   5,
		WantCurrency: "usd",
		GiveCurrency: "eur",
	}
	tl, err := InitTrade(TradeDetails{myToffer, &newBankBalance, &myBalance})
	if err == nil {
		t.Fatalf(`Should error, but got "err == nil".`)
	}
	err = ExecuteTrade(tl)
	if err == nil {
		t.Fatalf(`Should error, but got "err == nil".`)
	}
}

func TestBankNoMoney(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"usd": 10, "eur": 5}
	ProvideLiquidity(&newBankBalance, newPair)

        myBalance := map[string]float64{"eur": 100, "jjt": 0, "gc":100}
	myToffer := TradeOffer{
		WantAmount:   15,
		WantCurrency: "usd",
		GiveCurrency: "eur",
	}
	_, err := InitTrade(TradeDetails{myToffer, &newBankBalance, &myBalance})
	if err == nil {
		t.Fatalf(`Should error, but got "err == nil".`)
	}
}
