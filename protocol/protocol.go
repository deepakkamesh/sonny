package protocol

import (
	"errors"
	"fmt"
)

const (
	PKT_SZ      = 16
	MAX_DEVICES = 16

	// Device Definitions.
	DEV_ADMIN byte = 0x0
	DEV_LED   byte = 0x1
	DEV_SERVO byte = 0x2

	// Command definitions.
	CMD_ON      byte = 0x1
	CMD_PING    byte = 0x2
	CMD_VERSION byte = 0x3
	CMD_OFF     byte = 0x4
	CMD_ROTATE  byte = 0x6

	// Error Codes.
	ERR_CHECKSUM_FAILURE byte = 0x1
	ERR_DEVICE_BUSY      byte = 0x2

	// Status Codes.
	ACK      byte = 0x8
	ACK_DONE byte = 0xC
	ERR      byte = 0x0
	DONE     byte = 0x4
)

func Error(errCode byte) error {
	err := map[byte]string{
		0x1: "checksum mismatch",
		0x2: "device busy",
		0x3: "unimplemented",
		0x4: "insufficient parameters",
	}

	if e, ok := err[errCode]; ok {
		return errors.New(e)
	}
	return errors.New("unknown")
}

func Checksum(packet []byte) byte {
	return 0x1
}

func VerifyChecksum(packet []byte) bool {
	return true
}

func PacketSz(packet []byte) byte {
	return packet[0] >> 4
}

func Header(packet []byte) byte {
	header := byte(len(packet) << 4)
	return header | Checksum(packet)

}
func StatusCode(b byte) byte {
	return b >> 4
}

func DeviceID(b byte) byte {
	return b & 0xF
}
func PrettyPrint(packet []byte) (logline string) {

	// Calculate len of packet
	//sz := (packet[0] >> 4) & 0xF
	// Print starting at header.
	for i := 0; i < len(packet); i++ {
		switch i {
		case 0:
			logline = logline + fmt.Sprintf("\n header - %08b\n", packet[i])
		case 1:
			logline = logline + fmt.Sprintf(" status/cmd - %08b\n", packet[i])
		default:
			logline = logline + fmt.Sprintf(" param%02d - 0x%X (%d)\n", i-1, packet[i], packet[i])
		}
	}
	return

}
