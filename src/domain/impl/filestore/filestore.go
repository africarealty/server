package impl

import (
	"context"
	"fmt"
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/errors/filestore"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
	"io"
	"path"
	"regexp"
	"strings"
)

type serviceImpl struct {
	storage domain.FileStorageRepository
}

func NewStoreService(storage domain.FileStorageRepository) domain.StoreService {
	return &serviceImpl{
		storage: storage,
	}
}

func (t *serviceImpl) l() log.CLogger {
	return service.L().Cmp("store-svc")
}

func (t *serviceImpl) validateFileInfo(ctx context.Context, fi *domain.FileInfo) error {

	// check filename
	if fi.Filename == "" {
		return errors.ErrStoreFieldRequired(ctx, "Filename")
	}
	matched, _ := regexp.MatchString(`^[^:%*]{4,}$`, fi.Filename)
	if !matched {
		return errors.ErrStoreFilenameInvalid(ctx, "[^:%*]{4,}$")
	}

	// check bucket name
	if fi.BucketName == "" {
		return errors.ErrStoreFieldRequired(ctx, "BucketName")
	}
	matched, _ = regexp.MatchString(`^[a-z0-9]{4,}$`, fi.BucketName)
	if !matched {
		return errors.ErrStoreBucketNameInvalid(ctx, "[a-z0-9]{4,}")
	}

	//check content
	if fi.ContentType == "" {
		return errors.ErrStoreFieldRequired(ctx, "ContentType")
	}
	return nil
}

func (t *serviceImpl) extractBucketName(ctx context.Context, fileID string) (string, error) {
	s := strings.Split(fileID, ":")
	if len(s) < 2 {
		return "", errors.ErrStoreInvalidFileID(ctx, fileID)
	}
	return s[0], nil
}

func (t *serviceImpl) BuildFileID(bucketName string, filename string) string {
	var extension = path.Ext(filename)
	return fmt.Sprintf("%s:%s%s", bucketName, kit.NewRandString(), extension)
}

func (t *serviceImpl) PutFile(ctx context.Context, file io.Reader, fi *domain.FileInfo) (string, error) {
	l := t.l().C(ctx).Mth("put-file").Dbg()

	if err := t.validateFileInfo(ctx, fi); err != nil {
		return "", err
	}

	// generate file ID, if not provided
	if fi.Id == "" {
		fi.Id = t.BuildFileID(fi.BucketName, fi.Filename)
	}
	l.F(log.FF{"fileId": fi.Id})

	if !t.storage.IsBucketExist(ctx, fi.BucketName) {
		err := t.storage.CreateBucket(ctx, fi.BucketName)
		if err != nil {
			return "", errors.ErrStoreCreateBucket(err, ctx)
		}
		l.Dbg("bucket created: ", fi.BucketName)
	}

	// put to storage
	err := t.storage.Put(ctx, fi, file)
	if err != nil {
		return "", err
	}
	l.Dbg("ok")
	return fi.Id, nil
}

func (t *serviceImpl) GetFile(ctx context.Context, fileID string) (*domain.FileContent, error) {
	t.l().C(ctx).Mth("get-file").Dbg()
	bucket, err := t.extractBucketName(ctx, fileID)
	if err != nil {
		return nil, err
	}
	if !t.storage.IsBucketExist(ctx, bucket) {
		return nil, errors.ErrStoreBucketNotExist(ctx, bucket)
	}
	md, err := t.GetMetadata(ctx, fileID)
	if err != nil {
		return nil, err
	}
	reader, err := t.storage.Get(ctx, bucket, fileID)
	if err != nil {
		return nil, err
	}
	r := &domain.FileContent{
		FileID:      fileID,
		Filename:    md.Filename,
		Extension:   md.Extension,
		ContentType: md.ContentType,
		Content:     reader,
	}
	return r, nil
}

func (t *serviceImpl) GetMetadata(ctx context.Context, fileID string) (*domain.FileInfo, error) {
	t.l().C(ctx).Mth("get-metadata").Dbg()
	bucket, err := t.extractBucketName(ctx, fileID)
	if err != nil {
		return nil, err
	}
	if !t.storage.IsBucketExist(ctx, bucket) {
		return nil, errors.ErrStoreBucketNotExist(ctx, bucket)
	}
	return t.storage.GetMetadata(ctx, bucket, fileID)
}

func (t *serviceImpl) DeleteFile(ctx context.Context, fileID string) error {
	t.l().C(ctx).Mth("delete-file").Dbg()
	// extract bucket
	bucket, err := t.extractBucketName(ctx, fileID)
	if err != nil {
		return err
	}
	// check bucket exists
	if !t.storage.IsBucketExist(ctx, bucket) {
		return errors.ErrStoreBucketNotExist(ctx, bucket)
	}
	if !t.storage.IsFileExist(ctx, bucket, fileID) {
		return errors.ErrStoreFileNotExists(ctx)
	}
	// delete file
	return t.storage.Delete(ctx, bucket, fileID)
}
