package schema

import (
	"time"
)

var (
	SchemaVersionRawData         = "2024-02-20"
	SchemaVersionTransformedData = "2024-02-20"
)

type RawData struct {
	Time          time.Time         `json:"time"`
	SchemaVersion string            `json:"schema_version"`
	Onwer         string            `json:"owner"`
	Project       string            `json:"project"`
	Service       string            `json:"service"`
	Tags          map[string]string `json:"tags"`
	Data          []byte            `json:"data"`
}

type TransformedData struct {
	Time              time.Time `json:"time"`
	SchemaVersion     string    `json:"schema_version"`
	Metadata          string    `json:"metadata"`
	Owner             string    `json:"owner"`
	Project           string    `json:"project"`
	Provider          string    `json:"provider"`
	Service           string    `json:"service"`
	SKU               string    `json:"sku"`
	CostAmount        float64   `json:"cost_amount"`
	CostAmountUnit    string    `json:"cost_amount_unit"`
	UsageQuantity     float64   `json:"usage_quantity"`
	UsageQuantityUnit string    `json:"usage_quantity_unit"`
	ExchangeRate      float64   `json:"exchange_rate"`
	Total             float64   `json:"total"`
}
