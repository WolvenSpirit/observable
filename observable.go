package observable

import (
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
)

/*

 */
type observableInterface interface {
	Next(interface{})
	Subscribe(func(*chan interface{}))
}

type observable struct {
	buffer      chan interface{}
	subscribers sync.Map
	i           *int
}

// Observable provides methods and hooks for transmitting data.
type Observable struct {
	observable
}

// New returns a pointer to a new Observable
func New() *Observable {
	var i int
	return &Observable{observable: observable{i: &i}}
}

// Next - pass a new value to be broadcasted.
func (o *Observable) Next(sl ...interface{}) {
	for k := range sl {
		/*o.subscribers.Range(func(key interface{}, value interface{}) bool {
			*value.(*chan interface{}) <- sl[k]
			return true
		})*/
		for i := 0; i < *o.i; i++ {
			if ch, ok := o.subscribers.Load(i); ok {
				*ch.(*chan interface{}) <- sl[k]
			}
		}
	}
}
func (o *Observable) allocate(id int) {
	io.Copy(os.Stderr, strings.NewReader("Unimplemented"))
}

// Subscribe with your own channel to the events of the observable.
func (o *Observable) Subscribe(ch *chan interface{}) error {
	if ch == nil {
		return errors.New("Argument to subscribe is nil not chan")
	}
	k := *o.i
	*o.i++
	o.subscribers.Store(k, ch)
	return nil
}

// On provides a hook to execute a callback when a certain value is passed through the channel.
func (o *Observable) On(value interface{}, ch *chan interface{}, fn func()) {
	for true {
		select {
		case v := <-*ch:
			if reflect.DeepEqual(v, value) {
				fn()
			}
		}
	}
}

// Close effectivelly calls close() on all channels used by subscribers to this observable.
func (o *Observable) Close() {
	for i := 0; i < *o.i; i++ {
		ch, ok := o.subscribers.Load(i)
		if ok {
			close(*ch.(*chan interface{}))
		}
	}
}
