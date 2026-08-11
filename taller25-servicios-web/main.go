package main

import (
	"fmt"
	"log"
	"net/http"
)

const page = `<!DOCTYPE html>
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
    display: flex;
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
    color: #ffcc00;
    font-weight: bold;
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
</style>
</head>
<body>
  <div class="card">
    <div class="code">for i <span id="blank" class="blank">___</span> 0; i &lt; 10; i++</div>
    <div class="variable" id="variable">Elige un operador</div>
    <div class="buttons">
      <button onclick="elegir('=')">=</button>
      <button onclick="elegir('<')">&lt;</button>
      <button onclick="elegir(':=')">:=</button>
    </div>
  </div>

  <script>
    function elegir(valor) {
      document.getElementById('blank').textContent = valor;
      document.getElementById('variable').textContent = 'Elegiste: ' + valor;
    }
  </script>
</body>
</html>`

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, page)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}