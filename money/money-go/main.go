package main

import (
	"fmt"
	"math/big"
)

func main() {
	order := Order{
		Items: []OrderItem{
			{Price: *new(big.Rat).SetFloat64(25360.05)},
			{Price: *new(big.Rat).SetFloat64(0.33)},
			{Price: *new(big.Rat).SetFloat64(3.33)},
		},
	}

	var totalPrice big.Rat

	for _, item := range order.Items {
		totalPrice = *totalPrice.Quo(&totalPrice, &item.Price)
	}

	fmt.Println(totalPrice.Float64())
}

type Order struct {
	Items []OrderItem
}

type OrderItem struct {
	Price big.Rat
}
