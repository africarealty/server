package listener

import "github.com/africarealty/server/src/kit/queue"

// Builder simplifies listener registration
type Builder interface {
	// AtMostOnce adds a listener with "at most once" delivery guaranty (implemented by NATS without durability)
	AtMostOnce() Builder
	// AtLeastOnce adds a listener with "at least once" delivery guaranty (implemented by STAN with durability)
	// read about durable here https://github.com/nats-io/stan.go
	AtLeastOnce(durableId string) Builder
	// WithHandler specifies a new handler
	WithHandler(QueueMessageHandler) Builder
	// WithLoadBalancing specifies a listener as load balanced
	// messages for multiple listeners with the same group value are balanced
	WithLoadBalancing(lbGroup string) Builder
	// Add adds a listener to queue listener
	Add()
}

type builderImpl struct {
	queueListener QueueListener
	queueType     queue.QueueType
	lbGroup       string
	topic         string
	handlers      []QueueMessageHandler
	durableId     string
}

func newBuilder(queueListener QueueListener, topic string) Builder {
	return &builderImpl{
		queueListener: queueListener,
		topic:         topic,
		queueType:     queue.QueueTypeAtMostOnce,
	}
}

func (b *builderImpl) AtMostOnce() Builder {
	b.queueType = queue.QueueTypeAtMostOnce
	return b
}

func (b *builderImpl) AtLeastOnce(durableId string) Builder {
	b.queueType = queue.QueueTypeAtLeastOnce
	b.durableId = durableId
	return b
}

func (b *builderImpl) WithHandler(handler QueueMessageHandler) Builder {
	b.handlers = append(b.handlers, handler)
	return b
}

func (b *builderImpl) WithLoadBalancing(lbGroup string) Builder {
	b.lbGroup = lbGroup
	return b
}

func (b *builderImpl) Add() {
	if b.lbGroup == "" {
		b.queueListener.Add(b.queueType, b.topic, b.durableId, b.handlers...)
	} else {
		b.queueListener.AddLb(b.queueType, b.topic, b.lbGroup, b.durableId, b.handlers...)
	}
}
