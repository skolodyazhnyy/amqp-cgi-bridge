package bridge

import (
	"testing"
	"time"
)

func TestBackOffTimeout(t *testing.T) {
	tests := []struct {
		name    string
		last    time.Time     // last time timeout was requested
		delay   time.Duration // last timeout value
		timeout time.Duration // expected timeout value
	}{
		{name: "initial timeout", last: time.Unix(0, 0), delay: 0, timeout: time.Second},
		{name: "delay reset", last: time.Now().Add(-71 * time.Second), delay: time.Minute, timeout: time.Second},
		{name: "delay increase", last: time.Now(), delay: time.Second, timeout: 2 * time.Second},
		{name: "delay max", last: time.Now(), delay: 55 * time.Second, timeout: time.Minute},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := &backOff{delay: test.delay, last: test.last}
			d := b.Timeout()

			if d != test.timeout {
				t.Errorf("BackOff timeout does not match expected value: want %v, got %v", test.timeout, d)
			}
		})
	}
}
