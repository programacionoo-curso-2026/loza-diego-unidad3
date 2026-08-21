// internal/api/dto.go

package api

// EvaluarRequest es el payload que envía el frontend
type EvaluarRequest struct {
	TeacherID   string `json:"teacher_id"`
	Descripcion string `json:"descripcion"`
	Disciplina  string `json:"disciplina"`
}

// CrearDocenteRequest para registrar un nuevo docente
type CrearDocenteRequest struct {
	Nombre     string `json:"nombre"`
	Disciplina string `json:"disciplina"`
}

// ErrorResponse formato uniforme de error para el frontend
type ErrorResponse struct {
	Error string `json:"error"`
}