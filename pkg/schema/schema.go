package schema

import (
	"time"
)

type RawData struct {
	Time    time.Time         `json:"time"`
	Onwer   string            `json:"owner"`
	Project string            `json:"project"`
	Service string            `json:"service"`
	Tags    map[string]string `json:"tags"`
	Data    []byte            `json:"data"`
}

type TransformedData struct {
	Time     time.Time `json:"time"`
	Onwer    string    `json:"owner"`
	Project  string    `json:"project"`
	Service  string    `json:"service"`
	SKU      string    `json:"sku"`
	Price    float64   `json:"price"`
	Quantity int       `json:"quantity"`
	Total    float64   `json:"total"`
}
