package views

type Entity struct {
	ID        int32
	Values    []string
	EditURL   string
	DeleteURL string
	Nombre    string
	Email     string
	Telefono  string
	Detalle   string
}

type Option struct {
	ID    int32
	Label string
}

type FormValues struct {
	Nombre            string
	Apellido          string
	Email             string
	Contrasena        string
	Telefono          string
	FechaInscripcion  string
	FechaVto          string
	TipoPlan          string
	Especialidad      string
	Descripcion       string
	GrupoMuscular     string
	DuracionSemanas   string
	ProfesorID        string
	MaxAlumnos        string
	RutinaID          string
	SelectedExercises []string
}

type PageData struct {
	Title             string
	Active            string
	Kind              string
	Description       string
	Action            string
	FormError         string
	IsEditing         bool
	ShowActions       bool
	FormValues        FormValues
	Entities          []Entity
	Headers           []string
	Options           []Option
	ExerciseOptions   []Option
	SelectedExercises []string
	UserName          string // usuario con sesión; vacío si no hay
}

// LoginData alimenta la pantalla de login. Rol vacío significa "todavía no
// eligió": se muestran los dos accesos. Con Rol definido se muestra el formulario.
type LoginData struct {
	Rol   string // "alumno" | "profesor" | ""
	Email string // se conserva si el intento falla; la contraseña nunca
	Error string
}

// RutinaAlumnoData alimenta la vista del alumno. Nombre vacío significa que
// todavía no tiene rutina asignada.
type RutinaAlumnoData struct {
	UserName   string
	Nombre     string
	Duracion   string // "8 semanas"
	Profesor   string
	Ejercicios []EjercicioItem
}

// EjercicioItem trae los textos ya armados para mostrar; los datos que no
// existen quedan vacíos y la vista los omite.
type EjercicioItem struct {
	Nombre      string
	Grupo       string
	Descripcion string
	Esquema     string // "4 series × 12 repeticiones"
	Descanso    string // "60 s de descanso"
}
