// Automatically generated from the message definition "geometry_msgs/TwistWithCovariance.msg"
package geometry_msgs

import (
	"bytes"
	"encoding/binary"

	"github.com/akio/rosgo/ros"
)

type _MsgTwistWithCovariance struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgTwistWithCovariance) Text() string {
	return t.text
}

func (t *_MsgTwistWithCovariance) Name() string {
	return t.name
}

func (t *_MsgTwistWithCovariance) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgTwistWithCovariance) NewMessage() ros.Message {
	m := new(TwistWithCovariance)
	m.Twist = Twist{}
	for i := 0; i < 36; i++ {
		m.Covariance[i] = 0.0
	}
	return m
}

var (
	MsgTwistWithCovariance = &_MsgTwistWithCovariance{
		`# This expresses velocity in free space with uncertainty.

Twist twist

# Row-major representation of the 6x6 covariance matrix
# The orientation parameters use a fixed-axis representation.
# In order, the parameters are:
# (x, y, z, rotation about X axis, rotation about Y axis, rotation about Z axis)
float64[36] covariance
`,
		"geometry_msgs/TwistWithCovariance",
		"eb9275f4ce74337e2a9cd750980e8d8c",
	}
)

type TwistWithCovariance struct {
	Twist      Twist       `rosmsg:"twist:Twist"`
	Covariance [36]float64 `rosmsg:"covariance:float64[36]"`
}

func (m *TwistWithCovariance) Type() ros.MessageType {
	return MsgTwistWithCovariance
}

func (m *TwistWithCovariance) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	if err = m.Twist.Serialize(buf); err != nil {
		return err
	}
	for _, e := range m.Covariance {
		binary.Write(buf, binary.LittleEndian, e)
	}
	return err
}

func (m *TwistWithCovariance) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = m.Twist.Deserialize(buf); err != nil {
		return err
	}
	{
		for i := 0; i < 36; i++ {
			if err = binary.Read(buf, binary.LittleEndian, &m.Covariance[i]); err != nil {
				return err
			}
		}
	}
	return err
}
