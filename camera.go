package visca

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/connpool"
)

type Camera struct {
	address string
	pool    *connpool.Pool

	logger Logger
	dialer net.Dialer
}

func New(addr string, opts ...Option) *Camera {
	options := options{
		ttl:    _defaultTTL,
		delay:  _defaultDelay,
		dialer: net.Dialer{},
	}

	for _, o := range opts {
		o.apply(&options)
	}

	cam := &Camera{
		address: addr,
		pool: &connpool.Pool{
			TTL:    options.ttl,
			Delay:  options.delay,
			Logger: options.logger,
		},
		logger: options.logger,
		dialer: options.dialer,
	}

	cam.pool.NewConnection = func(ctx context.Context) (net.Conn, error) {
		return cam.dialer.DialContext(ctx, "udp", cam.address)
	}

	return cam
}

func (c *Camera) sendPayload(ctx context.Context, p payload) (payload, error) {
	var resp payload

	bytes, err := p.MarshalBinary()
	if err != nil {
		return resp, err
	}

	err = c.pool.Do(ctx, func(conn connpool.Conn) error {
		c.debugf("Sending payload: %# x", bytes)

		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(3 * time.Second)
		}

		conn.SetWriteDeadline(deadline)

		n, err := conn.Write(bytes)
		switch {
		case err != nil:
			return fmt.Errorf("unable to write payload: %w", err)
		case n != len(bytes):
			return fmt.Errorf("unable to write payload: wrote %d/%d bytes", n, len(bytes))
		}

		data, err := conn.ReadUntil(_terminator, deadline)
		if err != nil {
			return fmt.Errorf("unable to read response: %w", err)
		}

		c.debugf("Got response: %# x", data)

		if err := resp.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("unable to parse response: %w", err)
		}

		return nil
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}
