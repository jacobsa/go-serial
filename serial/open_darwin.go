package serial

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// NCCS is the number of control character sequences used for c_cc
const (
	IOSSIOSPEED = 0x80045402
)

// makeTermios2 returns a pointer to an instantiates termios2 struct, based on the given
// OpenOptions. Termios2 is a Linux extension which allows arbitrary baud rates
// to be specified.
func makeTermios2(fd uintptr, options OpenOptions) (*unix.Termios, error) {

	t := &unix.Termios{}
	// unix.IoctlGetTermios(int(fd), )
	// unix.TIOCGETA

	err := unix.IoctlSetTermios(int(fd), unix.TIOCGETA, t)
	if err != nil {
		fmt.Println("TCGETS openInternal err")
		return nil, err
	}

	t.Cflag |= (syscall.CLOCAL | syscall.CREAD)
	t.Lflag &= ^uint64(
		syscall.ICANON | syscall.ECHO | syscall.ECHOE |
			syscall.ECHOK | syscall.ECHONL |
			syscall.ISIG | syscall.IEXTEN)
	t.Lflag &= ^uint64(syscall.ECHOCTL)
	t.Lflag &= ^uint64(syscall.ECHOKE)

	t.Oflag &= ^uint64(syscall.OPOST | syscall.ONLCR | syscall.OCRNL)
	t.Iflag &= ^uint64(syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IGNBRK)
	t.Iflag &= ^uint64(syscall.PARMRK)

	// character size
	t.Cflag &= ^uint64(syscall.CSIZE)
	t.Cflag |= uint64(syscall.CS8)

	// setup stop bits
	t.Cflag &= ^uint64(syscall.CSTOPB)

	// setup parity
	t.Iflag &= ^uint64(syscall.INPCK | syscall.ISTRIP)
	t.Cflag &= ^uint64(syscall.PARENB | syscall.PARODD)

	t.Iflag &= ^uint64(syscall.IXON | syscall.IXOFF | syscall.IXANY)

	// t.Cflag &= ^uint64(syscall.CRTSCTS)

	// t.Cc[syscall.VMIN] = 0
	// t.Oflag &= ^uint64(syscall.OPOST | syscall.ONLCR | syscall.OCRNL)

	// t.Oflag &= c

	// fmt.Println("makeTermios2 TCGETS")
	// errno := ioctl(TCGETS, fd, uintptr(unsafe.Pointer(t)))
	// if errno != nil {
	// 	fmt.Println("TCGETS openInternal err")
	// 	return nil, errno
	// }
	//
	// fmt.Printf("%v", t)

	// // Sanity check inter-character timeout and minimum read size options.
	// // See serial.go for more information on vtime/vmin -- these only work in non-canonical mode
	vtime := uint(round(float64(options.InterCharacterTimeout)/100.0) * 100)
	vmin := options.MinimumReadSize
	//
	// if vmin == 0 && vtime < 100 {
	// 	return nil, errors.New("invalid values for InterCharacterTimeout and MinimumReadSize")
	// }
	//
	// if vtime > 25500 {
	// 	return nil, errors.New("invalid value for InterCharacterTimeout")
	// }
	//
	t.Cc[syscall.VTIME] = uint8(vtime / 100)
	t.Cc[unix.VMIN] = uint8(vmin)
	//
	// // We set the flags for CLOCAL, CREAD and BOTHER
	// // CLOCAL : ignore modem control lines
	// // CREAD  : enable receiver
	// // BOTHER : allow generic BAUDRATE values
	// t2 := &syscall.Termios{
	// 	Cflag:  unix.CLOCAL | unix.CREAD,
	// 	Ispeed: uint64(options.BaudRate),
	// 	Ospeed: uint64(options.BaudRate),
	// 	Cc:     ccOpts,
	// }
	//
	// // Un-set the ICANON mode to allow non-canonical mode
	// // See: https://www.gnu.org/software/libc/manual/html_node/Canonical-or-Not.html
	// if !options.CanonicalMode {
	// 	t2.Lflag &= ^uint64(unix.ICANON)
	// }
	//
	// // Allow for setting 1 or 2 stop bits
	// switch options.StopBits {
	// case 1:
	// case 2:
	// 	t2.Cflag |= unix.CSTOPB
	//
	// default:
	// 	return nil, errors.New("invalid setting for StopBits")
	// }
	//
	// // If odd or even, enable parity generation (PARENB) and determine the type
	// switch options.ParityMode {
	// case Parity_None:
	// case Parity_Odd:
	// 	t2.Cflag |= unix.PARENB
	// 	t2.Cflag |= unix.PARODD
	//
	// case Parity_Even:
	// 	t2.Cflag |= unix.PARENB
	//
	// default:
	// 	return nil, errors.New("invalid setting for ParityMode")
	// }
	//
	// // Choose the databits per frame
	// switch options.DataBits {
	// case 5:
	// 	t2.Cflag |= unix.CS5
	// case 6:
	// 	t2.Cflag |= unix.CS6
	// case 7:
	// 	t2.Cflag |= unix.CS7
	// case 8:
	// 	t2.Cflag |= unix.CS8
	// default:
	// 	return nil, errors.New("invalid setting for DataBits")
	// }
	//
	return t, nil
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

	fd := file.Fd()

	// When we call Fd(), we make the file descriptor blocking, which we don't want
	// Let's unset the blocking flag and save the pointer for later.
	nonblockErr := unix.SetNonblock(int(fd), true)
	if nonblockErr != nil {
		return nil, nonblockErr
	}

	t, optErr := makeTermios2(fd, options)
	if optErr != nil {
		return nil, optErr
	}

	// Set our termios2 struct as the file descriptor's settings
	err := unix.IoctlSetTermios(int(fd), unix.TIOCSETA, t)
	if err != nil {
		return nil, err
	}
	// 189  ->	                buf = array.array('i', [baudrate])
	b := uint(options.BaudRate)
	errcode := ioctl(IOSSIOSPEED, fd, uintptr(unsafe.Pointer(&b)))
	if errcode != nil {
		return nil, errcode
	}

	return NewPort(file, fd, options), nil
}
