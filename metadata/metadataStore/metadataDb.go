package metadatastore

import (
	"database/sql"
	"log"
)

type MetadataDb struct {
	dbConn *sql.DB
	logger *log.Logger
	//additional field for communicating with storage engine
}
