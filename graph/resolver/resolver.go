package resolver

import "github.com/flyfy1/diarier/orm"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen@v0.17.31 generate

type ResolverOption func(*Resolver)

func ResolverOptionWithUserOrm(userOrm *orm.UserORM) ResolverOption {
	return func(r *Resolver) {
		r.userORM = userOrm
	}
}
func ResolverOptionWithTaskOrm(taskOrm *orm.TaskORM) ResolverOption {
	return func(r *Resolver) {
		r.taskOrm = taskOrm
	}
}

func NewResolver(options... ResolverOption) *Resolver {
	r := &Resolver{}
	for _, option := range options {
		option(r)
	}
	return r
}

type Resolver struct {
	userORM *orm.UserORM
	taskOrm *orm.TaskORM
}
