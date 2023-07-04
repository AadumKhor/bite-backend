package bite

import (
	"fmt"

	"github.com/AadumKhor/bitespeed-backend-task/src/app/handlers"
	"github.com/AadumKhor/bitespeed-backend-task/src/app/middleware"
	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/database"
	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/models"
	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Run is the function that starts the router and also does pre-liminary setup
func Run() {
	// Pre-requisites
	// NOTE: These can also include logging, queue setup etc.
	config, err := utils.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	err = utils.InitLogger()
	if err != nil {
		panic(err)
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

	// init pgStore
	err := database.Connect(*config)
	if err != nil {
		panic(fmt.Sprintf("could not connect with DB: %+v", err))
	}

	// setup routes and their handlers
	identifyHandler := handlers.IdentifyHandler{
		Store: *database.GetPGStore(),
	}
	router.POST(models.IdentifyRoute, middleware.ValidatePhoneNumber(), identifyHandler.Handle)

	// run the router
	router.Run(fmt.Sprintf(":%d", config.Port))
}
