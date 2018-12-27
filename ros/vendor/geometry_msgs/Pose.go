// Automatically generated from the message definition "geometry_msgs/Pose.msg"
package geometry_msgs

import (
	"bytes"

	"github.com/akio/rosgo/ros"
)

type _MsgPose struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgPose) Text() string {
	return t.text
}

func (t *_MsgPose) Name() string {
	return t.name
}

func (t *_MsgPose) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgPose) NewMessage() ros.Message {
	m := new(Pose)
	m.Position = Point{}
	m.Orientation = Quaternion{}
	return m
}

var (
	MsgPose = &_MsgPose{
		`# A representation of pose in free space, composed of position and orientation. 
Point position
Quaternion orientation
`,
		"geometry_msgs/Pose",
		"e45d45a5a1ce597b249e23fb30fc871f",
	}
)

type Pose struct {
	Position    Point      `rosmsg:"position:Point"`
	Orientation Quaternion `rosmsg:"orientation:Quaternion"`
}

func (m *Pose) Type() ros.MessageType {
	return MsgPose
}

func (m *Pose) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	if err = m.Position.Serialize(buf); err != nil {
		return err
	}
	if err = m.Orientation.Serialize(buf); err != nil {
		return err
	}
	return err
}

func (m *Pose) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = m.Position.Deserialize(buf); err != nil {
		return err
	}
	if err = m.Orientation.Deserialize(buf); err != nil {
		return err
	}
	return err
}
