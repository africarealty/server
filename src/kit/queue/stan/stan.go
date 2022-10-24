package stan

import (
	"context"
	"encoding/json"
	"fmt"
	kitContext "github.com/africarealty/server/src/kit/context"
	"github.com/africarealty/server/src/kit/goroutine"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/kit/queue"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"sync"
	"time"
)

const (
	PingIntervalSec        = 2
	MaxOutPingsCount       = 5
	WaitBeforeReconnectSec = 3
	ReconnectIntervalSec   = 5
)

type stanImpl struct {
	conn                   stan.Conn
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
	return &stanImpl{
		logger: logger,
	}
}

func (s *stanImpl) l() log.CLogger {
	return s.logger().Pr("queue").Cmp("stan")
}

func (s *stanImpl) SetLostConnectionHandler(handler queue.EventHandler) {
	s.handlerMtx.Lock()
	defer s.handlerMtx.Unlock()
	s.lostConnectionHandlers = append(s.lostConnectionHandlers, handler)
}

func (s *stanImpl) SetReconnectHandler(handler queue.EventHandler) {
	s.handlerMtx.Lock()
	defer s.handlerMtx.Unlock()
	s.reconnectHandlers = append(s.reconnectHandlers, handler)
}

func (s *stanImpl) Open(ctx context.Context, clientId string, config *queue.Config) error {
	l := s.l().Mth("open").F(log.FF{"client": clientId, "host": config.Host}).Dbg("connecting")

	s.connMtx.Lock()
	defer s.connMtx.Unlock()

	s.ctx = ctx
	s.clientId = clientId
	s.config = config

	url := fmt.Sprintf("nats://%s:%s", config.Host, config.Port)
	c, err := stan.Connect(config.ClusterId, clientId, stan.NatsURL(url),
		stan.Pings(PingIntervalSec, MaxOutPingsCount),
		stan.SetConnectionLostHandler(s.StanConnectionLostHandler))
	if err != nil {
		return ErrStanConnect(err)
	}
	s.conn = c

	l.Inf("ok")

	return nil
}

func (s *stanImpl) lostConnectionNotify() {
	s.handlerMtx.Lock()
	defer s.handlerMtx.Unlock()
	for _, handler := range s.lostConnectionHandlers {
		goroutine.New().
			WithLogger(s.l().Mth("lost-connection-notify")).
			Go(context.Background(), handler)
	}
}

func (s *stanImpl) reconnectNotify() {
	s.handlerMtx.Lock()
	defer s.handlerMtx.Unlock()
	for _, handler := range s.reconnectHandlers {
		goroutine.New().
			WithLogger(s.l().Mth("reconnect-notify")).
			Go(context.Background(), handler)
	}
}

func (s *stanImpl) reconnect() {
	s.l().Mth("reconnect").Dbg("started")
	// wait a bit before reconnecting
	time.Sleep(time.Second * WaitBeforeReconnectSec)
	// creat a ticker
	ticker := time.NewTicker(time.Second * ReconnectIntervalSec)
	defer ticker.Stop()
	// start the process
	for range ticker.C {
		if err := s.Open(s.ctx, s.clientId, s.config); err != nil {
			s.l().E(err).Err("reconnect error")
		} else {
			s.l().Dbg("reconnected")
			// notify subscribers
			s.reconnectNotify()
			return
		}
	}
}

func (s *stanImpl) StanConnectionLostHandler(connection stan.Conn, err error) {
	s.l().Mth("connection-lost-handler").E(err).Err("connection lost")
	// close current connection
	_ = s.Close()
	// notify subscribers
	s.lostConnectionNotify()
	// start reconnect process
	goroutine.New().
		WithLogger(s.l().Mth("reconnect")).
		Go(context.Background(), s.reconnect)
}

func (s *stanImpl) Declare(ctx context.Context, qt queue.QueueType, topic string) error {
	s.l().Mth("declare").F(log.FF{"topic": topic, "type": qt.String()}).Dbg()
	return nil
}

