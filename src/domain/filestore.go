package domain

import (
	"context"
	"io"
)

// FileInfo is a meta info about stored file
type FileInfo struct {
	Id           string            // Id file ID
	Filename     string            // Filename - file name
	BucketName   string            // BucketName - bucket name (like storage folder)
	Extension    string            // Extension - file ext
	LastModified string            // LastModified - last modified date
	Size         int64             // Size - length of a file in bytes
	ContentType  string            // ContentType - content type
	Metadata     map[string]string // Metadata - some additional user params attached to a stored file
}

// FileContent represents a stored file with content
type FileContent struct {
	FileID      string    // FileID - file id in store
	Filename    string    // Filename - file name
	Extension   string    // Extension - file extension
	ContentType string    // ContentType - file mimetype
	Content     io.Reader // Content - file content
}

// StoreService - is a domain interface
// This is the main contract for your business logic
type StoreService interface {
	// PutFile saves file and return generated filename which can be used to GetFile
	PutFile(ctx context.Context, file io.Reader, fi *FileInfo) (string, error)
	// GetFile returns io.Reader on requested filename
	GetFile(ctx context.Context, fileID string) (*FileContent, error)
	// GetMetadata returns *FileInfo
	GetMetadata(ctx context.Context, fileID string) (*FileInfo, error)
	// DeleteFile delete file by FileID
	DeleteFile(ctx context.Context, fileID string) error
	// BuildFileID builds a new file ID
	BuildFileID(bucketName string, filename string) string
}

// FileStorageRepository - contract for interaction with file storage
type FileStorageRepository interface {
	// Put persists file
	Put(ctx context.Context, fi *FileInfo, file io.Reader) error
	// Get retrieves file by file ID as io.Reader
	Get(ctx context.Context, bucketName string, fileID string) (io.Reader, error)
	// IsFileExist checks if file with the given fileId exists in a bucket
	IsFileExist(ctx context.Context, bucketName string, fileID string) bool
	// GetMetadata retrieves file's  metadata
	GetMetadata(ctx context.Context, bucketName string, fileID string) (*FileInfo, error)
	// IsBucketExist checks if a bucket exists
	IsBucketExist(ctx context.Context, bucketName string) bool
	// CreateBucket creates a new bucket if it's not exists
	CreateBucket(ctx context.Context, bucketName string) error
	// Delete deletes file
	Delete(ctx context.Context, bucketName string, fileID string) error
	// Copy copies file
	Copy(ctx context.Context, srcBucketName, srcFileID, destBucketName, destFileId string) error
}
