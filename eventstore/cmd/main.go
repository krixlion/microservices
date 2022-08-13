package main

import (
	"eventstore/cmd/repository"
	"eventstore/cmd/service"
)

func main() {
	repository.Ping()
	service.Run()
}
