package observable

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestObservable_Next(t *testing.T) {
	ch := make(chan interface{}, 1)
	o := New()
	o.Subscribe(&ch)
	go func() {
		o.Next("3")
	}()
	v := <-ch
	if v.(string) != "3" {
		t.FailNow()
	}
}

func BenchmarkObservable_Next(t *testing.B) {
	o := New()
	subscription := make(chan interface{}, 1)
	o.Subscribe(&subscription)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		o.Next("3")
		wg.Done()
	}()
	go func() {
		o.Once("3", &subscription, func() {
			o.Next("5")
			wg.Done()
		})
	}()
	wg.Wait()
	v := <-subscription
	if v.(string) != "5" {
		t.FailNow()
	}
}

func TestObservable_Subscribe(t *testing.T) {
	ch := make(chan interface{}, 1)
	type args struct {
		ch *chan interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Subscribe",
			args: args{ch: &ch},
		},
		{
			name: "Subscribe nil chan",
			args: args{ch: nil},
		},
	}
	for k := range tests {
		t.Run(tests[k].name, func(t *testing.T) {
			o := New()
			err := o.Subscribe(tests[k].args.ch)
			if k == 1 && err == nil {
				t.FailNow()
			} else {
				return
			}
			_, ok := o.subscribers.Load(0)
			if !ok {
				t.FailNow()
			}
			o.Next("3")
			ok = false
			val, ok := <-*tests[k].args.ch
			if !ok {
				t.Errorf("No value transmitted through stored chan")
			}
			fmt.Println(val.(string))
		})
	}
}

func TestObservable_allocate(t *testing.T) {
	type fields struct {
		observable observable
	}
	type args struct {
		id int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			args: args{id: 1},
		},
	}
	for k := range tests {
		t.Run(tests[k].name, func(t *testing.T) {
			o := New()
			o.allocate(tests[k].args.id)
		})
	}
}

func TestObservable_On(t *testing.T) {
	ch := make(chan interface{}, 1)
	o := New()
	defer o.Close()
	o.Subscribe(&ch)
	var i int
	go o.On("ev3", &ch, func() {
		i = 5
	})
	o.Next("ev3")
	time.Sleep(time.Second * 3)
	if i != 5 {
		t.Errorf("Expected %d and got %d", 5, i)
	}
}

func TestObservable_Once(t *testing.T) {
	ch := make(chan interface{}, 1)
	o := New()
	defer o.Close()
	o.Subscribe(&ch)
	var i int
	go o.Once("ev3", &ch, func() {
		i = 5
	})
	o.Next("ev3")
	time.Sleep(time.Second * 3)
	if i != 5 {
		t.Errorf("Expected %d and got %d", 5, i)
	}
}
