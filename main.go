package main

import (
	"github.com/aburluka/sgtask/internal/logger"
	"github.com/aburluka/sgtask/internal/server"
	"github.com/aburluka/sgtask/internal/storage"
)

func main() {
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}
	storage, err := storage.New(logger)
	if err != nil {
		panic(err)
	}
	defer storage.Close()
	server.New(logger, storage).Run()
}
