package devices

import (
	"fmt"
	"math"
	"os"
	"time"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/go-roomba/constants"
	"github.com/deepakkamesh/ydlidar"
	"github.com/golang/glog"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
)

type SensorData struct {
	data    map[byte]int16
	lastUpd time.Time
}

// Sonny is the struct that represents all the devices.
type Sonny struct {
	*Controller                                 // PIC controller.
	*ydlidar.YDLidar                            // YDLidar.
	*i2c.QMC5883Driver                          // Magnetometer QMC5883.
	*i2c.MPU6050Driver                          // MPU 6050 Accelerometer / Gryo.
	*roomba.Roomba                              // Roomba controller.
	i2cEn                 *gpio.DirectPinDriver // GPIO port control for I2C Bus.
	*gpio.PIRMotionDriver                       // PIR driver.
	*Video
	pirState        int           // State of PIR. 1=enabled, 0=disabled.
	i2cBusState     int           // State of I2CBus. 1=enabled, 0=disabled.
	auxPowerState   int           // Start of AuxPower. 1=enabled, 0=disabled.
	roombaMode      int           // Roomba mode: 1 = passive, 2=safe, 3=full.
	auxPowerOnInit  func() error  // initialization to execute after Aux Power is on.
	auxPowerOffInit func() error  // initialization to execute after Aux Power is off.
	sensorData      SensorData    // Sensor data updated regularly.
	magXmin         int16         // magnetometer X min value for calibration.
	magXmax         int16         // magnetometer X max value for calibration.
	magYmin         int16         // magnetometer Y min value for calibration.
	magYmax         int16         // magnetometer Y max value for calibration.
	killSensorData  chan struct{} // shutdown sensor loop.
}

func NewSonny(
	c *Controller,
	l *ydlidar.YDLidar,
	m *i2c.QMC5883Driver,
	a *i2c.MPU6050Driver,
	r *roomba.Roomba,
	i2cEn *gpio.DirectPinDriver,
	p *gpio.PIRMotionDriver,
	v *Video,
) *Sonny {

	return &Sonny{
		c, l, m, a, r, i2cEn, p, v, 0, 0, 0, 0,
		func() error { return nil },
		func() error { return nil },
		SensorData{
			make(map[byte]int16),
			time.Now(),
		},
		-4928, // min X. Default min/max values for a sane Magnetometer offset.
		2460,  // max X.
		-427,  // min Y.
		6844,  // max Y.
		make(chan struct{}, 2),
	}
}

// StartRoomba starts up the roomba platform.
func (s *Sonny) StartRoomba(keepAlive bool) error {
	return s.Roomba.Start(keepAlive)
}

func (s *Sonny) Startup() {
	s.PIREventLoop()
	go s.updateRoombaTelemetry()
}

func (s *Sonny) Shutdown() {
	if s == nil {
		glog.Warningf("Sonny not initialized")
		return
	}
	s.killSensorData <- struct{}{}
	if err := s.I2CBusEnable(false); err != nil {
		glog.Fatalf("Failed to disable I2C Bus: %v", err)
	}
	if err := s.AuxPower(false); err != nil {
		glog.Errorf("Failed to disable aux power: %v", err)
	}
	if err := s.SetRoombaMode(constants.OI_MODE_PASSIVE); err != nil { // Reset roomba turns it off.
		glog.Errorf("Failed to reset Roomba on shutdown")
	}
	if err := s.LidarPower(false); err != nil {
		glog.Errorf("Failed to shutdown lidar power")
	}
	s.Reset()
	if s.Video != nil {
		s.Video.StopVideoStream()
	}
}

// RoombaInitialized returns true if Roomba is initialized.
func (s *Sonny) RoombaInitialized() bool {
	return (s.Roomba != nil)
}

func (s *Sonny) LidarInitialized() bool {
	return (s.YDLidar != nil)
}
func (s *Sonny) MagnetometerInitialized() bool {
	return (s.QMC5883Driver != nil)
}

func (s *Sonny) ControllerInitialized() bool {
	return (s.Controller != nil)
}

// GetAuxPowerState returns the state of Aux Power.
func (s *Sonny) GetAuxPowerState() int {
	return s.auxPowerState
}

