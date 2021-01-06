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

	for i := 0; i < 10; i++ {
		id := rand.Intn(10) + 1
		wg.Add(2) // adding number of tasks for wait
		go func(id int, wg *sync.WaitGroup) {
			if b, ok := getFromCache(id); ok {
				fmt.Println("from cache")
				fmt.Println(b)

			}
			wg.Done() // signal that task was done
		}(id, wg)

		go func(id int, wg *sync.WaitGroup) {
			if b, ok := getFromDb(id); ok {
				fmt.Println("from db")
				fmt.Println(b)

			}
			wg.Done()
		}(id, wg)

		time.Sleep(150 * time.Millisecond)
	}

	wg.Wait() // wait for all routines are done
}

func getFromCache(id int) (Book, bool) {
	b, ok := cache[id]
	return b, ok
}

func getFromDb(id int) (Book, bool) {
	time.Sleep(100 * time.Millisecond)
	for _, b := range books {
		if b.ID == id {
			cache[id] = b
			return b, true
		}
	}
	return Book{}, false
}
