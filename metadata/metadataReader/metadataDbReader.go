package metadatareader

import (
	"database/sql"
	"log"
)

type MetadataDBReader struct {
	dbConn *sql.DB
	logger *log.Logger
}

func NewMetadataDBReader(db *sql.DB, logger *log.Logger) *MetadataDBReader {
	if logger == nil {
		// Fallback to standard logger if none provided
		logger = log.Default()
	}

	return &MetadataDBReader{
		dbConn: db,
		logger: logger,
	}
}
