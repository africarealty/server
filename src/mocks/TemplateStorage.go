// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/africarealty/server/src/domain"
	mock "github.com/stretchr/testify/mock"
)

// TemplateStorage is an autogenerated mock type for the TemplateStorage type
type TemplateStorage struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, template
func (_m *TemplateStorage) Create(ctx context.Context, template *domain.Template) error {
	ret := _m.Called(ctx, template)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Template) error); ok {
		r0 = rf(ctx, template)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *TemplateStorage) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id
func (_m *TemplateStorage) Get(ctx context.Context, id string) (*domain.Template, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.Template
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Template); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Search provides a mock function with given fields: ctx, query
func (_m *TemplateStorage) Search(ctx context.Context, query string) ([]*domain.Template, error) {
	ret := _m.Called(ctx, query)

	var r0 []*domain.Template
	if rf, ok := ret.Get(0).(func(context.Context, string) []*domain.Template); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Template)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, template
func (_m *TemplateStorage) Update(ctx context.Context, template *domain.Template) error {
	ret := _m.Called(ctx, template)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Template) error); ok {
		r0 = rf(ctx, template)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}