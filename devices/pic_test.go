package devices

import (
	"testing"
	"time"

	p "github.com/deepakkamesh/sonny/protocol"
	"github.com/tarm/serial"
)

func TestController(t *testing.T) {

	readCnt := 0

	serialOpen = func(c *serial.Config) (*serial.Port, error) {
		return &serial.Port{}, nil
	}

	for _, i := range []struct {
		read    func(s *serial.Port, b []byte) (int, error)
		write   func(s *serial.Port, b []byte) (int, error)
		wantErr bool
		info    string
	}{
		{
			func(s *serial.Port, b []byte) (int, error) {
				// Block till there is a write request.
				for readCnt == 0 {
					time.Sleep(time.Millisecond * 50)
					return 0, nil
				}
				switch readCnt {
				case 1:
					b[0] = 0x11
					readCnt += 1
					return 1, nil

				case 2:
					b[0] = p.ACK_DONE<<4 | p.DEV_ADMIN
					readCnt = 0
					return 1, nil
				}
				return 0, nil
			},
			func(s *serial.Port, b []byte) (int, error) {
				for _, i := range b {
					if i == (p.CMD_PING<<4 | p.DEV_ADMIN) {
						readCnt = 1
					}
				}
				return 1, nil
			},
			false,
			"ACK_DONE tests",
		},
		{
			func(s *serial.Port, b []byte) (int, error) {
				// Block till there is a write request.
				for readCnt == 0 {
					time.Sleep(time.Millisecond * 50)
					return 0, nil
				}
				switch readCnt {
				case 1:
					b[0] = 0x21
					readCnt += 1
					return 1, nil

				case 2:
					b[0] = p.ERR<<4 | p.DEV_ADMIN
					b[1] = p.ERR_UNIMPLEMENTED
					readCnt = 0
					return 2, nil
				}
				return 0, nil
			},
			func(s *serial.Port, b []byte) (int, error) {
				for _, i := range b {
					if i == (p.CMD_PING<<4 | p.DEV_ADMIN) {
						readCnt = 1
					}
				}
				return 1, nil
			},
			true,
			"device error tests",
		},
		{
			func(s *serial.Port, b []byte) (int, error) {
				// Block till there is a write request.
				for readCnt == 0 {
					time.Sleep(time.Millisecond * 50)
					return 0, nil
				}
				switch readCnt {
				case 1:
					b[0] = 0x11
					readCnt += 1
					return 1, nil

				case 2:
					b[0] = p.ACK<<4 | p.DEV_ADMIN
					readCnt += 1
					return 1, nil

				case 3:
					time.Sleep(time.Millisecond * 100)
					b[0] = 0x11
					readCnt += 1
					return 1, nil

				case 4:
					b[0] = p.DONE<<4 | p.DEV_ADMIN
					readCnt = 0
					return 1, nil
				}
				return 0, nil
			},
			func(s *serial.Port, b []byte) (int, error) {
				for _, i := range b {
					if i == (p.CMD_PING<<4 | p.DEV_ADMIN) {
						readCnt = 1
					}
				}
				return 0, nil
			},
			false,
			"ACK->DONE  tests",
		},
		{
			func(s *serial.Port, b []byte) (int, error) {
				for readCnt == 0 {
					time.Sleep(time.Millisecond * 50)
					return 0, nil
				}
				time.Sleep(time.Millisecond * 500)

				switch readCnt {
				case 1:
					b[0] = 0x11
					readCnt += 1
					return 1, nil

				case 2:
					b[0] = p.ACK_DONE<<4 | p.DEV_ADMIN
					readCnt = 0
					return 1, nil
				}

				return 0, nil
			},
			func(s *serial.Port, b []byte) (int, error) {
				if b[0] == (p.CMD_PING<<4 | p.DEV_ADMIN) {
					readCnt = 1
				}
				return 0, nil
			},
			true,
			"ACK/ACK_DONE timeout test",
		},
		{
			func(s *serial.Port, b []byte) (int, error) {
				for readCnt == 0 {
					time.Sleep(time.Millisecond * 50)
					return 0, nil
				}

				switch readCnt {
				case 1:
					b[0] = 0x11
					readCnt += 1
					return 1, nil

				case 2:
					b[0] = p.ACK<<4 | p.DEV_ADMIN
					readCnt += 1
					return 1, nil

				case 3:
					time.Sleep(time.Millisecond * 700)
					b[0] = 0x11
					readCnt += 1
					return 1, nil

				case 4:
					b[0] = p.DONE<<4 | p.DEV_ADMIN
					readCnt = 0
					return 1, nil
				}
				return 0, nil
			},
			func(s *serial.Port, b []byte) (int, error) {
				if b[0] == (p.CMD_PING<<4 | p.DEV_ADMIN) {
					readCnt = 1
				}
				return 0, nil
			},
			true,
			"DONE timeout test",
		},
	} {
		t.Logf("Executing %s\n", i.info)

		ctrl, err := NewController("/dev/ttyAMA0", 115211)
		if err != nil {
			t.Fatalf("Unable to open tty %v", err)
		}

		serialRead = i.read
		serialWrite = i.write
		ctrl.Start()
		err = ctrl.Ping()
		if (err == nil) == i.wantErr {
			t.Errorf("Expected error: %v got: %v", i.wantErr, err)
		}
		ctrl.Stop()
	}
}
