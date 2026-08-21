package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ==========================
// ESTRUCTURAS
// ==========================

type Producto struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
}

type Cliente struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Correo   string `json:"correo"`
	Telefono string `json:"telefono"`
}

type Pedido struct {
	ID         int     `json:"id"`
	ClienteID  int     `json:"cliente_id"`
	ProductoID int     `json:"producto_id"`
	Cantidad   int     `json:"cantidad"`
	Total      float64 `json:"total"`
	Estado     string  `json:"estado"`
}

// ==========================
// DATOS INICIALES
// ==========================

var productos = []Producto{
	{
		ID:          1,
		Nombre:      "Laptop Gamer",
		Descripcion: "Laptop para gaming y programación",
		Precio:      950,
		Stock:       10,
	},
	{
		ID:          2,
		Nombre:      "Mouse Gamer",
		Descripcion: "Mouse RGB",
		Precio:      35.50,
		Stock:       25,
	},
	{
		ID:          3,
		Nombre:      "Teclado Mecánico",
		Descripcion: "Teclado mecánico RGB",
		Precio:      65,
		Stock:       15,
	},
}

var clientes = []Cliente{
	{
		ID:       1,
		Nombre:   "Carlos Pérez",
		Correo:   "carlos@gmail.com",
		Telefono: "0991111111",
	},
	{
		ID:       2,
		Nombre:   "Ana López",
		Correo:   "ana@gmail.com",
		Telefono: "0982222222",
	},
}

var pedidos = []Pedido{}

// ==========================
// INICIO
// ==========================

func inicio(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "API Sistema de Gestión de E-Commerce",
	})
}

// ==========================
// PRODUCTOS
// ==========================

// GET /productos
// POST /productos
func manejarProductos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:

		json.NewEncoder(w).Encode(productos)

	case http.MethodPost:

		var nuevoProducto Producto

		err := json.NewDecoder(r.Body).Decode(&nuevoProducto)

		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		nuevoProducto.ID = obtenerNuevoIDProducto()

		productos = append(productos, nuevoProducto)

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(nuevoProducto)

	default:

		http.Error(
			w,
			"Método no permitido",
			http.StatusMethodNotAllowed,
		)
	}
}

// GET /productos/{id}
// PUT /productos/{id}
// DELETE /productos/{id}
func manejarProductoPorID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	idTexto := strings.TrimPrefix(
		r.URL.Path,
		"/productos/",
	)

	id, err := strconv.Atoi(idTexto)

	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:

		for _, producto := range productos {

			if producto.ID == id {

				json.NewEncoder(w).Encode(producto)
				return
			}
		}

		http.Error(
			w,
			"Producto no encontrado",
			http.StatusNotFound,
		)

	case http.MethodPut:

		var productoActualizado Producto

		err := json.NewDecoder(
			r.Body,
		).Decode(&productoActualizado)

		if err != nil {

			http.Error(
				w,
				"JSON inválido",
				http.StatusBadRequest,
			)

			return
		}

		for i, producto := range productos {

			if producto.ID == id {

				productoActualizado.ID = id

				productos[i] = productoActualizado

				json.NewEncoder(w).Encode(
					productoActualizado,
				)

				return
			}
		}

		http.Error(
			w,
			"Producto no encontrado",
			http.StatusNotFound,
		)

	case http.MethodDelete:

		for i, producto := range productos {

			if producto.ID == id {

				productos = append(
					productos[:i],
					productos[i+1:]...,
				)

				json.NewEncoder(w).Encode(
					map[string]string{
						"mensaje": "Producto eliminado correctamente",
					},
				)

				return
			}
		}

		http.Error(
			w,
			"Producto no encontrado",
			http.StatusNotFound,
		)

	default:

		http.Error(
			w,
			"Método no permitido",
			http.StatusMethodNotAllowed,
		)
	}
}

// ==========================
// CLIENTES
// ==========================

// GET /clientes
// POST /clientes
func manejarClientes(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:

		json.NewEncoder(w).Encode(clientes)

	case http.MethodPost:

		var nuevoCliente Cliente

		err := json.NewDecoder(
			r.Body,
		).Decode(&nuevoCliente)

		if err != nil {

			http.Error(
				w,
				"JSON inválido",
				http.StatusBadRequest,
			)

			return
		}

		nuevoCliente.ID = obtenerNuevoIDCliente()

		clientes = append(
			clientes,
			nuevoCliente,
		)

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(
			nuevoCliente,
		)

	default:

		http.Error(
			w,
			"Método no permitido",
			http.StatusMethodNotAllowed,
		)
	}
}

// ==========================
// PEDIDOS
// ==========================

