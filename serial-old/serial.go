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
	"io"
	"math"
)

// Valid parity values.
type ParityMode int

const (
	PARITY_NONE ParityMode = 0
	PARITY_ODD  ParityMode = 1
	PARITY_EVEN ParityMode = 2
)

var (
	// The list of standard baud-rates.
	StandardBaudRates = map[uint]bool{
		50:     true,
		75:     true,
		110:    true,
		134:    true,
		150:    true,
		200:    true,
		300:    true,
		600:    true,
		1200:   true,
		1800:   true,
		2400:   true,
		4800:   true,
		7200:   true,
		9600:   true,
		14400:  true,
		19200:  true,
		28800:  true,
		38400:  true,
		57600:  true,
		76800:  true,
		115200: true,
		230400: true,
	}
)

// IsStandardBaudRate checks whether the specified baud-rate is standard.
//
// Some operating systems may support non-standard baud-rates (OSX) via
// additional IOCTL.
func IsStandardBaudRate(baudRate uint) bool { return StandardBaudRates[baudRate] }

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
	//
	// InterCharacterTimeout = 0 and MinimumReadSize = 0 (the default):
	//     This arrangement is not legal; you must explicitly set at least one of
	//     these fields to a positive number. (If MinimumReadSize is zero then
	//     InterCharacterTimeout must be at least 100.)
	//
	// InterCharacterTimeout > 0 and MinimumReadSize = 0
	//     If data is already available on the read queue, it is transferred to
	//     the caller's buffer and the Read() call returns immediately.
	//     Otherwise, the call blocks until some data arrives or the
	//     InterCharacterTimeout milliseconds elapse from the start of the call.
	//     Note that in this configuration, InterCharacterTimeout must be at
	//     least 100 ms.
	//
	// InterCharacterTimeout > 0 and MinimumReadSize > 0
	//     Calls to Read() return when at least MinimumReadSize bytes are
	//     available or when InterCharacterTimeout milliseconds elapse between
	//     received bytes. The inter-character timer is not started until the
	//     first byte arrives.
	//
	// InterCharacterTimeout = 0 and MinimumReadSize > 0
	//     Calls to Read() return only when at least MinimumReadSize bytes are
	//     available. The inter-character timer is not used.
	//
	// For windows usage, these options (termios) do not conform well to the
	//     windows serial port / comms abstractions.  Please see the code in
	//		 open_windows setCommTimeouts function for full documentation.
	//   	 Summary:
	//			Setting MinimumReadSize > 0 will cause the serialPort to block until
	//			until data is available on the port.
	//			Setting IntercharacterTimeout > 0 and MinimumReadSize == 0 will cause
	//			the port to either wait until IntercharacterTimeout wait time is
	//			exceeded OR there is character data to return from the port.
	//

	InterCharacterTimeout uint
	MinimumReadSize       uint

	// Use to enable RS485 mode -- probably only valid on some Linux platforms
	Rs485Enable bool

	// Set to true for logic level high during send
	Rs485RtsHighDuringSend bool

	// Set to true for logic level high after send
	Rs485RtsHighAfterSend bool

	// set to receive data during sending
	Rs485RxDuringTx bool

	// RTS delay before send
	Rs485DelayRtsBeforeSend int

	// RTS delay after send
	Rs485DelayRtsAfterSend int
}

// Open creates an io.ReadWriteCloser based on the supplied options struct.
func Open(options OpenOptions) (io.ReadWriteCloser, error) {
	// Redirect to the OS-specific function.
	return openInternal(options)
}

// Rounds a float to the nearest integer.
func round(f float64) float64 {
	return math.Floor(f + 0.5)
}
