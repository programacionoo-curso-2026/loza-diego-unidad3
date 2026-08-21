// internal/domain/matrix.go

package domain

type Area string

const (
	AreaRecursosDigitales       Area = "recursos_digitales"
	AreaEvaluacion               Area = "evaluacion"
	AreaEmpoderamientoEstudiante Area = "empoderamiento_estudiante"
)

type Nivel int

const (
	NivelInicial Nivel = iota + 1
	NivelIntermedio
	NivelAvanzado
	NivelExperto
)

func (n Nivel) String() string {
	nombres := map[Nivel]string{
		NivelInicial:    "Inicial",
		NivelIntermedio: "Intermedio",
		NivelAvanzado:   "Avanzado",
		NivelExperto:    "Experto",
	}
	return nombres[n]
}

type AreaResultado struct {
	Area          Area   `json:"area"`
	Nivel         Nivel  `json:"nivel"`
	Justificacion string `json:"justificacion"`
	Sugerencia    string `json:"sugerencia"`
}

type MatrizResultado struct {
	Resultados []AreaResultado `json:"resultados"`
	NotaEtica  string          `json:"nota_etica"`
}