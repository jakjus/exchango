package main

import (
	"errors"
	"fmt"
	"log"
	"math"
)

var (
	tax         float64              = 0.01
	allowChange float64              = 0.05
)

type TradeOffer struct {
	wantAmount   float64
	wantCurrency string
	giveCurrency string
}

func getPrice(toffer TradeOffer, exchange map[string]float64) float64 {
	from := exchange[toffer.giveCurrency]
	to := exchange[toffer.wantCurrency]
	price_init := from / to
	price_closing := (from + toffer.wantAmount*price_init) / (to - toffer.wantAmount)
	price := (price_init + price_closing) / 2
	return price
}

func provideLiquidity(bankBalance *[]map[string]float64, exchange map[string]float64) error {
        for key := range exchange {
                if exchange[key] <= 0 {
                        return errors.New(fmt.Sprintf("Amount of both provided currencies must be positive (%v%v<=0).", exchange[key], key))
                }
        }
	if len(exchange) != 2 {
		return errors.New("You need mapping with exactly two currencies")
	}
	if len(exchange) != 2 {
		return errors.New("You need mapping with exactly two currencies")
	}
	for i, pairInBank := range *bankBalance {
		existingPairIndex := -1
		for currencyName := range exchange {
			if pairInBank[currencyName] != 0.0 {
				existingPairIndex = i
			} else {
				existingPairIndex = -1
				break
			}
		}
		if existingPairIndex != -1 {
			log.Println("There is an existing currency pair. Adding...")
			for currencyName := range exchange {
				pairInBank[currencyName] += exchange[currencyName]
			}
			log.Printf("%v\n", *bankBalance)
			return nil
		}
	}
	log.Println("Adding a new pair...")
	*bankBalance = append(*bankBalance, exchange)
	log.Printf("%v\n", *bankBalance)
	return nil
}

func checkExchange(toffer TradeOffer, bankBalance *[]map[string]float64) (map[string]float64, error) {
	for _, v := range *bankBalance {
		var wantExists, giveExists bool
		for i := range v {
			if toffer.wantCurrency == i {
				wantExists = true
			} else if toffer.giveCurrency == i {
				giveExists = true
			}
		}
		if wantExists && giveExists {
			log.Println("Found exchange.")
			return v, nil
		}
	}
	return map[string]float64{}, errors.New("Exchange not found.")
}

func trade(toffer TradeOffer, bankBalance *[]map[string]float64, userBalance *map[string]float64) error {
	exchange, err := checkExchange(toffer, bankBalance)
	bankBal := exchange[toffer.wantCurrency]
	price := getPrice(toffer, exchange)
	roundedPrice := price * toffer.wantAmount * (1 + tax)

        // Errors:
	if err != nil {
		return err
	}
        if (*userBalance)[toffer.giveCurrency] < roundedPrice {
                return errors.New(fmt.Sprintf("User does not have %.2f %v. It has currently %.2f %v.\n", roundedPrice, toffer.giveCurrency, (*userBalance)[toffer.giveCurrency], toffer.giveCurrency))
        }
	if bankBal < toffer.wantAmount {
		return errors.New(fmt.Sprintf("Bank does not have %.2f %v. It has currently %.2f %v.\n", toffer.wantAmount, toffer.wantCurrency, bankBal, toffer.wantCurrency))
	}

	log.Printf("\nYou give: %.2f %v\nYou get: %v %v\nAccept?\n", roundedPrice, toffer.giveCurrency, toffer.wantAmount, toffer.wantCurrency)
	err = executeTrade(toffer, exchange, roundedPrice, bankBalance, userBalance)
	if err != nil {
		return err
	}
	return nil
}

func executeTrade(toffer TradeOffer, exchange map[string]float64, roundedPricePre float64, bankBalance *[]map[string]float64, userBalance *map[string]float64) error {
	price := getPrice(toffer, exchange)
	roundedPrice := price * toffer.wantAmount * (1 + tax)
	bankBal := exchange[toffer.wantCurrency]

        // Errors:
        if (*userBalance)[toffer.giveCurrency] < roundedPrice {
                return errors.New(fmt.Sprintf("User does not have %.2f %v. It has currently %.2f %v.\n", roundedPrice, toffer.giveCurrency, (*userBalance)[toffer.giveCurrency], toffer.giveCurrency))
        }
	if bankBal < toffer.wantAmount {
		return errors.New(fmt.Sprintf("Bank does not have %.2f %v. It has currently %.2f %v.\n", toffer.wantAmount, toffer.wantCurrency, bankBal, toffer.wantCurrency))
	}
	if math.Abs((roundedPrice-roundedPricePre)/roundedPrice) > allowChange {
		return errors.New(fmt.Sprintf("Price has changed more than %v%%, before trade was accepted. Please try again.", allowChange*100))
	}

        (*userBalance)[toffer.giveCurrency] -= roundedPrice
        (*userBalance)[toffer.wantCurrency] += toffer.wantAmount
	exchange[toffer.wantCurrency] -= toffer.wantAmount
	exchange[toffer.giveCurrency] += roundedPrice
	log.Printf("Done. Bank: %v\n", *bankBalance)
	return nil
}
