package jetstream

import (
	"context"
	"encoding/json"
	"fmt"
	kitContext "github.com/africarealty/server/src/kit/context"
	"github.com/africarealty/server/src/kit/goroutine"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/kit/queue"
	"github.com/nats-io/nats.go"
	"regexp"
	"sync"
)

const (
	PingIntervalSec        = 2
	MaxOutPingsCount       = 5
	WaitBeforeReconnectSec = 3
	ReconnectIntervalSec   = 5
)

type jsImpl struct {
	conn                   *nats.Conn
	js                     nats.JetStreamContext
	connMtx                sync.RWMutex
	clientId               string
	logger                 log.CLoggerFunc
	ctx                    context.Context
	config                 *queue.Config
	lostConnectionHandlers []queue.EventHandler
	reconnectHandlers      []queue.EventHandler
	handlerMtx             sync.RWMutex
}

func New(logger log.CLoggerFunc) queue.Queue {
	return &jsImpl{
		logger: logger,
	}
}

func (s *jsImpl) l() log.CLogger {
	return s.logger().Pr("queue").Cmp("js")
}

func (s *jsImpl) SetLostConnectionHandler(handler queue.EventHandler) {
	s.handlerMtx.Lock()
	defer s.handlerMtx.Unlock()
	s.lostConnectionHandlers = append(s.lostConnectionHandlers, handler)
}

func (s *jsImpl) SetReconnectHandler(handler queue.EventHandler) {
	s.handlerMtx.Lock()
	defer s.handlerMtx.Unlock()
	s.reconnectHandlers = append(s.reconnectHandlers, handler)
}

func (s *jsImpl) Open(ctx context.Context, clientId string, config *queue.Config) error {
	l := s.l().Mth("open").F(log.FF{"client": clientId, "host": config.Host}).Dbg("connecting")

	s.connMtx.Lock()
	defer s.connMtx.Unlock()

	s.ctx = ctx
	s.clientId = clientId
	s.config = config

	url := fmt.Sprintf("nats://%s:%s", config.Host, config.Port)
	c, err := nats.Connect(url,
		nats.Name(clientId),
		nats.MaxPingsOutstanding(MaxOutPingsCount),
		nats.ReconnectHandler(s.ReconnectHandler),
		nats.DisconnectErrHandler(s.LostConnectionHandler))
	if err != nil {
		return ErrJsConnect(err)
	}
	s.conn = c
	s.js, err = c.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return ErrJsConnect(err)
	}

	l.Inf("ok")

	return nil
}

func (s *jsImpl) lostConnectionNotify() {
	s.handlerMtx.Lock()
	defer s.handlerMtx.Unlock()
	for _, handler := range s.lostConnectionHandlers {
		goroutine.New().
			WithLogger(s.l().Mth("lost-connection-notify")).
			Go(context.Background(), handler)
	}
}

func (s *jsImpl) reconnectNotify() {
	s.handlerMtx.Lock()
	defer s.handlerMtx.Unlock()
	for _, handler := range s.reconnectHandlers {
		goroutine.New().
			WithLogger(s.l().Mth("reconnect-notify")).
			Go(context.Background(), handler)
	}
}

func (s *jsImpl) ReconnectHandler(connection *nats.Conn) {
	s.l().Mth("reconnect-handler").Dbg()
	// notify subscribers
	s.reconnectNotify()
}

func (s *jsImpl) LostConnectionHandler(connection *nats.Conn, err error) {
	l := s.l().Mth("lost-connection-handler")
	if err == nil {
		return
	}
	l.E(err).Err("connection lost")
	// notify subscribers
	s.lostConnectionNotify()
}

var rg = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func (s *jsImpl) topicToStream(topic string) string {
	return rg.ReplaceAllString(topic, "")
}

func (s *jsImpl) Declare(ctx context.Context, qt queue.QueueType, topic string) error {
	s.l().Mth("declare").F(log.FF{"topic": topic, "type": qt.String()}).Dbg()
	_, err := s.js.AddStream(&nats.StreamConfig{
		Name:     s.topicToStream(topic),
		Subjects: []string{topic},
	})
	return err
}

