package main

import (
	"fmt"

	"github.com/deepakkamesh/sonny/devices"
)

func main() {

	ctrl := devices.NewController("ttyAMA0", 115200)
	ctrl.Start()

	if err := ctrl.RotateServo(10); err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("completed")

	for {
	}
}
