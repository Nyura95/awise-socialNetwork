package main

import (
	"awise-socialNetwork/models"
	"awise-socialNetwork/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init of the pool mysql
	models.InitDb(os.Getenv("USERNAME"), os.Getenv("PASSWORD"), os.Getenv("HOST"), os.Getenv("DATABASE"))
	defer models.Close()

	// Launch of the http server
	server.Start()

}

/*
{
  "User": "acourapied",
  "Host": "5.135.167.23",
  "Password": "Yfful95df",
  "Database": "SOCIAL_NETWORK",
  "Port": 3000,
  "basePathImage": "https://pictures-awise.co",
  "Version": "1.0.0"
}


*/
