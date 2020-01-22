package serial

import (
	"io"
	"time"
)

type port interface {
	io.ReadWriteCloser
	Inwaiting() int
	SetDeadline(time.Time) error
}
