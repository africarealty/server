//build integration

package jetstream

import (
	"context"
	"fmt"
	"github.com/africarealty/server/src/kit"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/kit/queue"
	"sync"
	"testing"
	"time"
)

var logger = log.Init(&log.Config{Level: log.TraceLevel})
var logf = func() log.CLogger {
	return log.L(logger)
}

var cfg = &queue.Config{
	Host: "localhost",
	Port: "14222",
}

func Test_Js_AtLeastOnce(t *testing.T) {

	ctx := context.Background()
	clientId := kit.NewRandString()
	topic := kit.NewRandString()

	jsQ := New(logf)
	err := jsQ.Open(ctx, clientId, cfg)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	recChan := make(chan []byte)
	go func() {
		defer close(recChan)
		for {
			select {
			case msg := <-recChan:
				logf().Dbg("received:" + string(msg))
				wg.Done()
				return
			case <-time.After(time.Second * 2):
				logf().Err("timeout")
				return
			}
		}
	}()
	err = jsQ.Subscribe(queue.QueueTypeAtLeastOnce, topic, clientId, recChan)
	if err != nil {
		t.Fatal(err)
	}

	err = jsQ.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: "test"})
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()

	err = jsQ.Close()
	if err != nil {
		t.Fatal(err)
	}

}

func Test_Js_AtLeastOnce_LoadBalancing(t *testing.T) {

	ctx := context.Background()
	clientId := kit.NewRandString()
	topic := kit.NewRandString()
	group := "group"

	jsSub1 := New(logf)
	err := jsSub1.Open(ctx, clientId+"1", cfg)
	if err != nil {
		t.Fatal(err)
	}

	jsSub2 := New(logf)
	err = jsSub2.Open(ctx, clientId+"2", cfg)
	if err != nil {
		t.Fatal(err)
	}

	jsPub := New(logf)
	err = jsPub.Open(ctx, clientId, cfg)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(20)

	ctxCancel, cancel := context.WithCancel(context.Background())

	rec1Chan := make(chan []byte)
	go func() {
		for {
			select {
			case msg := <-rec1Chan:
				logf().Dbg("received 1:" + string(msg))
				wg.Done()
			case <-ctxCancel.Done():
				return
			case <-time.After(time.Second * 2):
				logf().Err("timeout")
				return
			}
		}
	}()
	err = jsSub1.SubscribeLB(queue.QueueTypeAtLeastOnce, topic, group, clientId, rec1Chan)
	if err != nil {
		t.Fatal(err)
	}

	rec2Chan := make(chan []byte)
	go func() {
		for {
			select {
			case msg := <-rec2Chan:
				logf().Dbg("received 2:" + string(msg))
				wg.Done()
			case <-ctxCancel.Done():
				return
			case <-time.After(time.Second * 2):
				logf().Err("timeout")
				return
			}
		}
	}()
	err = jsSub2.SubscribeLB(queue.QueueTypeAtLeastOnce, topic, group, clientId, rec2Chan)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 20; i++ {
		err = jsPub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: fmt.Sprintf("%d", i)})
		if err != nil {
			t.Fatal(err)
		}
	}

	wg.Wait()
	cancel()

	_ = jsPub.Close()
	_ = jsSub1.Close()
	_ = jsSub2.Close()
}

