package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var cache = map[int]Book{}

func main() {

	rand.Seed(time.Now().UnixNano())

	wg := &sync.WaitGroup{}
	m := &sync.Mutex{}

	for i := 0; i < 20; i++ {
		id := rand.Intn(10) + 1
		wg.Add(2) // adding number of tasks for wait
		go func(id int, wg *sync.WaitGroup, m *sync.Mutex) {
			if b, ok := getFromCache(id, m); ok {
				fmt.Println("from cache")
				fmt.Println(b)

			}
			wg.Done() // signal that task was done
		}(id, wg, m)

		go func(id int, wg *sync.WaitGroup, m *sync.Mutex) {
			if b, ok := getFromDb(id, m); ok {
				fmt.Println("from db")
				fmt.Println(b)
			}
			wg.Done()
		}(id, wg, m)

		time.Sleep(150 * time.Millisecond)
	}

	wg.Wait() // wait for all routines are done
}

func getFromCache(id int, m *sync.Mutex) (Book, bool) {
	m.Lock()
	b, ok := cache[id]
	m.Unlock()
	return b, ok
}

func getFromDb(id int, m *sync.Mutex) (Book, bool) {
	time.Sleep(100 * time.Millisecond)
	for _, b := range books {
		if b.ID == id {
			m.Lock()
			cache[id] = b
			m.Unlock()
			return b, true
		}
	}
	return Book{}, false
}
