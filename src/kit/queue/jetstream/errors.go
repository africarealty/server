package jetstream

import "github.com/africarealty/server/src/kit/er"

const (
	ErrCodeJsNoOpenConn           = "JS-001"
	ErrCodeJsQtNotSupported       = "JS-002"
	ErrCodeJsConnect              = "JS-003"
	ErrCodeJsPublishAtLeastOnce   = "JS-005"
	ErrCodeJsPublishAtMostOnce    = "JS-006"
	ErrCodeJsSubscribeAtLeastOnce = "JS-007"
	ErrCodeJsSubscribeAtMostOnce  = "JS-008"
)

var (
	ErrJsNoOpenConn     = func() error { return er.WithBuilder(ErrCodeJsNoOpenConn, "no open connections").Err() }
	ErrJsQtNotSupported = func(qt int) error {
		return er.WithBuilder(ErrCodeJsQtNotSupported, "queue type not supported").F(er.FF{"qt": qt}).Err()
	}
	ErrJsConnect              = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeJsConnect, "").Err() }
	ErrJsPublishAtLeastOnce   = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeJsPublishAtLeastOnce, "").Err() }
	ErrJsPublishAtMostOnce    = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeJsPublishAtMostOnce, "").Err() }
	ErrJsSubscribeAtLeastOnce = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeJsSubscribeAtLeastOnce, "").Err() }
	ErrJsSubscribeAtMostOnce  = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeJsSubscribeAtMostOnce, "").Err() }
)
