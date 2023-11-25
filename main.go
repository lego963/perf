package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	bookStash := make(map[string]chan struct{}, 10)
	for i := 0; i < 10; i++ {
		bookStash[strconv.Itoa(i)] = make(chan struct{}, 2)
	}
	stopCh := time.After(10 * time.Second)

	for i := 0; i < len(bookStash); i++ {
		i := i
		go func(i int) {
			for {
				select {
				case <-bookStash[strconv.Itoa(i)]:
					fmt.Printf("you returned the %d book\n", i)
				default:
					time.Sleep(1 * time.Second)
				}
			}
		}(i)
	}

	go func() {
		for {
			randBookInt := rand.Intn(15-0) + 0
			randBookStr := strconv.Itoa(randBookInt)
			if _, ok := bookStash[randBookStr]; !ok {
				fmt.Println(fmt.Errorf("we don't have requested book: %s", randBookStr))
			}
			bookStash[randBookStr] <- struct{}{}
			fmt.Printf("you took the %s book\n", randBookStr)
		}
	}()

	select {
	case <-stopCh:
		fmt.Println("cya")
		return
	}
}