// SetAuxPostInit sets the initialization routines to call after aux power
// is turned on or off.
func (s *Sonny) SetAuxPostInit(fOn func() error, fOff func() error) {
	s.auxPowerOnInit = fOn
	s.auxPowerOffInit = fOff
}

// AuxPower enables/disables Auxillary power from main brush motor on Roomba.
func (s *Sonny) AuxPower(enable bool) error {
	if enable {
		if s.GetRoombaMode() < 2 {
			return fmt.Errorf("Roomba mode should be Safe or Full. In Mode %v", s.GetRoombaMode())
		}
		// Power up auxillary battery on main brush.
		time.Sleep(300 * time.Millisecond) // Not sure why, but a little time is needed.
		if err := s.MainBrush(true, true); err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond) // Time to power up Aux.
		s.auxPowerState = 1
		if err := s.auxPowerOnInit(); err != nil {
			return err
		}
		return nil
	}

	if err := s.MainBrush(false, true); err != nil {
		return err
	}
	time.Sleep(1000 * time.Millisecond) // Time to power down Aux.
	s.auxPowerState = 0
	if err := s.auxPowerOffInit(); err != nil {
		return err
	}
	return nil
}

// LidarPower starts the lidar motor.
func (s *Sonny) LidarPower(enable bool) error {
	if s.YDLidar == nil {
		return fmt.Errorf("lidar not initialized")
	}

	return s.SetDTR(enable)
}

func (s *Sonny) LidarData() ydlidar.Packet {
	return <-s.YDLidar.D
}

// PIREventLoop subscribes to events from the PIR gpio.
func (s *Sonny) PIREventLoop() {
	if s.PIRMotionDriver == nil {
		return
	}

	pirCh := s.PIRMotionDriver.Subscribe()
	go func() {
		for {
			evt := <-pirCh
			s.pirState = evt.Data.(int)
			glog.V(3).Infof("Got pir data %v %v", evt.Name, evt.Data.(int))
		}
	}()
}

// Returns PIR state.
func (s *Sonny) GetPIRState() int {
	return s.pirState
}

