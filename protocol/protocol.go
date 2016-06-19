package protocol

import "fmt"

const (
	PKT_SZ      = 16
	MAX_DEVICES = 16

	// Device Definitions.
	DEV_ADMIN = 0x0
	DEV_LED   = 0x1

	// Command definitions.
	CMD_ON      = 0x1
	CMD_PING    = 0x2
	CMD_VERSION = 0x3

	// Error Codes.
	ERR_CHECKSUM_FAILURE = 0x1
	ERR_DEVICE_BUSY      = 0x2
)

func CalcChecksum(packet []byte) byte {
	return 0x1
}

func VerifyChecksum(packet []byte, checksum byte) bool {
	return true
}

func Header(packet []byte) byte {
	header := byte(len(packet) << 4)
	return header | CalcChecksum(packet)
}

func PrettyPrint(packet []byte) (logline string) {

	// Calculate len of packet
	//sz := (packet[0] >> 4) & 0xF
	// Print starting at header.
	for i := 0; i < len(packet); i++ {
		switch i {
		case 0:
			logline = logline + fmt.Sprintf(" header - %08b |", packet[i])
		case 1:
			logline = logline + fmt.Sprintf(" status/cmd - %08b |", packet[i])
		default:
			logline = logline + fmt.Sprintf(" param%02d - 0x%X |", i+1, packet[i])
		}
	}
	return

}