func (s *jsImpl) Close() error {
	l := s.l().Mth("close")
	s.connMtx.Lock()
	defer s.connMtx.Unlock()
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
		l.Inf("closed")
	}
	return nil
}

func (s *jsImpl) Publish(ctx context.Context, qt queue.QueueType, topic string, msg *queue.Message) error {
	l := s.l().Mth("publish").F(log.FF{"topic": topic, "type": qt.String()})
	s.connMtx.RLock()
	defer s.connMtx.RUnlock()
	if msg.Ctx == nil {
		msg.Ctx = kitContext.NewRequestCtx().Queue().WithNewRequestId()
	}
	l.C(msg.Ctx.ToContext(context.Background()))
	if s.conn == nil {
		return ErrJsNoOpenConn()
	}
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	l.Dbg("ok").TrcF("%s\n", string(m))
	if qt == queue.QueueTypeAtLeastOnce {
		// publish async
		_, err = s.js.PublishAsync(topic, m)
		if err != nil {
			return ErrJsPublishAtLeastOnce(err)
		}
	} else if qt == queue.QueueTypeAtMostOnce {
		if err := s.conn.Publish(topic, m); err != nil {
			return ErrJsPublishAtMostOnce(err)
		}
	} else {
		return ErrJsQtNotSupported(int(qt))
	}
	return nil
}

func (s *jsImpl) Subscribe(qt queue.QueueType, topic, durableId string, receiverChan chan<- []byte) error {
	l := s.l().Mth("received").F(log.FF{"topic": topic, "type": qt.String()})
	s.connMtx.RLock()
	defer s.connMtx.RUnlock()
	if qt == queue.QueueTypeAtLeastOnce {
		// add stream if not exists
		_, err := s.js.AddStream(&nats.StreamConfig{
			Name:     s.topicToStream(topic),
			Subjects: []string{topic},
		})
		if err != nil {
			return ErrJsSubscribeAtLeastOnce(err)
		}
		// subscribe
		_, err = s.js.Subscribe(topic, func(m *nats.Msg) {
			l.TrcF("%s\n", string(m.Data))
			defer func() { _ = m.Ack() }()
			receiverChan <- m.Data
		}, nats.Durable(durableId), nats.ManualAck())
		if err != nil {
			return ErrJsSubscribeAtLeastOnce(err)
		}
		return nil
	} else if qt == queue.QueueTypeAtMostOnce {
		_, err := s.conn.Subscribe(topic, func(m *nats.Msg) {
			l.TrcF("%s\n", string(m.Data))
			receiverChan <- m.Data
		})
		if err != nil {
			return ErrJsSubscribeAtMostOnce(err)
		}
		return nil
	} else {
		return ErrJsQtNotSupported(int(qt))
	}
}

func (s *jsImpl) SubscribeLB(qt queue.QueueType, topic, loadBalancingGroup, durableId string, receiverChan chan<- []byte) error {
	l := s.l().Mth("received").F(log.FF{"topic": topic, "type": qt.String(), "lbGrp": loadBalancingGroup})
	s.connMtx.RLock()
	defer s.connMtx.RUnlock()
	if qt == queue.QueueTypeAtLeastOnce {
		// add stream if not exists
		_, err := s.js.AddStream(&nats.StreamConfig{
			Name:     s.topicToStream(topic),
			Subjects: []string{topic},
		})
		if err != nil {
			return ErrJsSubscribeAtLeastOnce(err)
		}
		// subscribe
		_, err = s.js.QueueSubscribe(topic, loadBalancingGroup, func(m *nats.Msg) {
			l.TrcF("%s\n", string(m.Data))
			defer func() { _ = m.Ack() }()
			receiverChan <- m.Data
		}, nats.Durable(durableId), nats.ManualAck())
		if err != nil {
			return ErrJsSubscribeAtLeastOnce(err)
		}
		return nil
	} else if qt == queue.QueueTypeAtMostOnce {
		_, err := s.conn.QueueSubscribe(topic, loadBalancingGroup, func(m *nats.Msg) {
			l.TrcF("%s\n", string(m.Data))
			receiverChan <- m.Data
		})
		if err != nil {
			return ErrJsSubscribeAtMostOnce(err)
		}
		return nil
	} else {
		return ErrJsQtNotSupported(int(qt))
	}
}
