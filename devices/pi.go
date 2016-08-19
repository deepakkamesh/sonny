package devices

type Pi struct {
}

func NewPi() *Pi {
	return &Pi{}
}

func (m *Pi) GetDistance() (uint16, error) {

	return 100, nil
}
