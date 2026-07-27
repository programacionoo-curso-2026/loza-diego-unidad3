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