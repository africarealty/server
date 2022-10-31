// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/africarealty/server/src/domain"
	mock "github.com/stretchr/testify/mock"

	usecase "github.com/africarealty/server/src/usecase"
)

// UserUseCases is an autogenerated mock type for the UserUseCases type
type UserUseCases struct {
	mock.Mock
}

// CreateActiveUser provides a mock function with given fields: ctx, rq
func (_m *UserUseCases) CreateActiveUser(ctx context.Context, rq *usecase.UserRegistrationRq) (*domain.User, error) {
	ret := _m.Called(ctx, rq)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(context.Context, *usecase.UserRegistrationRq) *domain.User); ok {
		r0 = rf(ctx, rq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *usecase.UserRegistrationRq) error); ok {
		r1 = rf(ctx, rq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, rq
func (_m *UserUseCases) Register(ctx context.Context, rq *usecase.UserRegistrationRq) (*domain.User, error) {
	ret := _m.Called(ctx, rq)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(context.Context, *usecase.UserRegistrationRq) *domain.User); ok {
		r0 = rf(ctx, rq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *usecase.UserRegistrationRq) error); ok {
		r1 = rf(ctx, rq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserUseCases interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserUseCases creates a new instance of UserUseCases. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserUseCases(t mockConstructorTestingTNewUserUseCases) *UserUseCases {
	mock := &UserUseCases{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
