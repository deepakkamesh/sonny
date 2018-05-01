package navigator

import (
	"testing"
)

func TestCell(t *testing.T) {
	x, y := calcCell(50, 20)

	if x != -5 || y != 2 {
		t.Errorf("Expected x,y = -5,2. Got %v %v", x, y)
	}
}

func TestImg(t *testing.T) {

	s := NewOgrid(nil)
	if err := s.UpdateMap(); err != nil {
		t.Errorf("Failed to update map")
	}
	if err := s.GenerateMap(); err != nil {
		t.Errorf("Failed")
	}
}
