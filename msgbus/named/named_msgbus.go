package named

import (
	"github.com/donpark/golang/msgbus"
	"github.com/donpark/golang/named"
)

var (
	registry = named.NewRegistry()
)

func GetNamedMsgBus(name string) *msgbus.MsgBus {
	v := registry.Get(name, func() interface{} {
		return msgbus.New()
	})
	if v == nil {
		panic("msgbus.named: GetNamedMsgBus failed unexpectedly")
	}
	switch v := v.(type) {
	case *msgbus.MsgBus:
		return v
	default:
		panic("msgbus.named: invalid registered type for name: " + name)
	}
}
