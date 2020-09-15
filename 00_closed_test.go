package go_channel__test

import (
	"testing"
)

func Test_closed(t *testing.T) {
	c := make(chan string)
	close(c)

	// panic: send on closed channel
	c <- "string 01"

	t.Log("never finish")
}

func Test_already_closed(t *testing.T) {
	c := make(chan string)
	close(c)

	// panic: close of closed channel
	close(c)

	t.Log("never finish")
}
