This is a package that allows you to read from and write to serial ports in Go.


OS support
==========

Currently this package works only on OS X (tested on version 10.6). It could
probably be ported to Linux simply by updating a few constants; get in touch if
you are interested in helping and have hardware to test with. Windows would
likely be a lot more work.


Installation
============

Simply use `goinstall`:

    goinstall github.com/jacobsa/go-serial/serial

To update later:

    goinstall -u github.com/jacobsa/go-serial/serial


Use
===

Set up a `serial.OpenOptions` struct, then call `serial.Open`. For example:

    import "fmt"
    import "github.com/jacobsa/go-serial/serial"

    ...

    // Set up options.
    var options serial.OpenOptions
    options.PortName = "/dev/tty.usbserial-A8008HlV"
    options.BaudRate = 19200
    options.DataBits = 8
    options.StopBits = 1
    options.MinimumReadSize = 4

    // Open the port.
    port, err := serial.Open(options)
    if err != nil {
      panic("serial.Open: " + err.String())
    }

    // Make sure to close it later.
    defer port.Close()

    // Write 4 bytes to the port.
    b := []byte{0x00, 0x01, 0x02, 0x03}
    n, err := port.Write(b)
    if err !+ nil {
      panic("port.Write: " + err.String())
    }

    fmt.Println("Wrote", n, "bytes.")

See the documentation for the `OpenOptions` struct in `serial.go` for more
information on the supported options.
