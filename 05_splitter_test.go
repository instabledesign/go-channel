package go_channel__test

import (
	"sync"
	"testing"
	"time"
)

type Splitter struct {
	from chan string
	tos  []chan string
	mu   sync.RWMutex
}

func (splitter *Splitter) start() *Splitter {
	go func() {
		defer splitter.close()
		for data := range splitter.from {
			splitter.mu.RLock()
			for _, to := range splitter.tos {
				to <- data
			}
			splitter.mu.RUnlock()
		}
	}()
	return splitter
}

func (splitter *Splitter) close() {
	for _, to := range splitter.tos {
		close(to)
	}
}

func (splitter *Splitter) Split() (chan string, func()) {
	to := make(chan string, cap(splitter.from))
	i := len(splitter.tos)
	splitter.mu.Lock()
	splitter.tos = append(splitter.tos, to)
	splitter.mu.Unlock()
	return to, func() {
		defer close(to)
		splitter.mu.Lock()
		splitter.tos = append(splitter.tos[:i], splitter.tos[i+1:]...)
		splitter.mu.Unlock()
	}
}

func NewSplitter(from chan string) *Splitter {
	return (&Splitter{from: from}).start()
}

// read c1 value "string 01"
// read c1 value "string 02"
// read c2 value "string 01"
// read c2 value "string 02"
// channel c has 0/5 items
// channel c1 has 0/5 items
// channel c2 has 0/5 items
// read c2 value "string 03"
// c1 closed
// c2 closed
// finish
func Test_splitter(t *testing.T) {
	c := make(chan string, 5)

	splitter := NewSplitter(c)
	c1, c1Closefn := splitter.Split()
	c2, _ := splitter.Split()

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

	time.Sleep(100 * time.Millisecond)
	c1Closefn()

	c <- "string 03"

	t.Logf("channel c has %d/%d items\n", len(c), cap(c))
	t.Logf("channel c1 has %d/%d items\n", len(c1), cap(c1))
	t.Logf("channel c2 has %d/%d items\n", len(c2), cap(c2))

	close(c)

	time.Sleep(100 * time.Millisecond)
	t.Log("finish")
}

func Benchmark_Splitter(b *testing.B) {
	source := make(chan string, 3)
	splitter := NewSplitter(source)
	out1, _ := splitter.Split()
	out2, _ := splitter.Split()

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
