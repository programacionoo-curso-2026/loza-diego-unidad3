// internal/evaluator/service.go

package evaluator

import (
	"fmt"
	"context"
	"digcompedu-api/internal/domain"
	"digcompedu-api/internal/llm"
	"digcompedu-api/internal/storage"
)

type Service struct {
	llmClient llm.LLMProvider
	repo      storage.Repository
}

func NewService(llmClient llm.LLMProvider, repo storage.Repository) *Service {
	return &Service{llmClient: llmClient, repo: repo}
}

func (s *Service) EvaluarActividad(ctx context.Context, descripcion, disciplina, teacherID string) (*domain.MatrizResultado, error) {
	prompt := llm.ConstruirPrompt(descripcion, disciplina)

	respuesta, err := s.llmClient.Analizar(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("error al analizar con IA: %w", err)
	}

	resultado, err := parsearRespuestaLLM(respuesta)
	if err != nil {
		return nil, fmt.Errorf("error al interpretar respuesta: %w", err)
	}

	resultado.NotaEtica = "Esta evaluación es una sugerencia generada por IA. La decisión pedagógica final corresponde al docente."

	eval := domain.Evaluation{
		TeacherID:   teacherID,
		Descripcion: descripcion,
		Disciplina:  disciplina,
		Resultado:   resultado,
	}
	if err := s.repo.Guardar(ctx, eval); err != nil {
		return nil, fmt.Errorf("error al guardar: %w", err)
	}

	return resultado, nil
}