package asr

import (
	"context"
	"errors"
	"net"
)

func isNetTimeout(err error) bool {
	var nerr net.Error
	if errors.As(err, &nerr) {
		return nerr.Timeout()
	}
	return errors.Is(err, context.DeadlineExceeded)
}
