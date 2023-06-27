package main

import (
	"log"
	"time"

	"github.com/crypto/app"
	"github.com/gin-gonic/gin"
)

func main() {

	ok := app.InitSyncMap()
	if !ok {
		log.Fatalln("Init Routes failed")
		return
	}

	app.StartRealTimeSync(100 * time.Millisecond)

	router := gin.Default()

	ok = app.InitRoutes(router)
	if !ok {
		log.Fatalln("Init Routes failed")
		return
	}

	router.Run(":8000")

	log.Println("Service Started Successfully")
}
