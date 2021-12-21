package main

import (
	"dbmsbackend/router"
	"dbmsbackend/util"
	"fmt"

	"github.com/gin-gonic/gin"
)

//Execution starts from main function
func main() {

	engine := gin.Default()
	engine.SetTrustedProxies(nil)

	configPath := "./.env"
	config, configErr := util.LoadConfig(configPath)

	if configErr != nil {
		panic(fmt.Errorf("loading config %v: %w", configPath, configErr))
	}

	router.SetupRouter(engine, &config)

	engine.Run(":" + config.ServerPort)

}
