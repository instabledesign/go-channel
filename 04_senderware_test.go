package go_channel__test

import (
	"strings"
	"testing"
	"time"
)

func prefixSenderware(prefix string, in chan<- string) chan<- string {
	out := make(chan string, cap(in))
	go func() {
		defer close(out)
		for data := range out {
			in <- prefix + data
		}
	}()
	return out
}

// This senderware will uppercase string from original input data
func uppercaseSenderware(in chan<- string) chan<- string {
	out := make(chan string, cap(in))
	go func() {
		defer close(out)
		for data := range out {
			in <- strings.ToUpper(data)
		}
	}()
	return out
}

// channel c has 0/5 items
// read c value "my custom prefixSTRING 01"
// read c value "my custom prefixSTRING 02"
// finish
func Test_senderware(t *testing.T) {
	c := make(chan string, 5)
	go func() {
		for s := range c {
			t.Logf("read c value \"%s\"\n", s)
		}
	}()

	t.Logf("channel c has %d/%d items\n", len(c), cap(c))

	wrappedC := prefixSenderware("my custom prefix", c)
	wrappedC = uppercaseSenderware(wrappedC)

	wrappedC <- "string 01"
	wrappedC <- "string 02"

	t.Log("finish")
	time.Sleep(100 * time.Millisecond)
}
