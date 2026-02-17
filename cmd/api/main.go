package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ivsanmendez/ControlDeGastos/db/migrations"
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

	if err := migrations.Run(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	repo := postgres.NewExpenseRepo(db)
	bus := eventbus.New()

	// Domain services
	expenseSvc := expense.NewService(repo, bus)

	// Inbound adapters
	mux := http.NewServeMux()
	httpapi.RegisterRoutes(mux, expenseSvc)

	// Serve static files (production React build)
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./web/dist"
	}
	serveSPA(mux, staticDir)

	// Start server
	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "8080"
	}
	log.Printf("API listening on :%s", listenPort)
	log.Printf("Serving static files from: %s", staticDir)
	if err := http.ListenAndServe(":"+listenPort, mux); err != nil {
		log.Fatal(err)
	}
}

// serveSPA serves the React SPA and handles client-side routing
func serveSPA(mux *http.ServeMux, staticDir string) {
	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Warning: static directory not found: %s (SPA serving disabled)", staticDir)
		return
	}

	indexPath := filepath.Join(staticDir, "index.html")
	fs := http.FileServer(http.Dir(staticDir))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to open the file in the static directory
		path := staticDir + r.URL.Path
		f, err := os.Open(path)

		if err == nil {
			defer f.Close()
			// Check if it's a directory
			stat, statErr := f.Stat()
			if statErr == nil && !stat.IsDir() {
				// It's a file and exists - let the file server handle it
				fs.ServeHTTP(w, r)
				return
			}
		}

		// File doesn't exist or is a directory - serve index.html
		http.ServeFile(w, r, indexPath)
	})
}