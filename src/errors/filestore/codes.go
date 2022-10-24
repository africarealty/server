package errors

const (
	ErrCodeStoreMaxFileSize       = "STORE-001"
	ErrCodeStoreFileNotFound      = "STORE-002"
	ErrCodeStoreFieldRequired     = "STORE-003"
	ErrCodeStoreInvalidFilename   = "STORE-004"
	ErrCodeStoreInvalidBucketName = "STORE-005"
	ErrCodeStoreCreateBucket      = "STORE-006"
	ErrCodeStorePutFile           = "STORE-007"
	ErrCodeStoreDeleteFile        = "STORE-008"
	ErrCodeStoreGetFile           = "STORE-009"
	ErrCodeStoreGetMetadata       = "STORE-010"
	ErrCodeStoreInvalidFileID     = "STORE-011"
	ErrCodeStoreCannotGetFile     = "STORE-012"
	ErrCodeStoreCannotStatFile    = "STORE-013"
	ErrCodeStoreNotExcelFileType  = "STORE-014"
	ErrCodeStoreReadExcel         = "STORE-015"
	ErrCodeStoreFetchExcelRows    = "STORE-016"
	ErrCodeStoreSheetOutOfIndex   = "STORE-017"
	ErrCodeStoreFileNotExists     = "STORE-018"
	ErrCodeStoreBucketNotExist    = "STORE-019"
	ErrCodeWriteToFileStreamFail  = "STORE-020"
	ErrCodeReadFromFileStreamFail = "STORE-021"
	ErrCodeStreamCloseFailed      = "STORE-022"
	ErrCodeFileReadFail           = "STORE-023"
	ErrCodeFileWriteFail          = "STORE-024"
	ErrCodeDecodeScaleFail        = "STORE-025"
	ErrCodeEncodeScaleFail        = "STORE-026"
	ErrCodeFileIsNotImage         = "STORE-027"
	ErrCodeNoScaleForImage        = "STORE-028"
	ErrCodeStoreEmptyFile         = "STORE-029"
	ErrCodeStorePutImage          = "STORE-030"
	ErrCodeClosingStreamFailed    = "STORE-031"
	ErrCodeStoreCopyStreamFailed  = "STORE-032"
)
