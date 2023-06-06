package auth

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}

	return next(ctx)
}
