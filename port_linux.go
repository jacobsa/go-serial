package serial

import (
	"os"
	"time"
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

func (p *Port) InWaiting() int {
	// Funky time
	return 0
}

func (p *Port) SetDeadline(t time.Time) error {
	// Funky Town
	return nil
}

func NewPort(f *os.File) *Port {
	return &Port{f}
}
