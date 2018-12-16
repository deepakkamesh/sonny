package navigator

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/deepakkamesh/sonny/mocks"
)

type datapt struct {
	posture float64
	x       int
	y       int
	reading []int32
}

func TestMoveForward(t *testing.T) {
	/*
		mockPl := &mocks.Platform{}

		for _, i := range []struct {
			d     []byte
			delta float64
			req   int
			act   int
		}{
			{[]byte{0x7F, 0xFE, 0x7F, 0xF8, 0x84, 0xB8, 0x84, 0xB8}, 1210, 500, 100},   // 32670, -32560
			{[]byte{0x7F, 0xFE, 0x7F, 0xF8, 0x84, 0xB8, 0x84, 0xB8}, 1210, 600, 100},   // 32670, -32560
			{[]byte{0x84, 0xB8, 0x7F, 0xF8, 0x7F, 0xFE, 0x84, 0xB8}, 1210, -500, -100}, // 32670, -32560
			{[]byte{0x84, 0xB8, 0x7F, 0xF8, 0x7F, 0xFE, 0x84, 0xB8}, 1210, -600, -100}, // 32670, -32560
			//		{[]byte{0, 200, 0, 200, 0, 250, 0, 240}, 16, -100, 0},
		} {

			// Starting encoder readings.
			mockPl.On("Sensors", byte(constants.SENSOR_LEFT_ENCODER)).
				Return([]byte{i.d[0], i.d[1]}, nil).Once()
			mockPl.On("Sensors", byte(constants.SENSOR_RIGHT_ENCODER)).
				Return([]byte{i.d[2], i.d[3]}, nil).Once()

			// Ending encoder readings.
			mockPl.On("Sensors", byte(constants.SENSOR_LEFT_ENCODER)).
				Return([]byte{i.d[4], i.d[5]}, nil).Once()
			mockPl.On("Sensors", byte(constants.SENSOR_RIGHT_ENCODER)).
				Return([]byte{i.d[6], i.d[7]}, nil).Once()

			mockPl.On("DirectDrive", int16(i.act), int16(i.act)).Return(nil)
			mockPl.On("DirectDrive", int16(0), int16(0)).Return(nil)

			//drive := NewAutoDrive(mockPl)
			//d, _ := drive.MoveForward(i.req)

			//t.Logf("Expected %v Got %v", i.delta*math.Pi*72/508.8-float64(math.Abs(float64(i.req))), d-int(math.Abs(float64(i.req))))
		}
		t.Log("Test concluded")
	*/

}

func TestMove(t *testing.T) {

	pf := &mocks.Platform{}
	driver := NewAutoDrive(pf)
	X, Y := driver.GetXY() // Default X,Y

	for _, i := range []struct {
		angle     float64
		vel       int
		dist      int
		distMoved int
		x         int // delta from X and Y.
		y         int
	}{
		{0, 50, 100, 100, -20, 0},
		{60, 50, 100, 100, -10, -17},
		{90, 50, 100, 100, 0, -20},
		{180, 50, 100, 100, 20, 0},
		{0, -50, 100, 100, 20, 0},
		{270, 50, 100, 100, 0, 20},
		{270, -50, 100, 100, 0, -20},
		{90, -50, 100, 100, 0, 20},
	} {
		pf.On("TiltHeading").Once().Return(float64(i.angle), nil)
		pf.On("MoveForward", i.dist, i.vel).Once().Return(float64(i.distMoved), nil)

		x0, y0 := driver.GetXY()
		driver.Move(i.dist, i.vel)
		x, y := driver.GetXY()
		t.Logf("Angle: %v Scale: Vel:%v  %v X: %v->%v, Y:%v->%v", i.angle, i.vel, mapScale, x0, x, y0, y)
		if X-x != i.x || Y-y != i.y {
			t.Errorf("X expected: %v got: %v Y expected:%v got %v", X-i.x, x, Y-i.y, y)
		}

		driver.ResetMap()
	}

}