func (s *stanImpl) Close() error {
	s.connMtx.Lock()
	defer s.connMtx.Unlock()
	if s.conn != nil {
		err := s.conn.Close()
		s.conn = nil
		if err != nil {
			return ErrStanClose(err)
		}
		s.l().Mth("close").Inf("closed")
	}
	return nil
}

func (s *stanImpl) Publish(ctx context.Context, qt queue.QueueType, topic string, msg *queue.Message) error {
	l := s.l().Mth("publish").F(log.FF{"topic": topic, "type": qt.String()})

	s.connMtx.RLock()
	defer s.connMtx.RUnlock()

	if msg.Ctx == nil {
		msg.Ctx = kitContext.NewRequestCtx().Queue().WithNewRequestId()
	}
	l.C(msg.Ctx.ToContext(context.Background()))

	if s.conn == nil {
		return ErrStanNoOpenConn()
	}

	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	l.Dbg("ok").TrcF("%s\n", string(m))

	if qt == queue.QueueTypeAtLeastOnce {
		if err := s.conn.Publish(topic, m); err != nil {
			return ErrStanPublishAtLeastOnce(err)
		}
	} else if qt == queue.QueueTypeAtMostOnce {
		if err := s.conn.NatsConn().Publish(topic, m); err != nil {
			return ErrStanPublishAtMostOnce(err)
		}
	} else {
		return ErrStanQtNotSupported(int(qt))
	}
	return nil

}

func (s *stanImpl) Subscribe(qt queue.QueueType, topic, durableId string, receiverChan chan<- []byte) error {

	l := s.l().Mth("received").F(log.FF{"topic": topic, "type": qt.String()})

	s.connMtx.RLock()
	defer s.connMtx.RUnlock()

	if qt == queue.QueueTypeAtLeastOnce {

		_, err := s.conn.Subscribe(topic, func(m *stan.Msg) {
			l.TrcF("%s\n", string(m.Data))
			receiverChan <- m.Data
		}, stan.DurableName(durableId))
		if err != nil {
			return ErrStanSubscribeAtLeastOnce(err)
		}
		return nil

	} else if qt == queue.QueueTypeAtMostOnce {
		_, err := s.conn.NatsConn().Subscribe(topic, func(m *nats.Msg) {
			l.TrcF("%s\n", string(m.Data))
			receiverChan <- m.Data
		})
		if err != nil {
			return ErrStanSubscribeAtMostOnce(err)
		}
		return nil
	} else {
		return ErrStanQtNotSupported(int(qt))
	}

}

func (s *stanImpl) SubscribeLB(qt queue.QueueType, topic, loadBalancingGroup, durableId string, receiverChan chan<- []byte) error {
	l := s.l().Mth("received").F(log.FF{"topic": topic, "type": qt.String(), "lbGrp": loadBalancingGroup})

	s.connMtx.RLock()
	defer s.connMtx.RUnlock()

	if qt == queue.QueueTypeAtLeastOnce {
		_, err := s.conn.QueueSubscribe(topic, loadBalancingGroup, func(m *stan.Msg) {
			l.TrcF("%s\n", string(m.Data))
			receiverChan <- m.Data
			// https://github.com/nats-io/stan.go
		}, stan.DurableName(durableId))
		if err != nil {
			return ErrStanSubscribeAtLeastOnce(err)
		}
		return nil
	} else if qt == queue.QueueTypeAtMostOnce {
		_, err := s.conn.NatsConn().QueueSubscribe(topic, loadBalancingGroup, func(m *nats.Msg) {
			l.TrcF("%s\n", string(m.Data))
			receiverChan <- m.Data
		})
		if err != nil {
			return ErrStanSubscribeAtMostOnce(err)
		}
		return nil
	} else {
		return ErrStanQtNotSupported(int(qt))
	}
}
