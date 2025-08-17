package queue

import "errors"

var (
	ErrOverflow  = errors.New("queue overflow")
	ErrUnderflow = errors.New("queue underflow")
)
