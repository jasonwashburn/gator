package main

import (
	"fmt"
	"log"

	"github.com/jasonwashburn/gator/internal/config"
)

func main() {
	user := "Jason"

	configFile, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	configFile.SetUser(user)
	configFile, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(configFile)
}
