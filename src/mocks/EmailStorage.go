// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/africarealty/server/src/domain"
	mock "github.com/stretchr/testify/mock"
)

// EmailStorage is an autogenerated mock type for the EmailStorage type
type EmailStorage struct {
	mock.Mock
}

// CreateEmail provides a mock function with given fields: ctx, requests
func (_m *EmailStorage) CreateEmail(ctx context.Context, requests *domain.Email) error {
	ret := _m.Called(ctx, requests)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Email) error); ok {
		r0 = rf(ctx, requests)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateEmail provides a mock function with given fields: ctx, requests
func (_m *EmailStorage) UpdateEmail(ctx context.Context, requests ...*domain.Email) error {
	_va := make([]interface{}, len(requests))
	for _i := range requests {
		_va[_i] = requests[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...*domain.Email) error); ok {
		r0 = rf(ctx, requests...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
