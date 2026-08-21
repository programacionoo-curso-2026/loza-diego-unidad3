# Sistema de Gestión de E-Commerce

## Descripción del proyecto

Este proyecto consiste en el desarrollo de un Sistema de Gestión de E-Commerce utilizando el lenguaje de programación Go.

El sistema implementa servicios web que permiten gestionar productos, clientes y pedidos mediante una API REST. La información es enviada y recibida utilizando el formato JSON.

El proyecto integra diferentes conocimientos revisados durante la asignatura, como estructuras de datos, funciones, programación modular, servicios web, métodos HTTP y serialización de información.

## Objetivo

Desarrollar un sistema de gestión de comercio electrónico que permita administrar productos, clientes y pedidos mediante servicios web, utilizando Go y JSON.

## Funcionalidades principales

El sistema permite:

* Registrar productos.
* Consultar todos los productos.
* Consultar productos mediante su ID.
* Actualizar información de productos.
* Eliminar productos.
* Registrar clientes.
* Consultar clientes.
* Crear pedidos.
* Consultar pedidos.
* Cambiar el estado de un pedido.
* Validar disponibilidad de stock.
* Calcular automáticamente el valor total de un pedido.

## Tecnologías utilizadas

* Go
* Servicios Web REST
* HTTP
* JSON
* Visual Studio Code
* Git
* GitHub

## Servicios web implementados

| Número | Método | Ruta                   | Funcionalidad                     |
| ------ | ------ | ---------------------- | --------------------------------- |
| 1      | GET    | `/productos`           | Consultar todos los productos     |
| 2      | POST   | `/productos`           | Registrar un nuevo producto       |
| 3      | GET    | `/productos/{id}`      | Consultar un producto por ID      |
| 4      | PUT    | `/productos/{id}`      | Actualizar un producto            |
| 5      | DELETE | `/productos/{id}`      | Eliminar un producto              |
| 6      | GET    | `/clientes`            | Consultar clientes                |
| 7      | POST   | `/clientes`            | Registrar un cliente              |
| 8      | POST   | `/pedidos`             | Crear un pedido                   |
| 9      | GET    | `/pedidos`             | Consultar pedidos                 |
| 10     | PUT    | `/pedidos/{id}/estado` | Actualizar el estado de un pedido |

## Ejecución del proyecto

Para ejecutar el programa se debe tener Go instalado.

Abrir una terminal dentro de la carpeta del proyecto y ejecutar:

```bash
go run main.go
```

El servidor se ejecutará en:

```text
http://localhost:8080
```

Para consultar los productos se puede acceder desde el navegador a:

```text
http://localhost:8080/productos
```

Para consultar los clientes:

```text
http://localhost:8080/clientes
```

Para consultar los pedidos:

```text
http://localhost:8080/pedidos
```

## Serialización JSON

La información del sistema es serializada mediante JSON.

Ejemplo de un producto:

```json
{
  "id": 1,
  "nombre": "Laptop Gamer",
  "descripcion": "Laptop para gaming y programación",
  "precio": 950,
  "stock": 10
}
```

Ejemplo de un pedido:

```json
{
  "id": 1,
  "cliente_id": 1,
  "producto_id": 1,
  "cantidad": 2,
  "total": 1900,
  "estado": "Pendiente"
}
```

## Integración de conocimientos de la asignatura

El proyecto integra diferentes conceptos revisados durante las unidades de la materia.

Se utilizan estructuras para representar productos, clientes y pedidos.

También se utilizan slices para almacenar temporalmente la información, funciones para organizar el programa y estructuras de control para validar las operaciones.

Además, se implementaron servicios web utilizando métodos HTTP como GET, POST, PUT y DELETE.

La serialización de datos se realiza mediante JSON, permitiendo el intercambio de información entre el cliente y el servidor.

## Visualización del futuro

En el futuro, el sistema podría evolucionar incorporando nuevas tecnologías y funcionalidades.

Entre las posibles mejoras se encuentran:

* Implementación de una base de datos.
* Desarrollo de una interfaz web.
* Desarrollo de una aplicación móvil.
* Sistema de autenticación de usuarios.
* Integración de pagos electrónicos.
* Carrito de compras.
* Gestión avanzada de inventario.
* Uso de servicios en la nube.
* Inteligencia artificial para recomendaciones de productos.
* Análisis de comportamiento de los clientes.

La inteligencia artificial podría utilizarse para analizar las compras de los usuarios y recomendar productos relacionados con sus intereses.

También sería posible implementar servicios en la nube para permitir que el sistema pueda ser utilizado desde diferentes dispositivos y lugares.

## Conclusión

El desarrollo de este proyecto permitió aplicar diferentes conocimientos de programación y servicios web en una aplicación práctica.

Durante el desarrollo se trabajó con estructuras de datos, funciones, métodos HTTP, rutas, validaciones y serialización mediante JSON.

El proyecto demuestra cómo Go puede utilizarse para desarrollar servicios web eficientes y cómo estos servicios pueden convertirse en la base de aplicaciones modernas de comercio electrónico.

Como mejora futura, se podrían incorporar bases de datos, interfaces gráficas, seguridad, sistemas de pago y tecnologías de inteligencia artificial.
