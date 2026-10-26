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
}
