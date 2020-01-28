package serial

import (
	"io"
	"time"
)

// Structs that implement these methods are considered ports
type port interface {
	io.ReadWriteCloser
	Inwaiting() (int, error)
	SetDeadline(time.Time) error
}
