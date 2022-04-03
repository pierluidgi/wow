package metrics

import (
	log "github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

type RateMeter struct {
	interval time.Duration
	size     int
	requests uint64
	rate     uint64
	values   []uint64
}

func (m *RateMeter) Add(n uint64) {
	atomic.AddUint64(&m.requests, n)
}

func (m *RateMeter) Rate() uint64 {
	return atomic.LoadUint64(&m.rate)
}

func (m *RateMeter) run() {
	ticker := time.NewTicker(m.interval)

	for range ticker.C {
		requests := atomic.SwapUint64(&m.requests, 0)

		m.values = append(m.values, requests/uint64(m.interval.Seconds()))

		if len(m.values) > m.size {
			for i := 0; i < m.size; i++ {
				m.values[i] = m.values[i+1]
			}

			m.values = m.values[:m.size]
		}

		var avgRate uint64

		for _, r := range m.values {
			avgRate += r
		}

		avgRate = avgRate / uint64(len(m.values))

		atomic.StoreUint64(&m.rate, avgRate)

		log.Debugf("Rate: %d rps", avgRate)
	}
}

func NewRateMeter(interval int, size int) *RateMeter {
	if interval < 1 {
		interval = 1
	}

	if size < 1 {
		size = 1
	}

	m := &RateMeter{
		interval: time.Duration(interval) * time.Second,
		size:     size,
		requests: 0,
		rate:     0,
		values:   make([]uint64, 0, size+1),
	}

	go m.run()

	return m
}
