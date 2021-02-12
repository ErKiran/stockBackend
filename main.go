package main

import (
	"fmt"
	"log"
	"stockwatch/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting Environment Variables %v", err)
	}
	api.Run()

	fmt.Println("hey")
}
