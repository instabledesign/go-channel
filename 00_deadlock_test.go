package go_channel__test

import (
	"testing"
)

// fatal error: all goroutines are asleep - deadlock!
//
// goroutine 1 [chan send (nil chan)]:
// main.nilChannel()
// ..../go-channel/deadlock/main.go:26 +0x5b
// main.main()
// ..../go-channel/deadlock/main.go:13 +0x20
//
// goroutine 6 [chan receive (nil chan)]:
// main.nilChannel.func1(0x0)
// ..../go-channel/deadlock/main.go:20 +0xd4
// created by main.nilChannel
// ..../go-channel/deadlock/main.go:19 +0x42
func Test_deadlock_nil_chan(t *testing.T) {
	var c chan string

	go func() {
		// ERROR chan receive (nil chan)
		for s := range c {
			t.Logf("received data \"%s\"\n", s)
		}
	}()

	// ERROR: chan send (nil chan)
	c <- "string 01"
	close(c)
	
	t.Log("never finish")
}

// fatal error: all goroutines are asleep - deadlock!
func Test_deadlock_channel_unread(t *testing.T) {
	c := make(chan string, 2)

	c <- "string 01"
	c <- "string 02"
	// fatal error: all goroutines are asleep - deadlock!
	// script deadlock because it cannot send "string 03" and the script will never end
	c <- "string 03"
	close(c)

	t.Log("never finish")
}
