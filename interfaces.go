package serial

import (
	"io"
	"time"
)

type port interface {
	io.ReadWriteCloser
	Inwaiting() (int, error)
	SetDeadline(time.Time) error
}
