//+build example

package listener

import (
	"context"
	"github.com/africarealty/server/src/kit"
	kitContext "github.com/africarealty/server/src/kit/context"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/kit/queue"
	"github.com/africarealty/server/src/kit/queue/stan"
	"sync"
	"testing"
)

var logger = log.Init(&log.Config{Level: log.TraceLevel})
var logf = func() log.CLogger {
	return log.L(logger)
}

func Test_PublishConsume_AtMostOnce_WithTwoListeners(t *testing.T) {

	ctxRq := kitContext.NewRequestCtx().WithNewRequestId().Test()
	ctx := ctxRq.ToContext(context.Background())
	topic := kit.NewRandString()

	wg := sync.WaitGroup{}

	handler := func(id string) QueueMessageHandler {
		return func(payload []byte) error {
			logf().DbgF("[%s] %s", id, string(payload))
			wg.Done()
			return nil
		}
	}

	pub := stan.New(logf)
	err := pub.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer pub.Close()

	sub1 := stan.New(logf)
	err = sub1.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub1.Close()
	ql1 := NewQueueListener(sub1, logf)
	ql1.New(topic).AtMostOnce().WithHandler(handler("1")).Add()
	ql1.ListenAsync()

	sub2 := stan.New(logf)
	err = sub2.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub2.Close()
	ql2 := NewQueueListener(sub2, logf)
	ql2.New(topic).AtMostOnce().WithHandler(handler("2")).Add()
	ql2.ListenAsync()

	wg.Add(2)

	err = pub.Publish(ctx, queue.QueueTypeAtMostOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg"}})
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

func Test_PublishConsume_AtMostOnce_WithTwoListenersAndTwoHandlers(t *testing.T) {

	ctxRq := kitContext.NewRequestCtx().WithNewRequestId().Test()
	ctx := ctxRq.ToContext(context.Background())
	topic := kit.NewRandString()

	wg := sync.WaitGroup{}

	handler := func(id string) QueueMessageHandler {
		return func(payload []byte) error {
			logf().DbgF("[%s] %s", id, string(payload))
			wg.Done()
			return nil
		}
	}

	pub := stan.New(logf)
	err := pub.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer pub.Close()

	sub1 := stan.New(logf)
	err = sub1.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub1.Close()
	ql1 := NewQueueListener(sub1, logf)
	ql1.New(topic).AtMostOnce().WithHandler(handler("11")).WithHandler(handler("12")).Add()
	ql1.ListenAsync()

	sub2 := stan.New(logf)
	err = sub2.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub2.Close()
	ql2 := NewQueueListener(sub2, logf)
	ql2.New(topic).AtMostOnce().WithHandler(handler("21")).WithHandler(handler("22")).Add()
	ql2.ListenAsync()

	wg.Add(4)

	err = pub.Publish(ctx, queue.QueueTypeAtMostOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg"}})
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

func Test_PublishConsume_AtLeastOnce_WithLB_WithTwoListeners(t *testing.T) {

	ctxRq := kitContext.NewRequestCtx().WithNewRequestId().Test()
	ctx := ctxRq.ToContext(context.Background())
	topic := kit.NewRandString()
	durableId := kit.NewRandString()

	wg := sync.WaitGroup{}

	handler := func(id string) QueueMessageHandler {
		return func(payload []byte) error {
			logf().DbgF("[%s] %s", id, string(payload))
			wg.Done()
			return nil
		}
	}

	pub := stan.New(logf)
	err := pub.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer pub.Close()

	sub1 := stan.New(logf)
	err = sub1.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	ql1 := NewQueueListener(sub1, logf)
	ql1.New(topic).AtLeastOnce(durableId).WithLoadBalancing("group").WithHandler(handler("1")).Add()
	ql1.ListenAsync()

	sub2 := stan.New(logf)
	err = sub2.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	ql2 := NewQueueListener(sub2, logf)
	ql2.New(topic).AtLeastOnce(durableId).WithLoadBalancing("group").WithHandler(handler("2")).Add()
	ql2.ListenAsync()

	wg.Add(2)

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg1"}})
	if err != nil {
		t.Fatal(err)
	}

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg2"}})
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()

	ql1.Stop()
	ql2.Stop()
	_ = sub1.Close()
	_ = sub2.Close()

	wg.Add(3)

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg3"}})
	if err != nil {
		t.Fatal(err)
	}

	sub1 = stan.New(logf)
	err = sub1.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	ql1 = NewQueueListener(sub1, logf)
	ql1.New(topic).AtLeastOnce(durableId).WithLoadBalancing("group").WithHandler(handler("1")).Add()
	ql1.ListenAsync()

	sub2 = stan.New(logf)
	err = sub2.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	ql2 = NewQueueListener(sub2, logf)
	ql2.New(topic).AtLeastOnce(durableId).WithLoadBalancing("group").WithHandler(handler("2")).Add()
	ql2.ListenAsync()

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg4"}})
	if err != nil {
		t.Fatal(err)
	}

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg4"}})
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()

}

func Test_PublishConsume_AtLeastOnce_WithoutLB_WithTwoListeners(t *testing.T) {

	ctxRq := kitContext.NewRequestCtx().WithNewRequestId().Test()
	ctx := ctxRq.ToContext(context.Background())
	topic := kit.NewRandString()

	wg := sync.WaitGroup{}

	handler := func(id string) QueueMessageHandler {
		return func(payload []byte) error {
			logf().DbgF("[%s] %s", id, string(payload))
			wg.Done()
			return nil
		}
	}

	pub := stan.New(logf)
	err := pub.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer pub.Close()

	sub1 := stan.New(logf)
	subClientId1 := kit.NewRandString()
	err = sub1.Open(ctx, subClientId1, &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	ql1 := NewQueueListener(sub1, logf)
	// durableId must be unique
	durableId1 := kit.NewRandString()
	ql1.New(topic).AtLeastOnce(durableId1).WithHandler(handler("1")).Add()
	ql1.ListenAsync()

	sub2 := stan.New(logf)
	subClientId2 := kit.NewRandString()
	err = sub2.Open(ctx, subClientId2, &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	ql2 := NewQueueListener(sub2, logf)
	// durableId must be unique
	durableId2 := kit.NewRandString()
	ql2.New(topic).AtLeastOnce(durableId2).WithLoadBalancing("group").WithHandler(handler("2")).Add()
	ql2.ListenAsync()

	wg.Add(4)

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg1"}})
	if err != nil {
		t.Fatal(err)
	}

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg2"}})
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()

	// check durability: close subscribers, publish messages and then open subscriber with the same durableId
	ql1.Stop()
	ql2.Stop()
	_ = sub1.Close()
	_ = sub2.Close()

	wg.Add(1)

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg3"}})
	if err != nil {
		t.Fatal(err)
	}

	sub1 = stan.New(logf)
	err = sub1.Open(ctx, subClientId1, &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	ql1 = NewQueueListener(sub1, logf)
	// durableId must be unique
	ql1.New(topic).AtLeastOnce(durableId1).WithHandler(handler("1")).Add()
	ql1.ListenAsync()

	wg.Wait()
}

func Test_PublishConsume_AtLeastOnce_WithLB_SameGroupDifferentTopic(t *testing.T) {

	ctxRq := kitContext.NewRequestCtx().WithNewRequestId().Test()
	ctx := ctxRq.ToContext(context.Background())
	topic1 := kit.NewRandString()
	topic2 := kit.NewRandString()
	durableId := kit.NewRandString()

	wg := sync.WaitGroup{}

	handler := func(id string) QueueMessageHandler {
		return func(payload []byte) error {
			logf().DbgF("[%s] %s", id, string(payload))
			wg.Done()
			return nil
		}
	}

	pub := stan.New(logf)
	err := pub.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer pub.Close()

	sub1 := stan.New(logf)
	err = sub1.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	defer sub1.Close()
	if err != nil {
		t.Fatal(err)
	}
	ql1 := NewQueueListener(sub1, logf)
	ql1.New(topic1).AtLeastOnce(durableId).WithLoadBalancing("group").WithHandler(handler("1")).Add()
	ql1.ListenAsync()

	sub2 := stan.New(logf)
	err = sub2.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	defer sub2.Close()
	if err != nil {
		t.Fatal(err)
	}
	ql2 := NewQueueListener(sub2, logf)
	ql2.New(topic1).AtLeastOnce(durableId).WithLoadBalancing("group").WithHandler(handler("2")).Add()
	ql2.ListenAsync()

	sub3 := stan.New(logf)
	err = sub3.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	defer sub3.Close()
	if err != nil {
		t.Fatal(err)
	}
	ql3 := NewQueueListener(sub3, logf)
	ql3.New(topic2).AtLeastOnce(durableId).WithLoadBalancing("group").WithHandler(handler("3")).Add()
	ql3.ListenAsync()

	sub4 := stan.New(logf)
	err = sub4.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	ql4 := NewQueueListener(sub4, logf)
	ql4.New(topic2).AtLeastOnce(durableId).WithLoadBalancing("group").WithHandler(handler("4")).Add()
	ql4.ListenAsync()

	wg.Add(4)

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic1, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg11"}})
	if err != nil {
		t.Fatal(err)
	}

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic2, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg21"}})
	if err != nil {
		t.Fatal(err)
	}

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic1, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg12"}})
	if err != nil {
		t.Fatal(err)
	}

	err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic2, &queue.Message{Payload: struct {
		Val string
	}{Val: "msg22"}})
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

func Test_PublishConsume_AtLeastOnce_WithLB_Multiple(t *testing.T) {

	ctxRq := kitContext.NewRequestCtx().WithNewRequestId().Test()
	ctx := ctxRq.ToContext(context.Background())
	topic1 := "test.topic"

	wg := sync.WaitGroup{}
	wg.Add(100)
	handler := func(id string) QueueMessageHandler {
		return func(payload []byte) error {
			wg.Done()
			logf().DbgF("[%s] %s", id, string(payload))
			return nil
		}
	}

	pub := stan.New(logf)
	err := pub.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer pub.Close()

	sub1 := stan.New(logf)
	err = sub1.Open(ctx, kit.NewRandString(), &queue.Config{
		Host:      "localhost",
		Port:      "4222",
		ClusterId: "test-cluster",
	})
	defer sub1.Close()
	if err != nil {
		t.Fatal(err)
	}
	ql1 := NewQueueListener(sub1, logf)
	ql1.New(topic1).AtLeastOnce("durable").WithLoadBalancing("group").WithHandler(handler("1")).Add()
	ql1.ListenAsync()

	for i := 0; i < 100; i++ {
		err = pub.Publish(ctx, queue.QueueTypeAtLeastOnce, topic1, &queue.Message{Payload: struct {
			Val string
		}{Val: "msg11"}})
		if err != nil {
			t.Fatal(err)
		}
	}

	wg.Wait()
}
