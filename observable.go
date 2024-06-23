package observable

import (
	"errors"
	"reflect"
)

type observable struct {
	subscribers map[int]*chan any
	i           *int
}

// Observable provides methods and hooks for transmitting data.
type Observable struct {
	observable
}

// New returns a pointer to a new Observable
func New() *Observable {
	var i int
	return &Observable{observable: observable{i: &i, subscribers: make(map[int]*chan any)}}
}

// Next - pass a new value to be broadcasted.
func (o *Observable) Next(sl ...interface{}) {
	for k := range sl {
		/*o.subscribers.Range(func(key interface{}, value interface{}) bool {
			*value.(*chan interface{}) <- sl[k]
			return true
		})*/
		for i := 0; i < *o.i; i++ {
			if ch, ok := o.subscribers[i]; ok {
				*ch <- sl[k]
			}
		}
	}
}

// Subscribe with your own channel to the events of the observable. The function returns an id that can be used to unsubscribe.
func (o *Observable) Subscribe(ch *chan interface{}) (int, error) {
	if ch == nil {
		return 0, errors.New("argument to subscribe is nil not chan")
	}
	k := *o.i
	*o.i++
	o.subscribers[k] = ch
	return k, nil
}

// On provides a hook to execute a callback when a certain value is passed through the channel.
func (o *Observable) On(value interface{}, ch *chan interface{}, fn func()) {
	inf := true
	for inf {
		select {
		case v := <-*ch:
			if reflect.DeepEqual(v, value) {
				fn()
			}
		default:

		}
	}
}

// Close effectivelly calls close() on all channels used by subscribers to this observable.
func (o *Observable) Close() {
	for i := 0; i < *o.i; i++ {
		ch, ok := o.subscribers[i]
		if ok {
			close(*ch)
		}
	}
}

// Unsubscribe removes the channel from the list of active subscribers, does not call close on them.
func (o *Observable) Unsubscribe(id int) (ok bool) {
	delete(o.subscribers, id)
	return true
}
