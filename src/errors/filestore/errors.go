package errors

import (
	"context"
	"github.com/africarealty/server/src/kit/er"
)

var (
	ErrStoreMaxFileSize = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeStoreMaxFileSize, "max file size exceed").Business().C(ctx).Err()
	}
	ErrStoreFieldRequired = func(ctx context.Context, fieldName string) error {
		return er.WithBuilder(ErrCodeStoreFieldRequired, "param should be filled").Business().F(er.FF{"param": fieldName}).C(ctx).Err()
	}
	ErrStoreFilenameInvalid = func(ctx context.Context, regexp string) error {
		return er.WithBuilder(ErrCodeStoreInvalidFilename, "filename should pass regexp").Business().F(er.FF{"regexp ": regexp}).C(ctx).Err()
	}
	ErrStoreBucketNameInvalid = func(ctx context.Context, regexp string) error {
		return er.WithBuilder(ErrCodeStoreInvalidBucketName, "bucketName should pass regexp").Business().F(er.FF{"regexp ": regexp}).C(ctx).Err()
	}
	ErrStoreCreateBucket = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeStoreCreateBucket, "").C(ctx).Err()
	}
	ErrStorePutFile = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeStorePutFile, "").C(ctx).Err()
	}
	ErrStoreGetFile = func(cause error, ctx context.Context, fileID string) error {
		return er.WrapWithBuilder(cause, ErrCodeStoreGetFile, "").F(er.FF{"FileID ": fileID}).C(ctx).Err()
	}
	ErrStoreInvalidFileID = func(ctx context.Context, fileID string) error {
		return er.WithBuilder(ErrCodeStoreInvalidFileID, "invalid FileID").Business().F(er.FF{"FileID ": fileID}).C(ctx).Err()
	}
	ErrStoreNotExcelFileType = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeStoreNotExcelFileType, "not excel file content-type").Business().C(ctx).Err()
	}
	ErrStoreReaderExcel = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeStoreReadExcel, "").C(ctx).Err()
	}
	ErrStoreFetchExcelRows = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeStoreFetchExcelRows, "").C(ctx).Err()
	}
	ErrStoreExcelSheetOutOfIndex = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeStoreSheetOutOfIndex, "sheet index is out of range").Business().C(ctx).Err()
	}
	ErrStoreFileNotExists = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeStoreFileNotExists, "file not exists").Business().C(ctx).Err()
	}
	ErrStoreBucketNotExist = func(ctx context.Context, bucketName string) error {
		return er.WithBuilder(ErrCodeStoreBucketNotExist, "bucket not exist").Business().F(er.FF{"BucketName ": bucketName}).C(ctx).Err()
	}
	ErrWriteToFileStreamFailed = func(ctx context.Context, cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeWriteToFileStreamFail, "write to file stream failed").C(ctx).Err()
	}
	ErrReadFromFileStreamFailed = func(ctx context.Context, cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeReadFromFileStreamFail, "read from file stream failed").C(ctx).Err()
	}
	ErrStreamCloseFailed = func(ctx context.Context, cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeStreamCloseFailed, "close stream failed").C(ctx).Err()
	}
	ErrFileReadFailed = func(ctx context.Context, cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeFileReadFail, "file reading failed").C(ctx).Err()
	}
	ErrFileWriteFailed = func(ctx context.Context, cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeFileWriteFail, "file reading failed").C(ctx).Err()
	}
	ErrStoreDecodeScale = func(ctx context.Context, cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeDecodeScaleFail, "image decode file for scale failed").C(ctx).Err()
	}
	ErrStoreEncodeScale = func(ctx context.Context, cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeEncodeScaleFail, "image encode file for scale failed").C(ctx).Err()
	}
	ErrStoreFileIsNotImage = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeFileIsNotImage, "file type isn't image").Business().C(ctx).Err()
	}
	ErrStoreNoScaleForImage = func(ctx context.Context, fileID, scale string) error {
		return er.WithBuilder(ErrCodeNoScaleForImage, "no scale for image with the given size").Business().F(er.FF{"FileID": fileID, "scale": scale}).C(ctx).Err()
	}
	ErrStoreCopyStreamFailed = func(ctx context.Context, cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeStoreCopyStreamFailed, "copy stream failed").C(ctx).Err()
	}
)
