package errors

import (
	"context"
	"github.com/africarealty/server/src/kit/er"
)

var (
	ErrEmailValidationInvalidEmail = func(ctx context.Context, p string) error {
		return er.WithBuilder(ErrCodeEmailValidationInvalidEmail, "invalid email").Business().C(ctx).F(er.FF{"email": p}).Err()
	}
	ErrEmailValidationInvalidFrom = func(ctx context.Context, p string) error {
		return er.WithBuilder(ErrCodeEmailValidationInvalidFrom, "invalid from").Business().C(ctx).F(er.FF{"from": p}).Err()
	}
	ErrEmailSmtpInvalidConfig = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeEmailSmtpInvalidConfig, "invalid config").Business().C(ctx).Err()
	}
	ErrEmailStorageCreateEmailDb = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeEmailStorageCreateEmailDb, "").C(ctx).Err()
	}
	ErrEmailStorageUpdateEmailDb = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeEmailStorageUpdateEmailDb, "").C(ctx).Err()
	}
	ErrEmailSmtpSend = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeEmailSmtpSend, "").C(ctx).Err()
	}
	ErrTemplateEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeTemplateEmpty, "template empty").Business().C(ctx).Err()
	}
	ErrTemplateIdEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeTemplateIdEmpty, "template id empty").Business().C(ctx).Err()
	}
	ErrTemplateTitleEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeTemplateTitleEmpty, "template title empty").Business().C(ctx).Err()
	}
	ErrTemplateBodyEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeTemplateBodyEmpty, "template body empty").Business().C(ctx).Err()
	}
	ErrTemplateNotFound = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeTemplateNotFound, "template not found").Business().C(ctx).Err()
	}
	ErrTemplateAlreadyExists = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeTemplateAlreadyExists, "template already exists").Business().C(ctx).Err()
	}
	ErrTemplateGenerator = func(cause error, ctx context.Context, id string) error {
		return er.WrapWithBuilder(cause, ErrCodeTemplateGenerator, "").C(ctx).F(er.FF{"templId": id}).Err()
	}
	ErrTemplateStorageDbSearch = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeTemplateStorageDbSearch, "").C(ctx).Err()
	}
	ErrTemplateStorageDbCreate = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeStorageTemplateDbCreate, "").C(ctx).Err()
	}
	ErrTemplateStorageDbUpdate = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeStorageTemplateDbUpdate, "").C(ctx).Err()
	}
	ErrTemplateStorageDbDelete = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeStorageTemplateDbDelete, "").C(ctx).Err()
	}
	ErrTemplateStorageAeroKey = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeTemplateStorageAeroKey, "").C(ctx).Err()
	}
	ErrTemplateStorageClearCache = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeTemplateStorageClearCache, "").C(ctx).Err()
	}
	ErrTemplateStorageGetCache = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeTemplateStorageGetCache, "").C(ctx).Err()
	}
	ErrTemplateStoragePutCache = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeTemplateStoragePutCache, "").C(ctx).Err()
	}
	ErrTemplateStorageGetDb = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeTemplateStorageGetDb, "").C(ctx).Err()
	}
)
