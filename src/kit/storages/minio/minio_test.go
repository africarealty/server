//go:build integration
// +build integration

package minio

import (
	"bytes"
	"context"
	"github.com/africarealty/server/src/kit"
	kitContext "github.com/africarealty/server/src/kit/context"
	"github.com/africarealty/server/src/kit/log"
	"io/ioutil"
	"testing"
)

var Logger = log.Init(&log.Config{Level: log.TraceLevel})

func lf() log.CLoggerFunc {
	return func() log.CLogger {
		return log.L(Logger).Srv("test")
	}
}

var (
	FileBytes = []byte{76, 111, 114}
)

func Test(t *testing.T) {
	client, err := New(&Config{
		Host:      "localhost",
		Port:      "8062",
		AccessKey: "minioaccesskey",
		SecretKey: "miniosecretkey",
		Ssl:       false,
	}, lf())

	if err != nil {
		t.Fatal(err)
	}

	fi := &FileInfo{
		Id:         kit.NewRandString(),
		BucketName: "testbucket",
		Metadata: map[string]string{
			"Key": "value",
		},
	}

	ctx := kitContext.NewRequestCtx().Test().ToContext(context.Background())
	if !client.IsBucketExist(ctx, fi.BucketName) {
		err := client.CreateBucket(ctx, fi.BucketName)
		if err != nil {
			t.Fatal(err)
		}
	}

	// put file
	file := bytes.NewReader(FileBytes)
	err = client.Put(ctx, fi, file)
	if err != nil {
		t.Fatal(err)
	}

	// get file
	fileRead, err := client.Get(ctx, fi.BucketName, fi.Id)
	if err != nil {
		t.Fatal(err)
	}
	buf, err := ioutil.ReadAll(fileRead)
	if err != nil {
		t.Fatal(err)
	}
	if len(buf) <= 0 {
		t.Fatal("File was not downloaded")
	}

	// get metadata
	meta, err := client.GetMetadata(ctx, fi.BucketName, fi.Id)
	if err != nil {
		t.Fatal(err, "Cannot get metadata for file")
		return
	}
	val, ok := meta.Metadata["Key"]

	if !ok || val != "value" {
		t.Fatal("Metadata was not downloaded")
	}
}
