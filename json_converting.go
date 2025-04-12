package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadData(filename string) (*Resolver, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}
	var jsonData jsonData
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON: %w", err)
	}
	resolver := &Resolver{
		Sellers:    make(map[int]*Seller),
		MCatalog:   nil,
		UserCart:   make(map[string][]*CartItem),
		TotalItems: make(map[int]*Item),
	}
	for _, s := range jsonData.Sellers {
		resolver.Sellers[*s.ID] = &Seller{
			ID:    s.ID,
			Name:  s.Name,
			Deals: s.Deals,
		}
	}
	catalogMap := make(map[int]*Catalog)
	convertCatalog(resolver, jsonData.JCatalog, nil, catalogMap)
	resolver.MCatalog = catalogMap
	return resolver, nil
}

type JCatalog struct {
	ID     *int        `json:"id,omitempty"`
	Name   *string     `json:"name,omitempty"`
	Parent *JCatalog   `json:"parent"`
	Childs []*JCatalog `json:"childs"`
	Items  []*JItem    `json:"items,omitempty"`
}

type JItem struct {
	ID          *int            `json:"id,omitempty"`
	Name        *string         `json:"name,omitempty"`
	Parent      *JCatalog       `json:"parent,omitempty"`
	Seller      int             `json:"seller_id"`
	InCart      int             `json:"inCart"`
	InStockText json.RawMessage `json:"in_stock"`
}

func convertCatalog(resolver *Resolver, jc *JCatalog, parent *Catalog, catalogMap map[int]*Catalog) *Catalog {
	if jc == nil {
		return nil
	}
	c := &Catalog{
		ID:     jc.ID,
		Name:   jc.Name,
		Parent: parent,
	}
	catalogMap[*jc.ID] = c
	for _, child := range jc.Childs {
		c.Childs = append(c.Childs, convertCatalog(resolver, child, c, catalogMap))
	}
	for _, item := range jc.Items {
		var inStockStr string
		if item.InStockText != nil {
			inStockStr = string(item.InStockText)
		}
		seller := resolver.Sellers[item.Seller]
		if seller == nil {
			continue
		}
		newItem := &Item{
			ID:          item.ID,
			Name:        item.Name,
			Parent:      c,
			Seller:      seller,
			InCart:      0,
			InStockText: inStockStr,
		}
		c.Items = append(c.Items, newItem)
		resolver.TotalItems[*newItem.ID] = newItem
		seller.Items = append(seller.Items, newItem)
	}
	return c
}

type jsonData struct {
	JCatalog *JCatalog `json:"catalog"`
	Sellers  []*Seller `json:"sellers"`
}
