package devices

import (
	"testing"
	"time"

	p "github.com/deepakkamesh/sonny/protocol"
	"github.com/tarm/serial"
)

func TestController(t *testing.T) {

	m := make(chan struct{})
	n := make(chan struct{})
	serialOpen = func(c *serial.Config) (*serial.Port, error) {
		return &serial.Port{}, nil
	}

	ctrl, err := NewController("/dev/ttyAMA0", 115211)
	if err != nil {
		t.Fatalf("Unable to open tty %v", err)
	}
	ctrl.Start()
	for _, i := range []struct {
		read    func(s *serial.Port, b []byte) (int, error)
		write   func(s *serial.Port, b []byte) (int, error)
		wantErr bool
		info    string
	}{
		{

			func(s *serial.Port, b []byte) (int, error) {
				<-n
				b[0] = 0x11
				b[1] = p.ACK_DONE<<4 | p.DEV_ADMIN
				return 2, nil
			},
			func(s *serial.Port, b []byte) (int, error) {
				for _, i := range b {
					if i == (p.CMD_PING<<4 | p.DEV_ADMIN) {
						n <- struct{}{}
					}
				}
				return 0, nil
			},
			false,
			"all pass test",
		},
		{
			func(s *serial.Port, b []byte) (int, error) {
				<-m
				time.Sleep(700 * time.Millisecond)
				b[0] = 0x11
				b[1] = p.ACK_DONE<<4 | p.DEV_ADMIN
				return 2, nil
			},
			func(s *serial.Port, b []byte) (int, error) {
				if b[0] == (p.CMD_PING<<4 | p.DEV_ADMIN) {
					m <- struct{}{}
				}
				return 0, nil
			},
			true,
			"> 500ms timeout test",
		},
		{
			func(s *serial.Port, b []byte) (int, error) {
				<-m
				time.Sleep(300 * time.Millisecond)
				b[0] = 0x11
				b[1] = p.ACK_DONE<<4 | p.DEV_ADMIN
				return 2, nil
			},
			func(s *serial.Port, b []byte) (int, error) {
				if b[0] == (p.CMD_PING<<4 | p.DEV_ADMIN) {
					m <- struct{}{}
				}
				return 0, nil
			},
			false,
			"< 500ms timeout test",
		},
	} {
		t.Logf("Executing %s", i.info)
		serialRead = i.read
		serialWrite = i.write
		ctrl.Start()
		err := ctrl.Ping()
		if i.wantErr && err == nil {
			t.Errorf("expected error %v got %v", i.wantErr, err)
		}
	}
}
