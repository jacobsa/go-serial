package main

import (
	"fmt"
	"log"
	serial "simple-go-serial/serial"
)

func main() {
	ser, error := serial.Open(serial.OpenOptions{
		PortName:              "/dev/ttyUSB0",
		BaudRate:              3000000,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       1,
		InterCharacterTimeout: 10,
	})
	if error != nil {
		log.Fatal(error)
	}
	defer ser.Close()

	b := make([]byte, 2048)
	fmt.Printf("READING")
	ret, err := ser.Read(b)
	if ret != 0 {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read %d bytes", ret)

	bytes, errs := ser.InWaiting()
	if errs != nil {
		fmt.Printf("OI")
	}
	fmt.Printf("Data bits in waiting: %d", bytes)
}
