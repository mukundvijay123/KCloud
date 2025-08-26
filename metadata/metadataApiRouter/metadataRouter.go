package metadatarouter

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	metadatastore "github.com/mukundvijay123/KCloud/metadata/metadataStore"
)

type MetadataRouter struct {
	dbConn        *sql.DB
	logger        *log.Logger
	MdataStore    *metadatastore.MetadataDb
	JWTMiddleWare *JWTMiddleWare
	Router        *mux.Router
}

func NewMetadataRouter(dbConn *sql.DB, logger *log.Logger) *MetadataRouter {
	if logger == nil {
		logger = log.Default()
	}

	return &MetadataRouter{
		dbConn:     dbConn,
		logger:     logger,
		MdataStore: metadatastore.NewMetadataDb(dbConn, logger),
	}
}

//Function to add JWT middleware here

func (m *MetadataRouter) CreateRouter() error {
	m.Router = mux.NewRouter()
	err := m.AddRoutes()
	if err != nil {
		return fmt.Errorf("error initialising router")
	}

	return nil
}
