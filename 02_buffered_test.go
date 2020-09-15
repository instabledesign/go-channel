package go_channel__test

import (
	"testing"
)

// channel c has 0/5
// A) read c value "string 01"
// B) read c value "string 02"
// finish
func Test_buffered(t *testing.T) {
	c := make(chan string, 5)
	t.Logf("channel c has %d/%d items\n", len(c), cap(c))

	c <- "string 01"
	c <- "string 02"
	close(c)

	t.Logf("A) read c value \"%s\"\n", <-c)
	t.Logf("B) read c value \"%s\"\n", <-c)

	t.Log("finish")
}
