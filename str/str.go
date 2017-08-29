package str

import (
	"math/rand"
	"time"

	"go.uber.org/atomic"
)

// Str ...
type Str struct {
	alph   []rune
	f      atomic.Int32
	t      atomic.Int32
	length int
	from   int
	to     int
	delay  time.Duration
	source rand.Source
	strs   chan string
	done   chan struct{}
}

// Option ...
type Option func(*Str)

// Source ...
func Source(param rand.Source) Option {
	return func(r *Str) {
		r.source = param
	}
}

// Alphabet ...
func Alphabet(alph string) Option {
	return func(r *Str) {
	}
}

// Upper ...
func Upper() Option {
	return func(r *Str) {
	}
}

// Lower ...
func Lower() Option {
	return func(r *Str) {
	}
}

// Digits ...
func Digits() Option {
	return func(r *Str) {
	}
}

// Length ...
func Length(param int) Option {
	return func(r *Str) {
		r.length = param
	}
}

// From ...
func From(param int) Option {
	return func(r *Str) {
		r.from = param
	}
}

// To ...
func To(param int) Option {
	return func(r *Str) {
		r.to = param
	}
}

// Size ...
func Size(param int) Option {
	return func(r *Str) {
		r.strs = make(chan string, param)
	}
}

// Delay ...
func Delay(param time.Duration) Option {
	return func(r *Str) {
		r.delay = param
	}
}

// NewStr ...
func NewStr(Options ...Option) *Str {
	r := &Str{
		done: make(chan struct{}),
	}

	defaults := []Option{
		Alphabet(`abcdefghijklmnopqrstuvwxyz0123456789`),
		Source(rand.NewSource(time.Now().UnixNano())),
		From(10),
		To(10),
		Size(64),
		Delay(500 * time.Millisecond),
	}
	Options = append(defaults, Options...)
	for _, o := range Options {
		o(r)
	}

	go r.run()
	return r
}

// Next ...
func (r *Str) Next() string {
	res := <-r.strs
	return res
}

func (r *Str) new() string {
	res := make([]rune, r.length)
	for i := 0; i < r.length; i++ {
		n := r.source.Int63()
		res[i] = r.alph[int(n)%len(r.alph)]
	}
	return string(res)
}

func (r *Str) fill() {
	for {
		select {
		case r.strs <- r.new():
		default:
			break
		}
	}
}

func (r *Str) run() {
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
