// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/africarealty/server/src/domain"
	mock "github.com/stretchr/testify/mock"
)

// Adapter is an autogenerated mock type for the Adapter type
type Adapter struct {
	mock.Mock
}

// Close provides a mock function with given fields: ctx
func (_m *Adapter) Close(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEmailStorage provides a mock function with given fields:
func (_m *Adapter) GetEmailStorage() domain.EmailStorage {
	ret := _m.Called()

	var r0 domain.EmailStorage
	if rf, ok := ret.Get(0).(func() domain.EmailStorage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.EmailStorage)
		}
	}

	return r0
}

// GetTemplateStorage provides a mock function with given fields:
func (_m *Adapter) GetTemplateStorage() domain.TemplateStorage {
	ret := _m.Called()

	var r0 domain.TemplateStorage
	if rf, ok := ret.Get(0).(func() domain.TemplateStorage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.TemplateStorage)
		}
	}

	return r0
}

// Init provides a mock function with given fields: ctx, cfg
func (_m *Adapter) Init(ctx context.Context, cfg interface{}) error {
	ret := _m.Called(ctx, cfg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, cfg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAdapter interface {
	mock.TestingT
	Cleanup(func())
}

// NewAdapter creates a new instance of Adapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAdapter(t mockConstructorTestingTNewAdapter) *Adapter {
	mock := &Adapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
