// internal/api/handlers.go

package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"digcompedu-api/internal/domain"
	"digcompedu-api/internal/evaluator"
	"digcompedu-api/internal/storage"
)

type Handler struct {
	service *evaluator.Service
	repo    storage.Repository
}

func NewHandler(service *evaluator.Service, repo storage.Repository) *Handler {
	return &Handler{service: service, repo: repo}
}

// respondJSON centraliza la escritura de respuestas — evita repetir esto en cada handler
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("error al codificar respuesta JSON: %v", err)
	}
}

func respondError(w http.ResponseWriter, status int, mensaje string) {
	respondJSON(w, status, ErrorResponse{Error: mensaje})
}

// POST /api/docentes
func (h *Handler) CrearDocente(w http.ResponseWriter, r *http.Request) {
	var req CrearDocenteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "cuerpo de solicitud inválido")
		return
	}

	if req.Nombre == "" || req.Disciplina == "" {
		respondError(w, http.StatusBadRequest, "nombre y disciplina son obligatorios")
		return
	}

	teacher := domain.Teacher{
		ID:         uuid.NewString(),
		Nombre:     req.Nombre,
		Disciplina: req.Disciplina,
	}

	if err := h.repo.CrearDocente(r.Context(), teacher); err != nil {
		log.Printf("error al crear docente: %v", err)
		respondError(w, http.StatusInternalServerError, "no se pudo crear el docente")
		return
	}

	respondJSON(w, http.StatusCreated, teacher)
}

// POST /api/evaluaciones
func (h *Handler) EvaluarActividad(w http.ResponseWriter, r *http.Request) {
	var req EvaluarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "cuerpo de solicitud inválido")
		return
	}

	if req.Descripcion == "" || req.TeacherID == "" || req.Disciplina == "" {
		respondError(w, http.StatusBadRequest, "teacher_id, descripcion y disciplina son obligatorios")
		return
	}

	resultado, err := h.service.EvaluarActividad(r.Context(), req.Descripcion, req.Disciplina, req.TeacherID)
	if err != nil {
		log.Printf("error al evaluar actividad: %v", err)
		respondError(w, http.StatusInternalServerError, "no se pudo procesar la evaluación")
		return
	}

	respondJSON(w, http.StatusOK, resultado)
}

// GET /api/evaluaciones/{teacherID}
func (h *Handler) ObtenerHistorial(w http.ResponseWriter, r *http.Request, teacherID string) {
	if teacherID == "" {
		respondError(w, http.StatusBadRequest, "teacher_id es obligatorio")
		return
	}

	evaluaciones, err := h.repo.ObtenerPorTeacher(r.Context(), teacherID)
	if err != nil {
		log.Printf("error al obtener historial: %v", err)
		respondError(w, http.StatusInternalServerError, "no se pudo obtener el historial")
		return
	}

	respondJSON(w, http.StatusOK, evaluaciones)
}

// GET /api/matriz
func (h *Handler) ObtenerMatriz(w http.ResponseWriter, r *http.Request) {
	// Devuelve la definición estática de la matriz para que el frontend
	// pueda renderizar las 3 áreas × 4 niveles sin hardcodearlo en el cliente
	matriz := map[string]interface{}{
		"areas": []string{
			string(domain.AreaRecursosDigitales),
			string(domain.AreaEvaluacion),
			string(domain.AreaEmpoderamientoEstudiante),
		},
		"niveles": []map[string]interface{}{
			{"valor": int(domain.NivelInicial), "nombre": domain.NivelInicial.String()},
			{"valor": int(domain.NivelIntermedio), "nombre": domain.NivelIntermedio.String()},
			{"valor": int(domain.NivelAvanzado), "nombre": domain.NivelAvanzado.String()},
			{"valor": int(domain.NivelExperto), "nombre": domain.NivelExperto.String()},
		},
	}
	respondJSON(w, http.StatusOK, matriz)
}

// GET /api/metricas
func (h *Handler) ObtenerMetricas(w http.ResponseWriter, r *http.Request) {
	metricas, err := h.repo.MetricasPorArea(r.Context())
	if err != nil {
		log.Printf("error al obtener métricas: %v", err)
		respondError(w, http.StatusInternalServerError, "no se pudieron obtener las métricas")
		return
	}
	respondJSON(w, http.StatusOK, metricas)
}

// GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}