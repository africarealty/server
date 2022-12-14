// Code generated by mockery 2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/africarealty/server/src/domain"
	mock "github.com/stretchr/testify/mock"
)

// TemplateGenerator is an autogenerated mock type for the TemplateGenerator type
type TemplateGenerator struct {
	mock.Mock
}

// Generate provides a mock function with given fields: ctx, rq
func (_m *TemplateGenerator) Generate(ctx context.Context, rq *domain.TemplateRequest) (*domain.TemplateResponse, error) {
	ret := _m.Called(ctx, rq)

	var r0 *domain.TemplateResponse
	if rf, ok := ret.Get(0).(func(context.Context, *domain.TemplateRequest) *domain.TemplateResponse); ok {
		r0 = rf(ctx, rq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.TemplateResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.TemplateRequest) error); ok {
		r1 = rf(ctx, rq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTemplateGenerator interface {
	mock.TestingT
	Cleanup(func())
}

// NewTemplateGenerator creates a new instance of TemplateGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTemplateGenerator(t mockConstructorTestingTNewTemplateGenerator) *TemplateGenerator {
	mock := &TemplateGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
