package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"system-design/internal/handlers"
	"system-design/internal/macros"
	"system-design/internal/registry"
)

//go:embed web/static/*
var staticFS embed.FS

//go:embed web/templates/*
var templateFS embed.FS

//go:embed content/*
var contentFS embed.FS

func main() {
	// Load registry
	reg, err := registry.Load(contentFS, "content/_registry.yaml")
	if err != nil {
		log.Fatalf("Failed to load registry: %v", err)
	}
	log.Printf("Loaded %d problems, %d fundamentals", len(reg.Problems), len(reg.AllFundamentals()))

	// Build template macros
	funcMap := macros.FuncMap()

	// Create handler with dependencies
	h := handlers.New(reg, templateFS, contentFS, funcMap)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	r.Handle("/static/*", http.FileServer(http.FS(staticFS)))

	// Page routes
	r.Get("/", h.Dashboard)
	r.Get("/problem/{slug}", h.ProblemDetail)
	r.Get("/fund/*", h.FundamentalDetail)

	// API routes
	r.Post("/api/generate/{slug}", h.GeneratePDF)
	r.Get("/api/status/{taskID}", h.GenerateStatus)

	// PDF serving
	r.Get("/pdf/{filename}", h.ServePDF)

	port := 8080
	log.Printf("Dashboard running at http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
