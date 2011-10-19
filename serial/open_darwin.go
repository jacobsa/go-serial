// Copyright 2011 Aaron Jacobs. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file contains OS-specific constants and types that work on OS X (tested
// on version 10.6.8).

package serial

import "io"
import "os"
import "unsafe"

// termios types
type cc_t byte
type speed_t uint64
type tcflag_t uint64

// sys/termios.h
const (
	B9600 = 9600
	B14400 = 14400
	B19200 = 19200

	CS5 = 0x00000000
	CS6 = 0x00000100
	CS7 = 0x00000200
	CS8 = 0x00000300
	CLOCAL = 0x00008000
	CREAD = 0x00000800
	IGNPAR = 0x00000004

	NCCS = 20;

	VMIN = tcflag_t(16);
	VTIME = tcflag_t(17);
)

// sys/ttycom.h
const (
	TIOCGETA = 1078490131
	TIOCSETA = 2152231956
)

// sys/termios.h
type termios struct {
	c_iflag tcflag_t
	c_oflag tcflag_t
	c_cflag tcflag_t
	c_lflag tcflag_t
	c_cc [NCCS]cc_t
	c_ispeed speed_t
	c_ospeed speed_t
}

// setTermios updates the termios struct associated with a serial port file
// descriptor. This sets appropriate options for how the OS interacts with the
// port.
func setTermios(fd syscall.Handle, src termios) os.Error {
	// Make the ioctl syscall that sets the termios struct.
	r1, _, errno :=
		syscall.Syscall(
			syscall.SYS_IOCTL,
			uintptr(fd),
			uintptr(TIOCSETA),
			uintptr(unsafe.Pointer(&src)))

	// Did the syscall return an error?
	if err := os.NewSyscallError("SYS_IOCTL", int(errno)); err != nil {
		return err
	}

	// Just in case, check the return value as well.
	if r1 != 0 {
		return os.NewError("Unknown error from SYS_IOCTL.")
	}

	return nil
}

func openInternal(options OpenOptions) (io.ReadWriteCloser, os.Error) {
	// Open the serial port in non-blocking mode, since that seems to be required
	// for OS X for some reason (otherwise it just blocks forever).
	file, err :=
		os.OpenFile(
			options.PortName,
			os.O_RDWR | os.O_NOCTTY | os.O_NONBLOCK,
			0600)

	if err != nil {
		return nil, err
	}

	// We want to do blocking I/O, so clear the non-blocking flag set above.
	r1, _, errno :=
		syscall.Syscall(
			syscall.SYS_FCNTL,
			uintptr(file.Fd()),
			uintptr(syscall.F_SETFL),
			uintptr(0))

	if err := os.NewSyscallError("SYS_IOCTL", int(errno)); err != nil {
		return err
	}

	if r1 != 0 {
		return os.NewError("Unknown error from SYS_FCNTL.")
	}

	// Set appropriate options.
	terminalOptions := convertOptions(options)

	err = setTermios(file.Fd(), terminalOptions)
	if err != nil {
		return nil, err
	}

	// We're done.
	return file, nil
}

