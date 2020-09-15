package go_channel__test

import (
	"testing"
	"time"
)

func split(from chan string) (chan string, chan string) {
	to1 := make(chan string, cap(from))
	to2 := make(chan string, cap(from))
	go func() {
		defer close(to1)
		defer close(to2)
		for data := range from {
			to1 <- data
			to2 <- data
		}
	}()
	return to1, to2
}

// read c1 value "string 01"
// read c1 value "string 02"
// channel c has 2/5 items
// channel c1 has 0/5 items
// read c2 value "string 01"
// read c2 value "string 02"
// channel c2 has 1/5 items
// c1 closed
// c2 closed
// finish
func Test_split(t *testing.T) {
	c := make(chan string, 5)

	c1, c2 := split(c)

	go func() {
		for s := range c1 {
			t.Logf("read c1 value \"%s\"\n", s)
		}
		t.Log("c1 closed")
	}()

	go func() {
		for s := range c2 {
			t.Logf("read c2 value \"%s\"\n", s)
		}
		t.Log("c2 closed")
	}()

	c <- "string 01"
	c <- "string 02"

	t.Logf("channel c has %d/%d items\n", len(c), cap(c))
	t.Logf("channel c1 has %d/%d items\n", len(c1), cap(c1))
	t.Logf("channel c2 has %d/%d items\n", len(c2), cap(c2))

	close(c)

	time.Sleep(100 * time.Millisecond)
	t.Log("finish")
}

func Benchmark_Split(b *testing.B) {
	source := make(chan string, 3)
	out1, out2 := split(source)
	go func() {
		for _ = range out1 {
		}
	}()
	go func() {
		for _ = range out2 {
		}
	}()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		source <- "aaa"
	}
	close(source)
}