func (s *Sonny) updateRoombaTelemetry() {

	if s.Roomba == nil {
		glog.Errorf("Failed to update telemetry because Roomba is not initialized")
		return
	}

	for {
		select {
		case <-s.killSensorData:
			return

		default:
			data := make(map[byte]int16)
			d, e := s.Roomba.QueryList(constants.PACKET_GROUP_100)
			if e != nil {
				glog.Errorf("Failed to read roomba sensors:%v", e)
			}

			for i, p := range d {
				pktID := constants.PACKET_GROUP_100[i]
				if len(p) == 1 {
					data[pktID] = int16(p[0])
					continue
				}
				data[pktID] = int16(p[0])<<8 | int16(p[1])
			}
			s.sensorData.data = data
			s.sensorData.lastUpd = time.Now()

			// Inspect roomba mode. If different, reset aux power if passive mode.
			prevMode := s.roombaMode
			s.roombaMode = int(data[constants.SENSOR_OI_MODE])
			if s.roombaMode != prevMode && s.roombaMode == 1 {
				s.AuxPower(false)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// GetRoombaTelemetry returns the current value of the roomba sensors.
func (s *Sonny) GetRoombaTelemetry() (data map[byte]int16, err error) {
	return s.sensorData.data, nil
}

// GetRoombaMode returns the current roomba mode from the sensor reading.
func (s *Sonny) GetRoombaMode() int {
	return s.roombaMode
}

// SetRoombaMode sets the mode for Roomba.
func (s *Sonny) SetRoombaMode(mode byte) error {
	if s.Roomba == nil {
		return fmt.Errorf("roomba not initialized")
	}
	switch mode {
	case constants.OI_MODE_OFF:
		if err := s.Roomba.Power(); err != nil {
			return err
		}
		s.roombaMode = 0
	case constants.OI_MODE_PASSIVE:
		if err := s.Roomba.Passive(); err != nil {
			return err
		}
		s.roombaMode = 1
	case constants.OI_MODE_SAFE:
		if err := s.Roomba.Safe(); err != nil {
			return err
		}
		s.roombaMode = 2
	case constants.OI_MODE_FULL:
		if err := s.Roomba.Full(); err != nil {
			return err
		}
		s.roombaMode = 3
	default:
		return fmt.Errorf("unknown mode %v requested", mode)
	}

	return nil
}

// I2CBusEnable enables/disables the I2C buffer chip.
// Connects the rest of the I2C devices with Pi.
func (s *Sonny) I2CBusEnable(b bool) error {
	if s.i2cEn == nil {
		return fmt.Errorf("gpio I2C not initialized")
	}
	if b {
		s.i2cBusState = 1
		return s.i2cEn.DigitalWrite(1)
	}
	s.i2cBusState = 0
	return s.i2cEn.DigitalWrite(0)
}

// GetI2CBusState return 1 if I2C bus is enabled otherwise it returns 0.
func (s *Sonny) GetI2CBusState() int {
	return s.i2cBusState
}

// CalibrateCompass runs the calibration routine on the compass and sets
// the offsets. In test mode it recalibrations after apply existing offets
// without storing the calibration values. Useful to see if recalib is needed.
func (s *Sonny) CalibrateCompass(test bool) error {

	// Lock controller to avoid any I2C contention.
	s.LockController()
	defer s.UnlockController()

	if s.QMC5883Driver == nil {
		return fmt.Errorf("Compass not enabled")
	}

	var minX, maxX, minY, maxY int16
	f, err := os.OpenFile("calibrationReading.csv", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	glog.Infof("Calibration in test mode: %v", test)

	vel := int16(50) // vel in mm/s.
	circum := constants.RoombaRadius * 2 * math.Pi
	driveTime := float32(circum) / float32(vel) // in secs.

	go func() {
		// Turn the bot through two 360 rotation.
		if err := s.DirectDrive(vel, -vel); err != nil {
			glog.Errorf("%v", err)
			return
		}
		time.Sleep(time.Duration(driveTime*1000*2) * time.Millisecond)
		if err := s.DirectDrive(0, 0); err != nil {
			glog.Errorf("%v", err)
			return
		}

	}()

	for i := 0; i < int(driveTime*1000*2/50); i++ {
		x, y, z, err := s.RawHeading()
		if err != nil {
			return err
		}

		// If this is a test then apply offsets to the raw readings.
		if test {
			x -= (s.magXmin + s.magXmax) / 2
			y -= (s.magYmin + s.magYmax) / 2
		}

		f.WriteString(fmt.Sprintf(" Reading, %v,%v,%v\n", x, y, z))

		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
		time.Sleep(50 * time.Millisecond)
	}
	f.Sync()

	offX := (minX + maxX) / 2
	offY := (minY + maxY) / 2
	glog.Infof("Calibration Complete: X:(%v - %v)/2=%v, Y:(%v - %v)/2=%v", minX, maxX, offX, minY, maxY, offY)
	if !test {
		s.magXmin = minX
		s.magXmax = maxX
		s.magYmin = minY
		s.magYmax = maxY
	}
	return nil
}

// Accelerometer returns the X,Y,Z data.
func (s *Sonny) Accelerometer() (x, y, z int16, err error) {

	if err = s.MPU6050Driver.GetData(); err != nil {
		return
	}

	x = s.MPU6050Driver.Accelerometer.X
	y = s.MPU6050Driver.Accelerometer.Y
	z = s.MPU6050Driver.Accelerometer.Z
	return
}

// Gryo returns the X,Y,Z data from Gryo.
func (s *Sonny) Gyro() (x, y, z int16, err error) {
	if err = s.MPU6050Driver.GetData(); err != nil {
		return
	}

	x = s.MPU6050Driver.Gyroscope.X
	y = s.MPU6050Driver.Gyroscope.Y
	z = s.MPU6050Driver.Gyroscope.Z
	return
}

// TiltHeading returns tilt compensated compass reading
// with a low pass filter.
func (s *Sonny) TiltHeading() (float64, error) {

	// Lock controller to avoid any I2C contention.
	s.LockController()
	defer s.UnlockController()

	magLPF := 0.4
	accLPF := 0.1

	var accX, accY, accZ, xR, yR, zR int16
	var paccX, paccY, paccZ, pxR, pyR, pzR int16
	var err error

	// Apply Low Pass Filter to stablize value.
	for i := 0; i < 50; i++ {
		accX, accY, accZ, err = s.Accelerometer()
		if err != nil {
			return 0, nil
		}
		xR, yR, zR, err = s.RawHeading()
		if err != nil {
			return 0, nil
		}

		accX = int16(float64(accX)*accLPF + float64(paccX)*(1-accLPF))
		accY = int16(float64(accY)*accLPF + float64(paccY)*(1-accLPF))
		accZ = int16(float64(accZ)*accLPF + float64(paccZ)*(1-accLPF))

		paccX = accX
		paccY = accY
		paccZ = accZ

		xR = int16(float64(xR)*magLPF + float64(pxR)*(1-magLPF))
		yR = int16(float64(yR)*magLPF + float64(pyR)*(1-magLPF))
		zR = int16(float64(zR)*magLPF + float64(pzR)*(1-magLPF))

		pxR = xR
		pyR = yR
		pzR = zR
		time.Sleep(5 * time.Millisecond)
	}

	// Apply hard iron offsets.
	xR -= (s.magXmin + s.magXmax) / 2
	yR -= (s.magYmin + s.magYmax) / 2

	// TODO: Apply soft iron offets.
	//	x := float64((xR-s.magXmin)/(s.magXmax-s.magXmin))*2 - 1
	//	y := float64((yR-s.magXmin)/(s.magXmax-s.magXmin))*2 - 1

	x := float64(xR)
	y := float64(yR)
	z := float64(zR)

	// Normalize accelerometer reading.
	div := math.Sqrt(float64(accX)*float64(accX) +
		float64(accY)*float64(accY) +
		float64(accZ)*float64(accZ))
	accXn := float64(accX) / div
	accYn := float64(accY) / div

	// Tilt compensation.
	pitch := math.Asin(accXn)
	roll := -math.Asin(accYn / math.Cos(pitch))

	magX := x*math.Cos(pitch) + (z)*math.Sin(pitch)
	magY := x*math.Sin(roll)*math.Sin(pitch) +
		y*math.Cos(roll) -
		z*math.Sin(roll)*math.Cos(pitch)

	heading := 180 * math.Atan2(magY, magX) / math.Pi
	if heading < 0 {
		heading += 360
	}

	// Declination.
	declination := 231.8 / 1000
	heading += declination * 180 / math.Pi
	if heading > 360 {
		heading -= 360
	}

	return heading, nil
}

// MoveForward moves the bot in mm. Negative vel moves backward. It returns the
// actual distance moved.
func (s *Sonny) MoveForward(desiredDist, vel int) (float64, error) {

	if desiredDist == 0 {
		return 0, nil
	}

	// Just in case , always shutdown motors.
	defer s.DirectDrive(0, 0)

	// Get starting readings.
	encLStart, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0, err
	}
	encRStart, err := s.getEncoderReading(constants.SENSOR_RIGHT_ENCODER)
	if err != nil {
		return 0, err
	}

	driveTime := float64(desiredDist) * 1000 / float64(math.Abs(float64(vel))) // # of ms to engage motor.
	if err := s.DirectDrive(int16(vel), int16(vel)); err != nil {
		return 0, err
	}
	// Goroutine applies correct drive power to motors equalized by
	// applying differential power.
	fin := make(chan struct{})
	go func() {
		tDelta := 30 // Check every x ms.
		t := time.NewTicker(time.Duration(tDelta) * time.Millisecond)
		var encStL, encStR int16 = encLStart, encRStart
		kP := 0.6 // Constant factor for power differential.

		for {
			select {
			case <-t.C:
				encL, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
				if err != nil {
					glog.Errorf("Failed to get encoder: %v", err)
				}
				encR, err := s.getEncoderReading(constants.SENSOR_RIGHT_ENCODER)
				if err != nil {
					glog.Errorf("Failed to get encoder: %v", err)
				}

				dL := getDistFromEncoder(encStL, encL)
				dR := getDistFromEncoder(encStR, encR)

				// Power differential is difference in dist travelled in mm converted
				// to speed in mm/s with a constant applied for gradual change.
				pDiff := (dL - dR) * (1000 / float64(tDelta)) * kP
				encStL = encL
				encStR = encR

				lPower := int16(vel)
				rPower := int16(vel + int(pDiff))
				if vel < 0 {
					rPower = int16(vel - int(pDiff))
				}
				if err := s.DirectDrive(rPower, lPower); err != nil {
					glog.Errorf("Failed to set drive power: %v", err)
				}

			case <-fin:
				return
			}
		}
	}()

	time.Sleep(time.Duration(driveTime) * time.Millisecond)
	fin <- struct{}{}
	if err := s.DirectDrive(0, 0); err != nil {
		return 0, err
	}
	time.Sleep(300 * time.Millisecond)

	// End Encoder readings.
	encLEnd, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0, err
	}
	encREnd, err := s.getEncoderReading(constants.SENSOR_RIGHT_ENCODER)
	if err != nil {
		return 0, err
	}

	distL := getDistFromEncoder(encLStart, encLEnd)
	distR := getDistFromEncoder(encRStart, encREnd)
	if math.Abs(distR-distL) > 2 {
		glog.Warningf("Wheels reporting different distance travelled: R:%0.2f mm, L:%0.2f mm", distR, distL)
		glog.Infof("EncRight st:%v, end:%v | EncLeft st:%v end: %v", encRStart, encREnd, encLStart, encLEnd)
	}

	return distR, nil
}

// getEncoderReading returns the current encoder reading.
func (s *Sonny) getEncoderReading(enc byte) (int16, error) {
	p, err := s.Sensors(enc)
	if err != nil {
		return 0, err
	}
	return int16(p[0])<<8 | int16(p[1]), nil
}

// getDistFromEncoder calculates dist travelled in mm from encoder reading.
func getDistFromEncoder(start, end int16) float64 {
	delta := math.Abs(float64(end - start))

	// Calculate dist from encoder reading.
	return float64(delta) * math.Pi * 72.0 / 508.8

}

// Turn rotates the bot by angle in degrees and returns the delta in degrees.
// positve angle rotates clockwise.
func (s *Sonny) Turn(angle float64) (float64, error) {
	if angle == 0 {
		return 0.0, nil
	}
	vel := 50 //speed in mm/s.

	// Get starting encoder reading.
	encStart, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0.0, err
	}

	// Calculate cirumference to drive. θ radian = circumfrence/radius.
	c := constants.RoombaRadius * math.Abs(angle) * math.Pi / 180
	driveTime := float32(c) / float32(vel)

	glog.V(3).Infof("driveTime(s)%v Circumfrence(mm):%v Angle:%v", driveTime, c, angle)

	rvel := -1 * vel
	lvel := vel
	if angle < 0 {
		rvel = vel
		lvel = -1 * vel
	}

	st := make(chan struct{})

	// Get gyro readings when turning.
	var yaw float64
	go func() {
		for {
			select {
			case <-st:
				glog.Infof("Gyroscopc Rotation: %0.2f", yaw)
				return
			default:
				x, y, z, err := s.Gyro()
				if err != nil {
					glog.Errorf("%v", err)
				}
				rX := float64(x) * 0.00875
				rY := float64(y) * 0.00875
				rZ := float64(z) * 0.00875
				yaw += rZ * 10 / 1000
				_ = rX
				_ = rY
				//	glog.Infof("%0.3f %0.3f %0.3f", rX, rY, rZ)
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Turn the bot.
	if err := s.DirectDrive(int16(rvel), int16(lvel)); err != nil {
		return 0.0, err
	}
	time.Sleep(time.Duration(driveTime*1000) * time.Millisecond)
	if err := s.DirectDrive(0, 0); err != nil {
		return 0.0, err
	}
	time.Sleep(300 * time.Millisecond)
	st <- struct{}{}

	// Calculate if we overshot or undershot landing and return delta.
	encEnd, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0.0, err
	}

	// TODO: need to return the right delta from compass.
	_ = encStart
	_ = encEnd
	return 0.0, nil
}
