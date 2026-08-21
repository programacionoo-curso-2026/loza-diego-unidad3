// internal/evaluator/parser.go

package evaluator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"digcompedu-api/internal/domain"
)

var jsonBlockRegex = regexp.MustCompile(`(?s)\{.*\}`)

// areasValidas mapea variantes comunes que puede devolver el LLM al valor canónico exacto
var areasValidas = map[string]domain.Area{
	"recursos_digitales":       domain.AreaRecursosDigitales,
	"recursos digitales":       domain.AreaRecursosDigitales,
	"evaluacion":               domain.AreaEvaluacion,
	"evaluación":               domain.AreaEvaluacion,
	"empoderamiento_estudiante": domain.AreaEmpoderamientoEstudiante,
	"empoderamiento del estudiante": domain.AreaEmpoderamientoEstudiante,
	"empoderamiento estudiante": domain.AreaEmpoderamientoEstudiante,
}

func normalizarArea(raw string) (domain.Area, error) {
	clave := strings.ToLower(strings.TrimSpace(raw))
	if area, ok := areasValidas[clave]; ok {
		return area, nil
	}
	return "", fmt.Errorf("área no reconocida: %q", raw)
}

func parsearRespuestaLLM(respuesta string) (*domain.MatrizResultado, error) {
	match := jsonBlockRegex.FindString(respuesta)
	if match == "" {
		return nil, fmt.Errorf("no se encontró JSON en la respuesta del LLM")
	}

	var raw struct {
		Resultados []struct {
			Area          string `json:"area"`
			Nivel         int    `json:"nivel"`
			Justificacion string `json:"justificacion"`
			Sugerencia    string `json:"sugerencia"`
		} `json:"resultados"`
	}

	if err := json.Unmarshal([]byte(match), &raw); err != nil {
		return nil, fmt.Errorf("error al deserializar resultado: %w", err)
	}

	if len(raw.Resultados) != 3 {
		return nil, fmt.Errorf("se esperaban 3 áreas evaluadas, se obtuvieron %d", len(raw.Resultados))
	}

	resultado := &domain.MatrizResultado{}
	for _, r := range raw.Resultados {
		areaNormalizada, err := normalizarArea(r.Area)
		if err != nil {
			return nil, fmt.Errorf("error al normalizar área: %w (valor crudo del LLM: %q)", err, r.Area)
		}

		nivel := domain.Nivel(r.Nivel)
		if nivel < domain.NivelInicial || nivel > domain.NivelExperto {
			return nil, fmt.Errorf("nivel fuera de rango: %d", r.Nivel)
		}

		resultado.Resultados = append(resultado.Resultados, domain.AreaResultado{
			Area:          areaNormalizada,
			Nivel:         nivel,
			Justificacion: r.Justificacion,
			Sugerencia:    r.Sugerencia,
		})
	}

	return resultado, nil
}