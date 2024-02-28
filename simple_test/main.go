package main

import (
	"log"
	"sync"
	"time"
)

func worker(idx int, wg *sync.WaitGroup, ch chan<- int) {
	defer func() {
		ch <- idx
		wg.Done()
	}()

	time.Sleep(1000 * time.Millisecond)

}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go worker(i, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for idx := range ch {
		log.Printf("Received signal from worker %d\n", idx)
	}
}
