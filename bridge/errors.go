package bridge

import "errors"

var ErrProcessorInternal = errors.New("processor was not able to perform request")
var ErrUnknownStatus = errors.New("processor was not able to read response status code")
var ErrProcessingError = errors.New("request to processing backend has failed (response status code 3xx or 4xx)")
var ErrProcessingFailed = errors.New("message processing failed (response status code 5xx)")
