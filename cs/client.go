package cs

import (
	"context"
	"net"
	"time"
)

// NetDialer A Dialer is a means to establish a connection.
type NetDialer struct {
	Timeout      time.Duration
	Forward      Dialer
	Chains       AdornConnsChain
	BeforeChains AdornConnsChain
	AfterChains  AdornConnsChain
}

// Dial connects to the address on the named network.
func (sf *NetDialer) Dial(network, addr string) (net.Conn, error) {
	return sf.DialContext(context.Background(), network, addr)
}

// DialContext connects to the address on the named network using the provided context.
func (sf *NetDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	var d Dialer = &net.Dialer{Timeout: sf.Timeout}

	if sf.Forward != nil {
		d = sf.Forward
	}

	contextDial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return DialContext(ctx, d, network, addr)
	}
	if f, ok := d.(ContextDialer); ok {
		contextDial = f.DialContext
	}

	conn, err := contextDial(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	for _, chain := range sf.BeforeChains {
		conn = chain(conn)
	}
	for _, chain := range sf.Chains {
		conn = chain(conn)
	}
	for _, chain := range sf.BeforeChains {
		conn = chain(conn)
	}
	return conn, nil
}

// DialContext dial context with dialer
// WARNING: this can leak a goroutine for as long as the underlying Dialer implementation takes to timeout
// A Conn returned from a successful Dial after the context has been cancelled will be immediately closed.
func DialContext(ctx context.Context, d Dialer, network, address string) (net.Conn, error) {
	var conn net.Conn
	var err error

	done := make(chan struct{}, 1)
	go func() {
		conn, err = d.Dial(network, address)
		close(done)
		if conn != nil && ctx.Err() != nil {
			conn.Close()
		}
	}()
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-done:
	}
	return conn, err
}