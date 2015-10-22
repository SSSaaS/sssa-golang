package main

import (
	"./sssaas"
	"fmt"
)

func main() {
	shares := sssaas.Create(3, 4, "asdf")
	fmt.Println("Shares: ", shares)
	fmt.Println(sssaas.Combine(shares))
}
