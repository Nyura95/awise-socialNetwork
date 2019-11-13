package main

import (
	"awise-socialNetwork/config"
	"awise-socialNetwork/models"
	"awise-socialNetwork/server"
)

func main() {
	// Init of the config
	config.Start("dev")
	configuration, _ := config.GetConfig()
	// Init of the pool mysql
	models.InitDb(configuration.User, configuration.Password, configuration.Host, configuration.Database)
	defer models.Close()

	// Launch of the http server
	server.Start()

}
