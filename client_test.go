package nanoleaf

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGet_StatusCodeMapping exercises get()'s HTTP status -> error mapping (via GetPanel, the
// simplest get() caller) against a synthetic server. No live panel is needed for this - it's
// pure wire-protocol behaviour.
func TestGet_StatusCodeMapping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    error
	}{
		{name: "ok", statusCode: http.StatusOK, body: `{"name":"Test Panel"}`},
		{name: "bad request", statusCode: http.StatusBadRequest, wantErr: ErrBadRequest},
		{name: "unauthorized", statusCode: http.StatusUnauthorized, wantErr: ErrUnauthorized},
		{name: "not found", statusCode: http.StatusNotFound, wantErr: ErrNotFound},
		{name: "too many requests", statusCode: http.StatusTooManyRequests, wantErr: ErrTooManyRequests},
		{name: "unexpected status", statusCode: http.StatusInternalServerError, wantErr: ErrUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.body != "" {
					w.Write([]byte(tt.body))
				}
			}))
			defer srv.Close()

			c := newTestClient(t, srv)
			panel, err := c.GetPanel(context.Background())

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if panel.Name != "Test Panel" {
				t.Fatalf("expected decoded panel name %q, got %q", "Test Panel", panel.Name)
			}
		})
	}
}

// TestPut_StatusCodeMapping exercises put()'s HTTP status -> error mapping via SetOn, whose real
// panels respond 204 (not 200) to - SetOn passes a nil respType, which only put()'s 204 branch
// (and every error branch, which returns before touching respType) is compatible with.
func TestPut_StatusCodeMapping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    error
	}{
		{name: "no content", statusCode: http.StatusNoContent},
		{name: "bad request", statusCode: http.StatusBadRequest, wantErr: ErrBadRequest},
		{name: "unauthorized", statusCode: http.StatusUnauthorized, wantErr: ErrUnauthorized},
		{name: "not found", statusCode: http.StatusNotFound, wantErr: ErrNotFound},
		{name: "unprocessable entity", statusCode: http.StatusUnprocessableEntity, wantErr: ErrBadRequest},
		{name: "too many requests", statusCode: http.StatusTooManyRequests, wantErr: ErrTooManyRequests},
		{name: "unexpected status", statusCode: http.StatusInternalServerError, wantErr: ErrUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("expected PUT, got %s", r.Method)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer srv.Close()

			c := newTestClient(t, srv)
			err := c.SetOn(context.Background(), true)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestPut_DecodesSuccessBody covers put()'s 200 branch, which decodes the response body into the
// caller-supplied respType - GetEffect is the put() caller that actually uses one (SetOn and
// friends pass nil, since real panels never answer their PUTs with a 200+body).
func TestPut_DecodesSuccessBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"animName":"Nemo"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	effect, err := c.GetEffect(context.Background(), "Nemo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if effect.Name != "Nemo" {
		t.Fatalf("expected decoded effect name %q, got %q", "Nemo", effect.Name)
	}
}
