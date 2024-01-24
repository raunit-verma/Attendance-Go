package main

import (
	"attendance/api/router"
	"attendance/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading environment variable.")
	}

	serverConfig := router.ServerConfig{
		Port: ":1025",
	}

	r := router.NewMUXRouter()

	db := repository.GetDB()
	defer db.Close()

	fmt.Printf("Server Starting at Port : %v\n", serverConfig.Port)
	log.Fatal(http.ListenAndServe(serverConfig.Port, r))

}
