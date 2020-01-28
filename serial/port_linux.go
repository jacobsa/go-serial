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

func (p *Port) SetTimeout(t time.Time) error {
	// Funky Town
	// todo(ahollist): Implement
	return nil
}

// Get the port's DTR pin state
func (p *Port) DTR() (bool, error) {
	var status int
	_, _, err := unix.Syscall(unix.SYS_IOCTL, p.f.Fd(), unix.TIOCMGET, uintptr(unsafe.Pointer(&status)))
	if err != 0 {
		return false, err
	}
	if status&unix.TIOCM_DTR > 0 {
		return true, nil
	}
	return false, nil
}

// Set the port's DTR pin state
func (p *Port) SetDTR(state bool) error {
	// todo(ahollist): Implement
	var command int
	if state {
		command = unix.TIOCMBIS
	} else {
		command = unix.TIOCMBIC
	}
	_, _, err := unix.Syscall(unix.SYS_IOCTL, p.f.Fd(), uintptr(command), 0)
	if err != 0 {
		return err
	}
	return nil
}

func NewPort(f *os.File) *Port {
	return &Port{f}
}
