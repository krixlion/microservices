package main

import (
	"eventstore/cmd/service"
	"eventstore/pkg/repository"
)

func main() {
	repository.Ping()
	service.Run()
}
