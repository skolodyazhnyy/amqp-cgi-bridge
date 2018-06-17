package bridge

import (
	"context"
	"testing"
)

// ProcessorWithEnv should add additional headers to request
func TestProcessorWithEnv_Merge(t *testing.T) {
	done := make(chan struct{})

	p := func(c context.Context, h map[string]string, b []byte) error {
		if x := h["foo"]; x != "bar" {
			t.Errorf("Environment variables are not injected")
		}

		close(done)

		return nil
	}

	p = ProcessorWithEnv(p, map[string]string{"foo": "bar"})
	p(context.Background(), nil, nil)

	select {
	case <-done:
	default:
		t.Errorf("Inner processor was not executed")
	}
}

// ProcessWithEnv should not override any headers of the request
func TestProcessorWithEnv_Override(t *testing.T) {
	done := make(chan struct{})

	p := func(c context.Context, h map[string]string, b []byte) error {
		if x := h["foo"]; x != "bar" {
			t.Errorf("Environment variables are not injected")
		}

		close(done)

		return nil
	}

	p = ProcessorWithEnv(p, map[string]string{"foo": "overriden"})
	p(context.Background(), map[string]string{"foo": "bar"}, nil)

	select {
	case <-done:
	default:
		t.Errorf("Inner processor was not executed")
	}
}
