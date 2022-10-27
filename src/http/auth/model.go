package auth

import (
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`    // Email - login
	Password string `json:"password"` // Password - password
}

// SessionToken specifies a session token
type SessionToken struct {
	SessionId             string    // SessionId - session ID
	AccessToken           string    // AccessToken
	AccessTokenExpiresAt  time.Time // AccessTokenExpiresAt - when access token expires
	RefreshToken          string    // RefreshToken
	RefreshTokenExpiresAt time.Time // RefreshToken - when refresh token expires
}

type LoginResponse struct {
	Token  *SessionToken `json:"token"`  // Token - auth token must be passed as  "Authorization Bearer" header for all the requests (except ones which don't require authorization)
	UserId string        `json:"userId"` // UserId - ID of account
}

type RegistrationRequest struct {
	Email     string `json:"email"`     // Email - user email
	Password  string `json:"password"`  // Password - password
	FirstName string `json:"firstName"` // FirstName - user first name
	LastName  string `json:"lastName"`  // LastName - user last name
	UserType  string `json:"userType"`  // UserType - user type
}

type OwnerProfile struct {
	Avatar string `json:"avatar,omitempty"` // Avatar avatar
}

type AgentProfile struct {
	Avatar string `json:"avatar,omitempty"` // Avatar avatar
}

type User struct {
	Id        string        `json:"id"`                  // Id - user ID
	Email     string        `json:"email"`               // Email - email
	FirstName string        `json:"firstName,omitempty"` // FirstName - user's first name
	LastName  string        `json:"lastName,omitempty"`  // LastName - user's last name
	Owner     *OwnerProfile `json:"owner,omitempty"`     // Owner - owner profile
	Agent     *AgentProfile `json:"agent,omitempty"`     // Agent - agent profile
}

type SetPasswordRequest struct {
	PrevPassword string `json:"prevPassword"` // PrevPassword - current password
	NewPassword  string `json:"newPassword"`  // NewPassword - new password
}
