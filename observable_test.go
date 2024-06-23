package observable

import (
	"fmt"
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
			_, err := o.Subscribe(tests[k].args.ch)
			if k == 1 && err == nil {
				t.FailNow()
			} else {
				return
			}
			_, ok := o.subscribers[0]
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

func TestObservable_Unsubscribe(t *testing.T) {
	type args struct {
		id int
	}
	o := New()
	c := make(chan interface{}, 1)
	o.Subscribe(&c)
	tests := []struct {
		name   string
		O      *Observable
		args   args
		wantOk bool
	}{
		{
			name:   "#1",
			O:      o,
			args:   args{id: 0},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := o.Unsubscribe(tt.args.id); gotOk != tt.wantOk {
				t.Errorf("Observable.Unsubscribe() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
