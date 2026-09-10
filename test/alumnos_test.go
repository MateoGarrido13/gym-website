import (
	"testing"
)

func TestListAlumnos(t *testing.T) {
	alumnos, err := db.ListAlumnos()
	if err != nil {
		t.Fatalf("Error al listar alumnos: %v", err)
	}
	fmt.Println(alumnos)
}