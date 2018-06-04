package navigator

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"testing"
)

/*
func TestPattern(t *testing.T) {

	// Test data.
	//	reading := []int32{137, 139, 58, 57, 61, 60, 154, 159, 142, 143, 136, 43, 48, 46, 113, 112, 105, 109, 98, 97, 101, 99, 97, 94, 83, 83, 90, 89, 89, 87, 88, 87, 84, 89, 85, 85, 89, 84, 87, 86, 85, 88, 86, 85, 86, 89, 85, 88, 89, 89, 95, 92, 95, 96, 94, 86, 85, 89, 89, 91, 108, 109, 99, 98, 106, 109, 73, 72, 69, 67, 68}

	reading := []int32{40, 38, 36, 33, 30, 27, 29, 28, 27, 19, 23, 19, 19, 15, 14, 12, 13, 13, 14, 12, 9, 9, 10, 12, 7, 10, 7, 12, 8, 13, 12, 8, 11, 9, 8, 12, 9, 12, 14, 9, 12, 14, 12, 13, 16, 15, 16, 18, 22, 22, 19, 23, 20, 24, 22, 24, 27, 29, 30, 33, 35, 39, 37, 47, 47, 52, 44, 48, 40, 33, 32}
	//	reading := []int32{43, 28, 20, 16, 10, 10, 8, 13, 11, 17, 22, 27, 35, 51, 70}
	//	reading := []int32{46, 27, 20, 19, 12, 13, 9, 9, 13, 16, 23, 30, 39, 49, 31}
	minAngle := 20
	shiftAngle := 2
	posture := 300
	fmt.Println(len(reading))
	s := NewOgrid()
	s.ResetMap()
	if err := s.UpdateMap(reading, minAngle, shiftAngle, float64(posture)); err != nil {
		t.Errorf("error updating map %v", err)
	}
	bytes1, err := s.GenerateMap()
	bytes2, err := s.GenCompressMap()
	if err != nil {
		t.Errorf("Failed %v", err)
	}

	if err := ioutil.WriteFile("out1.png", bytes1.Bytes(), os.ModePerm); err != nil {
		t.Errorf("failed %v", err)
	}
	if err := ioutil.WriteFile("out2.png", bytes2.Bytes(), os.ModePerm); err != nil {
		t.Errorf("failed %v", err)
	}
} */

type datapt struct {
	posture float64
	x       int
	y       int
	reading []int32
}

/*
func Test360(t *testing.T) {

	reading := map[float64][]int32{
		182: []int32{59, 61, 68, 72, 67, 71, 71, 78, 306, 289, 174,
			155, 142, 128, 121, 112, 152, 152, 147, 133, 131, 157, 137, 132, 138, 135,
			134, 153, 175},
		135: []int32{337, 167, 150, 129, 127, 155, 154, 153, 134, 132, 158,
			146, 135, 135, 143, 139, 137, 152, 162, 179, 190, 174, 191, 187, 206, 78, 82, 78, 81},
		88: []int32{137, 155, 154, 173, 179, 175, 185, 205, 94, 75, 74, 98, 82, 83, 91, 97, 105,
			110, 112, 116, 123, 143, 96, 208, 195, 189, 183, 172, 107},
	}
	order := []float64{182, 135, 88}
	//Posture: 34.33551113931616 Reading:[99 115 119 167 145 174 104 165 163 163 162 161 159 112 169 174 181 188 177 198 153 238 199 241 253 253 263 270 61]


	colors := []color.RGBA{
		color.RGBA{244, 26, 26, 255},
		color.RGBA{26, 26, 244, 255},
		color.RGBA{0, 51, 0, 255},
	}

	minAngle := 20
	shiftAngle := 5
	s := NewOgrid()
	s.ResetMap()

	for i := 0; i < len(order); i++ {
		ang := order[i]
		//fmt.Println("Readings", reading[ang], len(reading[ang]))
		//	if i != 135.40983201669695 {
		//	continue
		//	}
		if err := s.UpdateMap(reading[ang], minAngle, shiftAngle, 360-ang, colors[i]); err != nil {
			t.Errorf("error updating map %v", err)
		}
		bytes1, err := s.GenerateMap()
		bytes2, err := s.GenCompressMap()
		if err != nil {
			t.Errorf("Failed %v", err)
		}
		fn1 := fmt.Sprintf("out1-%v.png", ang)
		fn2 := fmt.Sprintf("out2-%v.png", ang)

		if err := ioutil.WriteFile(fn1, bytes1.Bytes(), os.ModePerm); err != nil {
			t.Errorf("failed %v", err)
		}
		if err := ioutil.WriteFile(fn2, bytes2.Bytes(), os.ModePerm); err != nil {
			t.Errorf("failed %v", err)
		}

	}
}*/

func Test3602(t *testing.T) {

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

		datapt{
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
		},
	}

	colors := []color.RGBA{
		color.RGBA{244, 26, 26, 255},
		color.RGBA{26, 26, 244, 255},
		color.RGBA{80, 51, 30, 255},
		color.RGBA{0, 51, 0, 255},
		color.RGBA{24, 166, 26, 255},
		color.RGBA{26, 0, 100, 255},
		color.RGBA{0, 100, 30, 255},
		color.RGBA{125, 51, 49, 255},
	}

	minAngle := 20
	shiftAngle := 5
	s := NewOgrid()
	s.ResetMap()

	for i := 0; i < len(readings); i++ {
		pos := readings[i].posture
		x := readings[i].x
		y := readings[i].y
		reading := readings[i].reading

		s.SetPos(x, y)
		s.SetPos(500, 500)

		if err := s.UpdateMap(reading, minAngle, shiftAngle, 360-pos, colors[i]); err != nil {
			t.Errorf("error updating map %v", err)
		}

		bytes2, err := s.scaledMap()
		bytes1, err := s.normalMap()
		if err != nil {
			t.Errorf("Failed %v", err)
		}
		fn1 := fmt.Sprintf("out1-%v-%v-%v.png", int(pos), x, y)
		fn2 := fmt.Sprintf("out2-%v-%v-%v.png", int(pos), x, y)

		if err := ioutil.WriteFile(fn1, bytes1.Bytes(), os.ModePerm); err != nil {
			t.Errorf("failed %v", err)
		}

		if err := ioutil.WriteFile(fn2, bytes2.Bytes(), os.ModePerm); err != nil {
			t.Errorf("failed %v", err)
		}

	}

}
