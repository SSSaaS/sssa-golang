package main

import (
	"./sssaas"
	"fmt"
)

func main() {
	shares := sssaas.Create(1, 2, "asdf")
	fmt.Println("Shares: ", shares)
	fmt.Println(sssaas.Combine(shares))
}
