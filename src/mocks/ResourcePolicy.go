// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/africarealty/server/src/kit/auth"

	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// ResourcePolicy is an autogenerated mock type for the ResourcePolicy type
type ResourcePolicy struct {
	mock.Mock
}

// Resolve provides a mock function with given fields: ctx, r
func (_m *ResourcePolicy) Resolve(ctx context.Context, r *http.Request) (*auth.AuthorizationResource, error) {
	ret := _m.Called(ctx, r)

	var r0 *auth.AuthorizationResource
	if rf, ok := ret.Get(0).(func(context.Context, *http.Request) *auth.AuthorizationResource); ok {
		r0 = rf(ctx, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.AuthorizationResource)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *http.Request) error); ok {
		r1 = rf(ctx, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewResourcePolicy interface {
	mock.TestingT
	Cleanup(func())
}

// NewResourcePolicy creates a new instance of ResourcePolicy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewResourcePolicy(t mockConstructorTestingTNewResourcePolicy) *ResourcePolicy {
	mock := &ResourcePolicy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
