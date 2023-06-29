package bite

import (
	"fmt"

	"github.com/AadumKhor/bitespeed-backend-task/src/app/handlers"
	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	IdentifyRoute = "/identify"
)

func Run() {
	// Pre-requisites
	// NOTE: These can also include logging, queue setup etc. 
	config, err := utils.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	err = utils.InitTimeZone(config.DefaultTimezone)
	if err != nil {
		panic(err)
	}

	// start router on the port specified in the config
	startRouter(config)
}

func startRouter(config *utils.Config) {
	// for prettier logs
	gin.ForceConsoleColor()
	
	// set mode for gin
	gin.SetMode(config.Mode)

	// init the router
	router := gin.Default()

	// setup routes and their handlers
	router.POST(IdentifyRoute, handlers.HandleIdentify)

	// run the router
	router.Run(fmt.Sprintf(":%d", config.Port))
}