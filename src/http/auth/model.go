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

type ClientRegistrationRequest struct {
	Email     string `json:"email"`     // Email - user's email
	Password  string `json:"password"`  // Password - password
	FirstName string `json:"firstName"` // FirstName - user's first name
	LastName  string `json:"lastName"`  // LastName - user's last name
}

type ClientUser struct {
	Id        string `json:"id"`                  // Id - user ID
	Email     string `json:"email"`               // Email - email
	FirstName string `json:"firstName,omitempty"` // FirstName - user's first name
	LastName  string `json:"lastName,omitempty"`  // LastName - user's last name
}

type SetPasswordRequest struct {
	PrevPassword string `json:"prevPassword"` // PrevPassword - current password
	NewPassword  string `json:"newPassword"`  // NewPassword - new password
}
