// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/africarealty/server/src/kit/auth"

	mock "github.com/stretchr/testify/mock"
)

// AuthorizeSession is an autogenerated mock type for the AuthorizeSession type
type AuthorizeSession struct {
	mock.Mock
}

// AuthorizeSession provides a mock function with given fields: ctx, rq
func (_m *AuthorizeSession) AuthorizeSession(ctx context.Context, rq *auth.AuthorizationRequest) (bool, error) {
	ret := _m.Called(ctx, rq)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *auth.AuthorizationRequest) bool); ok {
		r0 = rf(ctx, rq)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *auth.AuthorizationRequest) error); ok {
		r1 = rf(ctx, rq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRolesForGroups provides a mock function with given fields: ctx, groups
func (_m *AuthorizeSession) GetRolesForGroups(ctx context.Context, groups []string) ([]string, error) {
	ret := _m.Called(ctx, groups)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, []string) []string); ok {
		r0 = rf(ctx, groups)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, groups)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAuthorizeSession interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthorizeSession creates a new instance of AuthorizeSession. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthorizeSession(t mockConstructorTestingTNewAuthorizeSession) *AuthorizeSession {
	mock := &AuthorizeSession{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
