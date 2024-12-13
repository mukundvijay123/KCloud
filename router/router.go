package router

import (
	"database/sql"
	"net/http"
)

const ProvisioningMaxBodySize = 1024

type MainServer struct {
	server *http.ServeMux // ServeMux for the Main server
	dbPool *sql.DB        // dbPool
}

func InitServer(db *sql.DB) MainServer {
	var newServer MainServer
	newServer.server = http.NewServeMux()
	newServer.dbPool = db

	return newServer

}

func (m *MainServer) AddRoutes() {

}
