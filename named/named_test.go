package named

import (
	"fmt"
	"testing"
)

func TestEmptyName(t *testing.T) {
	r := NewRegistry()
	if c := r.Get("", nil); c != nil {
		t.Error("null named should fail", c)
	}
}

func TestNonExistantName(t *testing.T) {
	r := NewRegistry()
	if c := r.Get("test/string/should/not/exist", nil); c != nil {
		t.Error("non-existant should fail", c)
	}
}

func TestCreateString(t *testing.T) {
	r := NewRegistry()
	c := r.Get("test/string", func() interface{} {
		return "hello world"
	})
	if c == nil {
		t.Error("expected named string")
	}
	switch vt := c.(type) {
	case string: // expected type
	default:
		t.Error("unexpected value type", vt)
	}
}

func TestDeleteString(t *testing.T) {
	r := NewRegistry()
	c := r.Get("test/string", func() interface{} {
		return "hello world"
	})
	if c == nil {
		t.Error("expected named string")
	}
	r.Delete("test/string")
	if r.Get("test/string", nil) != nil {
		t.Error("delete failed")
	}
}

func TestIntChannel(t *testing.T) {
	r := NewRegistry()
	c := r.Get("test/channel", func() interface{} {
		return make(chan int)
	})
	if c == nil {
		t.Error("expected named channel")
	}
	ch := c.(chan int)
	go func() {
		select {
		case v := <-ch:
			fmt.Println("got int", v)
		default:
			t.Error("unexpected select")
		}
	}()

	ch <- 5
}

func TestAnonymousChannel(t *testing.T) {
	r := NewRegistry()
	c := r.Get("test/channel", func() interface{} {
		return make(chan interface{})
	})
	if c == nil {
		t.Error("expected named channel")
	}
	ch := c.(chan interface{})
	go func() {
		select {
		case v := <-ch:
			switch vt := v.(type) {
			case int: // expected type
				fmt.Println("got value", vt)
			default:
				t.Error("unexpected value type", vt)
			}
		default:
			t.Error("unexpected select")
		}
	}()

	ch <- 5
}

func TestDeleteAll(t *testing.T) {
	r := NewRegistry()
	c1 := r.Get("test/string", func() interface{} {
		return "hello world"
	})
	if c1 == nil {
		t.Error("expected named string")
	}
	if r.Get("test/string", nil) != c1 {
		t.Error("string value don't match")
	}

	c2 := r.Get("test/channel", func() interface{} {
		return make(chan string, 1)
	})
	if c2 == nil {
		t.Error("expected named channel")
	}
	if r.Get("test/channel", nil) != c2 {
		t.Error("channel value don't match")
	}

	r.DeleteAll()
	if r.Get("test/string", nil) != nil {
		t.Error("delete all failed")
	}
	if r.Get("test/channel", nil) != nil {
		t.Error("delete all failed")
	}
}
