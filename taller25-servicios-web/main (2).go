package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"sync"
)

// ---------- Preguntas ----------

type Pregunta struct {
	Pregunta string `json:"pregunta"`
	OpcionA  string `json:"opcion_a"`
	OpcionB  string `json:"opcion_b"`
	OpcionC  string `json:"opcion_c"`
	Solucion string `json:"solucion"`
}

var preguntas = []Pregunta{
	{
		Pregunta: "for i &blank 0; i &lt; 10; i++",
		OpcionA:  "=",
		OpcionB:  "<",
		OpcionC:  ":=",
		Solucion: "C",
	},
	{
		Pregunta: `var nombre &blank "Sofia"`,
		OpcionA:  "=",
		OpcionB:  "<",
		OpcionC:  ":=",
		Solucion: "A",
	},
}

// ---------- Resultados (ranking en memoria) ----------

type Resultado struct {
	Nombre      string `json:"nombre"`
	TiempoMs    int64  `json:"tiempo_ms"`
	TiempoTexto string `json:"tiempo_texto"`
}

var (
	mu         sync.Mutex
	resultados []Resultado
)

func formatearTiempo(ms int64) string {
	totalSeg := ms / 1000
	min := totalSeg / 60
	seg := totalSeg % 60
	return fmt.Sprintf("%02d:%02d", min, seg)
}

func guardarResultadoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var entrada struct {
		Nombre   string `json:"nombre"`
		TiempoMs int64  `json:"tiempo_ms"`
	}
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}

	nuevo := Resultado{
		Nombre:      entrada.Nombre,
		TiempoMs:    entrada.TiempoMs,
		TiempoTexto: formatearTiempo(entrada.TiempoMs),
	}

	mu.Lock()
	resultados = append(resultados, nuevo)
	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].TiempoMs < resultados[j].TiempoMs
	})
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func resultadosHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultados)
}

// ---------- Página ----------

