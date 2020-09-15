package go_channel__test

import (
	"testing"
	"time"
)


func splitN(from chan string, n uint) []chan string {
	tos := make([]chan string, int(n))
	for i := uint(0); i < n; i++ {
		tos[i] = make(chan string, cap(from))
	}
	go func() {
		defer func() {
			for _, to := range tos {
				close(to)
			}
		}()
		for data := range from {
			for _, to := range tos {
				// TODO IN GO ROUTINE???
				to <- data
			}
		}
	}()
	return tos
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
func Test_splitN(t *testing.T) {
	c := make(chan string, 5)

	cs := splitN(c, 2)

	go func() {
		for s := range cs[0] {
			t.Logf("read cs[0] value \"%s\"\n", s)
		}
		t.Log("cs[0] closed")
	}()

	go func() {
		for s := range cs[1] {
			t.Logf("read cs[1] value \"%s\"\n", s)
		}
		t.Log("cs[1] closed")
	}()


	c <- "string 01"
	c <- "string 02"

	t.Logf("channel c has %d/%d items\n", len(c), cap(c))
	t.Logf("channel cs[0] has %d/%d items\n", len(cs[0]), cap(cs[0]))
	t.Logf("channel cs[1] has %d/%d items\n", len(cs[1]), cap(cs[1]))

	close(c)

	time.Sleep(100 * time.Millisecond)
	t.Log("finish")
}

func Benchmark_SplitN(b *testing.B) {
	source := make(chan string, 3)
	tos := splitN(source, 2)
	go func() {
		for _ = range tos[0] {
		}
	}()
	go func() {
		for _ = range tos[1] {
		}
	}()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		source <- "aaa"
	}
	close(source)
}