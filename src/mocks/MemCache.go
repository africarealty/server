// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// MemCache is an autogenerated mock type for the MemCache type
type MemCache struct {
	mock.Mock
}

// Delete provides a mock function with given fields: key
func (_m *MemCache) Delete(key string) {
	_m.Called(key)
}

// Get provides a mock function with given fields: key
func (_m *MemCache) Get(key string) (interface{}, bool) {
	ret := _m.Called(key)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Set provides a mock function with given fields: key, v, ttl
func (_m *MemCache) Set(key string, v interface{}, ttl time.Duration) {
	_m.Called(key, v, ttl)
}