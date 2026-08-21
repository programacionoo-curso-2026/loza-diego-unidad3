// internal/domain/teacher.go

package domain

import "time"

type Teacher struct {
	ID         string    `json:"id"`
	Nombre     string    `json:"nombre"`
	Disciplina string    `json:"disciplina"`
	CreatedAt  time.Time `json:"created_at"`
}