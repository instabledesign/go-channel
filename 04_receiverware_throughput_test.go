package go_channel__test

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func throughputReceiverware(in <-chan string, throughput *uint64) <-chan string {
	previousThroughput := uint64(0)

	// a ticker to capture throughput each second
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				atomic.StoreUint64(throughput, previousThroughput)
				atomic.StoreUint64(&previousThroughput, 0)
			case <-quit:
				return
			}
		}
	}()

	// we decorate the input channel
	out := make(chan string, len(in))
	go func() {
		defer close(quit)
		defer close(out)
		for data := range in {
			atomic.AddUint64(&previousThroughput, 1)
			out <- data
		}
	}()
	return out
}

func Test_throughputReceiverware(t *testing.T) {
	c := make(chan string, 10000)
	t.Logf("channel c has %d/%d items\n", len(c), cap(c))

	instantThroughput := uint64(0)

	// print throughput every second
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				t.Logf("instant throughput %d", instantThroughput)
			case <-quit:
				return
			}
		}
	}()

	wrappedC := throughputReceiverware(c, &instantThroughput)

	// this simulate your application usage of the channel data
	var result []string
	go func() {
		for s := range wrappedC {
			result = append(result, s)
		}
	}()

	// publish data on the source channel
	go func() {
		for i := 0; i < 100000; i++ {
			c <- ""
			// add entropy in sleep time in order to have a different throughput over the time
			time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
		}
	}()

	time.Sleep(5 * time.Second)
	t.Logf("finish with %d results", len(result))
}
