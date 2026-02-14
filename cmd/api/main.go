package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ivsanmendez/ControlDeGastos/internal/adapter/eventbus"
	"github.com/ivsanmendez/ControlDeGastos/internal/adapter/httpapi"
	"github.com/ivsanmendez/ControlDeGastos/internal/adapter/postgres"
	"github.com/ivsanmendez/ControlDeGastos/internal/domain/expense"
)

func main() {
	// Outbound adapters
	dbURL := os.Getenv("DATABASE_URL")
	db, err := postgres.Connect(dbURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewExpenseRepo(db)
	bus := eventbus.New()

	// Domain services
	expenseSvc := expense.NewService(repo, bus)

	// Inbound adapters
	mux := http.NewServeMux()
	httpapi.RegisterRoutes(mux, expenseSvc)

	// Start server
	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "8080"
	}
	log.Printf("API listening on :%s", listenPort)
	if err := http.ListenAndServe(":"+listenPort, mux); err != nil {
		log.Fatal(err)
	}
}