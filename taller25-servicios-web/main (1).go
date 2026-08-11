package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// Pregunta representa un ejercicio con su enunciado, sus 3 opciones y la solución correcta (A, B o C).
type Pregunta struct {
	Pregunta string `json:"pregunta"`
	OpcionA  string `json:"opcion_a"`
	OpcionB  string `json:"opcion_b"`
	OpcionC  string `json:"opcion_c"`
	Solucion string `json:"solucion"`
}

// Aquí se parametriza el arreglo de preguntas.
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
    gap: 14px;
  }
  #pantalla-final h2 {
    margin: 0;
    color: #333;
  }
  #listaFinal {
    width: 100%;
    box-sizing: border-box;
    text-align: left;
    background: #f7f7f7;
    border-radius: 8px;
    padding: 16px 16px 16px 36px;
    margin: 0;
  }
  #listaFinal li {
    font-family: 'Courier New', monospace;
    margin-bottom: 8px;
    color: #333;
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
      <ul id="listaFinal"></ul>
    </div>

  </div>

  <script>
    const preguntas = {{.Preguntas}};
    let nombre = '';
    let indice = 0;
    const respondidas = [];

    function iniciar() {
      const valor = document.getElementById('inputNombre').value.trim();
      if (!valor) {
        document.getElementById('inputNombre').focus();
        return;
      }
      nombre = valor;
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
            mostrarFinal();
          }
        }, 700);
      } else {
        blank.classList.remove('correcto');
        blank.classList.add('incorrecto');
        document.getElementById('variable').textContent = 'Incorrecto, intenta de nuevo';
      }
    }

    function mostrarFinal() {
      document.getElementById('pantalla-pregunta').style.display = 'none';
      document.getElementById('pantalla-final').style.display = 'flex';
      document.getElementById('tituloFinal').textContent =
        '¡Felicidades, ' + nombre + '!';

      const lista = document.getElementById('listaFinal');
      lista.innerHTML = '';
      respondidas.forEach(texto => {
        const li = document.createElement('li');
        li.textContent = texto;
        lista.appendChild(li);
      });
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
	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}