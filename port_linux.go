package serial

import "os"

type port struct {
	f os.File
}

func (p *port) Read(b []byte) (int, error) {
	return p.f.Read(b)
}

func (p *port) Write(b []byte) error {
	return p.f.Write(b)
}

func (p *port) Close() error {
	return p.f.Close()
}

func (p *port) InWaiting() int {
	// Funky time
}

funct (p *port) SetDeadline(t time.Time) error {
	// Funky Town
}