package main

import (
	"./sssaas"
	"fmt"
)

func main() {
	shares := sssaas.Create(1, 2, "asdffdsaasdffdsaasdffdsaasdffdsaasdffdsaasdffdsaasdffdsaasdffdsa")
	fmt.Println(sssaas.Combine(shares))
}
