package domain

import "fmt"

// EvaluacionDocente representa una evaluación complementaria del docente.
type EvaluacionDocente struct {
	DominioTema int `json:"dominio_tema"`
	Claridad    int `json:"claridad"`
	Puntualidad int `json:"puntualidad"`
}

// CalcularPromedio valida los criterios y devuelve el promedio.
func (e EvaluacionDocente) CalcularPromedio() (float64, error) {
	valores := []int{e.DominioTema, e.Claridad, e.Puntualidad}
	for _, valor := range valores {
		if valor < 1 || valor > 5 {
			return 0, fmt.Errorf("cada criterio debe estar entre 1 y 5")
		}
	}

	total := e.DominioTema + e.Claridad + e.Puntualidad
	return float64(total) / float64(len(valores)), nil
}
