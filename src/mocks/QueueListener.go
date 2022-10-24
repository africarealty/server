// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	listener "github.com/africarealty/server/src/kit/queue/listener"
	mock "github.com/stretchr/testify/mock"

	queue "github.com/africarealty/server/src/kit/queue"
)

// QueueListener is an autogenerated mock type for the QueueListener type
type QueueListener struct {
	mock.Mock
}

// Add provides a mock function with given fields: qt, topic, durableId, h
func (_m *QueueListener) Add(qt queue.QueueType, topic string, durableId string, h ...listener.QueueMessageHandler) {
	_va := make([]interface{}, len(h))
	for _i := range h {
		_va[_i] = h[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, qt, topic, durableId)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// AddLb provides a mock function with given fields: qt, topic, lbGroup, durableId, h
func (_m *QueueListener) AddLb(qt queue.QueueType, topic string, lbGroup string, durableId string, h ...listener.QueueMessageHandler) {
	_va := make([]interface{}, len(h))
	for _i := range h {
		_va[_i] = h[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, qt, topic, lbGroup, durableId)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Clear provides a mock function with given fields:
func (_m *QueueListener) Clear() {
	_m.Called()
}

// ListenAsync provides a mock function with given fields:
func (_m *QueueListener) ListenAsync() {
	_m.Called()
}

// New provides a mock function with given fields: topic
func (_m *QueueListener) New(topic string) listener.Builder {
	ret := _m.Called(topic)

	var r0 listener.Builder
	if rf, ok := ret.Get(0).(func(string) listener.Builder); ok {
		r0 = rf(topic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(listener.Builder)
		}
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *QueueListener) Stop() {
	_m.Called()
}