package serial

import (
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type Port struct {
	f *os.File
}

// Pass directly through to the file pointer and read the data stream
func (p *Port) Read(b []byte) (int, error) {
	return p.f.Read(b)
}

// Pass directly through to the file pointer and write to the stream
func (p *Port) Write(b []byte) (int, error) {
	return p.f.Write(b)
}

// Close the file in our Port
func (p *Port) Close() error {
	return p.f.Close()
}

// Return the number of bytes waiting in the stream, using ioctl
func (p *Port) InWaiting() (int, error) {
	// Funky time
	var waiting int
	err := ioctl(unix.TIOCINQ, p.f.Fd(), uintptr(unsafe.Pointer(&waiting)))
	if err != nil {
		return 0, err
	}
	return waiting, nil
}

func (p *Port) SetDeadline(t time.Time) error {
	// Funky Town
	err := p.f.SetDeadline(t)
	if err != nil {
		return err
	}
	return nil
}

func NewPort(f *os.File) *Port {
	return &Port{f}
}
