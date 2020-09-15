package go_channel__test

import (
	"testing"
	"time"
)
////////////////////////////////////////////////////////////
// UNBUFFERED CHANNEL
// With unbuffered channel you need a goroutine to write or read channel data
////////////////////////////////////////////////////////////

// A) read c value "string 01"
// B) read c value "string 02"
// finish
func Test_simple_goroutine_send_data(t *testing.T) {
	// UNbuffered channel
	c := make(chan string)

	// goroutine send data in channel
	go func() {
		c <- "string 01"
		// it will block here until first value was read
		c <- "string 02"
	}()

	// sync read first value (it wait until value was send)
	t.Logf("A) read c value \"%s\"\n", <-c)

	// sync read second value (it wait until value was send)
	t.Logf("B) read c value \"%s\"\n", <-c)

	t.Log("finish")
}

// A) read c value "string 01"
// B) read c value "string 02"
// finish
func Test_simple_goroutine_read_data(t *testing.T) {
	// UNbuffered channel
	c := make(chan string)

	// goroutine read channel data
	go func() {
		// sync read first value (it wait until value was send)
		t.Logf("A) read c value \"%s\"\n", <-c)

		// sync read second value (it wait until value was send)
		t.Logf("B) read c value \"%s\"\n", <-c)
	}()

	c <- "string 01"
	// it will block here until first value was read
	c <- "string 02"
	close(c)

	t.Log("finish")

	// give time to goroutine to read second value
	time.Sleep(100 * time.Millisecond)
}
