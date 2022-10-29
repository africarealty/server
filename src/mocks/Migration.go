// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Migration is an autogenerated mock type for the Migration type
type Migration struct {
	mock.Mock
}

// Up provides a mock function with given fields:
func (_m *Migration) Up() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMigration interface {
	mock.TestingT
	Cleanup(func())
}

// NewMigration creates a new instance of Migration. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMigration(t mockConstructorTestingTNewMigration) *Migration {
	mock := &Migration{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
