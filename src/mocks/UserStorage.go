// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/africarealty/server/src/kit/auth"

	domain "github.com/africarealty/server/src/domain"

	mock "github.com/stretchr/testify/mock"
)

// UserStorage is an autogenerated mock type for the UserStorage type
type UserStorage struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, u
func (_m *UserStorage) CreateUser(ctx context.Context, u *domain.User) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: ctx, u
func (_m *UserStorage) DeleteUser(ctx context.Context, u *domain.User) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetActivationToken provides a mock function with given fields: ctx, userId
func (_m *UserStorage) GetActivationToken(ctx context.Context, userId string) (string, error) {
	ret := _m.Called(ctx, userId)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, userId)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *UserStorage) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUsername provides a mock function with given fields: ctx, username
func (_m *UserStorage) GetByUsername(ctx context.Context, username string) (*auth.User, error) {
	ret := _m.Called(ctx, username)

	var r0 *auth.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *auth.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, userId
func (_m *UserStorage) GetUser(ctx context.Context, userId string) (*domain.User, error) {
	ret := _m.Called(ctx, userId)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.User); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
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

// GetUserByIds provides a mock function with given fields: ctx, userIds
func (_m *UserStorage) GetUserByIds(ctx context.Context, userIds []string) ([]*domain.User, error) {
	ret := _m.Called(ctx, userIds)

	var r0 []*domain.User
	if rf, ok := ret.Get(0).(func(context.Context, []string) []*domain.User); ok {
		r0 = rf(ctx, userIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, userIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetActivationToken provides a mock function with given fields: ctx, userId, token, ttl
func (_m *UserStorage) SetActivationToken(ctx context.Context, userId string, token string, ttl uint32) error {
	ret := _m.Called(ctx, userId, token, ttl)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, uint32) error); ok {
		r0 = rf(ctx, userId, token, ttl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUser provides a mock function with given fields: ctx, u
func (_m *UserStorage) UpdateUser(ctx context.Context, u *domain.User) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserStorage creates a new instance of UserStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserStorage(t mockConstructorTestingTNewUserStorage) *UserStorage {
	mock := &UserStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
