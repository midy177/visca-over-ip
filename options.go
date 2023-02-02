package visca

import (
	"net"
	"time"
)

var (
	_defaultTTL   = 30 * time.Second
	_defaultDelay = 75 * time.Millisecond
)

type options struct {
	ttl    time.Duration
	delay  time.Duration
	logger Logger
	dialer net.Dialer
}

// Option configures how we create the Camera.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithTTL changes the TTL for the underlying UDP connection to the Camera.
// The default value is 30 seconds.
// See more details about TTL at https://github.com/byuoitav/connpool.
func WithTTL(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.ttl = t
	})
}

// WithDelay changes the delay between sending commands to the Camera.
// The default value is 250 milliseconds.
// See more details about delay at https://github.com/byuoitav/connpool.
func WithDelay(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.delay = t
	})
}

// WithLogger adds a logger to Camera.
// Camera will log appropriate information about the underlying connection and the commands being sent.
// The default value is nil, meaning that no logs are written.
func WithLogger(l Logger) Option {
	return optionFunc(func(o *options) {
		o.logger = l
	})
}

// WithDialer sets the dialer to use when opening connections with the Camera.
// The default value is net.Dialer{}.
// See more details at https://golang.org/pkg/net/#Dialer
func WithDialer(d net.Dialer) Option {
	return optionFunc(func(o *options) {
		o.dialer = d
	})
}
