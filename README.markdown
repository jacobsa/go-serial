This is a package that allows you to read from and write to serial ports in Go.


OS support
==========

Currently this package works only on OS X, Linux and Windows. It could probably be ported
to other Unix-like platforms simply by updating a few constants; get in touch if
you are interested in helping and have hardware to test with. Windows would
likely be a lot more work.


Installation
============

Simply use `go install`:

    go install github.com/jacobsa/go-serial/serial

To update later:

    go install -u github.com/jacobsa/go-serial/serial


Use
===

Set up a `serial.OpenOptions` struct, then call `serial.Open`. For example:

````go
    import "fmt"
    import "log"
    import "github.com/jacobsa/go-serial/serial"

    ...

    // Set up options.
    options := serial.OpenOptions{
      PortName: "/dev/tty.usbserial-A8008HlV",
      BaudRate: 19200,
      DataBits: 8,
      StopBits: 1,
      MinimumReadSize: 4,
    }

    // Open the port.
    port, err := serial.Open(options)
    if err != nil {
      log.Fatalf("serial.Open: %v", err)
    }

    // Make sure to close it later.
    defer port.Close()

    // Write 4 bytes to the port.
    b := []byte{0x00, 0x01, 0x02, 0x03}
    n, err := port.Write(b)
    if err != nil {
      log.Fatalf("port.Write: %v", err)
    }

    fmt.Println("Wrote", n, "bytes.")
````

See the documentation for the `OpenOptions` struct in `serial.go` for more
information on the supported options.
