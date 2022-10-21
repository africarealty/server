package domain

import (
	"context"
	"github.com/africarealty/server/src/kit/auth"
)

const (
	UserTypeAdmin  = "admin"
	UserTypeClient = "client"

	AuthGroupSysAdmin = "sysadmin"
	AuthGroupClient   = "client"

	AuthRoleSysAdmin        = "sysadmin"
	AuthRoleArbitrageClient = "arbitrage.client"

	AuthResUserProfileAll     = "users.all"
	AuthResUserProfileMy      = "users.my"
	AuthResArbitrageChainsAll = "arbitrage.chains.all"
)

type UserService interface {
	// Create creates a new user
	Create(ctx context.Context, user *auth.User) (*auth.User, error)
	// GetByEmail gets user by email
	GetByEmail(ctx context.Context, email string) (*auth.User, error)
	// Get gets user by id
	Get(ctx context.Context, userId string) (*auth.User, error)
	// GetByIds retrieves users by IDs
	GetByIds(ctx context.Context, userIds []string) ([]*auth.User, error)
	// SetPassword updates user password
	SetPassword(ctx context.Context, userId, prevPassword, newPassword string) error
}

type UserStorage interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, u *auth.User) error
	// UpdateUser updates an user
	UpdateUser(ctx context.Context, u *auth.User) error
	// GetByEmail retrieves an user by email
	GetByUsername(ctx context.Context, email string) (*auth.User, error)
	// GetUser retrieves a user by id
	GetUser(ctx context.Context, userId string) (*auth.User, error)
	// GetByIds retrieves users by IDs
	GetUserByIds(ctx context.Context, userIds []string) ([]*auth.User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, u *auth.User) error
}
