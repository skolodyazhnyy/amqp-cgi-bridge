package bridge

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func NewExecProcessor(cmd string, args ...string) Processor {
	return func(ctx context.Context, headers map[string]interface{}, body []byte) error {
		c := exec.CommandContext(ctx, cmd, args...)
		c.Stdin = bytes.NewReader(body)
		c.Stderr = ioutil.Discard
		c.Stdout = ioutil.Discard
		c.Env = make([]string, 0, len(headers))

		for k, v := range headers {
			c.Env = append(c.Env, fmt.Sprintf("%v=%v", k, v))
		}

		return c.Run()
	}
}
