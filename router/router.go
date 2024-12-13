package router

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/mukundvijay123/KCloud/services/metadata"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	metadataHandler := metadata.NewHandler(s.db)
	metadataHandler.RegisterRoutes(subrouter)

	return nil
}
