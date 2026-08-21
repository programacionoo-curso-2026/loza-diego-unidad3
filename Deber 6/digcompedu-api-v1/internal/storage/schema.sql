-- internal/storage/schema.sql

CREATE TABLE IF NOT EXISTS teachers (
    id          TEXT PRIMARY KEY,
    nombre      TEXT NOT NULL,
    disciplina  TEXT NOT NULL,          -- ej. "Informática", "Automotriz"
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS evaluations (
    id          TEXT PRIMARY KEY,
    teacher_id  TEXT NOT NULL REFERENCES teachers(id),
    descripcion TEXT NOT NULL,          -- texto libre de la actividad
    disciplina  TEXT NOT NULL,
    nota_etica  TEXT NOT NULL,          -- disclaimer de transparencia
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS area_resultados (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    evaluation_id TEXT NOT NULL REFERENCES evaluations(id),
    area          TEXT NOT NULL CHECK (area IN ('recursos_digitales', 'evaluacion', 'empoderamiento_estudiante')),
    nivel         INTEGER NOT NULL CHECK (nivel BETWEEN 1 AND 4),
    justificacion TEXT,
    sugerencia    TEXT
);

CREATE INDEX IF NOT EXISTS idx_evaluations_teacher ON evaluations(teacher_id);
CREATE INDEX IF NOT EXISTS idx_area_resultados_eval ON area_resultados(evaluation_id);