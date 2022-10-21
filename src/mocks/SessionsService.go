// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/africarealty/server/src/kit/auth"

	mock "github.com/stretchr/testify/mock"
)

// SessionsService is an autogenerated mock type for the SessionsService type
type SessionsService struct {
	mock.Mock
}

// AuthSession provides a mock function with given fields: ctx, token
func (_m *SessionsService) AuthSession(ctx context.Context, token string) (*auth.Session, error) {
	ret := _m.Called(ctx, token)

	var r0 *auth.Session
	if rf, ok := ret.Get(0).(func(context.Context, string) *auth.Session); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, sid
func (_m *SessionsService) Get(ctx context.Context, sid string) (*auth.Session, error) {
	ret := _m.Called(ctx, sid)

	var r0 *auth.Session
	if rf, ok := ret.Get(0).(func(context.Context, string) *auth.Session); ok {
		r0 = rf(ctx, sid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, sid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUser provides a mock function with given fields: ctx, userId
func (_m *SessionsService) GetByUser(ctx context.Context, userId string) ([]*auth.Session, error) {
	ret := _m.Called(ctx, userId)

	var r0 []*auth.Session
	if rf, ok := ret.Get(0).(func(context.Context, string) []*auth.Session); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*auth.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields: cfg
func (_m *SessionsService) Init(cfg *auth.Config) {
	_m.Called(cfg)
}

// LoginPassword provides a mock function with given fields: ctx, rq
func (_m *SessionsService) LoginPassword(ctx context.Context, rq *auth.LoginRequest) (*auth.Session, *auth.SessionToken, error) {
	ret := _m.Called(ctx, rq)

	var r0 *auth.Session
	if rf, ok := ret.Get(0).(func(context.Context, *auth.LoginRequest) *auth.Session); ok {
		r0 = rf(ctx, rq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Session)
		}
	}

	var r1 *auth.SessionToken
	if rf, ok := ret.Get(1).(func(context.Context, *auth.LoginRequest) *auth.SessionToken); ok {
		r1 = rf(ctx, rq)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*auth.SessionToken)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *auth.LoginRequest) error); ok {
		r2 = rf(ctx, rq)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Logout provides a mock function with given fields: ctx, sid
func (_m *SessionsService) Logout(ctx context.Context, sid string) error {
	ret := _m.Called(ctx, sid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, sid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RefreshToken provides a mock function with given fields: ctx, refreshToken
func (_m *SessionsService) RefreshToken(ctx context.Context, refreshToken string) (*auth.SessionToken, error) {
	ret := _m.Called(ctx, refreshToken)

	var r0 *auth.SessionToken
	if rf, ok := ret.Get(0).(func(context.Context, string) *auth.SessionToken); ok {
		r0 = rf(ctx, refreshToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.SessionToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, refreshToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
