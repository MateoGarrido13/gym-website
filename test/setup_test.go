package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// TODOS USARAN ESTA CONEXION 
const connStr = "host=localhost port=5432 user=postgres password=postgres dbname=tp2_db sslmode=disable"

var (
	queries *db.Queries
	ctx     = context.Background()
	dbConn  *sql.DB
)

func setup() {
	var err error
	dbConn, err = sql.Open("postgres", connStr) // LEVANTA LA CONEXION 
	if err != nil {
		log.Fatalf("no se pudo preparar la conexión: %v", err)
	}

	if err = dbConn.Ping(); err != nil {
		_ = dbConn.Close()
		log.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}

	queries = db.New(dbConn)
}

func TestMain(m *testing.M) {
	setup() // LEVANTA LA CONEXION A LA BASE DE DATOS
	code := m.Run()
	_ = dbConn.Close()
	os.Exit(code)
}
