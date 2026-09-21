package nanoleaf

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient points a Client at the given httptest server's address.
func newTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parsing server URL: %v", err)
	}
	host, portStr, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Fatalf("splitting host/port: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parsing port: %v", err)
	}

	return NewClient(srv.Client(), host, port, "test-key")
}

func writeStateOnEvent(w http.ResponseWriter, on bool) {
	fmt.Fprintf(w, "id: 1\ndata: {\"events\":[{\"attr\":1,\"Value\":%t}]}\n\n", on)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func TestSubscribeEvents_ReceivesUpdate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/events") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		writeStateOnEvent(w, true)
		<-r.Context().Done()
	}))
	defer srv.Close()

	c := newTestClient(t, srv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updates, errs := c.SubscribeEvents(ctx, 1)

	select {
	case u := <-updates:
		if u.TypeID != 1 {
			t.Fatalf("expected TypeID 1, got %d", u.TypeID)
		}
		if u.State == nil || u.State.On == nil || !u.State.On.Value {
			t.Fatalf("expected state.on = true, got %+v", u.State)
		}
	case err := <-errs:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for update")
	}
}

func TestSubscribeEvents_ReconnectsAfterDrop(t *testing.T) {
	var connCount int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&connCount, 1)
		w.Header().Set("Content-Type", "text/event-stream")
		writeStateOnEvent(w, n > 1)
		if n == 1 {
			// Drop the first connection immediately after one event, forcing a reconnect.
			return
		}
		<-r.Context().Done()
	}))
	defer srv.Close()

	c := newTestClient(t, srv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updates, errs := c.SubscribeEvents(ctx, 1)

	seen := 0
	deadline := time.After(10 * time.Second)
	for seen < 2 {
		select {
		case u := <-updates:
			seen++
			if seen == 2 && (u.State == nil || u.State.On == nil || !u.State.On.Value) {
				t.Fatalf("expected second connection's update to report on=true, got %+v", u.State)
			}
		case err := <-errs:
			t.Logf("got advisory error (expected around the reconnect): %v", err)
		case <-deadline:
			t.Fatalf("timed out after %d update(s), wanted 2 (reconnect never happened)", seen)
		}
	}

	if got := atomic.LoadInt32(&connCount); got < 2 {
		t.Fatalf("expected at least 2 connections, got %d", got)
	}
}

func TestSubscribeEvents_StopsOnContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		writeStateOnEvent(w, true)
		<-r.Context().Done()
	}))
	defer srv.Close()

	c := newTestClient(t, srv)

	ctx, cancel := context.WithCancel(context.Background())
	updates, errs := c.SubscribeEvents(ctx, 1)

	// Drain the first update so the goroutine is past its initial connect.
	select {
	case <-updates:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for initial update")
	}

	cancel()

	updatesClosed, errsClosed := false, false
	deadline := time.After(5 * time.Second)
	for !updatesClosed || !errsClosed {
		select {
		case _, ok := <-updates:
			if !ok {
				updatesClosed = true
			}
		case _, ok := <-errs:
			if !ok {
				errsClosed = true
			}
		case <-deadline:
			t.Fatalf("timed out waiting for channels to close (updatesClosed=%v errsClosed=%v)", updatesClosed, errsClosed)
		}
	}
}
