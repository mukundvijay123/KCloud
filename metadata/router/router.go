package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mukundvijay123/KCloud/metadata/middleware"
	authentication "github.com/mukundvijay123/KCloud/metadata/middleware/auth"
	"github.com/mukundvijay123/KCloud/metadata/services/metadata"
	"github.com/mukundvijay123/KCloud/metadata/services/user"
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
	subrouter.Use(middleware.CorsMiddleware)

	metadataSubRouter := subrouter.PathPrefix("/resources").Subrouter()
	metadataSubRouter.Use(authentication.JWTMiddleware)
	metadataHandler := metadata.NewHandler(s.db)
	metadataHandler.RegisterRoutes(metadataSubRouter)

	userHandler := user.NewHandler(s.db)
	userHandler.RegisterRoutes(subrouter)

	return http.ListenAndServe(s.addr, router)
}
