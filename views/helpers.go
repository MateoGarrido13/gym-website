package views

import "strconv"

func navClass(active, page string) string {
	if active == page {
		return "nav-link active"
	}
	return "nav-link"
}

func singularKind(kind string) string {
	switch kind {
	case "ejercicios":
		return "ejercicio"
	case "rutinas":
		return "rutina"
	case "clases":
		return "clase"
	case "profesores":
		return "profesor"
	default:
		return "alumno"
	}
}

func stringID(id int32) string {
	return strconv.FormatInt(int64(id), 10)
}

func catalogPath(kind string) string {
	if kind == "alumnos" {
		return "/"
	}
	return "/" + kind
}

func intString(value int) string {
	return strconv.Itoa(value)
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func isSelected(values []string, value int32) bool {
	for _, selected := range values {
		if selected == stringID(value) {
			return true
		}
	}
	return false
}
