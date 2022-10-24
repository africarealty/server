// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	domain "github.com/africarealty/server/src/domain"

	mock "github.com/stretchr/testify/mock"
)

// FileStorageRepository is an autogenerated mock type for the FileStorageRepository type
type FileStorageRepository struct {
	mock.Mock
}

// Copy provides a mock function with given fields: ctx, srcBucketName, srcFileID, destBucketName, destFileId
func (_m *FileStorageRepository) Copy(ctx context.Context, srcBucketName string, srcFileID string, destBucketName string, destFileId string) error {
	ret := _m.Called(ctx, srcBucketName, srcFileID, destBucketName, destFileId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) error); ok {
		r0 = rf(ctx, srcBucketName, srcFileID, destBucketName, destFileId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateBucket provides a mock function with given fields: ctx, bucketName
func (_m *FileStorageRepository) CreateBucket(ctx context.Context, bucketName string) error {
	ret := _m.Called(ctx, bucketName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, bucketName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, bucketName, fileID
func (_m *FileStorageRepository) Delete(ctx context.Context, bucketName string, fileID string) error {
	ret := _m.Called(ctx, bucketName, fileID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, bucketName, fileID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, bucketName, fileID
func (_m *FileStorageRepository) Get(ctx context.Context, bucketName string, fileID string) (io.Reader, error) {
	ret := _m.Called(ctx, bucketName, fileID)

	var r0 io.Reader
	if rf, ok := ret.Get(0).(func(context.Context, string, string) io.Reader); ok {
		r0 = rf(ctx, bucketName, fileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.Reader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, bucketName, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMetadata provides a mock function with given fields: ctx, bucketName, fileID
func (_m *FileStorageRepository) GetMetadata(ctx context.Context, bucketName string, fileID string) (*domain.FileInfo, error) {
	ret := _m.Called(ctx, bucketName, fileID)

	var r0 *domain.FileInfo
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *domain.FileInfo); ok {
		r0 = rf(ctx, bucketName, fileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.FileInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, bucketName, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsBucketExist provides a mock function with given fields: ctx, bucketName
func (_m *FileStorageRepository) IsBucketExist(ctx context.Context, bucketName string) bool {
	ret := _m.Called(ctx, bucketName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, bucketName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsFileExist provides a mock function with given fields: ctx, bucketName, fileID
func (_m *FileStorageRepository) IsFileExist(ctx context.Context, bucketName string, fileID string) bool {
	ret := _m.Called(ctx, bucketName, fileID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, bucketName, fileID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Put provides a mock function with given fields: ctx, fi, file
func (_m *FileStorageRepository) Put(ctx context.Context, fi *domain.FileInfo, file io.Reader) error {
	ret := _m.Called(ctx, fi, file)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FileInfo, io.Reader) error); ok {
		r0 = rf(ctx, fi, file)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}