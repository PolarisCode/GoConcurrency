package main

import (
	"fmt"
	"math/rand"
	"time"
)

var cache = map[int]Book{}

func main() {

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		id := rand.Intn(10) + 1

		go func(id int) {
			if b, ok := getFromCache(id); ok {
				fmt.Println("from cache")
				fmt.Println(b)

			}
		}(id)

		go func(id int) {
			if b, ok := getFromDb(id); ok {
				fmt.Println("from db")
				fmt.Println(b)

			}
		}(id)

		time.Sleep(150 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
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
