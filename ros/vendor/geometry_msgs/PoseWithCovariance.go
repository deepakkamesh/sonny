// Automatically generated from the message definition "geometry_msgs/PoseWithCovariance.msg"
package geometry_msgs

import (
	"bytes"
	"encoding/binary"

	"github.com/akio/rosgo/ros"
)

type _MsgPoseWithCovariance struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgPoseWithCovariance) Text() string {
	return t.text
}

func (t *_MsgPoseWithCovariance) Name() string {
	return t.name
}

func (t *_MsgPoseWithCovariance) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgPoseWithCovariance) NewMessage() ros.Message {
	m := new(PoseWithCovariance)
	m.Pose = Pose{}
	for i := 0; i < 36; i++ {
		m.Covariance[i] = 0.0
	}
	return m
}

var (
	MsgPoseWithCovariance = &_MsgPoseWithCovariance{
		`# This represents a pose in free space with uncertainty.

Pose pose

# Row-major representation of the 6x6 covariance matrix
# The orientation parameters use a fixed-axis representation.
# In order, the parameters are:
# (x, y, z, rotation about X axis, rotation about Y axis, rotation about Z axis)
float64[36] covariance
`,
		"geometry_msgs/PoseWithCovariance",
		"620fbdfad1a8bc89c4d876fdb4ffa7f8",
	}
)

type PoseWithCovariance struct {
	Pose       Pose        `rosmsg:"pose:Pose"`
	Covariance [36]float64 `rosmsg:"covariance:float64[36]"`
}

func (m *PoseWithCovariance) Type() ros.MessageType {
	return MsgPoseWithCovariance
}

func (m *PoseWithCovariance) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	if err = m.Pose.Serialize(buf); err != nil {
		return err
	}
	for _, e := range m.Covariance {
		binary.Write(buf, binary.LittleEndian, e)
	}
	return err
}

func (m *PoseWithCovariance) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = m.Pose.Deserialize(buf); err != nil {
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
