// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	config "github.com/africarealty/server/src/kit/config"
	mock "github.com/stretchr/testify/mock"
)

// Loader is an autogenerated mock type for the Loader type
type Loader struct {
	mock.Mock
}

// Load provides a mock function with given fields: target
func (_m *Loader) Load(target interface{}) error {
	ret := _m.Called(target)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(target)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithConfigPath provides a mock function with given fields: configPath
func (_m *Loader) WithConfigPath(configPath string) config.Loader {
	ret := _m.Called(configPath)

	var r0 config.Loader
	if rf, ok := ret.Get(0).(func(string) config.Loader); ok {
		r0 = rf(configPath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(config.Loader)
		}
	}

	return r0
}

// WithConfigPathFromEnv provides a mock function with given fields: env
func (_m *Loader) WithConfigPathFromEnv(env string) config.Loader {
	ret := _m.Called(env)

	var r0 config.Loader
	if rf, ok := ret.Get(0).(func(string) config.Loader); ok {
		r0 = rf(env)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(config.Loader)
		}
	}

	return r0
}

// WithEnvPath provides a mock function with given fields: envPath
func (_m *Loader) WithEnvPath(envPath string) config.Loader {
	ret := _m.Called(envPath)

	var r0 config.Loader
	if rf, ok := ret.Get(0).(func(string) config.Loader); ok {
		r0 = rf(envPath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(config.Loader)
		}
	}

	return r0
}
