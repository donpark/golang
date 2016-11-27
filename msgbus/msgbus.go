package msgbus

import (
	"reflect"
	"sync"
)

type MsgBusCallback func(msg interface{})

// MsgBus

type MsgBus struct {
	sync.RWMutex
	topics map[string]*msgTopicHandler
}

func New() *MsgBus {
	bus := &MsgBus{
		topics: make(map[string]*msgTopicHandler),
	}
	return bus
}

func (bus *MsgBus) Publish(topic string, arg interface{}) {
	bus.RLock()
	h, exists := bus.topics[topic]
	bus.RUnlock()
	if !exists {
		return
	}
	h.publish(arg)
}

func (bus *MsgBus) Subscribe(topic string, fn MsgBusCallback) {
	bus.RLock()
	h, exists := bus.topics[topic]
	bus.RUnlock()
	if !exists {
		bus.Lock()
		h, exists = bus.topics[topic]
		if !exists {
			h = newMsgTopicHandler()
			bus.topics[topic] = h
		}
		bus.Unlock()
	}
	h.subscribe(fn)
}

func (bus *MsgBus) Unsubscribe(topic string, fn MsgBusCallback) {
	bus.RLock()
	h, exists := bus.topics[topic]
	bus.RUnlock()
	if !exists {
		return
	}
	h.unsubscribe(fn)
	if h.hasSubscribers() {
		return
	}
	bus.Lock()
	defer bus.Unlock()
	if h.hasSubscribers() {
		return
	}
	delete(bus.topics, topic)
	h.stop <- true
}

func (bus *MsgBus) UnsubscribeAll() {
	bus.Lock()
	defer bus.Unlock()
	for topic, h := range bus.topics {
		h.unsubscribeAll()
		h.stop <- true
		delete(bus.topics, topic)
	}
}

// msgTopicHandler

type msgTopicHandler struct {
	sync.RWMutex
	msgs      chan interface{}
	stop      chan bool
	callbacks []reflect.Value
}

func newMsgTopicHandler() *msgTopicHandler {
	h := &msgTopicHandler{
		msgs:      make(chan interface{}, 10),
		stop:      make(chan bool),
		callbacks: make([]reflect.Value, 0),
	}
	go h.dispatchMsgs()
	return h
}

func (h *msgTopicHandler) dispatchMsgs() {
	for {
		select {
		case msg := <-h.msgs:
			h.deliver(msg)
		case <-h.stop:
			return
		}
	}
}

func (h *msgTopicHandler) hasSubscribers() bool {
	h.RLock()
	defer h.RUnlock()
	return len(h.callbacks) > 0
}

func (h *msgTopicHandler) publish(msg interface{}) {
	h.msgs <- msg
}

func (h *msgTopicHandler) deliver(msg interface{}) {
	h.RLock()
	defer h.RUnlock()
	args := []reflect.Value{reflect.ValueOf(msg)}
	for _, callback := range h.callbacks {
		callback.Call(args)
	}
}

func (h *msgTopicHandler) subscribe(fn MsgBusCallback) {
	h.Lock()
	defer h.Unlock()
	fnv := h.getCallbackValue(fn)
	if index := h.indexOfCallbackValue(fnv); index < 0 {
		h.callbacks = append(h.callbacks, fnv)
	}
}

func (h *msgTopicHandler) unsubscribe(fn MsgBusCallback) {
	h.Lock()
	defer h.Unlock()
	fnv := h.getCallbackValue(fn)
	if i := h.indexOfCallbackValue(fnv); i >= 0 {
		h.callbacks = append(h.callbacks[:i], h.callbacks[i+1:]...)
	}
}

func (h *msgTopicHandler) unsubscribeAll() {
	h.Lock()
	defer h.Unlock()
	h.callbacks = make([]reflect.Value, 0)
}

func (h *msgTopicHandler) getCallbackValue(fn MsgBusCallback) reflect.Value {
	if fn == nil || !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		panic("fn is not a function")
	}
	return reflect.ValueOf(fn)
}

func (h *msgTopicHandler) indexOfCallbackValue(fnv reflect.Value) int {
	for index, callback := range h.callbacks {
		if callback == fnv {
			return index
		}
	}
	return -1
}
