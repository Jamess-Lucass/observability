package model

import "github.com/google/uuid"

type BasketItem struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"productId"`
	Price     float64   `json:"price"`
	Quantity  uint      `json:"quantity"`
}