// Tests the update of a single sweep of lidar.
func TestUpdateOgridMap(t *testing.T) {
	readings := []datapt{
		/*	datapt{
			posture: 186.49,
			x:       500,
			y:       500,
			reading: []int32{115, 111, 110, 112, 114, 111, 113, 114, 118, 128, 127, 134, 140, 135, 121, 115, 107, 98, 99, 97, 76},
		}, */
		/*	datapt{
				posture: 150.31703686231944,
				x:       500,
				y:       500,
				reading: []int32{68, 100, 102, 113, 123, 142, 144, 139, 133, 131, 128, 127, 126, 128, 128,
					132, 134, 140, 142, 150, 156, 140, 133, 128, 144, 157, 177, 228, 200},
			},
			datapt{

				posture: 106.1509437176877,
				x:       500,
				y:       500,
				reading: []int32{128, 131, 136, 139, 145, 152, 161, 135, 125, 148, 155, 175, 177, 180,
					181, 210, 127, 126, 125, 182, 93, 93, 95, 97, 103, 103, 107, 112, 119},
			},
			datapt{
				posture: 73.49941462208987,
				x:       500,
				y:       500,
				reading: []int32{155, 235, 230, 224, 161, 132, 136, 90, 92, 94, 96, 96, 99, 103, 108,
					120, 148, 127, 123, 231, 217, 206, 194, 193, 186, 186, 182, 186, 133},
			},
			datapt{
				posture: 56.559022197907986,
				x:       500,
				y:       500,
				reading: []int32{135, 92, 92, 96, 95, 100, 102, 105, 111, 121, 124, 233, 220, 217, 202,
					195, 192, 186, 180, 183, 185, 128, 184, 190, 195, 201, 197, 215, 240},
			},*/
		datapt{
			posture: 0,
			x:       500,
			y:       500,
			reading: []int32{121, 112, 111, 106, 111, 132, 101, 86, 54, 51, 87, 61, 48},
		},
	}

	minAngle := 30
	shiftAngle := 10
	s := NewOgrid()
	s.ResetMap()

	for i := 0; i < len(readings); i++ {
		pos := readings[i].posture
		reading := readings[i].reading

		//	s.SetPos(x, y)
		if err := s.UpdateMap(reading, minAngle, shiftAngle, pos-90); err != nil {
			t.Errorf("error updating map %v", err)
		}

		bytes1, err := s.PrintMap()
		if err != nil {
			t.Errorf("Failed %v", err)
		}
		if err := ioutil.WriteFile("/Users/dkg/Downloads/ogrid.png", bytes1.Bytes(), os.ModePerm); err != nil {
			t.Errorf("failed %v", err)
		}
	}
}

