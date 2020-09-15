package go_channel__test

import (
	"testing"
	"time"
)

// channel c has 0/5 items
// receiver 1 read c(2/5) value "data 1"
// receiver 1 read c(0/5) value "data 3"
// receiver 1 goroutine end (channel closed)
// receiver 2 read c(1/5) value "data 2"
// receiver 2 goroutine end (channel closed)
// finish
func Test_multiple_receiver(t *testing.T) {
	c := make(chan string, 5)
	t.Logf("channel c has %d/%d items\n", len(c), cap(c))

	go func() {
		for s := range c {
			t.Logf("receiver 1 read c(%d/%d) value \"%s\" \n", len(c), cap(c), s)
		}
		t.Log("receiver 1 goroutine end (channel closed)")
	}()

	go func() {
		for s := range c {
			t.Logf("receiver 2 read c(%d/%d) value \"%s\" \n", len(c), cap(c), s)
		}
		t.Log("receiver 2 goroutine end (channel closed)")
	}()

	// values are randomly dispatch to the receiver
	c <- "data 1"
	c <- "data 2"
	c <- "data 3"
	close(c)

	time.Sleep(100 * time.Millisecond)

	t.Log("finish")
}
