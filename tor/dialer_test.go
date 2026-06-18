package tor

import (
	"context"
	"errors"
	"net"
	"testing"
)

var (
	errPlainDial   = errors.New("plain dial called")
	errContextDial = errors.New("context dial called")
)

type contextDialerStub struct{}

func (contextDialerStub) Dial(network, addr string) (net.Conn, error) {
	return nil, errPlainDial
}

func (contextDialerStub) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return nil, errContextDial
}

type plainDialerStub struct{}

func (plainDialerStub) Dial(network, addr string) (net.Conn, error) {
	return nil, errPlainDial
}

func TestDialContextUsesUnderlyingContextDialer(t *testing.T) {
	dialer := &Dialer{Dialer: contextDialerStub{}}

	_, err := dialer.DialContext(context.Background(), "tcp", "example.com:80")
	if !errors.Is(err, errContextDial) {
		t.Fatalf("DialContext error = %v, want %v", err, errContextDial)
	}
}

func TestDialContextFallsBackToPlainDialer(t *testing.T) {
	dialer := &Dialer{Dialer: plainDialerStub{}}

	_, err := dialer.DialContext(context.Background(), "tcp", "example.com:80")
	if !errors.Is(err, errPlainDial) {
		t.Fatalf("DialContext error = %v, want %v", err, errPlainDial)
	}
}
