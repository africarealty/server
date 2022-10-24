package impl

import (
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/errors/filestore"
	"github.com/africarealty/server/src/kit/er"
	kitTestSuite "github.com/africarealty/server/src/kit/test/suite"
	"github.com/africarealty/server/src/mocks"
	"github.com/africarealty/server/src/service"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type storeTestSuite struct {
	kitTestSuite.Suite
	storage *mocks.FileStorageRepository
	service domain.StoreService
}

func (s *storeTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
}

func (s *storeTestSuite) SetupTest() {
	s.storage = &mocks.FileStorageRepository{}
	s.service = NewStoreService(s.storage)
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(storeTestSuite))
}

func (s *storeTestSuite) Test_ValidateFileInfo_WhenFileNameEmpty_Fail() {
	svc := s.service.(*serviceImpl)
	fi := &domain.FileInfo{}
	err := svc.validateFileInfo(s.Ctx, fi)
	s.AssertAppErr(err, errors.ErrCodeStoreFieldRequired)
	appErr, _ := er.Is(err)
	s.Equal("Filename", appErr.Fields()["param"])
}

func (s *storeTestSuite) Test_ValidateFileInfo_WhenFileNameTooShort_Fail() {
	svc := s.service.(*serviceImpl)
	fi := &domain.FileInfo{Filename: "ert"}
	err := svc.validateFileInfo(s.Ctx, fi)
	s.AssertAppErr(err, errors.ErrCodeStoreInvalidFilename)
}

func (s *storeTestSuite) Test_ValidateFileInfo_WhenFileNameHasIlligalSymbols_Fail() {
	svc := s.service.(*serviceImpl)
	fi := &domain.FileInfo{Filename: "ertsfs%gjhg*"}
	err := svc.validateFileInfo(s.Ctx, fi)
	s.AssertAppErr(err, errors.ErrCodeStoreInvalidFilename)
}

func (s *storeTestSuite) Test_ValidateFileInfo_WhenFileNameHasLegalSymbols_Ok() {
	svc := s.service.(*serviceImpl)
	fi := &domain.FileInfo{
		Filename:    "filenam_filenam-.filename",
		BucketName:  "bucket",
		ContentType: "content",
	}
	err := svc.validateFileInfo(s.Ctx, fi)
	s.Nil(err)
}

func (s *storeTestSuite) Test_ValidateFileInfo_WhenBucketEmpty_Fail() {
	svc := s.service.(*serviceImpl)
	fi := &domain.FileInfo{
		Filename: "filename",
	}
	err := svc.validateFileInfo(s.Ctx, fi)
	s.AssertAppErr(err, errors.ErrCodeStoreFieldRequired)
	appErr, _ := er.Is(err)
	s.Equal("BucketName", appErr.Fields()["param"])
}

func (s *storeTestSuite) Test_ValidateFileInfo_WhenBucketNotMatch_Fail() {
	svc := s.service.(*serviceImpl)
	fi := &domain.FileInfo{Filename: "filename", BucketName: "hjdkfs25764$%^"}
	err := svc.validateFileInfo(s.Ctx, fi)
	s.AssertAppErr(err, errors.ErrCodeStoreInvalidBucketName)
}

func (s *storeTestSuite) Test_ValidateFileInfo_WhenContentTypeEmpty_Fail() {
	svc := s.service.(*serviceImpl)
	fi := &domain.FileInfo{Filename: "filename", BucketName: "bucket"}
	err := svc.validateFileInfo(s.Ctx, fi)
	s.AssertAppErr(err, errors.ErrCodeStoreFieldRequired)
	appErr, _ := er.Is(err)
	s.Equal("ContentType", appErr.Fields()["param"])
}

func (s *storeTestSuite) Test_ExtractBucketName() {
	svc := s.service.(*serviceImpl)

	bn, err := svc.extractBucketName(s.Ctx, "bnname:goodexample.jpg")
	s.Nil(err)
	s.Equal("bnname", bn)
	s.Nil(err)

	bn, err = svc.extractBucketName(s.Ctx, "bnname:hard:example.jpg")
	s.Nil(err)
	s.Equal("bnname", bn)

	bn, err = svc.extractBucketName(s.Ctx, "bnname:very:hard:example:")
	s.Equal("bnname", bn)
	s.Nil(err)

	bn, err = svc.extractBucketName(s.Ctx, "bnname_goodexample.jpg")
	s.Equal("", bn)
	s.AssertAppErr(err, errors.ErrCodeStoreInvalidFileID)
}