const pageTemplate = `<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="UTF-8">
<title>Completa el for</title>
<style>
  body {
    font-family: -apple-system, Arial, sans-serif;
    background: #f2f2f2;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    margin: 0;
  }
  .card {
    background: #fff;
    padding: 32px 40px;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
    width: 480px;
    box-sizing: border-box;
  }

  /* ---------- Pantalla inicio ---------- */
  #pantalla-inicio {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 18px;
  }
  #pantalla-inicio h2 {
    margin: 0;
    color: #333;
  }
  #inputNombre {
    width: 100%;
    box-sizing: border-box;
    font-size: 20px;
    padding: 12px;
    border-radius: 8px;
    border: 1px solid #ccc;
    text-align: center;
  }
  #btnIniciar {
    width: 100%;
    font-size: 22px;
    padding: 14px;
    border: none;
    border-radius: 8px;
    background: #007acc;
    color: #fff;
    cursor: pointer;
  }
  #btnIniciar:hover {
    background: #005f99;
  }

  /* ---------- Pantalla pregunta ---------- */
  #pantalla-pregunta {
    display: none;
    flex-direction: column;
    align-items: center;
  }
  .code {
    width: 100%;
    box-sizing: border-box;
    height: 160px;
    font-family: 'Courier New', monospace;
    font-size: 24px;
    background: #1e1e1e;
    color: #d4d4d4;
    border-radius: 8px;
    margin-bottom: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
  }
  .blank {
    font-weight: bold;
    color: #ffcc00;
  }
  .blank.correcto {
    color: #4caf50;
  }
  .blank.incorrecto {
    color: #e74c3c;
  }
  .variable {
    width: 100%;
    font-size: 26px;
    color: #555;
    margin-bottom: 20px;
    min-height: 40px;
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
  }
  .buttons {
    display: flex;
    justify-content: center;
    gap: 14px;
  }
  .buttons button {
    width: 143px;
    height: 143px;
    display: flex;
    justify-content: center;
    align-items: center;
    font-family: 'Courier New', monospace;
    font-size: 42px;
    padding: 0;
    border: none;
    border-radius: 8px;
    background: #007acc;
    color: #fff;
    cursor: pointer;
  }
  .buttons button:hover {
    background: #005f99;
  }

  /* ---------- Pantalla final ---------- */
  #pantalla-final {
    display: none;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 10px;
  }
  #pantalla-final h2 {
    margin: 0;
    color: #333;
  }
  #tiempoFinal {
    font-size: 20px;
    color: #007acc;
    font-weight: bold;
    margin-bottom: 6px;
  }
  #listaFinal {
    width: 100%;
    box-sizing: border-box;
    text-align: left;
    background: #f7f7f7;
    border-radius: 8px;
    padding: 16px 16px 16px 36px;
    margin: 0 0 10px 0;
  }
  #listaFinal li {
    font-family: 'Courier New', monospace;
    margin-bottom: 8px;
    color: #333;
  }
  #tituloRanking {
    margin: 10px 0 0 0;
    color: #333;
  }
  #ranking {
    width: 100%;
    box-sizing: border-box;
    text-align: left;
    background: #f7f7f7;
    border-radius: 8px;
    padding: 12px 20px;
    margin: 0;
    list-style: none;
  }
  #ranking li {
    display: flex;
    justify-content: space-between;
    padding: 6px 0;
    border-bottom: 1px solid #e0e0e0;
    font-size: 16px;
    color: #333;
  }
  #ranking li:last-child {
    border-bottom: none;
  }
  #ranking li.mejor {
    font-weight: bold;
    color: #4caf50;
  }
</style>
</head>
<body>
  <div class="card">

    <!-- Pantalla 1: pedir nombre -->
    <div id="pantalla-inicio">
      <h2>¿Cuál es tu nombre?</h2>
      <input id="inputNombre" type="text" placeholder="Escribe tu nombre">
      <button id="btnIniciar" onclick="iniciar()">Iniciar</button>
    </div>

    <!-- Pantalla 2: pregunta -->
    <div id="pantalla-pregunta">
      <div class="code" id="code"></div>
      <div class="variable" id="variable">Elige un operador</div>
      <div class="buttons">
        <button id="btnA"></button>
        <button id="btnB"></button>
        <button id="btnC"></button>
      </div>
    </div>

    <!-- Pantalla 3: final -->
    <div id="pantalla-final">
      <h2 id="tituloFinal"></h2>
      <div id="tiempoFinal"></div>
      <!--ul id="listaFinal"></ul-->
      <h3 id="tituloRanking">Ranking (menor tiempo primero)</h3>
      <ol id="ranking"></ol>
    </div>

  </div>

  <script>
    const preguntas = {{.Preguntas}};
    let nombre = '';
    let indice = 0;
    let inicioTimestamp = 0;
    const respondidas = [];

    function iniciar() {
      const valor = document.getElementById('inputNombre').value.trim();
      if (!valor) {
        document.getElementById('inputNombre').focus();
        return;
      }
      nombre = valor;
      inicioTimestamp = Date.now();
      document.getElementById('pantalla-inicio').style.display = 'none';
      document.getElementById('pantalla-pregunta').style.display = 'flex';
      cargarPregunta(indice);
    }

    function cargarPregunta(i) {
      const p = preguntas[i];

      const html = p.pregunta.replace(
        '&blank',
        '<span id="blank" class="blank">___</span>'
      );
      document.getElementById('code').innerHTML = html;
      document.getElementById('variable').textContent = 'Elige un operador';

      const btnA = document.getElementById('btnA');
      const btnB = document.getElementById('btnB');
      const btnC = document.getElementById('btnC');

      btnA.textContent = p.opcion_a;
      btnB.textContent = p.opcion_b;
      btnC.textContent = p.opcion_c;

      btnA.onclick = () => elegir(p.opcion_a, 'A');
      btnB.onclick = () => elegir(p.opcion_b, 'B');
      btnC.onclick = () => elegir(p.opcion_c, 'C');
    }

    function elegir(valor, letra) {
      const p = preguntas[indice];
      const blank = document.getElementById('blank');
      blank.textContent = valor;

      if (letra === p.solucion) {
        blank.classList.remove('incorrecto');
        blank.classList.add('correcto');
        document.getElementById('variable').textContent = '¡Correcto!';

        const textoCompleto = p.pregunta
          .replace('&blank', valor)
          .replace(/&lt;/g, '<');
        respondidas.push(textoCompleto);

        setTimeout(() => {
          indice++;
          if (indice < preguntas.length) {
            cargarPregunta(indice);
          } else {
            finalizar();
          }
        }, 700);
      } else {
        blank.classList.remove('correcto');
        blank.classList.add('incorrecto');
        document.getElementById('variable').textContent = 'Incorrecto, intenta de nuevo';
      }
    }

    function formatearTiempo(ms) {
      const totalSeg = Math.floor(ms / 1000);
      const min = Math.floor(totalSeg / 60);
      const seg = totalSeg % 60;
      const pad = n => String(n).padStart(2, '0');
      return pad(min) + ':' + pad(seg);
    }

    async function finalizar() {
      document.getElementById('pantalla-pregunta').style.display = 'none';
      document.getElementById('pantalla-final').style.display = 'flex';

      const tiempoMs = Date.now() - inicioTimestamp;

      document.getElementById('tituloFinal').textContent =
        '¡Felicidades, ' + nombre + '!';
      document.getElementById('tiempoFinal').textContent =
        'Tu tiempo: ' + formatearTiempo(tiempoMs);
      /*
      const lista = document.getElementById('listaFinal');
      lista.innerHTML = '';
      respondidas.forEach(texto => {
        const li = document.createElement('li');
        li.textContent = texto;
        lista.appendChild(li);
      });
      */

      // Guarda el resultado en el servidor
      try {
        await fetch('/api/resultado', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ nombre: nombre, tiempo_ms: tiempoMs })
        });
      } catch (e) {
        console.error('No se pudo guardar el resultado', e);
      }

      // Trae y muestra el ranking ya ordenado por el servidor
      try {
        const resp = await fetch('/api/resultados');
        const datos = await resp.json();
        const ranking = document.getElementById('ranking');
        ranking.innerHTML = '';
        datos.forEach((r, idx) => {
          const li = document.createElement('li');
          if (idx === 0) li.classList.add('mejor');
          li.innerHTML = '<span>' + (idx + 1) + '. ' + r.nombre + '</span><span>' + r.tiempo_texto + '</span>';
          ranking.appendChild(li);
        });
      } catch (e) {
        console.error('No se pudo obtener el ranking', e);
      }
    }
  </script>
</body>
</html>`

var tmpl = template.Must(template.New("page").Parse(pageTemplate))

func handler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(preguntas)
	if err != nil {
		http.Error(w, "error generando preguntas", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]template.JS{
		"Preguntas": template.JS(data),
	})
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/resultado", guardarResultadoHandler)
	http.HandleFunc("/api/resultados", resultadosHandler)
	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}