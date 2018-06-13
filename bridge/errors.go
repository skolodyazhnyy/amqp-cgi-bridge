package bridge

import "errors"

var ErrProcessorInternal = errors.New("processor was not able to perform request")
var ErrUnknownStatus = errors.New("processor was not able to read response status code")
var ErrRequestFailed = errors.New("request to processing backend has failed")
var ErrProcessingFailed = errors.New("message processing failed")
