package minio

import (
	"context"
	"github.com/africarealty/server/src/kit/er"
)

const (
	ErrCodeErrMinioPutObject        = "S3-001"
	ErrCodeMinioCannotGetObject     = "S3-002"
	ErrCodeMinioCannotGetStatObject = "S3-003"
	ErrCodeMinioObjectNotFound      = "S3-004"
	ErrCodeMinioCreateBucket        = "S3-005"
	ErrCodeMinioRemoveObject        = "S3-006"
	ErrCodeMinioNew                 = "S3-007"
	ErrCodeMinioCopyObject          = "S3-008"
)

var (
	ErrMinioPutObject = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeErrMinioPutObject, "").C(ctx).Err()
	}
	ErrMinioCannotGetObject = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeMinioCannotGetObject, "").C(ctx).Err()
	}
	ErrMinioCannotGetStatObject = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeMinioCannotGetStatObject, "").C(ctx).Err()
	}
	ErrMinioObjectNotFound = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeMinioObjectNotFound, "").C(ctx).Err()
	}
	ErrMinioCreateBucket = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeMinioCreateBucket, "").C(ctx).Err()
	}
	ErrMinioRemoveObject = func(cause error, ctx context.Context, fileId string) error {
		return er.WrapWithBuilder(cause, ErrCodeMinioRemoveObject, "").C(ctx).F(er.FF{"FileID ": fileId}).Err()
	}
	ErrMinioCopyObject = func(cause error, ctx context.Context, fileId string) error {
		return er.WrapWithBuilder(cause, ErrCodeMinioCopyObject, "").C(ctx).F(er.FF{"FileID ": fileId}).Err()
	}
	ErrMinioNew = func(cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeMinioNew, "").Err()
	}
)
