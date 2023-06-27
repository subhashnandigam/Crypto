package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	syncMap     map[string]Price
	supportList []string
	lock        sync.RWMutex
	listLock    sync.RWMutex

	symbolsAPICall  string = "https://api.hitbtc.com/api/3/public/symbol/%v"
	currencyAPICall string = "https://api.hitbtc.com/api/3/public/currency/%v"
	priceApiCall    string = "https://api.hitbtc.com/api/3/public/ticker/%v"
)

type Currency struct {
	Id          string `json:"id"`
	FullName    string `json:"fullName"`
	FeeCurrency string `json:"feeCurrency"`
}

type Price struct {
	*Currency
	Ask  string `json:"ask"`
	Bid  string `json:"bid"`
	Last string `json:"last"`
	Open string `json:"open"`
	High string `json:"high"`
	Low  string `json:"low"`
}

func InitSyncMap() bool {
	lock.Lock()
	defer lock.Unlock()

	listLock.Lock()
	defer listLock.Unlock()

	syncMap = make(map[string]Price)
	supportList = []string{"ETHBTC", "BTCUSDC"}
	for _, symbol := range supportList {
		curency, err := GetCurrencyDataFromSymbol(symbol)
		if err != nil {
			log.Printf("failed to get the currency and symbol data: %v data: %v", symbol, err)
			return false
		}

		price, err := GetLatestPriceData(symbol)
		if err != nil {
			log.Printf("failed to get the price of symbol: %v data: %v", symbol, err)
			return false
		}
		price.Currency = curency
		syncMap[symbol] = *price
	}

	return true
}

func StartRealTimeSync(interval time.Duration) {
	go func() {
		currentIndex := 0
		currentSymbol := ""
		for {
			listLock.RLock()
			if currentIndex >= len(supportList) {
				currentIndex = 0
				currentSymbol = supportList[0]
			} else {
				currentSymbol = supportList[currentIndex]
			}
			listLock.RUnlock()

			log.Println("Fetching the new info for the symbol: ", currentSymbol, currentIndex)

			price, err := GetLatestPriceData(currentSymbol)
			if err != nil {
				log.Printf("Failed to fetch the current price of symbol: %v, err: %v SKIPPING this iteration", currentSymbol, err)
			} else {
				lock.Lock()

				val, ok := syncMap[currentSymbol]
				if !ok {
					log.Printf("Missing %v symbol data in sync map", currentSymbol)
				} else {
					price.Currency = val.Currency
					syncMap[currentSymbol] = *price
				}

				lock.Unlock()
			}

			currentIndex++
			time.Sleep(interval)
		}
	}()
}

func GetCurrencyDataFromSymbol(symbol string) (*Currency, error) {
	c := http.Client{Timeout: time.Duration(2) * time.Second}
	resp, err := c.Get(fmt.Sprintf(symbolsAPICall, symbol))
	if err != nil {
		log.Printf("Get Symbol %v API Error %s", symbol, err)
		return nil, fmt.Errorf("Internal Server Error")
	}
	type symbolInfo struct {
		BaseCurrency string `json:"base_currency"`
		FeeCurrency  string `json:"fee_currency"`
		Status       string `json:"status"`
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Get Symbol %v body Error %s", symbol, err)
		return nil, fmt.Errorf("Internal Server Error")
	}
	var info symbolInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("Get Symbol %v unmarshall Error %s", symbol, err)
		return nil, fmt.Errorf("Internal Server Error")
	}

	if info.Status != "working" {
		return nil, fmt.Errorf("Symbol is not working state")
	}

	resp2, err := c.Get(fmt.Sprintf(currencyAPICall, info.BaseCurrency))
	if err != nil {
		log.Printf("Get Currency %v of symbol %v API Error %s", info.BaseCurrency, symbol, err)
		return nil, fmt.Errorf("Internal Server Error")
	}
	defer resp2.Body.Close()
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Printf("Get Currency %v of symbol %v body Error %s", info.BaseCurrency, symbol, err)
		return nil, fmt.Errorf("Internal Server Error")
	}

	type currencyInfo struct {
		FullName string `json:"full_name"`
	}

	var cinfo currencyInfo
	err = json.Unmarshal(body2, &cinfo)
	if err != nil {
		log.Printf("Get Currency %v of symbol %v unmarshall Error %s", info.BaseCurrency, symbol, err)
		return nil, fmt.Errorf("Internal Server Error")
	}

	return &Currency{
		Id:          info.BaseCurrency,
		FullName:    cinfo.FullName,
		FeeCurrency: info.FeeCurrency,
	}, nil
}

func GetLatestPriceData(symbol string) (*Price, error) {
	c := http.Client{Timeout: time.Duration(2) * time.Second}
	resp, err := c.Get(fmt.Sprintf(priceApiCall, symbol))
	if err != nil {
		log.Printf("Get Price Symbol %v API Error %s", symbol, err)
		return nil, fmt.Errorf("Internal Server Error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Get Price Symbol %v body Error %s", symbol, err)
		return nil, fmt.Errorf("Internal Server Error")
	}
	var info Price
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("Get Price Symbol %v unmarshall Error %s", symbol, err)
		return nil, fmt.Errorf("Internal Server Error")
	}
	return &info, nil
}
