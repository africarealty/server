// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/africarealty/server/src/domain"
	mock "github.com/stretchr/testify/mock"
)

// EmailRepository is an autogenerated mock type for the EmailRepository type
type EmailRepository struct {
	mock.Mock
}

// Send provides a mock function with given fields: ctx, email
func (_m *EmailRepository) Send(ctx context.Context, email *domain.Email) error {
	ret := _m.Called(ctx, email)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Email) error); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewEmailRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewEmailRepository creates a new instance of EmailRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEmailRepository(t mockConstructorTestingTNewEmailRepository) *EmailRepository {
	mock := &EmailRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
