package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ashkarin/ashkarin-api-test/pkg/recipe/gateways"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/ashkarin/ashkarin-api-test/internal/config"
	"github.com/ashkarin/ashkarin-api-test/internal/services/recipes"
)

// Server is the app container
type Server struct {
	recipesService *recipes.Service
	Router         *mux.Router
	server         *http.Server
}

// Initialize init the server
func (s *Server) Initialize(cfg *config.Config) {
	// Log
	log.Infof("Initialize server with: %v", cfg)

	// Open a gateway to the MongoDB storage
	collection := "recipes"
	recipesStorage, err := gateways.NewMongoDbGateway(cfg.DB.Server, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, collection)
	if err != nil {
		log.Fatalf("Connection to the recipes storage: %v", err)
	}
	log.Infof("Connected to the recipes storage: %s@%s:%s/%s/%s", cfg.DB.Username, cfg.DB.Server, cfg.DB.Port, cfg.DB.DBName, collection)

	// Create route and service
	s.Router = mux.NewRouter()
	s.recipesService = recipes.NewService(recipesStorage, s.Router)

	// Create the server
	s.server = &http.Server{
		Handler:      s.Router,
		Addr:         fmt.Sprintf("%s:%s", cfg.Address, cfg.Port),
		WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Timeout) * time.Second,
	}
}

// ListenAndServe start the server
func (s *Server) ListenAndServe() {
	log.Fatal(s.server.ListenAndServe())
}
