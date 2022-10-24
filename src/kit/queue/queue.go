package queue

import (
	"context"
)

type QueueType int

func (q QueueType) String() string {
	if q == QueueTypeAtLeastOnce {
		return "at-least-once"
	} else if q == QueueTypeAtMostOnce {
		return "at-most-once"
	}
	return ""
}

// delivery guaranties
const (
	QueueTypeAtLeastOnce = iota
	QueueTypeAtMostOnce
)

// Config queue configuration
type Config struct {
	Host      string
	Port      string
	ClusterId string
	Username  string
	Password  string
}

// EventHandler func type an event handler must correspond
type EventHandler func()

// Queue allows async communication with a message queue
type Queue interface {
	// Open opens connection
	// clientId must be unique
	Open(ctx context.Context, clientId string, options *Config) error
	// Declare allow preliminary declaration
	Declare(ctx context.Context, qt QueueType, topic string) error
	// Close closes connection
	Close() error
	// Publish publishes a message to topic
	Publish(ctx context.Context, qt QueueType, topic string, msg *Message) error
	// Subscribe subscribes on topic (pub-sub pattern: all subscribers receive the same message)
	Subscribe(qt QueueType, topic, durableId string, receiverChan chan<- []byte) error
	// SubscribeLB subscribes on topic with load balancing (producer - consumer pattern: one of the consumers receive a message in round-robin manner)
	// if more than one subscribers specify the same loadBalancingGroup, messages are balanced among all subscribers within the group
	// so that the only one subscriber gets the message
	SubscribeLB(qt QueueType, topic, loadBalancingGroup, durableId string, receiverChan chan<- []byte) error
	// SetLostConnectionHandler allows subscription on "queue lost connection" event
	SetLostConnectionHandler(EventHandler)
	// SetReconnectHandler allows subscription on "queue reconnected after losing connection" event
	SetReconnectHandler(EventHandler)
}
