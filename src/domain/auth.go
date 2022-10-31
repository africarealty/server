package domain

import (
	"context"
	"github.com/africarealty/server/src/kit/auth"
)

const (
	UserTypeAdmin = "admin"
	UserTypeOwner = "owner"
	UserTypeAgent = "agent"

	AuthGroupSysAdmin = "sysadmin"
	AuthGroupOwner    = "owner"
	AuthGroupAgent    = "agent"

	AuthRoleSysAdmin       = "sysadmin"
	AuthRoleProfileOwner   = "profile.owner"
	AuthRoleProfileManager = "profile.manager"

	AuthResUserProfileAll = "profile.all"
	AuthResUserProfileMy  = "profile.my"
)

// Profile common profile attrs
type Profile struct {
	Avatar string // Avatar avatar link
}

// OwnerProfile owner profile
type OwnerProfile struct {
	Profile
}

// AgentProfile agent profile
type AgentProfile struct {
	Profile
}

type User struct {
	auth.User               // basic profile attributes
	Owner     *OwnerProfile // Owner profile
	Agent     *AgentProfile // Agent profile
}

type UserService interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) (*User, error)
	// GetByEmail gets user by email
	GetByEmail(ctx context.Context, email string) (*User, error)
	// Get gets user by id
	Get(ctx context.Context, userId string) (*User, error)
	// GetByIds retrieves users by IDs
	GetByIds(ctx context.Context, userIds []string) ([]*User, error)
	// SetPassword updates user password
	SetPassword(ctx context.Context, userId, newPasswordHash string) error
	// SetActivationToken sets token for the given user with the given ttl
	SetActivationToken(ctx context.Context, userId, token string, ttl uint32) error
	// ActivateByToken activates a user by token
	ActivateByToken(ctx context.Context, userId, token string) (*User, error)
}

type UserStorage interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, u *User) error
	// UpdateUser updates an user
	UpdateUser(ctx context.Context, u *User) error
	// GetByEmail retrieves an user by email
	GetByEmail(ctx context.Context, email string) (*User, error)
	// GetUser retrieves a user by id
	GetUser(ctx context.Context, userId string) (*User, error)
	// GetByIds retrieves users by IDs
	GetUserByIds(ctx context.Context, userIds []string) ([]*User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, u *User) error
	// SetActivationToken sets token for the given user with the given ttl
	SetActivationToken(ctx context.Context, userId, token string, ttl uint32) error
	// GetActivationToken retrieves activation token
	GetActivationToken(ctx context.Context, userId string) (string, error)
	// GetByUsername retrieves an user by username
	GetByUsername(ctx context.Context, username string) (*auth.User, error)
}
