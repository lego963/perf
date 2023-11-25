package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type bookStash struct {
	stash map[string]chan struct{}
	mu    *sync.Mutex
}

func main() {
	bs := bookStash{
		stash: make(map[string]chan struct{}, 10),
		mu:    &sync.Mutex{},
	}
	for i := 0; i < 10; i++ {
		bs.stash[strconv.Itoa(i)] = make(chan struct{}, 10)
	}
	stopCh := time.After(10 * time.Second)

	for i := 0; i < len(bs.stash); i++ {
		i := i
		go func(i int) {
			for {
				select {
				case <-bs.stash[strconv.Itoa(i)]:
					fmt.Printf("you returned the %d book\n", i)
				}
			}
		}(i)
	}

	go func() {
		for {
			randBookInt := rand.Intn(10-0) + 0
			randBookStr := strconv.Itoa(randBookInt)
			if _, ok := bs.stash[randBookStr]; !ok {
				fmt.Println(fmt.Errorf("we don't have requested book: %s", randBookStr))
			}
			bs.mu.Lock()
			bs.stash[randBookStr] <- struct{}{}
			bs.mu.Unlock()
			fmt.Printf("you took the %s book\n", randBookStr)
		}
	}()

	select {
	case <-stopCh:
		fmt.Println("cya")
		return
	}
}
