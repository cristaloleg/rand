package floats

import (
	"math/rand"
	"time"
)

// Rand ...
type Rand struct {
	source rand.Source64
	from   int
	to     int
	nums   chan float64
	delay  time.Duration
	done   chan struct{}
}

// Option ...
type Option func(*Rand)

// Source ...
func Source(param rand.Source) Option {
	return func(r *Rand) {
		r.source = param.(rand.Source64)
	}
}

// Source64 ...
func Source64(param rand.Source64) Option {
	return func(r *Rand) {
		r.source = param
	}
}

// From ...
func From(param int) Option {
	return func(r *Rand) {
		r.from = param
	}
}

// To ...
func To(param int) Option {
	return func(r *Rand) {
		r.to = param
	}
}

// Size ...
func Size(param int) Option {
	return func(r *Rand) {
		r.nums = make(chan float64, param)
	}
}

// Delay ...
func Delay(param time.Duration) Option {
	return func(r *Rand) {
		r.delay = param
	}
}

// New ...
func New(opts ...Option) *Rand {
	r := &Rand{
		done: make(chan struct{}),
	}

	Size(1024)(r)
	Delay(500 * time.Millisecond)(r)

	for _, o := range opts {
		o(r)
	}
	go r.run()
	return r
}

// Int ...
func (r *Rand) Int() int {
	res := <-r.nums
	return int(res)
}

// float64 ...
func (r *Rand) float64() float64 {
	res := <-r.nums
	return res
}

// Close ...
func (r *Rand) Close() {
	r.done <- struct{}{}
}

func (r *Rand) new() float64 {
	return float64(r.source.Uint64())
}

func (r *Rand) fill() {
	for {
		select {
		case r.nums <- r.new():
		default:
			break
		}
	}
}

func (r *Rand) run() {
	ticker := time.NewTicker(r.delay)
	for {
		select {
		case <-ticker.C:
			r.fill()
		case <-r.done:
			return
		}
	}
}
