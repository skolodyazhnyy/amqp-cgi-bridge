package bridge

import (
	"context"
	"os"
	"testing"
)

const TestScript = "/amqp-cgi-bridge/processor_fastcgi_test.php"

func TestFastCGIProcessor_Accept(t *testing.T) {
	addr := os.Getenv("TEST_PHPFPM_ADDR")
	if addr == "" {
		t.Skip("This test requires PHP-FPM server, use environment variable TEST_PHPFPM_ADDR to set PHP-FPM address.")
	}

	p := NewFastCGIProcessor("tcp", addr, TestScript, &nilLogger{})

	if err := p(context.Background(), map[string]string{"TEST": "ACCEPT"}, nil); err != nil {
		t.Fatalf("An error occurred while processing request: %v", err)
	}
}

func TestFastCGIProcessor_BodySize(t *testing.T) {
	addr := os.Getenv("TEST_PHPFPM_ADDR")
	if addr == "" {
		t.Skip("This test requires PHP-FPM server, use environment variable TEST_PHPFPM_ADDR to set PHP-FPM address.")
	}

	p := NewFastCGIProcessor("tcp", addr, TestScript, &nilLogger{})

	err := p(context.Background(), map[string]string{"TEST": "BODYSIZE10"}, []byte("1234567890"))
	if err == ErrRequestFailed {
		t.Fatalf("It seems body size does not match expected value, check PHP-FPM logs for more details")
	}

	if err != nil {
		t.Fatalf("An error occurred while processing request: %v", err)
	}
}

func TestFastCGIProcessor_EnvVariables(t *testing.T) {
	addr := os.Getenv("TEST_PHPFPM_ADDR")
	if addr == "" {
		t.Skip("This test requires PHP-FPM server, use environment variable TEST_PHPFPM_ADDR to set PHP-FPM address.")
	}

	p := NewFastCGIProcessor("tcp", addr, TestScript, &nilLogger{})

	err := p(context.Background(), map[string]string{"TEST": "ENVVAR", "HTTP_FOO": "BAR"}, nil)
	if err == ErrRequestFailed {
		t.Fatalf("It seems environment variables are not passed to PHP script, check PHP-FPM logs for more details")
	}

	if err != nil {
		t.Fatalf("An error occurred while processing request: %v", err)
	}
}

func TestFastCGIProcessor_InternalError(t *testing.T) {
	p := NewFastCGIProcessor("tcp", "0.0.0.0:0", TestScript, &nilLogger{})

	err := p(context.Background(), nil, nil)
	if err != ErrProcessorInternal {
		t.Fatalf("Dialing invalid network address should cause ErrProcessorInternal, got %v instead", err)
	}
}
