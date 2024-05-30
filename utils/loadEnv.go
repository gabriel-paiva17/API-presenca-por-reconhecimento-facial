package utils

import (
	
	"log"

    "github.com/joho/godotenv"

)

func LoadEnv (env string) {

	envFile := "env/.env." + env
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}
	
}