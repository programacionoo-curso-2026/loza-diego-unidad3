// internal/storage/repository.go

package storage

import (
	"context"
	"digcompedu-api/internal/domain"
)

type Repository interface {
	CrearDocente(ctx context.Context, t domain.Teacher) error
	Guardar(ctx context.Context, eval domain.Evaluation) error
	ObtenerPorTeacher(ctx context.Context, teacherID string) ([]domain.Evaluation, error)
	MetricasPorArea(ctx context.Context) (map[domain.Area]float64, error) // promedio de nivel por área
}