//
//func initDurable(t *testing.T, topic string) {
//
//	sb, err := nats.Connect("test-cluster", "init")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer sb.Close()
//	s, err := sb.Subscribe(topic, func(m *nats.Msg) {
//		log.Printf("received: %v\n", m)
//	}, nats.DurableName("durable"))
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = s.Unsubscribe()
//	if err != nil {
//		t.Fatal(err)
//	}
//}
//
//func Test_Durable(t *testing.T) {
//
//	subj := "subj"
//	//initDurable(t, subj)
//
//	sendChan := make(chan string)
//	errChan := make(chan error)
//	quitPublisher := make(chan interface{})
//	quitSubscriber := make(chan interface{})
//
//	var wg sync.WaitGroup
//
//	go func() {
//		t.Fatal(<-errChan)
//	}()
//
//	// publisher
//	go func() {
//		pb, err := nats.Connect("test-cluster", "publisher")
//		if err != nil {
//			errChan <- err
//			return
//		}
//		defer pb.Close()
//		for {
//			select {
//			case msg := <-sendChan:
//				log.Printf("send: %s", msg)
//				err := pb.Publish(subj, []byte(msg))
//				if err != nil {
//					errChan <- err
//				}
//			case <-quitPublisher:
//				return
//			}
//		}
//
//	}()
//
//	// subscriber
//	subscriber := func() {
//		sb, err := nats.Connect("test-cluster", "subscriber")
//		if err != nil {
//			errChan <- err
//			return
//		}
//		defer sb.Close()
//		_, err = sb.Subscribe(subj, func(m *nats.Msg) {
//			log.Printf("received: %v\n", m)
//			wg.Done()
//		}, nats.DurableName("durable"))
//		if err != nil {
//			errChan <- err
//			return
//		}
//		//defer subscription.Unsubscribe()
//		<-quitSubscriber
//		log.Println("subscriber closed")
//	}
//
//	wg.Add(2)
//	sendChan <- "msg-1"
//	sendChan <- "msg-2"
//	log.Println("run subscriber")
//	go subscriber()
//	wg.Wait()
//	log.Println("sending quit subscriber")
//	quitSubscriber <- true
//	wg.Add(2)
//	sendChan <- "msg-3"
//	sendChan <- "msg-4"
//	log.Println("run subscriber")
//	go subscriber()
//
//	log.Println("waiting for wg")
//	wg.Wait()
//
//	log.Println("quit all")
//	quitSubscriber <- true
//	quitPublisher <- true
//
//}
//
//func Test_Js_Queue(t *testing.T) {
//
//	sendChan := make(chan string)
//	errChan := make(chan error)
//	quitPublisher := make(chan interface{})
//	quitSubscriber := make(chan interface{})
//
//	var wg sync.WaitGroup
//
//	go func() {
//		t.Fatal(<-errChan)
//	}()
//
//	// publisher
//	go func() {
//		pb, err := nats.Connect("test-cluster", "publisher")
//		if err != nil {
//			errChan <- err
//			return
//		}
//		defer pb.Close()
//		for {
//			select {
//			case msg := <-sendChan:
//				log.Printf("send: %s", msg)
//				err := pb.Publish("test-queue", []byte(msg))
//				if err != nil {
//					errChan <- err
//				}
//			case <-quitPublisher:
//				return
//			}
//		}
//
//	}()
//
//	// subscriber
//	subscriber := func(id string, quitChan <-chan interface{}) {
//		sb, err := nats.Connect("test-cluster", "subscriber"+id)
//		if err != nil {
//			errChan <- err
//			return
//		}
//		defer sb.Close()
//		_, err = sb.QueueSubscribe("test-queue", "queue", func(m *nats.Msg) {
//			log.Printf("received (%s): %v\n", id, m)
//			wg.Done()
//		}, nats.DurableName("durable"))
//		if err != nil {
//			errChan <- err
//			return
//		}
//		//defer subscription.Unsubscribe()
//		<-quitChan
//		log.Printf("subscriber %s closed", id)
//	}
//
//	log.Println("run subscribers")
//	go subscriber("1", quitSubscriber)
//	go subscriber("2", quitSubscriber)
//
//	time.Sleep(time.Second)
//
//	wg.Add(2)
//	sendChan <- "msg-1"
//	time.Sleep(time.Millisecond * 10)
//	sendChan <- "msg-2"
//
//	wg.Wait()
//	log.Println("sending quit subscriber")
//	quitSubscriber <- true
//	quitSubscriber <- true
//	wg.Add(2)
//	sendChan <- "msg-3"
//	time.Sleep(time.Millisecond * 10)
//	sendChan <- "msg-4"
//	log.Println("run subscriber")
//	go subscriber("3", quitSubscriber)
//	go subscriber("4", quitSubscriber)
//
//	log.Println("waiting for wg")
//	wg.Wait()
//
//	log.Println("quit all")
//	quitSubscriber <- true
//	quitSubscriber <- true
//	quitPublisher <- true
//
//}
