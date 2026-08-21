// internal/storage/sqlite.go

package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"digcompedu-api/internal/domain"
)

//go:embed schema.sql
var schemaSQL string

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(path string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("error al abrir sqlite: %w", err)
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, fmt.Errorf("error al migrar esquema: %w", err)
	}

	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) CrearDocente(ctx context.Context, t domain.Teacher) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO teachers (id, nombre, disciplina) VALUES (?, ?, ?)`,
		t.ID, t.Nombre, t.Disciplina,
	)
	return err
}

func (r *SQLiteRepository) Guardar(ctx context.Context, eval domain.Evaluation) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO evaluations (id, teacher_id, descripcion, disciplina, nota_etica) VALUES (?, ?, ?, ?, ?)`,
		eval.ID, eval.TeacherID, eval.Descripcion, eval.Disciplina, eval.Resultado.NotaEtica,
	)
	if err != nil {
		return fmt.Errorf("error al guardar evaluación: %w", err)
	}

	for _, r := range eval.Resultado.Resultados {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO area_resultados (evaluation_id, area, nivel, justificacion, sugerencia) VALUES (?, ?, ?, ?, ?)`,
			eval.ID, r.Area, r.Nivel, r.Justificacion, r.Sugerencia,
		)
		if err != nil {
			return fmt.Errorf("error al guardar resultado de área: %w", err)
		}
	}

	return tx.Commit()
}

func (r *SQLiteRepository) ObtenerPorTeacher(ctx context.Context, teacherID string) ([]domain.Evaluation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, teacher_id, descripcion, disciplina, created_at FROM evaluations WHERE teacher_id = ? ORDER BY created_at DESC`,
		teacherID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evals []domain.Evaluation
	for rows.Next() {
		var e domain.Evaluation
		var createdAt time.Time
		if err := rows.Scan(&e.ID, &e.TeacherID, &e.Descripcion, &e.Disciplina, &createdAt); err != nil {
			return nil, err
		}
		e.CreatedAt = createdAt
		evals = append(evals, e)
	}
	return evals, nil
}

func (r *SQLiteRepository) MetricasPorArea(ctx context.Context) (map[domain.Area]float64, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT area, AVG(nivel) FROM area_resultados GROUP BY area`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metricas := make(map[domain.Area]float64)
	for rows.Next() {
		var area domain.Area
		var promedio float64
		if err := rows.Scan(&area, &promedio); err != nil {
			return nil, err
		}
		metricas[area] = promedio
	}
	return metricas, nil
}