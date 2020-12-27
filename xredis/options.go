package xredis

import "time"

type options struct {
	maxActive   int
	maxIdle     int
	wait        bool
	idleTimeout time.Duration
}

var defaultOptions = options{
	maxActive:   5,
	maxIdle:     2,
	wait:        false,
	idleTimeout: 30 * time.Second,
}

type Option func(*options)

func WithMaxActive(maxActive int) Option {
	return func(o *options) {
		o.maxActive = maxActive
	}
}

func WithMaxIdle(maxIdle int) Option {
	return func(o *options) {
		o.maxIdle = maxIdle
	}
}

func WithWait(wait bool) Option {
	return func(o *options) {
		o.wait = wait
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *options) {
		o.idleTimeout = idleTimeout
	}
}
