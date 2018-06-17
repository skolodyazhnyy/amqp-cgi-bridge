package bridge

import "time"

const maxBackOffTimeout = time.Minute

type backOff struct {
	last  time.Time
	delay time.Duration
}

// Timeout returns backoff delay
func (b *backOff) Timeout() time.Duration {
	// reset delay if last run took longer than 10 seconds (without backoff delay)
	if time.Since(b.last)-b.delay > 10*time.Second {
		b.delay = 0
	}

	b.last = time.Now()
	b.delay = b.increase(b.delay)

	return b.delay
}

// increase backoff delay
func (b *backOff) increase(d time.Duration) time.Duration {
	if d <= 0 {
		return time.Second
	}

	if d*2 > maxBackOffTimeout {
		return maxBackOffTimeout
	}

	return d * 2
}
