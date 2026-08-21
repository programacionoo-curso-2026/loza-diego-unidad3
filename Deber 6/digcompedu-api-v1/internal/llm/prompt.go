// internal/llm/prompt.go

package llm

import "fmt"   // ← solo esto, quita el import de domain

func ConstruirPrompt(descripcion, disciplina string) string {
	return fmt.Sprintf(`Eres un asistente pedagógico. Analiza la siguiente actividad docente 
según el marco DigCompEdu, evaluando SOLO estas 3 áreas: recursos digitales, evaluación, 
y empoderamiento del estudiante. Para cada área, asigna un nivel (1-Inicial a 4-Experto) 
con justificación breve y una sugerencia concreta de mejora.

IMPORTANTE - principios éticos obligatorios:
- Esta es una SUGERENCIA de apoyo, nunca una calificación definitiva
- El juicio pedagógico final es siempre del docente humano
- Sé transparente: explica el razonamiento, no des veredictos sin justificar
- No hagas juicios sobre la persona, solo sobre la actividad descrita

Disciplina: %s
Actividad descrita por el docente: %s

Responde en formato JSON con la estructura: 
{"resultados": [{"area": "...", "nivel": N, "justificacion": "...", "sugerencia": "..."}]}`,
		disciplina, descripcion)
}