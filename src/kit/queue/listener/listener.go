package listener

import (
	"context"
	"github.com/africarealty/server/src/kit/goroutine"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/kit/queue"
	"go.uber.org/atomic"
	"sync"
)

const (
	GoroutinesSubscriberPerTopic = 4
)

// QueueMessageHandler is a func which acts as message handler
type QueueMessageHandler func(payload []byte) error

// QueueListener supports multiple subscriptions on Queue topics
type QueueListener interface {
	// Add adds handlers
	Add(qt queue.QueueType, topic, durableId string, h ...QueueMessageHandler)
	// AddLb adds handlers with load balancing
	AddLb(qt queue.QueueType, topic, lbGroup, durableId string, h ...QueueMessageHandler)
	// ListenAsync starts goroutine which is listening incoming messages and calls proper handlers
	ListenAsync()
	// Stop stops listening
	Stop()
	// Clear clears all handlers
	Clear()
	// New registers a new listener
	New(topic string) Builder
}

// topicKey used as a key for handlers
type topicKey struct {
	Topic     string // Topic - Queue topic
	LbGroup   string // LbGroup - load balancing group
	DurableId string // DurableId - durable ID
}

type queueListener struct {
	sync.RWMutex
	queue         queue.Queue
	topicHandlers map[queue.QueueType]map[topicKey][]QueueMessageHandler
	quit          chan struct{}
	listening     *atomic.Bool
	logger        log.CLoggerFunc
}

func NewQueueListener(q queue.Queue, logger log.CLoggerFunc) QueueListener {

	th := map[queue.QueueType]map[topicKey][]QueueMessageHandler{}
	th[queue.QueueTypeAtLeastOnce] = make(map[topicKey][]QueueMessageHandler)
	th[queue.QueueTypeAtMostOnce] = make(map[topicKey][]QueueMessageHandler)

	ql := &queueListener{
		topicHandlers: th,
		listening:     atomic.NewBool(false),
		queue:         q,
		logger:        logger,
	}
	ql.queue.SetReconnectHandler(ql.queueReconnectHandler)
	ql.queue.SetLostConnectionHandler(ql.queueLostConnectionHandler)
	return ql
}

func (q *queueListener) l() log.CLogger {
	return q.logger().Pr("queue").Cmp("queue-listener")
}

func (q *queueListener) New(topic string) Builder {
	return newBuilder(q, topic)
}

func (q *queueListener) add(qt queue.QueueType, topic, lbGroup, durableId string, h ...QueueMessageHandler) {
	q.l().Mth("add").F(log.FF{"topic": topic, "lb": lbGroup, "durId": durableId}).Inf()

	q.Stop()

	q.Lock()
	defer q.Unlock()

	key := topicKey{Topic: topic, LbGroup: lbGroup, DurableId: durableId}

	var handlers []QueueMessageHandler
	handlers, ok := q.topicHandlers[qt][key]
	if !ok {
		handlers = []QueueMessageHandler{}
	}

	handlers = append(handlers, h...)
	q.topicHandlers[qt][key] = handlers

}

func (q *queueListener) Add(qt queue.QueueType, topic, durableId string, h ...QueueMessageHandler) {
	q.add(qt, topic, "", durableId, h...)
}

func (q *queueListener) AddLb(qt queue.QueueType, topic, lbGroup, durableId string, h ...QueueMessageHandler) {
	q.add(qt, topic, lbGroup, durableId, h...)
}

func (q *queueListener) worker(topic string, handlers []QueueMessageHandler, receiverChan chan []byte) {
	goroutine.New().
		WithLoggerFn(q.logger).
		Cmp("queue-listener").
		Mth("listen").
		WithRetry(goroutine.Unrestricted).
		Go(context.Background(),
			func() {
				l := q.l().Mth("worker").F(log.FF{"topic": topic}).Dbg("started")
				for {
					select {
					case msg := <-receiverChan:
						l := l.Dbg("handle message").TrcF("%s", string(msg))
						// run handler
						for _, handler := range handlers {
							if err := handler(msg); err != nil {
								l.E(err).St().Err()
							}
						}
					case <-q.quit:
						l.Dbg("stopped")
						return
					}
				}
			},
		)
}

func (q *queueListener) ListenAsync() {

	l := q.l().Mth("listen-async")

	q.quit = make(chan struct{})

	// go through all queue types
	for queueType, topicHandlers := range q.topicHandlers {
		// go through handlers of the queue type
		for key, hnd := range topicHandlers {

			qt := queueType
			topic := key.Topic
			loadBalancingGroup := key.LbGroup
			durableId := key.DurableId
			handlers := hnd

			l := l.F(log.FF{"topic": topic, "queue-type": queueType, "lb": loadBalancingGroup}).Dbg("prepare workers")

			receiverChan := make(chan []byte, 256)

			// if load balancing group specified, make load balanced subscription
			if loadBalancingGroup == "" {
				if err := q.queue.Subscribe(qt, topic, durableId, receiverChan); err != nil {
					l.E(err).St().Err()
					return
				}
			} else {
				if err := q.queue.SubscribeLB(qt, topic, loadBalancingGroup, durableId, receiverChan); err != nil {
					l.E(err).St().Err()
					return
				}
			}

			// run worker to handle messages
			for i := 0; i < GoroutinesSubscriberPerTopic; i++ {
				q.worker(topic, handlers, receiverChan)
			}

			l.Dbg("workers run")
		}
	}
	q.listening.Store(true)
}

func (q *queueListener) Stop() {
	q.l().Mth("stop").Inf()
	if q.listening.Load() {
		q.listening.Store(false)
		close(q.quit)
	}
}

func (q *queueListener) queueLostConnectionHandler() {
	q.l().Mth("lost-connection-handler").Inf()
	q.Stop()
}

func (q *queueListener) queueReconnectHandler() {
	q.l().Mth("reconnect-handler").Inf()
	q.Stop()
	q.ListenAsync()
}

func (q *queueListener) Clear() {
	q.l().Mth("clear").Inf()
	q.Stop()
	q.Lock()
	defer q.Unlock()
	q.topicHandlers[queue.QueueTypeAtLeastOnce] = make(map[topicKey][]QueueMessageHandler)
	q.topicHandlers[queue.QueueTypeAtMostOnce] = make(map[topicKey][]QueueMessageHandler)
}
