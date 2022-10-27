package usecase

import (
	"context"
	"github.com/africarealty/server/src/domain"
)

type UserRegistrationRq struct {
	Email     string // Email user email
	Password  string // Password password
	UserType  string // UserType user type
	FirstName string // FirstName first name
	LastName  string // LastName last name
}

type UserRegistrationUseCase interface {
	// Register registers a new user
	Register(ctx context.Context, rq *UserRegistrationRq) (*domain.User, error)
}
