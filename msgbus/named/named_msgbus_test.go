package named

import (
	"testing"
)

func TestGetMsgBus(t *testing.T) {
	if bus := GetNamedMsgBus("bus/test"); bus == nil {
		t.Error("GetMsgBus failed")
	}
	registry.DeleteAll()
}

func TestSameMsgBus(t *testing.T) {
	bus1 := GetNamedMsgBus("bus/test")
	bus2 := GetNamedMsgBus("bus/test")
	if bus1 != bus2 {
		t.Error("GetMsgBus don't match")
	}
	registry.DeleteAll()
}

func TestOneToOne(t *testing.T) {
	bus := GetNamedMsgBus("bus/test")
	stop := make(chan bool)
	bus.Subscribe("foo", func(e interface{}) {
		switch e := e.(type) {
		case string:
			// fmt.Println("string msg received", e)
		default:
			t.Error("unexpected msg received", e)
		}
		stop <- true
	})

	bus.Publish("foo", "test msg")

	for {
		select {
		case <-stop:
			registry.DeleteAll()
			return
		}
	}
}
