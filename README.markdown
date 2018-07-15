go-serial
=========

This is a package that allows you to read from and write to serial ports in Go.


OS support
----------

Currently this package works only on OS X, Linux and Windows. It could probably be ported
to other Unix-like platforms simply by updating a few constants; get in touch if
you are interested in helping and have hardware to test with.


Installation
------------

Simply use `go get`:

    go get github.com/philipgreat/go-serial/serial

To update later:

    go get -u github.com/philipgreat/go-serial/serial


Use
---

Set up a `serial.OpenOptions` struct, then call `serial.Open`. For example:

````go
    package main

//env GOOS=linux GOARCH=mips go build -ldflags "-s -w" mem.go
import (
	"fmt"
	"log"

	"github.com/philipgreat/go-serial/serial"
)

func main() {
	// Below is an example of using our PrintMemUsage() function
	// Print our starting memory usage (should be around 0mb)
	options := serial.OpenOptions{
		PortName:        "/dev/ttyUSB0",
		BaudRate:        38400,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer port.Close()
	b := []byte{0x28, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00,
		0x00, 0x00, 0x80, 0x05, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x1F, 0x00, 0x08, 0x12, 0x21, 0x00, 0x00, 0x64, 0x00, 0x64, 0x00,
		0x64, 0x00, 0x64, 0x00, 0x64, 0x00, 0xFC, 0x63}
	n, err := port.Write(b)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}

	fmt.Println("Wrote", n, "bytes.")
	buf := make([]byte, 128)
	n, err = port.Read(buf)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}

}

````

See the documentation for the `OpenOptions` struct in `serial.go` for more
information on the supported options.
