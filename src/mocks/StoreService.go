// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	domain "github.com/africarealty/server/src/domain"

	mock "github.com/stretchr/testify/mock"
)

// StoreService is an autogenerated mock type for the StoreService type
type StoreService struct {
	mock.Mock
}

// BuildFileID provides a mock function with given fields: bucketName, filename
func (_m *StoreService) BuildFileID(bucketName string, filename string) string {
	ret := _m.Called(bucketName, filename)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(bucketName, filename)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// DeleteFile provides a mock function with given fields: ctx, fileID
func (_m *StoreService) DeleteFile(ctx context.Context, fileID string) error {
	ret := _m.Called(ctx, fileID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, fileID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFile provides a mock function with given fields: ctx, fileID
func (_m *StoreService) GetFile(ctx context.Context, fileID string) (*domain.FileContent, error) {
	ret := _m.Called(ctx, fileID)

	var r0 *domain.FileContent
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.FileContent); ok {
		r0 = rf(ctx, fileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.FileContent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMetadata provides a mock function with given fields: ctx, fileID
func (_m *StoreService) GetMetadata(ctx context.Context, fileID string) (*domain.FileInfo, error) {
	ret := _m.Called(ctx, fileID)

	var r0 *domain.FileInfo
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.FileInfo); ok {
		r0 = rf(ctx, fileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.FileInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutFile provides a mock function with given fields: ctx, file, fi
func (_m *StoreService) PutFile(ctx context.Context, file io.Reader, fi *domain.FileInfo) (string, error) {
	ret := _m.Called(ctx, file, fi)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, io.Reader, *domain.FileInfo) string); ok {
		r0 = rf(ctx, file, fi)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, io.Reader, *domain.FileInfo) error); ok {
		r1 = rf(ctx, file, fi)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}