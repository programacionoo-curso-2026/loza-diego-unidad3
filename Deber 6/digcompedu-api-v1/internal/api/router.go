// internal/api/router.go

package api

import (
	"io/fs"
	"net/http"
)

func NewRouter(h *Handler, staticFS fs.FS) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("POST /api/docentes", h.CrearDocente)
	mux.HandleFunc("POST /api/evaluaciones", h.EvaluarActividad)
	mux.HandleFunc("GET /api/evaluaciones/{teacherID}", func(w http.ResponseWriter, r *http.Request) {
		h.ObtenerHistorial(w, r, r.PathValue("teacherID"))
	})
	mux.HandleFunc("GET /api/matriz", h.ObtenerMatriz)
	mux.HandleFunc("GET /api/metricas", h.ObtenerMetricas)

	mux.Handle("/", http.FileServerFS(staticFS)) // Go 1.22+: http.FileServerFS

	return withCORS(mux)
}

// withCORS permite que el frontend (probablemente en otro puerto/origen durante desarrollo) llame a la API
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}