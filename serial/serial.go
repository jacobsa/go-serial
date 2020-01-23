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

import (
	"math"
)

// Valid parity values.
type ParityMode int

const (
	Parity_None ParityMode = 0
	Parity_Odd  ParityMode = 1
	Parity_Even ParityMode = 2
)

// OpenOptions is the struct containing all of the options necessary for
// opening a serial port.
type OpenOptions struct {
	// The name of the port, e.g. "/dev/tty.usbserial-A8008HlV".
	PortName string

	// The baud rate for the port.
	BaudRate uint

	// The number of data bits per frame. Legal values are 5, 6, 7, and 8.
	DataBits uint

	// The number of stop bits per frame. Legal values are 1 and 2.
	StopBits uint

	// The type of parity bits to use for the connection. Currently parity errors
	// are simply ignored; that is, bytes are delivered to the user no matter
	// whether they were received with a parity error or not.
	ParityMode ParityMode

	// Enable RTS/CTS (hardware) flow control.
	RTSCTSFlowControl bool

	CanonicalMode bool

	// An inter-character timeout value, in milliseconds, and a minimum number of
	// bytes to block for on each read. A call to Read() that otherwise may block
	// waiting for more data will return immediately if the specified amount of
	// time elapses between successive bytes received from the device or if the
	// minimum number of bytes has been exceeded.
	//
	// Note that the inter-character timeout value may be rounded to the nearest
	// 100 ms on some systems, and that behavior is undefined if calls to Read
	// supply a buffer whose length is less than the minimum read size.
	//
	// Behaviors for various settings for these values are described below. For
	// more information, see the discussion of VMIN and VTIME here:
	//
	//     http://www.unixwiz.net/techtips/termios-vmin-vtime.html

	// For windows usage, these options (termios) do not conform well to the
	//     windows serial port / comms abstractions.  Please see the code in
	//		 open_windows setCommTimeouts function for full documentation.
	//   	 Summary:
	//			Setting MinimumReadSize > 0 will cause the serialPort to block until
	//			until data is available on the port.
	//			Setting IntercharacterTimeout > 0 and MinimumReadSize == 0 will cause
	//			the port to either wait until IntercharacterTimeout wait time is
	//			exceeded OR there is character data to return from the port.

	InterCharacterTimeout uint
	MinimumReadSize       uint
}

// Open creates an io.ReadWriteCloser based on the supplied options struct.
func Open(options OpenOptions) (*Port, error) {
	// Redirect to the OS-specific function.
	return openInternal(options)
}

// Rounds a float to the nearest integer.
func round(f float64) float64 {
	return math.Floor(f + 0.5)
}
