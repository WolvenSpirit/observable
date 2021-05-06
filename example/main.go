package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/WolvenSpirit/observable"
)

const parralelRunners = 300

func produce(input *chan interface{}, o *observable.Observable, wg *sync.WaitGroup, shutdown *chan bool) {
	for true {
		select {
		case val := <-*input:
			switch val.(type) {
			case int:
				v := val.(int)
				v++
				o.Next(v)
			case float32:
				v := val.(float32)
				v++
				o.Next(v)
			case float64:
				v := val.(float64)
				v++
				o.Next(v)
			}
		case q := <-*shutdown:
			if q {
				log.Println("Received shutdown")
				wg.Done()
				break
			}
		}
	}
}

func consume(o *observable.Observable, shutdown *chan bool, name string) {
	ch := make(chan interface{}, 1)
	o.Subscribe(&ch)
	for true {
		select {
		case val := <-ch:
			log.Printf("Received %+v on %s", val, name)
		case q := <-*shutdown:
			if q {
				break
			}
		}
	}
	close(ch)
}

func main() {
	var q chan bool
	q = make(chan bool, parralelRunners*4) // we need to send signal to 300 producers and 900 consumers
	// shutdown
	go func() {
		time.Sleep(time.Second * 3)
		log.Println("Closing goroutines")
		for i := 0; i < parralelRunners*4; i++ {
			q <- true
		}
	}()
	var wg sync.WaitGroup
	wg.Add(parralelRunners)
	input := make(chan interface{}, parralelRunners)
	for i := 0; i < parralelRunners; i++ {
		o := observable.New()
		go produce(&input, o, &wg, &q)
		// one producer and three consumers pair
		for j := 0; j < 3; j++ {
			go consume(o, &q, fmt.Sprintf("@consumer#%d_%d", i, j))
		}
	}
	// let's get them some work
	for i := 0; i < parralelRunners; i++ {
		input <- i
	}
	wg.Wait()
	close(q)
}
