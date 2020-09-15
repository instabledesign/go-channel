package go_channel__test

import (
	"sync"
	"testing"
	"time"
)


//https://medium.com/justforfunc/two-ways-of-merging-n-channels-in-go-43c0b57cd1de
func merge(cs ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan string) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// out data receive data 1 c3
// out data receive data 1 c1
// out data receive data 1 c2
// channel c1 has 0/3 items
// channel c2 has 0/3 items
// channel c3 has 0/3 items
// channel out has 0/0 items
// out data receive data 2 c2
// out data receive data 2 c3
// out data receive data 3 c3
// c1 c2 c3 was closed
// finish
func Test_merge(t *testing.T) {
	c1 := make(chan string, 3)
	c2 := make(chan string, 3)
	c3 := make(chan string, 3)

	out := merge(c1, c2, c3)
	go func() {
		for data := range out {
			t.Log("out data receive", data)
		}
		t.Log("c1 c2 c3 was closed")
	}()

	c1 <- "data 1 c1"
	c2 <- "data 1 c2"
	c3 <- "data 1 c3"

	time.Sleep(100 * time.Millisecond)
	t.Logf("channel c1 has %d/%d items\n", len(c1), cap(c1))
	t.Logf("channel c2 has %d/%d items\n", len(c2), cap(c2))
	t.Logf("channel c3 has %d/%d items\n", len(c3), cap(c3))
	t.Logf("channel out has %d/%d items\n", len(out), cap(out))

	close(c1)
	c2 <- "data 2 c2"
	c3 <- "data 2 c3"

	close(c2)

	c3 <- "data 3 c3"
	close(c3)

	time.Sleep(100 * time.Millisecond)
	t.Log("finish")
}
