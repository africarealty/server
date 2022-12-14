// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ErrorHook is an autogenerated mock type for the ErrorHook type
type ErrorHook struct {
	mock.Mock
}

// Error provides a mock function with given fields: err
func (_m *ErrorHook) Error(err error) {
	_m.Called(err)
}

type mockConstructorTestingTNewErrorHook interface {
	mock.TestingT
	Cleanup(func())
}

// NewErrorHook creates a new instance of ErrorHook. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewErrorHook(t mockConstructorTestingTNewErrorHook) *ErrorHook {
	mock := &ErrorHook{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
