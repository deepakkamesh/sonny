// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import mjpeg "github.com/saljam/mjpeg"
import mock "github.com/stretchr/testify/mock"

// Platform is an autogenerated mock type for the Platform type
type Platform struct {
	mock.Mock
}

// Accelerometer provides a mock function with given fields:
func (_m *Platform) Accelerometer() (int16, int16, int16, error) {
	ret := _m.Called()

	var r0 int16
	if rf, ok := ret.Get(0).(func() int16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int16)
	}

	var r1 int16
	if rf, ok := ret.Get(1).(func() int16); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(int16)
	}

	var r2 int16
	if rf, ok := ret.Get(2).(func() int16); ok {
		r2 = rf()
	} else {
		r2 = ret.Get(2).(int16)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func() error); ok {
		r3 = rf()
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// AuxPower provides a mock function with given fields: enable
func (_m *Platform) AuxPower(enable bool) error {
	ret := _m.Called(enable)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool) error); ok {
		r0 = rf(enable)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BattState provides a mock function with given fields:
func (_m *Platform) BattState() (float32, error) {
	ret := _m.Called()

	var r0 float32
	if rf, ok := ret.Get(0).(func() float32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CalibrateCompass provides a mock function with given fields: _a0
func (_m *Platform) CalibrateCompass(_a0 bool) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ControllerInitialized provides a mock function with given fields:
func (_m *Platform) ControllerInitialized() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// DHT11 provides a mock function with given fields:
func (_m *Platform) DHT11() (uint8, uint8, error) {
	ret := _m.Called()

	var r0 uint8
	if rf, ok := ret.Get(0).(func() uint8); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint8)
	}

	var r1 uint8
	if rf, ok := ret.Get(1).(func() uint8); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(uint8)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DirectDrive provides a mock function with given fields: _a0, _a1
func (_m *Platform) DirectDrive(_a0 int16, _a1 int16) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(int16, int16) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Distance provides a mock function with given fields:
func (_m *Platform) Distance() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Drive provides a mock function with given fields: _a0, _a1
func (_m *Platform) Drive(_a0 int16, _a1 int16) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(int16, int16) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ForwardSweep provides a mock function with given fields: _a0, _a1, _a2
func (_m *Platform) ForwardSweep(_a0 int, _a1 int, _a2 int) ([]int32, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []int32
	if rf, ok := ret.Get(0).(func(int, int, int) []int32); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int32)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, int) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Full provides a mock function with given fields:
func (_m *Platform) Full() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAuxPowerState provides a mock function with given fields:
func (_m *Platform) GetAuxPowerState() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetI2CBusState provides a mock function with given fields:
func (_m *Platform) GetI2CBusState() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetPIRState provides a mock function with given fields:
func (_m *Platform) GetPIRState() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetRoombaMode provides a mock function with given fields:
func (_m *Platform) GetRoombaMode() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetRoombaTelemetry provides a mock function with given fields:
func (_m *Platform) GetRoombaTelemetry() (map[byte]int16, error) {
	ret := _m.Called()

	var r0 map[byte]int16
	if rf, ok := ret.Get(0).(func() map[byte]int16); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[byte]int16)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVideoStream provides a mock function with given fields:
func (_m *Platform) GetVideoStream() *mjpeg.Stream {
	ret := _m.Called()

	var r0 *mjpeg.Stream
	if rf, ok := ret.Get(0).(func() *mjpeg.Stream); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mjpeg.Stream)
		}
	}

	return r0
}

// Gyro provides a mock function with given fields:
func (_m *Platform) Gyro() (int16, int16, int16, error) {
	ret := _m.Called()

	var r0 int16
	if rf, ok := ret.Get(0).(func() int16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int16)
	}

	var r1 int16
	if rf, ok := ret.Get(1).(func() int16); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(int16)
	}

	var r2 int16
	if rf, ok := ret.Get(2).(func() int16); ok {
		r2 = rf()
	} else {
		r2 = ret.Get(2).(int16)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func() error); ok {
		r3 = rf()
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// I2CBusEnable provides a mock function with given fields: _a0
func (_m *Platform) I2CBusEnable(_a0 bool) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LDR provides a mock function with given fields:
func (_m *Platform) LDR() (uint16, error) {
	ret := _m.Called()

	var r0 uint16
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LEDBlink provides a mock function with given fields: _a0, _a1
func (_m *Platform) LEDBlink(_a0 uint16, _a1 byte) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint16, byte) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LEDOn provides a mock function with given fields: _a0
func (_m *Platform) LEDOn(_a0 bool) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LidarInitialized provides a mock function with given fields:
func (_m *Platform) LidarInitialized() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// LidarPower provides a mock function with given fields: _a0
func (_m *Platform) LidarPower(_a0 bool) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MagnetometerInitialized provides a mock function with given fields:
func (_m *Platform) MagnetometerInitialized() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MoveForward provides a mock function with given fields: _a0, _a1
func (_m *Platform) MoveForward(_a0 int, _a1 int) (float64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 float64
	if rf, ok := ret.Get(0).(func(int, int) float64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Passive provides a mock function with given fields:
func (_m *Platform) Passive() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Ping provides a mock function with given fields:
func (_m *Platform) Ping() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Power provides a mock function with given fields:
func (_m *Platform) Power() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Reset provides a mock function with given fields:
func (_m *Platform) Reset() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoombaInitialized provides a mock function with given fields:
func (_m *Platform) RoombaInitialized() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Safe provides a mock function with given fields:
func (_m *Platform) Safe() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SeekDock provides a mock function with given fields:
func (_m *Platform) SeekDock() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Sensors provides a mock function with given fields: _a0
func (_m *Platform) Sensors(_a0 byte) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(byte) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(byte) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServoRotate provides a mock function with given fields: _a0, _a1
func (_m *Platform) ServoRotate(_a0 byte, _a1 int) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(byte, int) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetRoombaMode provides a mock function with given fields: _a0
func (_m *Platform) SetRoombaMode(_a0 byte) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StartRoomba provides a mock function with given fields: _a0
func (_m *Platform) StartRoomba(_a0 bool) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TiltHeading provides a mock function with given fields:
func (_m *Platform) TiltHeading() (float64, error) {
	ret := _m.Called()

	var r0 float64
	if rf, ok := ret.Get(0).(func() float64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Turn provides a mock function with given fields: _a0
func (_m *Platform) Turn(_a0 float64) (float64, error) {
	ret := _m.Called(_a0)

	var r0 float64
	if rf, ok := ret.Get(0).(func(float64) float64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(float64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}