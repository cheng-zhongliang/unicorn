package unicorn

import (
	"errors"
	"syscall"
)

var (
	ErrEventExists              = errors.New("event already exists")
	ErrEventNotExists           = errors.New("event does not exist")
	ErrInvalidDemultiplexerType = errors.New("invalid demultiplexer type")
	ErrQueueFull                = errors.New("queue is full")
	ErrQueueEmpty               = errors.New("queue is empty")
)

func TemporaryErr(err error) bool {
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	return errno.Temporary()
}
