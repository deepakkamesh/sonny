package main

import (
	"fmt"
	"time"

	"github.com/deepakkamesh/sonny/devices"
)

func main() {

	ctrl := devices.NewController("/dev/ttyAMA0", 115200)
	ctrl.Start()

	if err := ctrl.Ping(); err != nil {
		fmt.Println("Error", err)
	}

	if err := ctrl.LedOn(true); err != nil {
		fmt.Println("Error", err)
	}
	time.Sleep(2 * time.Second)
	if err := ctrl.LedOn(false); err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println("completed")

	for {
	}
}
