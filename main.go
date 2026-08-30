package main

import (

	// Importación apuntando al módulo de EjercicioDos
	//db "EjercioDos/db/sqlc"

	_ "github.com/lib/pq" // Driver de Postgres
)

/*
func main() {
	// 1. Conexión a la base de datos (formato Clave-Valor)
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=tp2_db sslmode=disable"
	// pgx registers the driver name "pgx"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al preparar la conexión: %v", err)
	}
	defer dbConn.Close() // cuando finalice el rpograma cerrara la coexion

	if err = dbConn.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	fmt.Println("🔌 Conexión exitosa a Postgres desde EjercicioDos...")

	// Opcional: Limpiamos la tabla antes de arrancar para tener una demostración limpia
	_, _ = dbConn.Exec("DELETE FROM users")

	// 2. Instanciar el Queries generado por sqlc
	queries := db.New(dbConn)

	// Contexto base para las consultas
	ctx := context.Background()

	// ==========================================
	// 🟢 OPERACIÓN 1: CREATE (Crear Usuario)
	// ==========================================

	params := db.CreateUserParams{
		Name:  "Mateo Compilado",
		Email: "mateo.compilado@example.com",
	}

	nuevoUsuario, err := queries.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("Error al crear usuario: %v", err)
	}
	// Handle possible nullable CreatedAt (sql.NullTime)
	createdAtStr := "<nil>"
	if nuevoUsuario.CreatedAt.Valid {
		createdAtStr = nuevoUsuario.CreatedAt.Time.Format(time.RFC3339)
	}
	fmt.Printf("¡Usuario creado exitosamente!\nID: %d | Nombre: %s | Email: %s | Creado: %s\n",
		nuevoUsuario.ID, nuevoUsuario.Name, nuevoUsuario.Email, createdAtStr)

	// ==========================================
	// 🔵 OPERACIÓN 2: READ ONE (Obtener por ID)
	// ==========================================

	// sqlc may generate "GetUser" instead of "GetUserByID" depending on the query name
	usuarioBuscado, err := queries.GetUser(ctx, nuevoUsuario.ID)
	if err != nil {
		log.Fatalf("Error al buscar usuario: %v", err)
	}
	fmt.Printf("Usuario encontrado -> ID: %d, Nombre: %s, Email: %s\n",
		usuarioBuscado.ID, usuarioBuscado.Name, usuarioBuscado.Email)

	// ==========================================
	// 🟡 OPERACIÓN 3: UPDATE (Actualizar Datos)
	// ==========================================
	fmt.Println("\n--- [UPDATE] Modificando datos del usuario ---")

	updateParams := db.UpdateUserParams{
		ID:    nuevoUsuario.ID,
		Name:  "Mateo Compilado Actualizado",
		Email: "mateo.actualizado@example.com",
	}

	err = queries.UpdateUser(ctx, updateParams) // paso un contexto y un struct del tipo  UserPArams{}
	if err != nil {
		log.Fatalf("Error al actualizar usuario: %v", err)
	}
	fmt.Println("¡Usuario actualizado correctamente en la base de datos!")

	// ==========================================
	// 🔵 OPERACIÓN 4: READ MANY (Listar Todos)
	// ==========================================
	fmt.Println("\n--- [READ MANY] Listando todos los usuarios ---")

	usuarios, err := queries.ListUsers(ctx) // solo paso el contexto
	if err != nil {
		log.Fatalf("Error al listar usuarios: %v", err)
	}

	fmt.Printf("Se encontraron %d usuario(s) en la tabla:\n", len(usuarios))
	for _, u := range usuarios {
		fmt.Printf("- ID: %d | %s (%s)\n", u.ID, u.Name, u.Email)
	}

	// ==========================================
	// 🔴 OPERACIÓN 5: DELETE (Eliminar Usuario)
	// ==========================================
	fmt.Println("\n--- [DELETE] Eliminando al usuario ---")

	err = queries.DeleteUser(ctx, nuevoUsuario.ID) // contexto y un int32
	if err != nil {
		log.Fatalf("Error al eliminar usuario: %v", err)
	}
	fmt.Println("¡Usuario eliminado de la base de datos!")

	// ==========================================
	// 🔍 VERIFICACIÓN DE ELIMINACIÓN
	// ==========================================
	fmt.Println("\n--- [VERIFY] Verificando que el usuario ya no exista ---")

	_, err = queries.GetUser(ctx, nuevoUsuario.ID) // debe retonar un erorr ya que no esta
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("✅ Confirmado: El usuario ya no existe (retornó sql.ErrNoRows correctamente).")
		} else {
			log.Fatalf("Error inesperado al verificar la eliminación: %v", err)
		}
	} else {
		log.Fatalf("❌ ERROR: El usuario no debió haber sido encontrado.")
	}
}*/
