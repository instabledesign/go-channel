package go_channel__test

import (
	"strings"
	"testing"
)

// This receiverware will append string from original input data
func prefixReceiverware(prefix string, in <-chan string) <-chan string {
	out := make(chan string, len(in))
	go func() {
		defer close(out)
		// it will forward in data channel to out channel with prefix
		for data := range in {
			out <- prefix + data
		}
	}()
	return out
}

// This receiverware will uppercase string from original input data
func uppercaseReceiverware(in <-chan string) <-chan string {
	out := make(chan string, len(in))
	go func() {
		defer close(out)
		// it will forward in data channel to out channel with prefix
		for data := range in {
			out <- strings.ToUpper(data)
		}
	}()
	return out
}

// channel c has 0/5 items
// read c value "MY CUSTOM PREFIXSTRING 01"
// read c value "MY CUSTOM PREFIXSTRING 02"
// finish
func Test_receiverware(t *testing.T) {
	c := make(chan string, 5)
	t.Logf("channel c has %d/%d items\n", len(c), cap(c))

	go func() {
		defer close(c)
		c <- "string 01"
		c <- "string 02"
	}()

	// we decorate channel at
	wrappedC := prefixReceiverware("my custom prefix", c)
	wrappedC = uppercaseReceiverware(wrappedC)

	for s := range wrappedC {
		t.Logf("read c value \"%s\"\n", s)
	}

	t.Log("finish")
}
