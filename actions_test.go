package exchango

import (
	"testing"
)

func TestProvideLiquidity(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"coins": 15000, "jjt": 10}
	ProvideLiquidity(&newBankBalance, newPair)
	for key := range newPair {
		if newBankBalance[0][key] != newPair[key]{
			t.Fatalf(`BankBalance should be [%v], but got %v.`, newPair, newBankBalance)
		}
	}
}

func TestProvideLiquidityAdd(t *testing.T) {
	newBankBalance := make([]map[string]float64, 0)
	newPair := map[string]float64{"abc": 10, "def": 5}
	ProvideLiquidity(&newBankBalance, newPair)
	newPair2 := map[string]float64{"abc": 20, "def": 10}
	ProvideLiquidity(&newBankBalance, newPair2)

	want := map[string]float64{"abc": 30, "def": 15}
	for key := range want {
		if newBankBalance[0][key] != want[key] {
			t.Fatalf(`newBankBalance[0] should be %v, but got %v.`, want, newBankBalance)
		}
	}
}
