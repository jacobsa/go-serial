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

func (p *Port) Read(b []byte) (int, error) {
	return p.f.Read(b)
}

func (p *Port) Write(b []byte) (int, error) {
	return p.f.Write(b)
}

func (p *Port) Close() error {
	return p.f.Close()
}

func (p *Port) InWaiting() (int, error) {
	// Funky time
	var waiting int
	_, _, err := unix.Syscall(unix.SYS_IOCTL, p.f.Fd(), unix.TIOCINQ, uintptr(unsafe.Pointer(&waiting)))
	if err != 0 {
		return 0, err
	}
	return waiting, nil
}

func (p *Port) SetDeadline(t time.Time) error {
	// Funky Town
	return nil
}

func NewPort(f *os.File) *Port {
	return &Port{f}
}
