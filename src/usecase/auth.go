package usecase

import (
	"context"
	"github.com/africarealty/server/src/domain"
)

type UserRegistrationRq struct {
	Email                string // Email user email
	Password             string // Password password
	PasswordConfirmation string // PasswordConfirmation password confirmation
	UserType             string // UserType user type
	FirstName            string // FirstName first name
	LastName             string // LastName last name
}

type UserUseCases interface {
	// Register registers a new user
	Register(ctx context.Context, rq *UserRegistrationRq) (*domain.User, error)
	// CreateActiveUser creates a new active user
	CreateActiveUser(ctx context.Context, rq *UserRegistrationRq) (*domain.User, error)
}
