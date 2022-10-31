package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	kitContext "github.com/africarealty/server/src/kit/context"
	"github.com/africarealty/server/src/kit/er"
	kitHttp "github.com/africarealty/server/src/kit/http"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/service"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"time"
)

type Sdk struct {
	cfg         *service.CfgSdk
	baseUrl     string
	accessToken string
	logger      *log.Logger
}

func New(cfg *service.CfgSdk) *Sdk {
	rand.Seed(time.Now().UnixNano())
	sdk := &Sdk{
		cfg: cfg,
	}
	sdk.logger = log.Init(cfg.Log)
	return sdk
}

func (s *Sdk) l() log.CLogger {
	return s.logFn()().Srv("ar-sdk").Cmp("sdk")
}

func (s *Sdk) logFn() log.CLoggerFunc {
	return func() log.CLogger {
		return log.L(s.logger).Srv("africarealty").Node("africarealty")
	}
}

func (s *Sdk) Close(ctx context.Context) {}

func (s *Sdk) do(ctx context.Context, url, token, verb string, payload []byte) ([]byte, error) {
	l := s.l().C(ctx).Mth("do").F(log.FF{"url": url, "verb": verb, "pl": string(payload)}).Trc()
	client := &http.Client{}
	var rqReader io.Reader
	if payload != nil {
		rqReader = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(verb, url, rqReader)
	if err != nil {
		return nil, err
	}

	// setup separate connections for each call
	req.Close = true

	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	if rCtx, err := kitContext.MustRequest(ctx); err == nil {
		req.Header.Add("RequestId", rCtx.GetRequestId())
	} else {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, _ := ioutil.ReadAll(resp.Body)
	l.C(ctx).F(log.FF{"resp": string(data), "status": resp.StatusCode}).Dbg()

	// check app error
	httpErr := &kitHttp.Error{}
	_ = json.Unmarshal(data, &httpErr)
	if httpErr != nil && httpErr.Code != "" && httpErr.Message != "" {
		return nil, er.WithBuilder(httpErr.Code, httpErr.Message).Err()
	}

	return data, nil
}

func (s *Sdk) UploadFile(ctx context.Context, url, token string, fileParamName string, fileToUpload string, fields map[string]string) ([]byte, error) {
	l := s.l().C(ctx).Mth("upload-file").F(log.FF{"fileName": fileToUpload}).Trc()
	// Prepare a form that you will submit to that URL.
	client := &http.Client{}
	// open file to upload
	file, err := os.Open(fileToUpload)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer

	// detect correct mimetype (SHOULD USE DetectFile! DetectReader break the file)
	mtype, _ := mimetype.DetectFile(fileToUpload)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fileParamName, file.Name()))
	h.Set("Content-Type", mtype.String())
	l.Dbg("file mimetype detected: ", mtype.String())
	fw, _ = w.CreatePart(h)
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}
	for key, value := range fields {
		_ = w.WriteField(key, value)
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	_ = w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	if rCtx, err := kitContext.MustRequest(ctx); err == nil {
		req.Header.Add("RequestId", rCtx.GetRequestId())
	} else {
		return nil, err
	}

	// Submit the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(resp.Body)

	// check app error
	httpErr := &kitHttp.Error{}
	_ = json.Unmarshal(data, &httpErr)
	if httpErr != nil && httpErr.Code != "" && httpErr.Message != "" {
		return nil, er.WithBuilder(httpErr.Code, httpErr.Message).Err()
	}

	return data, nil

}

func (s *Sdk) POST(ctx context.Context, url, token string, payload []byte) ([]byte, error) {
	return s.do(ctx, url, token, "POST", payload)
}

func (s *Sdk) PUT(ctx context.Context, url, token string, payload []byte) ([]byte, error) {
	return s.do(ctx, url, token, "PUT", payload)
}

func (s *Sdk) DELETE(ctx context.Context, url, token string, payload []byte) ([]byte, error) {
	return s.do(ctx, url, token, "DELETE", payload)
}

func (s *Sdk) GET(ctx context.Context, url, token string) ([]byte, error) {
	return s.do(ctx, url, token, "GET", nil)
}

func (s *Sdk) HealthCheck(ctx context.Context) (string, error) {
	rs, err := s.GET(ctx, fmt.Sprintf("%s/api/ready", s.baseUrl), s.accessToken)
	if err != nil {
		return "", err
	}
	return string(rs), nil
}
