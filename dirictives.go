package main

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func AuthorizedDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user, err := GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("User not authorized")
	}
	fmt.Println("User authorized:", user.Username)
	return next(ctx)
}

func GetUserFromContext(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(sessionKey).(*User)
	if !ok {
		return nil, fmt.Errorf("User not authorized")
	}
	return user, nil
}
