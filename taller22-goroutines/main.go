package taller22goroutines
package main

import (
	"fmt"
)

func main() {
	fmt.Print("Todas las operaciones completadas. Finalizando\n")
}
package main

import (
  "fmt"
  "math/rand"
  "time"
)

type Order struct {
  ID     int
  Status string
}

func main() {
  orders := generateOrders(20)
  fmt.Print("Todas las operaciones completadas. Finalizando\n")
}
func generateOrders(count int) []*Order {
  orders := make([]*Order, count)
  for i := 0; i < count; i++ {
    orders[i] = &Order{
      ID: i + 1, Status: "pending",
    }
  }
  return orders
}
func processOrders(orders []*Order) {
  for _, order := range orders {
    time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    fmt.Printf("Procesando orden %d\n", order.ID)
  }
}
func updateOrderStatuses(orders []*Order) {
  for _, order := range orders {
    time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    status := []string{
      "Procesando", "Despachando", "Entregado",
    }[rand.Intn(3)]
    order.Status = status
    fmt.Printf("Actualizando orden %d con estado: %s\n",
      order.ID, status)
  }
}
func reportOrderStatus(orders []*Order) {
  for i := 0; i < 5; i++ {
    time.Sleep(1 * time.Second)
    fmt.Printf("\n--- Reporte Estado de las Ordenes ---\n")
    for _, order := range orders {
      fmt.Printf("Orden %d %s\n",
        order.ID, order.Status)
    }
    fmt.Printf("---------------------------------------\n")
  }
}
func main() {
  orders := generateOrders(20)
  processOrders(orders)
  updateOrderStatuses(orders)
  reportOrderStatus(orders)
  fmt.Print("Todas las operaciones completadas. Finalizando\n")
}
func main() {

  orders := generateOrders(20)

    go processOrders(orders)

  go updateOrderStatuses(orders)

    go reportOrderStatus(orders)

  fmt.Print("Todas las operaciones completadas. Saliendo")

}