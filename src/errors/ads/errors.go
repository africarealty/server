package errors

import (
	"context"
	"github.com/africarealty/server/src/kit/er"
)

var (
	ErrAdsNotFound = func(ctx context.Context, adsId string) error {
		return er.WithBuilder(ErrCodeAdsNotFound, "not found").Business().F(er.FF{"adsId": adsId}).C(ctx).Err()
	}
	ErrAdsIdEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeAdsIdEmpty, "id empty").Business().C(ctx).Err()
	}
	ErrAdsStorageAeroKey = func(err error, ctx context.Context) error {
		return er.WrapWithBuilder(err, ErrCodeAdsStorageAeroKey, "").C(ctx).Err()
	}
	ErrAdsStorageClearCache = func(err error, ctx context.Context) error {
		return er.WrapWithBuilder(err, ErrCodeAdsStorageClearCache, "").C(ctx).Err()
	}
	ErrAdsStorageGetCache = func(err error, ctx context.Context) error {
		return er.WrapWithBuilder(err, ErrCodeAdsStorageGetCache, "").C(ctx).Err()
	}
	ErrAdsStoragePutCache = func(err error, ctx context.Context) error {
		return er.WrapWithBuilder(err, ErrCodeAdsStoragePutCache, "").C(ctx).Err()
	}
	ErrAdsStorageCreate = func(err error, ctx context.Context) error {
		return er.WrapWithBuilder(err, ErrCodeAdsStorageCreate, "").C(ctx).Err()
	}
	ErrAdsStorageUpdate = func(err error, ctx context.Context) error {
		return er.WrapWithBuilder(err, ErrCodeAdsStorageUpdate, "").C(ctx).Err()
	}
	ErrAdsStorageDelete = func(err error, ctx context.Context) error {
		return er.WrapWithBuilder(err, ErrCodeAdsStorageDelete, "").C(ctx).Err()
	}
	ErrAdsStorageGetDb = func(err error, ctx context.Context) error {
		return er.WrapWithBuilder(err, ErrCodeAdsStorageGetDb, "").C(ctx).Err()
	}
	ErrAdsStorageDbNextSeqVal = func(err error, ctx context.Context) error {
		return er.WrapWithBuilder(err, ErrCodeAdsStorageDbNextSeqVal, "").C(ctx).Err()
	}
)
