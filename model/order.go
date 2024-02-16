package model

import (
	"time"

	"github.com/google/uuid"
)

// `CustomerID` gets uuidv4 as type
// we wont be setting this field
// because it is external to system

// storing as a string requires us to validate again later
// use a custom type to represent the uuidv4 for CustomerID
// uuid package

type Order struct {
	OrderID     uint64     `json:"order_id"`
	CustomerID  uuid.UUID  `json:"cust_id"`
	LineItems   []LineItem `json:"line_items"`
	CreatedAt   *time.Time `json:"created_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type LineItem struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity uint      `json:"quantity"`
	Price    uint      `json:"price"`
}
