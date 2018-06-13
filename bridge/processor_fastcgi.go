package bridge

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tomasen/fcgi_client"
	"strings"
)

func NewFastCGIProcessor(net, addr, script string, log logger) Processor {
	return func(ctx context.Context, headers map[string]interface{}, body []byte) error {
		conn, err := fcgiclient.Dial(net, addr)
		if err != nil {
			log.Errorf("Unable to connect to FastCGI server: %v", err)
			return ErrProcessorInternal
		}

		env := map[string]string{
			"SERVER_PROTOCOL": "AMQP/0.9",
			"SERVER_SOFTWARE": "AMQP CGI Bridge",
			"REQUEST_METHOD":  "POST",
			"REMOTE_ADDR":     "127.0.0.1",
		}

		for k, v := range headers {
			env[strings.ToUpper(k)] = fmt.Sprint(v)
		}

		env["SCRIPT_FILENAME"] = script

		resp, err := conn.Request(env, bytes.NewReader(body))
		if err != nil {
			log.Errorf("An error occurred while making FastCGI request: %v", err)
			return ErrProcessorInternal
		}

		defer resp.Body.Close()

		if resp.StatusCode == 0 {
			return ErrUnknownStatus
		}

		if resp.StatusCode/100 == 2 {
			return nil
		}

		if resp.StatusCode/100 == 3 || resp.StatusCode/100 == 4 {
			log.Errorf("Request to FastCGI server has returned %v status code which probably means request configuration is invalid", resp.StatusCode)
			return ErrRequestFailed
		}

		return ErrProcessingFailed
	}
}
