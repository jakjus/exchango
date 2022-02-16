package main

import (
	"errors"
	"fmt"
	"log"
	"math"
)

var (
	tax         float64              = 0.01
	BankBalance []map[string]float64 = make([]map[string]float64, 0)
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

func provideLiquidity(exchange map[string]float64) error {
	if len(exchange) != 2 {
		return errors.New("You need mapping with exactly two currencies")
	}
	for i, pairInBank := range BankBalance {
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
			log.Printf("%v\n", BankBalance)
			return nil
		}
	}
	log.Println("Adding a new pair...")
	BankBalance = append(BankBalance, exchange)
	log.Printf("%v\n", BankBalance)
	return nil
}

func checkExchange(toffer TradeOffer) (map[string]float64, error) {
	for _, v := range BankBalance {
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

func trade(toffer TradeOffer, userBalance *map[string]float64) error {
	exchange, err := checkExchange(toffer)
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
	err = executeTrade(toffer, exchange, roundedPrice, userBalance)
	if err != nil {
		return err
	}
	return nil
}

func executeTrade(toffer TradeOffer, exchange map[string]float64, roundedPricePre float64, userBalance *map[string]float64) error {
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

	exchange[toffer.wantCurrency] -= toffer.wantAmount
	exchange[toffer.giveCurrency] += roundedPrice
	log.Printf("Done. Bank: %v\n", BankBalance)
	return nil
}

func main() {
	existingPair := map[string]float64{"coins": 15000, "jjt": 10}
	provideLiquidity(existingPair)
	existingPair2 := map[string]float64{"gc": 222, "coins": 2}
	provideLiquidity(existingPair2)
	newPair := map[string]float64{"buba": 15000, "jjt": 10}
	provideLiquidity(newPair)
	mytoffer := TradeOffer{
		wantAmount:   5,
		wantCurrency: "jjt",
		giveCurrency: "coins",
	}
        myBalance := map[string]float64{"coins": 20000, "jjt": 0, "gc":100}
	err := trade(mytoffer, &myBalance)
	if err != nil {
		log.Fatal(err)
	}
	mytoffer = TradeOffer{
		wantAmount:   370,
		wantCurrency: "jjt",
		giveCurrency: "coins",
	}
	err = trade(mytoffer, &myBalance)
	if err != nil {
		log.Fatal(err)
	}
}
