package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	p "github.com/deepakkamesh/sonny/protocol"
	"github.com/tarm/serial"
)

func main() {
	fmt.Println("Welcome to Controller CLI")
	fmt.Println("Format: <Command Hex Code> <Device Hex Code> <optional params>")
	log.SetFlags(log.Lmicroseconds)
	in := bufio.NewReader(os.Stdin)
	//c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600}
	c := &serial.Config{Name: "/dev/ttyS0", Baud: 115200}
	ser, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Error opening %s", err)
	}

	// Read the serial line in a goroutine.
	go read(ser)

	// Command read and process loop.
Start:
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read stdin: %v", err)
		}
		// Split input into string slice.
		inputs := strings.Split(strings.Trim(line, "\n "), " ")
		// Validate size of input.
		if len(inputs) > p.PKT_SZ || len(inputs) < 2 {
			log.Printf("Bad packet size")
			continue
		}
		// Convert string into bytes slice.
		bytes := []byte{}
		for i := 0; i < len(inputs); i++ {
			c, err := strconv.ParseUint(inputs[i], 16, 8)
			if err != nil {
				log.Printf("Error converting input: %v", err)
				continue Start
			}
			bytes = append(bytes, byte(c))
		}

		// Create the command byte.
		packet := append([]byte{}, bytes[0]<<4|bytes[1])
		// Add optional param to the packet.
		packet = append(packet, bytes[2:]...)
		// Prepend header.
		packet = append(append([]byte{}, p.Header(packet)), packet...)

		// Prepend header and write to serial line.
		if _, err := ser.Write(packet); err != nil {

			log.Printf("failed to send to serial: %v", err)
		}
		log.Printf("%v", packet)
		log.Printf("Sent %s\n", p.PktPrint(packet))
	}

}

func read(ser *serial.Port) {

	for {

		buf := make([]byte, 16)
		_, err := ser.Read(buf)
		if err != nil {
			fmt.Printf("Error reading %s", err)
			continue
		}
		//log.Printf("Got %d bytes Binary:%08b  Hex:%X", n, buf[0], buf[0])
		log.Printf("Got %s\n\n", p.PktPrint(buf))
	}

}
