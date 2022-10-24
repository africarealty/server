package service

import (
	"context"
	kitContext "github.com/africarealty/server/src/kit/context"
	"github.com/africarealty/server/src/kit/er"
	"github.com/africarealty/server/src/kit/queue"
)

// BaseService can be used as a base service providing some helpers
type BaseService struct {
	Queue queue.Queue
}

const (
	ErrCodeBaseModelCannotPublishToQueue = "CMN-001"
)

var (
	ErrBaseModelCannotPublishToQueue = func(ctx context.Context, topic string) error {
		return er.WithBuilder(ErrCodeBaseModelCannotPublishToQueue, "cannot publish to topic").C(ctx).F(er.FF{"topic": topic}).Err()
	}
)

// Publish is helper method to publish a message to queue
// Note, it covers payload with &queue.Message, so you have to pass a pure payload object
func (s *BaseService) Publish(ctx context.Context, o interface{}, qt queue.QueueType, topic string) error {

	m := &queue.Message{Payload: o}

	if rCtx, ok := kitContext.Request(ctx); ok {
		m.Ctx = rCtx
	} else {
		return ErrBaseModelCannotPublishToQueue(ctx, topic)
	}

	return s.Queue.Publish(ctx, qt, topic, m)

}
