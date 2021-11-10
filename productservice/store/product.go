package store

import "time"

type Product struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	VendorName  string    `json:"vendor"`
	Inventory   int       `json:"inventory"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"create_at"`
}
