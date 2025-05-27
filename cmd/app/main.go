package main

import (
	"corpPR3/internal/transport"
	"corpPR3/internal/workers"
	"log"
)

func main() {
	workers.StartWorkerPool(10)

	router := transport.SetupRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
