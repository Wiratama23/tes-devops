package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ProductID       string          `json:"product_id"`
	ProductName     string          `json:"product_name"`
	ProductQuantity int             `json:"product_quantity"`
	ProductPrices   decimal.Decimal `json:"product_prices"`
	ProductType     string          `json:"product_type"`
	CreatedAt       time.Time       `json:"created_at"`
	CreatedBy       uuid.UUID       `json:"created_by"`
	ImagePath       string          `json:"image_path"`
}
