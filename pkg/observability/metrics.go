package observability

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Metrics stores application metrics
type Metrics struct {
	counters   map[string]*int64
	gauges     map[string]*uint64 // Changed to uint64 for atomic float64 operations
	histograms map[string]*Histogram
	mu         sync.RWMutex
}

// Histogram represents a simple histogram for tracking durations
type Histogram struct {
	count int64
	sum   uint64 // Changed to uint64 for atomic float64 operations
	min   uint64 // Changed to uint64 for atomic float64 operations
	max   uint64 // Changed to uint64 for atomic float64 operations
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	return &Metrics{
		counters:   make(map[string]*int64),
		gauges:     make(map[string]*uint64),
		histograms: make(map[string]*Histogram),
	}
}

// IncrementCounter increments a counter by 1
func (m *Metrics) IncrementCounter(name string) {
	m.mu.RLock()
	counter, exists := m.counters[name]
	m.mu.RUnlock()

	if !exists {
		m.mu.Lock()
		if _, exists := m.counters[name]; !exists {
			var c int64
			m.counters[name] = &c
			counter = &c
		}
		m.mu.Unlock()
	}

	atomic.AddInt64(counter, 1)
}

// SetGauge sets a gauge value
func (m *Metrics) SetGauge(name string, value float64) {
	m.mu.RLock()
	gauge, exists := m.gauges[name]
	m.mu.RUnlock()

	if !exists {
		m.mu.Lock()
		if _, exists := m.gauges[name]; !exists {
			var g uint64
			m.gauges[name] = &g
			gauge = &g
		}
		m.mu.Unlock()
	}

	atomic.StoreUint64(gauge, math.Float64bits(value))
}

// getFloat64 safely converts a uint64 to float64
func getFloat64(u uint64) float64 {
	return math.Float64frombits(u)
}

// RecordDuration records a duration in a histogram
func (m *Metrics) RecordDuration(name string, duration time.Duration) {
	m.mu.RLock()
	hist, exists := m.histograms[name]
	m.mu.RUnlock()

	if !exists {
		m.mu.Lock()
		if _, exists := m.histograms[name]; !exists {
			hist = &Histogram{}
			// Initialize min with the first value
			atomic.StoreUint64(&hist.min, math.Float64bits(float64(duration.Milliseconds())))
			m.histograms[name] = hist
		}
		m.mu.Unlock()
	}

	ms := float64(duration.Milliseconds())
	msBits := math.Float64bits(ms)

	atomic.AddInt64(&hist.count, 1)

	// Update sum atomically
	for {
		old := atomic.LoadUint64(&hist.sum)
		newSum := math.Float64bits(getFloat64(old) + ms)
		if atomic.CompareAndSwapUint64(&hist.sum, old, newSum) {
			break
		}
	}

	// Update min atomically
	for {
		old := atomic.LoadUint64(&hist.min)
		if ms >= getFloat64(old) {
			break
		}
		if atomic.CompareAndSwapUint64(&hist.min, old, msBits) {
			break
		}
	}

	// Update max atomically
	for {
		old := atomic.LoadUint64(&hist.max)
		if ms <= getFloat64(old) {
			break
		}
		if atomic.CompareAndSwapUint64(&hist.max, old, msBits) {
			break
		}
	}
}

// GetMetrics returns all current metrics
func (m *Metrics) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make(map[string]interface{})

	// Copy counters
	for name, value := range m.counters {
		metrics["counter."+name] = atomic.LoadInt64(value)
	}

	// Copy gauges
	for name, value := range m.gauges {
		metrics["gauge."+name] = getFloat64(atomic.LoadUint64(value))
	}

	// Copy histograms
	for name, hist := range m.histograms {
		count := atomic.LoadInt64(&hist.count)
		if count > 0 {
			sum := getFloat64(atomic.LoadUint64(&hist.sum))
			metrics["histogram."+name+".count"] = count
			metrics["histogram."+name+".avg"] = sum / float64(count)
			metrics["histogram."+name+".min"] = getFloat64(atomic.LoadUint64(&hist.min))
			metrics["histogram."+name+".max"] = getFloat64(atomic.LoadUint64(&hist.max))
		}
	}

	return metrics
}
