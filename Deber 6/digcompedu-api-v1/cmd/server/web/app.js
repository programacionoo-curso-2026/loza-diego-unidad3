// web/app.js

const API_BASE = "http://localhost:8080/api";

const NOMBRES_AREA = {
    recursos_digitales: "Recursos Digitales",
    evaluacion: "Evaluación",
    empoderamiento_estudiante: "Empoderamiento del Estudiante"
};

let teacherIdActual = null;

const form = document.getElementById("form-evaluacion");
const btnEnviar = document.getElementById("btn-enviar");
const loading = document.getElementById("loading");
const resultadosDiv = document.getElementById("resultados");
const errorDiv = document.getElementById("error");
const areasContainer = document.getElementById("areas-container");
const notaEticaDiv = document.getElementById("nota-etica");

form.addEventListener("submit", async (e) => {
    e.preventDefault();
    errorDiv.textContent = "";
    resultadosDiv.style.display = "none";

    const nombre = document.getElementById("nombre").value.trim();
    const disciplina = document.getElementById("disciplina").value;
    const descripcion = document.getElementById("descripcion").value.trim();

    btnEnviar.disabled = true;
    loading.style.display = "block";

    try {
        // Crear docente si aún no existe en esta sesión
        if (!teacherIdActual) {
            const docenteResp = await fetch(`${API_BASE}/docentes`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ nombre, disciplina })
            });
            if (!docenteResp.ok) throw new Error("No se pudo registrar el docente");
            const docente = await docenteResp.json();
            teacherIdActual = docente.id;
        }

        const evalResp = await fetch(`${API_BASE}/evaluaciones`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                teacher_id: teacherIdActual,
                descripcion,
                disciplina
            })
        });

        if (!evalResp.ok) {
            const errData = await evalResp.json();
            throw new Error(errData.error || "Error al procesar la evaluación");
        }

        const resultado = await evalResp.json();
        renderizarResultado(resultado);

    } catch (err) {
        errorDiv.textContent = err.message;
    } finally {
        btnEnviar.disabled = false;
        loading.style.display = "none";
    }
});

function renderizarResultado(resultado) {
    areasContainer.innerHTML = "";

    resultado.resultados.forEach(r => {
        const div = document.createElement("div");
        div.className = "resultado-area";
        div.innerHTML = `
            <h3>${NOMBRES_AREA[r.area] || r.area}</h3>
            <span class="nivel-badge">Nivel ${r.nivel}</span>
            <p style="font-size:0.9rem; margin-top:0.4rem;"><strong>Justificación:</strong> ${r.justificacion}</p>
            <p style="font-size:0.9rem; margin-top:0.3rem;"><strong>Sugerencia:</strong> ${r.sugerencia}</p>
        `;
        areasContainer.appendChild(div);
    });

    notaEticaDiv.textContent = "⚠️ " + resultado.nota_etica;
    resultadosDiv.style.display = "block";
}