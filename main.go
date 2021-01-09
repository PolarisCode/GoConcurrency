package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var cache = map[int]Book{}

func main() {
	waitGroupWithRWMutex()
	//waitGroupWithChannel()
}

func waitGroupWithChannel() {

	wg := &sync.WaitGroup{}
	ch := make(chan int)

	wg.Add(2) // adding number of tasks for wait

	go func(wg *sync.WaitGroup, ch <-chan int) {
		for msg := range ch {
			fmt.Println(msg)
		}
		wg.Done() // signal that task was done
	}(wg, ch)

	go func(wg *sync.WaitGroup, ch chan<- int) {

		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)

		wg.Done()
	}(wg, ch)

	time.Sleep(150 * time.Millisecond)

	wg.Wait() // wait for all routines are done
}

func waitGroupWithRWMutex() {
	rand.Seed(time.Now().UnixNano())

	wg := &sync.WaitGroup{}
	m := &sync.RWMutex{}
	cacheCh := make(chan Book)
	dbCh := make(chan Book)

	for i := 0; i < 10; i++ {
		id := rand.Intn(10) + 1
		wg.Add(2) // adding number of tasks for wait

		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, ch chan<- Book) {
			if b, ok := getFromCache(id, m); ok {
				ch <- b
			}
			wg.Done() // signal that task was done
		}(id, wg, m, cacheCh)

		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, ch chan<- Book) {
			if b, ok := getFromDb(id, m); ok {
				m.Lock()
				cache[id] = b
				m.Unlock()
				ch <- b
			}
			wg.Done()
		}(id, wg, m, dbCh)

		go func(cacheCh, dbCh <-chan Book) {
			select {
			case b := <-cacheCh:
				fmt.Println("from cache")
				fmt.Println(b)
				<-dbCh
			case b := <-dbCh:
				fmt.Println("from db")
				fmt.Println(b)
			}
		}(cacheCh, dbCh)

		time.Sleep(150 * time.Millisecond)
	}

	wg.Wait() // wait for all routines are done

}

func getFromCache(id int, m *sync.RWMutex) (Book, bool) {
	m.RLock()
	b, ok := cache[id]
	m.RUnlock()
	return b, ok
}

func getFromDb(id int, m *sync.RWMutex) (Book, bool) {
	time.Sleep(100 * time.Millisecond)
	for _, b := range books {
		if b.ID == id {
			return b, true
		}
	}
	return Book{}, false
}