func TestUpdateMap(t *testing.T) {

	_ = []datapt{
		datapt{
			posture: 150.31703686231944,
			x:       500,
			y:       500,
			reading: []int32{68, 100, 102, 113, 123, 142, 144, 139, 133, 131, 128, 127, 126, 128, 128,
				132, 134, 140, 142, 150, 156, 140, 133, 128, 144, 157, 177, 228, 200},
		},
		datapt{

			posture: 106.1509437176877,
			x:       500,
			y:       500,
			reading: []int32{128, 131, 136, 139, 145, 152, 161, 135, 125, 148, 155, 175, 177, 180,
				181, 210, 127, 126, 125, 182, 93, 93, 95, 97, 103, 103, 107, 112, 119},
		},
		datapt{
			posture: 73.49941462208987,
			x:       500,
			y:       500,
			reading: []int32{155, 235, 230, 224, 161, 132, 136, 90, 92, 94, 96, 96, 99, 103, 108,
				120, 148, 127, 123, 231, 217, 206, 194, 193, 186, 186, 182, 186, 133},
		},
		datapt{
			posture: 56.559022197907986,
			x:       500,
			y:       500,
			reading: []int32{135, 92, 92, 96, 95, 100, 102, 105, 111, 121, 124, 233, 220, 217, 202,
				195, 192, 186, 180, 183, 185, 128, 184, 190, 195, 201, 197, 215, 240},
		},
		datapt{
			posture: 154.07816426288912,
			x:       424,
			y:       460,
			reading: []int32{55, 94, 102, 109, 112, 137, 143, 132, 128, 123, 121, 120, 119, 118, 118,
				118, 119, 125, 124, 130, 137, 147, 155, 127, 119, 128, 148, 174, 203},
		},
		datapt{
			posture: 108.81609025448371,
			x:       424,
			y:       460,
			reading: []int32{120, 119, 121, 123, 126, 133, 142, 145, 136, 129, 119, 129, 147, 148,
				170, 173, 171, 195, 163, 122, 126, 136, 93, 96, 92, 99, 103, 105, 107},
		},
		datapt{
			posture: 76.39920001390699,
			x:       424,
			y:       460,
			reading: []int32{122, 147, 171, 216, 228, 203, 201, 129, 125, 91, 88, 94, 100, 100, 103,
				109, 134, 144, 165, 144, 130, 235, 226, 211, 205, 201, 123, 195, 194},
		},
		datapt{
			posture: 59.72262296121574,
			x:       424,
			y:       460,
			reading: []int32{205, 136, 196, 96, 94, 94, 100, 101, 106, 108, 116, 130, 134, 172, 229,
				227, 214, 204, 197, 152, 194, 196, 194, 194, 198, 202, 210, 203, 225},
		},
	}

	readings := []datapt{
		datapt{
			posture: 186.49,
			x:       500,
			y:       500,
			reading: []int32{115, 111, 110, 112, 114, 111, 113, 114, 118, 128, 127, 134, 140, 135, 121, 115, 107, 98, 99, 97, 76},
		},

		/*		datapt{
					posture: 160.35635040777913,
					x:       500,
					y:       500,
					reading: []int32{85, 118, 122, 134, 141, 159, 160, 145, 141, 134, 133, 129, 129, 128, 131, 129, 131, 132, 135, 141, 150, 122, 114, 116, 129, 137, 155, 174, 167},
				},
				datapt{
					posture: 110.31783509795588,
					x:       500,
					y:       500,
					reading: []int32{129, 129, 129, 132, 134, 140, 131, 122, 115, 127, 131, 134, 154, 152, 157, 160, 177, 108, 102, 107, 109, 70, 67, 67, 69, 75, 77, 82, 86},
				},

				datapt{
					posture: 79.17347188860884,
					x:       500,
					y:       500,
					reading: []int32{133, 155, 178, 205, 178, 110, 108, 104, 177, 63, 64, 66, 70, 73, 76, 101, 107, 115, 120, 140, 158, 106, 212, 206, 140, 192, 152, 191, 189},
				},

				datapt{
					posture: 64.35443304941813,
					x:       500,
					y:       500,
					reading: []int32{110, 105, 70, 62, 66, 67, 68, 76, 78, 84, 107, 120, 141, 163, 111, 116, 205, 200, 191, 159, 189, 185, 104, 187, 185, 140, 198, 207, 220},
				},

				/*	datapt{
						posture: 161.04086549421,
						x:       480,
						y:       494,
						reading: []int32{101, 106, 111, 119, 122, 136, 144, 135, 127, 122, 117, 115, 112, 111, 110, 107, 110, 110, 114, 114, 120, 125, 127, 116, 112, 102, 106, 123, 133},
					},
					datapt{
						posture: 114.500741415414,
						x:       480,
						y:       494,
						reading: []int32{112, 109, 109, 112, 114, 114, 120, 124, 131, 113, 108, 101, 104, 119, 130, 140, 150, 154, 147, 167, 166, 106, 104, 107, 117, 199, 79, 81, 90},
					},
					datapt{
						posture: 82.22251645431913,
						x:       480,
						y:       494,
						reading: []int32{112, 105, 123, 129, 151, 204, 155, 161, 172, 105, 103, 107, 165, 194, 183, 93, 100, 101, 112, 119, 129, 142, 167, 126, 131, 225, 170, 212, 210},
					},
					datapt{
						posture: 65.94051755124553,
						x:       480,
						y:       494,
						reading: []int32{206, 195, 193, 107, 104, 105, 194, 73, 74, 95, 107, 115, 120, 125, 139, 166, 182, 126, 176, 202, 201, 203, 200, 129, 132, 206, 169, 218, 223},
					},*/
	}

	minAngle := 60
	shiftAngle := 5
	s := NewOgrid()
	s.ResetMap()

	for i := 0; i < len(readings); i++ {
		pos := readings[i].posture
		//x := readings[i].x
		//y := readings[i].y
		reading := readings[i].reading

		//	s.SetPos(x, y)
		//	s.SetPos(500, 500)
		if err := s.UpdateMap(reading, minAngle, shiftAngle, pos-90); err != nil {
			t.Errorf("error updating map %v", err)
		}

		bytes1, err := s.PrintMap()
		if err != nil {
			t.Errorf("Failed %v", err)
		}
		if err := ioutil.WriteFile("updmap.png", bytes1.Bytes(), os.ModePerm); err != nil {
			t.Errorf("failed %v", err)
		}

	}

}
