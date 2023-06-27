package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) bool {
	router.GET("/currency/:symbol", getSymbols)
	router.POST("/internal/add/support/:symbol", addSymbol)
	return true
}

func getSymbols(c *gin.Context) {
	symbol := c.Param("symbol")
	lock.RLock()
	defer lock.RUnlock()

	if symbol == "all" {

		array := []Price{}
		for _, priceData := range syncMap {
			array = append(array, priceData)
		}

		c.JSON(http.StatusOK, array)
		return
	} else {
		price, ok := syncMap[symbol]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid symbol"})
			return
		}

		c.JSON(http.StatusOK, price)
	}

}

func addSymbol(c *gin.Context) {
	symbol := c.Param("symbol")

	lock.Lock()
	defer lock.Unlock()

	_, ok := syncMap[symbol]
	if ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Already Symbol is added"})
		return
	}

	currency, err := GetCurrencyDataFromSymbol(symbol)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	price, err := GetLatestPriceData(symbol)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	price.Currency = currency

	syncMap[symbol] = *price

	listLock.Lock()
	defer listLock.Unlock()

	supportList = append(supportList, symbol)

	c.JSON(http.StatusBadRequest, gin.H{"message": "Added Succefully"})
}
