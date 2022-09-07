package main

import (
	"eventstore/cmd/service"
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDir = "eventstore"

func main() {
	re := regexp.MustCompile(`^(.*` + projectDir + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		panic(fmt.Sprintf("Failed to load .env, err: %s", err))
	}

	service.Run()
}
