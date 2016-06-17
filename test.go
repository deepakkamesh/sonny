package main

import (
	"fmt"
	"log"

	"github.com/tarm/serial"
)

func main() {
	fmt.Println("Welcome to Sonny")
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600}
	ser, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Error opening %s", err)
	}
	pkt := []byte{0x10, 0x20}
	_, err = ser.Write(pkt)
	if err != nil {
		fmt.Printf("error transmit %s", err)
	}
	read(ser)
	for {
	}
}

func run(ser *serial.Port) {
	go read(ser)
}
func read(ser *serial.Port) {

	buf := make([]byte, 1)
	for {
		n, err := ser.Read(buf)
		if err != nil {
			fmt.Printf("Error reading %s", err)
			continue
		}
		fmt.Printf("Got %d bytes Binary:%08b  Hex:%X  Char:%c\n", n, buf[0], buf[0], buf[0])
	}

}
