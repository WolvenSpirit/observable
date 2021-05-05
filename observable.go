package observable

import (
	"reflect"
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
	subscribers map[int]*chan interface{}
	m           sync.Mutex
}

type Observable struct {
	observable
}

// Next - pass a new value to be broadcasted.
func (o *Observable) Next(sl ...interface{}) {
	o.m.Lock()
	for _, v := range sl {
		for k := range o.subscribers {
			*o.subscribers[k] <- v

		}
	}
	o.m.Unlock()
}
func (o *Observable) allocate(id int) {
	// ...
}

// Subscribe with your own channel to the events of the observable.
func (o *Observable) Subscribe(ch *chan interface{}) {
	o.m.Lock()
	if o.subscribers == nil {
		o.subscribers = make(map[int]*chan interface{})
	}
	o.subscribers[len(o.subscribers)+1] = ch
	o.allocate(len(o.subscribers) + 1)
	o.m.Unlock()
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
	for _, ch := range o.subscribers {
		close(*ch)
	}
}
