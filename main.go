package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	db "PRACTICO_DOS/db/sqlc"
	"PRACTICO_DOS/internal/handler"
	"PRACTICO_DOS/internal/service"

	_ "github.com/lib/pq"
)

// CONEXION POR DEFECTO A LA BASE DE DATOS
const defaultConnStr = "host=localhost port=5432 user=postgres password=postgres dbname=tp2_db sslmode=disable"

func main() {
	connStr := os.Getenv("DATABASE_URL") // VARIABLE DE ENTORNO PARA LA CONEXION A LA BASE DE DATOS
	if connStr == "" {
		connStr = defaultConnStr
	}

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("no se pudo preparar la conexión: %v", err)
	}
	defer dbConn.Close()

	if err = dbConn.Ping(); err != nil {
		log.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}

	queries := db.New(dbConn)
	ejercicios := service.NewEjercicioService(queries)
	alumnos := service.NewAlumnoService(queries)
	rutinas := service.NewRutinaService(queries)
	rutinas_ejercicio := service.NewRutinaEjercicioService(queries)
	clases := service.NewClaseService(queries)
	clases_horario := service.NewClaseHorarioService(queries)
	inscripciones := service.NewInscripcionService(queries)
	router := handler.NewRouter(ejercicios, alumnos, rutinas, rutinas_ejercicio, clases, clases_horario, inscripciones, "./static")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
