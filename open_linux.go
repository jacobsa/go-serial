package serial

import (
	"errors"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	NCCS = 19
)

//
// Types from asm-generic/termbits.h
//

type cc_t byte
type speed_t uint32
type tcflag_t uint32
type termios2 struct {
	c_iflag  tcflag_t   // input mode flags
	c_oflag  tcflag_t   // output mode flags
	c_cflag  tcflag_t   // control mode flags
	c_lflag  tcflag_t   // local mode flags
	c_line   cc_t       // line discipline
	c_cc     [NCCS]cc_t // control characters
	c_ispeed speed_t    // input speed
	c_ospeed speed_t    // output speed
}

//
// Returns a pointer to an instantiates termios2 struct, based on the given
// OpenOptions. Termios2 is a Linux extension which allows arbitrary baud rates
// to be specified.
//
func makeTermios2(options OpenOptions) (*termios2, error) {

	// Sanity check inter-character timeout and minimum read size options.

	vtime := uint(round(float64(options.InterCharacterTimeout)/100.0) * 100)
	vmin := options.MinimumReadSize

	if vmin == 0 && vtime < 100 {
		return nil, errors.New("invalid values for InterCharacterTimeout and MinimumReadSize")
	}

	if vtime > 25500 {
		return nil, errors.New("invalid value for InterCharacterTimeout")
	}

	ccOpts := [NCCS]cc_t{}
	ccOpts[unix.VTIME] = cc_t(vtime / 100)
	ccOpts[unix.VMIN] = cc_t(vmin)

	t2 := &termios2{
		c_cflag:  unix.CLOCAL | unix.CREAD | unix.BOTHER,
		c_ispeed: speed_t(options.BaudRate),
		c_ospeed: speed_t(options.BaudRate),
		c_cc:     ccOpts,
	}

	switch options.StopBits {
	case 1:
	case 2:
		t2.c_cflag |= unix.CSTOPB

	default:
		return nil, errors.New("invalid setting for StopBits")
	}

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

func openInternal(options OpenOptions) (*Port, error) {
	file, openErr :=
		os.OpenFile(
			options.PortName,
			unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK,
			0600)
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

	r, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(file.Fd()),
		uintptr(unix.TCSETS2),
		uintptr(unsafe.Pointer(t2)))

	if errno != 0 {
		return nil, os.NewSyscallError("SYS_IOCTL", errno)
	}

	if r != 0 {
		return nil, errors.New("unknown error from SYS_IOCTL")
	}

	return NewPort(file), nil
}
