// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	goroutine "github.com/africarealty/server/src/kit/goroutine"
	log "github.com/africarealty/server/src/kit/log"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Goroutine is an autogenerated mock type for the Goroutine type
type Goroutine struct {
	mock.Mock
}

// Cmp provides a mock function with given fields: component
func (_m *Goroutine) Cmp(component string) goroutine.Goroutine {
	ret := _m.Called(component)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(string) goroutine.Goroutine); ok {
		r0 = rf(component)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// Go provides a mock function with given fields: ctx, f
func (_m *Goroutine) Go(ctx context.Context, f func()) {
	_m.Called(ctx, f)
}

// Mth provides a mock function with given fields: method
func (_m *Goroutine) Mth(method string) goroutine.Goroutine {
	ret := _m.Called(method)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(string) goroutine.Goroutine); ok {
		r0 = rf(method)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// WithLogger provides a mock function with given fields: logger
func (_m *Goroutine) WithLogger(logger log.CLogger) goroutine.Goroutine {
	ret := _m.Called(logger)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(log.CLogger) goroutine.Goroutine); ok {
		r0 = rf(logger)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// WithLoggerFn provides a mock function with given fields: loggerFn
func (_m *Goroutine) WithLoggerFn(loggerFn log.CLoggerFunc) goroutine.Goroutine {
	ret := _m.Called(loggerFn)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(log.CLoggerFunc) goroutine.Goroutine); ok {
		r0 = rf(loggerFn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// WithRetry provides a mock function with given fields: retry
func (_m *Goroutine) WithRetry(retry int) goroutine.Goroutine {
	ret := _m.Called(retry)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(int) goroutine.Goroutine); ok {
		r0 = rf(retry)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

// WithRetryDelay provides a mock function with given fields: delay
func (_m *Goroutine) WithRetryDelay(delay time.Duration) goroutine.Goroutine {
	ret := _m.Called(delay)

	var r0 goroutine.Goroutine
	if rf, ok := ret.Get(0).(func(time.Duration) goroutine.Goroutine); ok {
		r0 = rf(delay)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(goroutine.Goroutine)
		}
	}

	return r0
}

type mockConstructorTestingTNewGoroutine interface {
	mock.TestingT
	Cleanup(func())
}

// NewGoroutine creates a new instance of Goroutine. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGoroutine(t mockConstructorTestingTNewGoroutine) *Goroutine {
	mock := &Goroutine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
