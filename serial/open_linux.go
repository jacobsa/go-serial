package serial

import (
	"errors"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

// NCCS is the number of control character sequences used for c_cc
const (
	nccs = 19
)

// Types from asm-generic/termbits.h
type cc_t byte
type speed_t uint32
type tcflag_t uint32
type termios2 struct {
	c_iflag  tcflag_t   // input mode flags
	c_oflag  tcflag_t   // output mode flags
	c_cflag  tcflag_t   // control mode flags
	c_lflag  tcflag_t   // local mode flags
	c_line   cc_t       // line discipline
	c_cc     [nccs]cc_t // control characters
	c_ispeed speed_t    // input speed
	c_ospeed speed_t    // output speed
}

// makeTermios2 returns a pointer to an instantiates termios2 struct, based on the given
// OpenOptions. Termios2 is a Linux extension which allows arbitrary baud rates
// to be specified.
func makeTermios2(options OpenOptions) (*termios2, error) {

	// Sanity check inter-character timeout and minimum read size options.
	// See serial.go for more information on vtime/vmin -- these only work in non-canonical mode
	vtime := uint(round(float64(options.InterCharacterTimeout)/100.0) * 100)
	vmin := options.MinimumReadSize

	if vmin == 0 && vtime < 100 {
		return nil, errors.New("invalid values for InterCharacterTimeout and MinimumReadSize")
	}

	if vtime > 25500 {
		return nil, errors.New("invalid value for InterCharacterTimeout")
	}

	ccOpts := [nccs]cc_t{}
	ccOpts[unix.VTIME] = cc_t(vtime / 100)
	ccOpts[unix.VMIN] = cc_t(vmin)

	// We set the flags for CLOCAL, CREAD and BOTHER
	// CLOCAL : ignore modem control lines
	// CREAD  : enable receiver
	// BOTHER : allow generic BAUDRATE values
	t2 := &termios2{
		c_cflag:  unix.CLOCAL | unix.CREAD | unix.BOTHER,
		c_ispeed: speed_t(options.BaudRate),
		c_ospeed: speed_t(options.BaudRate),
		c_cc:     ccOpts,
	}

	// Un-set the ICANON mode to allow non-canonical mode
	// See: https://www.gnu.org/software/libc/manual/html_node/Canonical-or-Not.html
	if !options.CanonicalMode {
		t2.c_lflag &= ^tcflag_t(unix.ICANON)
	}

	// Allow for setting 1 or 2 stop bits
	switch options.StopBits {
	case 1:
	case 2:
		t2.c_cflag |= unix.CSTOPB

	default:
		return nil, errors.New("invalid setting for StopBits")
	}

	// If odd or even, enable parity generation (PARENB) and determine the type
	switch options.ParityMode {
	case Parity_None:
	case Parity_Odd:
		t2.c_cflag |= unix.PARENB
		t2.c_cflag |= unix.PARODD

	case Parity_Even:
		t2.c_cflag |= unix.PARENB

	default:
		return nil, errors.New("invalid setting for ParityMode")
	}

	// Choose the databits per frame
	switch options.DataBits {
	case 5:
		t2.c_cflag |= unix.CS5
	case 6:
		t2.c_cflag |= unix.CS6
	case 7:
		t2.c_cflag |= unix.CS7
	case 8:
		t2.c_cflag |= unix.CS8
	default:
		return nil, errors.New("invalid setting for DataBits")
	}

	return t2, nil
}

// openInternal is the operating system specific port opening, given the OpenOptions
func openInternal(options OpenOptions) (*Port, error) {
	// Open the file with RDWR, NOCTTY, NONBLOCK flags
	// RDWR     : read/write
	// NOCTTY   : don't allow the port to become the controlling terminal
	// NONBLOCK : open with nonblocking so we don't stall
	file, openErr :=
		os.OpenFile(
			options.PortName,
			unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK,
			0777)
	if openErr != nil {
		return nil, openErr
	}

	// Clear the non-blocking flag set above.
	nonblockErr := unix.SetNonblock(int(file.Fd()), false)
	if nonblockErr != nil {
		return nil, nonblockErr
	}

	t2, optErr := makeTermios2(options)
	if optErr != nil {
		return nil, optErr
	}

	// Set our termios2 struct as the file descriptor's settings
	errno := ioctl(unix.TCSETS2, file.Fd(), uintptr(unsafe.Pointer(t2)))
	if errno != nil {
		return nil, errno
	}

	return NewPort(file), nil
}
