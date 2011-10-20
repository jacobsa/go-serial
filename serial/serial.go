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

// Package serial provides routines for interacting with serial ports.
// Currently it supports only OS X; see the readme file for details.

package serial

import "io"
import "os"

// Valid parity values.
type ParityMode int

const (
	PARITY_NONE = 0
	PARITY_ODD  = 1
	PARITY_EVEN = 2
)

// OpenOptions is the struct containing all of the options necessary for
// opening a serial port.
type OpenOptions struct {
	// The name of the port, e.g. "/dev/tty.usbserial-A8008HlV".
	PortName string

	// The baud rate for the port. The following values are supported:
	//
	//    50
	//    75
	//    110
	//    134
	//    150
	//    200
	//    300
	//    600
	//    1200
	//    1800
	//    2400
	//    4800
	//    7200
	//    9600
	//    14400
	//    19200
	//    28800
	//    38400
	//    57600
	//    76800
	//    115200
	//    230400
	//
	BaudRate uint

	// The number of data bits per frame. Legal values are 5, 6, 7, and 8.
	DataBits uint

	// The number of stop bits per frame. Legal values are 1 and 2.
	StopBits uint

	// The type of parity bits to use for the connection. Currently parity errors
	// are simply ignored; that is, bytes are delivered to the user no matter
	// whether they were received with a parity error or not.
	ParityMode ParityMode
}

// Open creates an io.ReadWriteCloser based on the supplied options struct.
func Open(options OpenOptions) (io.ReadWriteCloser, os.Error) {
	// Redirect to the OS-specific function.
	return openInternal(options)
}
