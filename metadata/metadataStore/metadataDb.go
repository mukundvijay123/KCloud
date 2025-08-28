package metadatastore

import (
	"database/sql"
	"log"

	metadatareader "github.com/mukundvijay123/KCloud/metadata/metadataReader"
)

type MetadataDb struct {
	dbConn           *sql.DB
	logger           *log.Logger
	MetadataDbReader *metadatareader.MetadataDBReader
}

func NewMetadataDb(db *sql.DB, logger *log.Logger) *MetadataDb {
	if logger == nil {
		logger = log.Default()
	}

	return &MetadataDb{
		dbConn:           db,
		logger:           logger,
		MetadataDbReader: metadatareader.NewMetadataDBReader(db, logger),
	}
}