func (s *storeTestSuite) getFileInfo() *domain.FileInfo {
	return &domain.FileInfo{
		Id:           "bucket:filename.png",
		Filename:     "filename.png",
		BucketName:   "bucket",
		Extension:    "png",
		LastModified: "",
		Size:         100,
		ContentType:  "content",
		Metadata:     map[string]string{},
	}
}

func (s *storeTestSuite) Test_PutFile_WhenNoBucket_Ok() {
	file := strings.NewReader("test content")
	fi := s.getFileInfo()
	s.storage.On("IsBucketExist", s.Ctx, fi.BucketName).Return(false)
	s.storage.On("CreateBucket", s.Ctx, fi.BucketName).Return(nil)
	s.storage.On("Put", s.Ctx, fi, file).Return(nil)
	fileId, err := s.service.PutFile(s.Ctx, file, fi)
	s.Nil(err)
	s.NotEmpty(fileId)
	s.AssertNumberOfCalls(&s.storage.Mock, "CreateBucket", 1)
	s.AssertNumberOfCalls(&s.storage.Mock, "Put", 1)
}

func (s *storeTestSuite) Test_GetFile_WhenNoBucket_Fail() {
	s.storage.On("IsBucketExist", s.Ctx, "bucket").Return(false)
	_, err := s.service.GetFile(s.Ctx, "bucket:filename")
	s.AssertAppErr(err, errors.ErrCodeStoreBucketNotExist)
}

func (s *storeTestSuite) Test_GetFile_WhenInvalidFileId_Fail() {
	s.storage.On("IsBucketExist", s.Ctx, "bucket").Return(true)
	_, err := s.service.GetFile(s.Ctx, "bucketfilename")
	s.AssertAppErr(err, errors.ErrCodeStoreInvalidFileID)
}

func (s *storeTestSuite) Test_GetFile_Ok() {
	file := strings.NewReader("test content")
	fi := s.getFileInfo()
	s.storage.On("IsBucketExist", s.Ctx, fi.BucketName).Return(true)
	s.storage.On("Get", s.Ctx, fi.BucketName, fi.Id).Return(file, nil)
	s.storage.On("GetMetadata", s.Ctx, fi.BucketName, fi.Id).Return(fi, nil)
	fileContent, err := s.service.GetFile(s.Ctx, fi.Id)
	s.Nil(err)
	s.NotNil(fileContent)
	s.NotNil(fileContent.Content)
	s.AssertNumberOfCalls(&s.storage.Mock, "Get", 1)
}

func (s *storeTestSuite) Test_GetMetadata_Ok() {
	fi := s.getFileInfo()
	s.storage.On("IsBucketExist", s.Ctx, fi.BucketName).Return(true)
	s.storage.On("GetMetadata", s.Ctx, fi.BucketName, fi.Id).Return(fi, nil)
	md, err := s.service.GetMetadata(s.Ctx, fi.Id)
	s.Nil(err)
	s.NotNil(md)
	s.AssertNumberOfCalls(&s.storage.Mock, "GetMetadata", 1)
}

func (s *storeTestSuite) Test_DeleteFile_Ok() {
	fi := s.getFileInfo()
	s.storage.On("IsBucketExist", s.Ctx, fi.BucketName).Return(true)
	s.storage.On("IsFileExist", s.Ctx, fi.BucketName, fi.Id).Return(true)
	s.storage.On("Delete", s.Ctx, fi.BucketName, fi.Id).Return(nil)
	err := s.service.DeleteFile(s.Ctx, fi.Id)
	s.Nil(err)
	s.AssertNumberOfCalls(&s.storage.Mock, "Delete", 1)
}

func (s *storeTestSuite) Test_DeleteFile_WhenNotExistsError() {
	fi := s.getFileInfo()
	s.storage.On("IsBucketExist", s.Ctx, fi.BucketName).Return(true)
	s.storage.On("IsFileExist", s.Ctx, fi.BucketName, fi.Id).Return(false)
	s.storage.On("Delete", s.Ctx, fi.BucketName, fi.Id).Return(nil)
	err := s.service.DeleteFile(s.Ctx, fi.Id)
	s.AssertAppErr(err, errors.ErrCodeStoreFileNotExists)
}
