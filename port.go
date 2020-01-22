package serial

import (
	"io"
	"time"
)

type Port interface {
	io.ReadWriteCloser
	Inwaiting() int
	SetDeadline(time.Time) error
}
