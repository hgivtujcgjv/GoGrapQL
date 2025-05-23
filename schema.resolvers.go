package main

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"context"
	"errors"
	"strconv"
)

// Parent is the resolver for the parent field.
func (r *catalogResolver) Parent(ctx context.Context, obj *Catalog) (*Catalog, error) {
	res := obj.Parent
	if (res) == nil {
		return nil, errors.New("Empty parent catalog")
	}
	return res, nil
}

// Childs is the resolver for the childs field.
func (r *catalogResolver) Childs(ctx context.Context, obj *Catalog) ([]*Catalog, error) {
	res := obj.Childs
	if (len(res)) == 1 {
		return nil, errors.New("Empty catalog childs")
	}
	return res, nil
}

// Items is the resolver for the items field.
func (r *catalogResolver) Items(ctx context.Context, obj *Catalog, limit *int, offset *int) ([]*Item, error) {
	if *offset >= len(obj.Items) {
		return []*Item{}, nil
	}
	end := *offset + *limit
	if end > len(obj.Items) {
		end = len(obj.Items)
	}
	res := obj.Items[*offset:end]
	return res, nil
}

// Parent is the resolver for the parent field.
func (r *itemResolver) Parent(ctx context.Context, obj *Item) (*Catalog, error) {
	if obj.Parent != nil {
		return obj.Parent, nil
	}
	return nil, nil
}

// Seller is the resolver for the seller field.
func (r *itemResolver) Seller(ctx context.Context, obj *Item) (*Seller, error) {
	if obj != nil {
		return obj.Seller, nil
	}
	return nil, errors.New("Emty object")
}

// InCart is the resolver for the inCart field.
func (r *itemResolver) InCart(ctx context.Context, obj *Item) (int, error) {
	user, err := ctx.Value(sessionKey).(*User)
	if !err {
		return 0, errors.New("Unknown user")
	}
	temp := r.UserCart[user.Username]
	for i := 0; i < len(temp); i++ {
		if temp[i].Item.ID == obj.ID {
			return temp[i].Quantity, nil
		}
	}
	return 0, nil
}

// InStockText is the resolver for the inStockText field.
func (r *itemResolver) InStockText(ctx context.Context, obj *Item) (string, error) {
	if obj == nil {
		return "", errors.New("Empty object")
	}
	switch obj.InStockText {
	case "1":
		return "мало", nil
	case "2", "3":
		return "хватает", nil
	case "4", "5":
		return "много", nil
	default:
		return obj.InStockText, nil // откат на исходное значение
	}
}

// AddToCart is the resolver for the AddToCart field.
func (r *mutationResolver) AddToCart(ctx context.Context, in *CartInput) ([]*CartItem, error) {
	user, ok := ctx.Value(sessionKey).(*User)
	if !ok || user == nil {
		return nil, errors.New("Internal server error")
	}
	cartItems := r.Resolver.UserCart[user.Username]
	if cartItems == nil {
		cartItems = []*CartItem{}
	}
	itm, exists := r.Resolver.TotalItems[in.ItemID]
	if !exists {
		return nil, errors.New("Item not found")
	}
	stock, err := strconv.Atoi(itm.InStockText)
	if err != nil {
		return nil, errors.New("Invalid stock data")
	}
	if stock < in.Quantity {
		return nil, errors.New("Not enough quantity")
	}
	found := false
	for _, ci := range cartItems {
		if *ci.Item.ID == in.ItemID {
			ci.Quantity += in.Quantity
			found = true
			break
		}
	}
	if !found {
		newItem := &CartItem{
			Quantity: in.Quantity,
			Item:     itm,
		}
		cartItems = append(cartItems, newItem)
	}
	itm.InStockText = strconv.Itoa(stock - in.Quantity)
	itm.InCart += in.Quantity
	r.Resolver.UserCart[user.Username] = cartItems
	return cartItems, nil
}

// RemoveFromCart is the resolver for the RemoveFromCart field.
func (r *mutationResolver) RemoveFromCart(ctx context.Context, in CartInput) ([]*CartItem, error) {
	user, ok := ctx.Value(sessionKey).(*User)
	if !ok || user == nil {
		return nil, errors.New("Internal server error")
	}
	temp := r.UserCart[user.Username]
	for i := 0; i < len(temp); i++ {
		if *temp[i].Item.ID == in.ItemID {
			temp[i].Item.InCart -= in.Quantity
			if temp[i].Item.InCart <= 0 {
				temp = append(temp[:i], temp[i+1:]...)
				i--
			}
		}
	}
	r.UserCart[user.Username] = temp
	return temp, nil
}

// Catalog is the resolver for the Catalog field.
func (r *queryResolver) Catalog(ctx context.Context, id *string) (*Catalog, error) {
	Temp, err := strconv.Atoi(*id)
	if err != nil {
		return nil, err
	}
	Res := r.Resolver.MCatalog[int(Temp)]
	return Res, nil
}

// Shop is the resolver for the Shop field.
func (r *queryResolver) Shop(ctx context.Context, parentID *string) ([]*Catalog, error) {
	Temp, err := strconv.Atoi(*parentID)
	if err != nil {
		return nil, err
	}
	Res := r.Resolver.MCatalog[int(Temp)]
	return Res.Childs, nil
}

// Seller is the resolver for the Seller field.
func (r *queryResolver) Seller(ctx context.Context, id *string) (*Seller, error) {
	Temp, err := strconv.Atoi(*id)
	if err != nil {
		return nil, err
	}
	Res := r.Resolver.Sellers[int(Temp)]
	return Res, nil
}

// MyCart is the resolver for the MyCart field.
func (r *queryResolver) MyCart(ctx context.Context) ([]*CartItem, error) {
	user, err := ctx.Value(sessionKey).(*User)
	if !err {
		return nil, errors.New("Don't know this user")
	}
	temp := r.UserCart[user.Username]
	return temp, nil
}

// Items is the resolver for the items field.
func (r *sellerResolver) Items(ctx context.Context, obj *Seller, limit *int, offset *int) ([]*Item, error) {
	if *offset >= len(obj.Items) {
		return []*Item{}, nil
	}
	end := *offset + *limit
	if end > len(obj.Items) {
		end = len(obj.Items)
	}
	res := obj.Items[*offset:end]
	return res, nil
}

// Catalog returns CatalogResolver implementation.
func (r *Resolver) Catalog() CatalogResolver { return &catalogResolver{r} }

// Item returns ItemResolver implementation.
func (r *Resolver) Item() ItemResolver { return &itemResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Seller returns SellerResolver implementation.
func (r *Resolver) Seller() SellerResolver { return &sellerResolver{r} }

type catalogResolver struct{ *Resolver }
type itemResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type sellerResolver struct{ *Resolver }
