package exchango

import (
	"errors"
	"fmt"
	"log"
	"math"
)

var (
	Tax         float64              = 0.01
	AllowChange float64              = 0.05
)

type TradeOffer struct {
	WantAmount   float64
	WantCurrency string
	GiveCurrency string
}

type TradeDetails struct {
	Toffer TradeOffer
	BankBalance *[]map[string]float64
	UserBalance *map[string]float64
}

type TradeLock struct {
	Td TradeDetails
	Exchange map[string]float64
	RoundedPrice float64
}

func getPrice(Toffer TradeOffer, Exchange map[string]float64) float64 {
	from := Exchange[Toffer.GiveCurrency]
	to := Exchange[Toffer.WantCurrency]
	price_init := from / to
	price_closing := (from + Toffer.WantAmount*price_init) / (to - Toffer.WantAmount)
	price := (price_init + price_closing) / 2
	return price
}

func ProvideLiquidity(BankBalance *[]map[string]float64, Exchange map[string]float64) error {
        for key := range Exchange {
                if Exchange[key] <= 0 {
                        return errors.New(fmt.Sprintf("Amount of both provided currencies must be positive (%v%v<=0).", Exchange[key], key))
                }
        }
	if len(Exchange) != 2 {
		return errors.New("You need mapping with exactly two currencies")
	}
	if len(Exchange) != 2 {
		return errors.New("You need mapping with exactly two currencies")
	}
	for i, pairInBank := range *BankBalance {
		existingPairIndex := -1
		for currencyName := range Exchange {
			if pairInBank[currencyName] != 0.0 {
				existingPairIndex = i
			} else {
				existingPairIndex = -1
				break
			}
		}
		if existingPairIndex != -1 {
			log.Println("There is an existing currency pair. Adding...")
			for currencyName := range Exchange {
				pairInBank[currencyName] += Exchange[currencyName]
			}
			log.Printf("%v\n", *BankBalance)
			return nil
		}
	}
	log.Println("Adding a new pair...")
	*BankBalance = append(*BankBalance, Exchange)
	log.Printf("%v\n", *BankBalance)
	return nil
}

func checkExchange(Toffer TradeOffer, BankBalance *[]map[string]float64) (map[string]float64, error) {
	for _, v := range *BankBalance {
		var wantExists, giveExists bool
		for i := range v {
			if Toffer.WantCurrency == i {
				wantExists = true
			} else if Toffer.GiveCurrency == i {
				giveExists = true
			}
		}
		if wantExists && giveExists {
			log.Println("Found Exchange.")
			return v, nil
		}
	}
	return map[string]float64{}, errors.New("Exchange not found.")
}

func InitTrade(Td TradeDetails) (TradeLock, error) {
	Toffer, BankBalance, UserBalance := Td.Toffer, Td.BankBalance, Td.UserBalance
	Exchange, err := checkExchange(Toffer, BankBalance)
	bankBal := Exchange[Toffer.WantCurrency]
	price := getPrice(Toffer, Exchange)
	RoundedPrice := price * Toffer.WantAmount * (1 + Tax)

	// Empty
	tl := TradeLock{Td, Exchange, RoundedPrice}
        // Errors:
	if err != nil {
		return tl, err
	}
        if (*UserBalance)[Toffer.GiveCurrency] < RoundedPrice {
                return tl, errors.New(fmt.Sprintf("User does not have %.2f %v. It has currently %.2f %v.\n", RoundedPrice, Toffer.GiveCurrency, (*UserBalance)[Toffer.GiveCurrency], Toffer.GiveCurrency))
        }
	if bankBal < Toffer.WantAmount {
		return tl, errors.New(fmt.Sprintf("Bank does not have %.2f %v. It has currently %.2f %v.\n", Toffer.WantAmount, Toffer.WantCurrency, bankBal, Toffer.WantCurrency))
	}

	log.Printf("\nYou give: %.2f %v\nYou get: %v %v\nAccept?\n", RoundedPrice, Toffer.GiveCurrency, Toffer.WantAmount, Toffer.WantCurrency)

	if err != nil {
		return tl, err
	}

	tl = TradeLock{
		Td,
		Exchange,
		RoundedPrice,
	}
	return tl, nil
}

func ExecuteTrade(tl TradeLock) error {
	Td, Exchange, RoundedPricePre := tl.Td, tl.Exchange, tl.RoundedPrice
	Toffer, BankBalance, UserBalance := Td.Toffer, Td.BankBalance, Td.UserBalance
	price := getPrice(Toffer, Exchange)
	RoundedPrice := price * Toffer.WantAmount * (1 + Tax)
	bankBal := Exchange[Toffer.WantCurrency]

        // Errors:
        if (*UserBalance)[Toffer.GiveCurrency] < RoundedPrice {
                return errors.New(fmt.Sprintf("User does not have %.2f %v. It has currently %.2f %v.\n", RoundedPrice, Toffer.GiveCurrency, (*UserBalance)[Toffer.GiveCurrency], Toffer.GiveCurrency))
        }
	if bankBal < Toffer.WantAmount {
		return errors.New(fmt.Sprintf("Bank does not have %.2f %v. It has currently %.2f %v.\n", Toffer.WantAmount, Toffer.WantCurrency, bankBal, Toffer.WantCurrency))
	}
	if math.Abs((RoundedPrice-RoundedPricePre)/RoundedPrice) > AllowChange {
		return errors.New(fmt.Sprintf("Price has changed more than %v%%, before trade was accepted. Please try again.", AllowChange*100))
	}

        (*UserBalance)[Toffer.GiveCurrency] -= RoundedPrice
        (*UserBalance)[Toffer.WantCurrency] += Toffer.WantAmount
	Exchange[Toffer.WantCurrency] -= Toffer.WantAmount
	Exchange[Toffer.GiveCurrency] += RoundedPrice
	log.Printf("Done. Bank: %v\n", *BankBalance)
	return nil
}
