// internal/domain/evaluation.go

package domain

import "time"

type Evaluation struct {
	ID          string           `json:"id"`
	TeacherID   string           `json:"teacher_id"`
	Descripcion string           `json:"descripcion"`
	Disciplina  string           `json:"disciplina"`
	Resultado   *MatrizResultado `json:"resultado,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}