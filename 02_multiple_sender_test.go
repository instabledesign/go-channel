package go_channel__test

import (
	"testing"
	"time"
)

// channel c has 4/5 items
// channel closed
// read c (3/5) value "sender 2 data 1"
// read c (2/5) value "sender 2 data 2"
// read c (1/5) value "sender 1 data 1"
// read c (0/5) value "sender 1 data 2"
// finish
func Test_multiple_sender(t *testing.T) {
	c := make(chan string, 5)
	time.AfterFunc(50*time.Millisecond, func() {
		t.Logf("channel c has %d/%d items\n", len(c), cap(c))
		t.Log("channel closed")
		close(c)
	})

	go func() {
		c <- "sender 1 data 1"
		c <- "sender 1 data 2"
	}()

	go func() {
		c <- "sender 2 data 1"
		c <- "sender 2 data 2"
	}()

	time.Sleep(100 * time.Millisecond)

	for s := range c {
		t.Logf("read c (%d/%d) value \"%s\" \n", len(c), cap(c), s)
	}

	t.Log("finish")
}
