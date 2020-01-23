package serial

import (
	"log"
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
	log.Println(p.f)
	log.Println("FIEL DESTERIPETOR ", p.f.Fd())
	log.Println("waiting before: ", waiting)
	a, b, err := unix.Syscall(unix.SYS_IOCTL, p.f.Fd(), unix.TIOCINQ, uintptr(unsafe.Pointer(&waiting)))
	log.Println("got a: ", a, "\ngot b: ", b, "\ngot err: ", err)
	log.Println("waiting after: ", waiting)
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