// GET /pedidos
// POST /pedidos
func manejarPedidos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	// SERVICIO 9
	case http.MethodGet:

		json.NewEncoder(w).Encode(pedidos)

	// SERVICIO 8
	case http.MethodPost:

		var nuevoPedido Pedido

		err := json.NewDecoder(
			r.Body,
		).Decode(&nuevoPedido)

		if err != nil {

			http.Error(
				w,
				"JSON inválido",
				http.StatusBadRequest,
			)

			return
		}

		// Verificar cliente
		clienteExiste := false

		for _, cliente := range clientes {

			if cliente.ID == nuevoPedido.ClienteID {

				clienteExiste = true
				break
			}
		}

		if !clienteExiste {

			http.Error(
				w,
				"Cliente no encontrado",
				http.StatusNotFound,
			)

			return
		}

		// Buscar producto
		productoEncontrado := false

		for i, producto := range productos {

			if producto.ID == nuevoPedido.ProductoID {

				productoEncontrado = true

				if nuevoPedido.Cantidad <= 0 {

					http.Error(
						w,
						"La cantidad debe ser mayor a 0",
						http.StatusBadRequest,
					)

					return
				}

				if producto.Stock < nuevoPedido.Cantidad {

					http.Error(
						w,
						"Stock insuficiente",
						http.StatusBadRequest,
					)

					return
				}

				// Calcular total
				nuevoPedido.Total =
					producto.Precio *
						float64(nuevoPedido.Cantidad)

				// Restar stock
				productos[i].Stock -= nuevoPedido.Cantidad

				break
			}
		}

		if !productoEncontrado {

			http.Error(
				w,
				"Producto no encontrado",
				http.StatusNotFound,
			)

			return
		}

		nuevoPedido.ID = obtenerNuevoIDPedido()

		nuevoPedido.Estado = "Pendiente"

		pedidos = append(
			pedidos,
			nuevoPedido,
		)

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(
			nuevoPedido,
		)

	default:

		http.Error(
			w,
			"Método no permitido",
			http.StatusMethodNotAllowed,
		)
	}
}

// PUT /pedidos/{id}/estado
func manejarPedidoPorID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {

		http.Error(
			w,
			"Método no permitido",
			http.StatusMethodNotAllowed,
		)

		return
	}

	ruta := strings.TrimPrefix(
		r.URL.Path,
		"/pedidos/",
	)

	partes := strings.Split(ruta, "/")

	if len(partes) != 2 || partes[1] != "estado" {

		http.Error(
			w,
			"Ruta inválida",
			http.StatusBadRequest,
		)

		return
	}

	id, err := strconv.Atoi(partes[0])

	if err != nil {

		http.Error(
			w,
			"ID inválido",
			http.StatusBadRequest,
		)

		return
	}

	var datos struct {
		Estado string `json:"estado"`
	}

	err = json.NewDecoder(
		r.Body,
	).Decode(&datos)

	if err != nil {

		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)

		return
	}

	if datos.Estado == "" {

		http.Error(
			w,
			"Debe indicar un estado",
			http.StatusBadRequest,
		)

		return
	}

	for i, pedido := range pedidos {

		if pedido.ID == id {

			pedidos[i].Estado = datos.Estado

			json.NewEncoder(w).Encode(
				pedidos[i],
			)

			return
		}
	}

	http.Error(
		w,
		"Pedido no encontrado",
		http.StatusNotFound,
	)
}

// ==========================
// FUNCIONES PARA IDS
// ==========================

func obtenerNuevoIDProducto() int {

	maxID := 0

	for _, producto := range productos {

		if producto.ID > maxID {
			maxID = producto.ID
		}
	}

	return maxID + 1
}

func obtenerNuevoIDCliente() int {

	maxID := 0

	for _, cliente := range clientes {

		if cliente.ID > maxID {
			maxID = cliente.ID
		}
	}

	return maxID + 1
}

func obtenerNuevoIDPedido() int {

	maxID := 0

	for _, pedido := range pedidos {

		if pedido.ID > maxID {
			maxID = pedido.ID
		}
	}

	return maxID + 1
}

// ==========================
// MAIN
// ==========================

func main() {

	http.HandleFunc(
		"/",
		inicio,
	)

	http.HandleFunc(
		"/productos",
		manejarProductos,
	)

	http.HandleFunc(
		"/productos/",
		manejarProductoPorID,
	)

	http.HandleFunc(
		"/clientes",
		manejarClientes,
	)

	http.HandleFunc(
		"/pedidos",
		manejarPedidos,
	)

	http.HandleFunc(
		"/pedidos/",
		manejarPedidoPorID,
	)

	fmt.Println("======================================")
	fmt.Println(" SISTEMA DE GESTIÓN DE E-COMMERCE")
	fmt.Println("======================================")

	fmt.Println("")
	fmt.Println("Servidor iniciado correctamente")
	fmt.Println("http://localhost:8080")

	fmt.Println("")
	fmt.Println("SERVICIOS WEB DISPONIBLES")
	fmt.Println("--------------------------------------")
	fmt.Println("1.  GET    /productos")
	fmt.Println("2.  POST   /productos")
	fmt.Println("3.  GET    /productos/{id}")
	fmt.Println("4.  PUT    /productos/{id}")
	fmt.Println("5.  DELETE /productos/{id}")
	fmt.Println("6.  GET    /clientes")
	fmt.Println("7.  POST   /clientes")
	fmt.Println("8.  POST   /pedidos")
	fmt.Println("9.  GET    /pedidos")
	fmt.Println("10. PUT    /pedidos/{id}/estado")
	fmt.Println("--------------------------------------")

	err := http.ListenAndServe(
		":8080",
		nil,
	)

	if err != nil {

		fmt.Println(
			"Error al iniciar servidor:",
			err,
		)
	}
}
