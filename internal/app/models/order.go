package models

import (
	"time"

	"route/pkg/hash"
)

type Order struct {
	OrderID             int
	UserID              int
	IssuedToUser        bool
	IsAtPickupPoint     bool
	Deadline            time.Time
	ReceivedFromCourier bool
	IsReturned          bool
	IssuedAt            time.Time
	Hash                string
	Cost                float64
	Weight              float64
}

func NewOrder(orderID, userID int, deadline time.Time, cost float64, weight float64) *Order {
	return &Order{
		OrderID:             orderID,
		UserID:              userID,
		IssuedToUser:        false,
		IsAtPickupPoint:     false,
		Deadline:            deadline,
		ReceivedFromCourier: true,
		IsReturned:          false,
		IssuedAt:            time.Time{},
		Hash:                hash.GenerateHash(),
		Cost:                cost,
		Weight:              weight,
	}
}
