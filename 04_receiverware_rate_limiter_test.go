package go_channel__test

import (
	"testing"
	"time"
)

func rateLimiterReceiverware(in <-chan string, limit uint, interval time.Duration) <-chan string {
	currentCount := uint64(0)

	// a ticker to capture throughput each second
	ticker := time.NewTicker(interval)

	// we decorate the input channel
	out := make(chan string, len(in))
	go func() {
		defer ticker.Stop()
		defer close(out)
		for data := range in {
			if currentCount >= uint64(limit) {
				<-ticker.C
				currentCount = 0
			}
			currentCount++
			out <- data
		}
	}()
	return out
}

func Test_rateLimiterReceiverware(t *testing.T) {
	c := make(chan string, 1000)
	t.Logf("channel c has %d/%d items\n", len(c), cap(c))

	wrappedC := rateLimiterReceiverware(c, 1, time.Second)

	// this simulate your application usage of the channel data
	var result []string
	go func() {
		for s := range wrappedC {
			result = append(result, s)
			t.Log(time.Now())
		}
	}()

	// publish data on the source channel
	go func() {
		for i := 0; i < 100000; i++ {
			c <- ""
			// add entropy in sleep time in order to have a different throughput over the time
			//time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
		}
	}()

	time.Sleep(5 * time.Second)
	t.Logf("finish with %d results", len(result))
}
