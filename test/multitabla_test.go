package test

import (
	"testing"
)

func TestClasesProfesor(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	otro := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	_ = createClaseDePrueba(t, otro.ID)
	horario := createHorarioDePrueba(t, clase.ID, "viernes", horaDePrueba(19, 30))

	filas, err := queries.ClasesProfesor(ctx, profesor.ID)
	if err != nil {
		t.Fatalf("ClasesProfesor falló: %v", err)
	}

	encontrada := false
	for _, fila := range filas {
		if fila.ClaseNombre == clase.Nombre && fila.DiaSemana == horario.DiaSemana {
			encontrada = true
			if fila.Nombre != profesor.Nombre || fila.Apellido != profesor.Apellido {
				t.Errorf("profesor = %q %q, se esperaba %q %q", fila.Nombre, fila.Apellido, profesor.Nombre, profesor.Apellido)
			}
			if fila.MaxAlumnos != clase.MaxAlumnos {
				t.Errorf("max_alumnos = %d, se esperaba %d", fila.MaxAlumnos, clase.MaxAlumnos)
			}
			if !mismaHora(fila.HoraInicio, horario.HoraInicio) {
				t.Errorf("hora_inicio = %s, se esperaba %s", fila.HoraInicio.Format("15:04:05"), horario.HoraInicio.Format("15:04:05"))
			}
		}
	}
	if !encontrada {
		t.Fatalf("no apareció la clase %q del profesor %d", clase.Nombre, profesor.ID)
	}

	ajenas, err := queries.ClasesProfesor(ctx, otro.ID)
	if err != nil {
		t.Fatalf("ClasesProfesor del otro profesor falló: %v", err)
	}
	for _, fila := range ajenas {
		if fila.ClaseNombre == clase.Nombre {
			t.Fatal("ClasesProfesor devolvió una clase de otro profesor")
		}
	}
}

func TestProfesoresdelDia(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	otroDia := createProfesorDePrueba(t)
	claseLunes := createClaseDePrueba(t, profesor.ID)
	claseMartes := createClaseDePrueba(t, otroDia.ID)
	_ = createHorarioDePrueba(t, claseLunes.ID, "lunes", horaDePrueba(7, 0))
	_ = createHorarioDePrueba(t, claseMartes.ID, "martes", horaDePrueba(7, 0))

	delLunes, err := queries.ProfesoresdelDia(ctx, "lunes")
	if err != nil {
		t.Fatalf("ProfesoresdelDia falló: %v", err)
	}

	estaElDeLunes := false
	estaElDeMartes := false
	for _, p := range delLunes {
		if p.Nombre == profesor.Nombre && p.Apellido == profesor.Apellido {
			estaElDeLunes = true
		}
		if p.Nombre == otroDia.Nombre && p.Apellido == otroDia.Apellido {
			estaElDeMartes = true
		}
	}
	if !estaElDeLunes {
		t.Fatalf("el profesor %s %s debió aparecer el lunes", profesor.Nombre, profesor.Apellido)
	}
	if estaElDeMartes {
		t.Fatal("el profesor del martes no debió aparecer en el listado del lunes")
	}
}

func TestListEjerciciosRutinaVacia(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	rutina := createRutinaDePrueba(t, profesor.ID)

	ejercicios, err := queries.ListEjerciciosRutina(ctx, rutina.ID)
	if err != nil {
		t.Fatalf("ListEjerciciosRutina falló: %v", err)
	}
	if len(ejercicios) != 0 {
		t.Fatalf("una rutina nueva no debería tener ejercicios, tiene %d", len(ejercicios))
	}
}
