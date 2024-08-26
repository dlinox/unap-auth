package infrastructure

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// NewDBConn creates a new connection to the database
func NewDBConn(dsn string) (*sql.DB, error) {
	// Configuración de la conexión a la base de datos
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
		return nil, err
	}
	// Verificar la conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
		return nil, err
	}
	return db, nil
}
