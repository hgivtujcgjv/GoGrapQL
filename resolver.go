package main

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	Sellers    map[int]*Seller        `json:"sellers"`
	MCatalog   map[int]*Catalog       `json:"catalog"`
	UserCart   map[string][]*CartItem `json:"cartitem,omitempty"`
	TotalItems map[int]*Item
	// для мутаций пока резолверы не написанны
}
