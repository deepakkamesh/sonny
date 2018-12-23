// Automatically generated from the message definition "geometry_msgs/Twist.msg"
package geometry_msgs

import (
	"bytes"

	"github.com/akio/rosgo/ros"
)

type _MsgTwist struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgTwist) Text() string {
	return t.text
}

func (t *_MsgTwist) Name() string {
	return t.name
}

func (t *_MsgTwist) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgTwist) NewMessage() ros.Message {
	m := new(Twist)
	m.Linear = Vector3{}
	m.Angular = Vector3{}
	return m
}

var (
	MsgTwist = &_MsgTwist{
		`# This expresses velocity in free space broken into its linear and angular parts.
Vector3  linear
Vector3  angular
`,
		"geometry_msgs/Twist",
		"9f195f881246fdfa2798d1d3eebca84a",
	}
)

type Twist struct {
	Linear  Vector3 `rosmsg:"linear:Vector3"`
	Angular Vector3 `rosmsg:"angular:Vector3"`
}

func (m *Twist) Type() ros.MessageType {
	return MsgTwist
}

func (m *Twist) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	if err = m.Linear.Serialize(buf); err != nil {
		return err
	}
	if err = m.Angular.Serialize(buf); err != nil {
		return err
	}
	return err
}

func (m *Twist) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = m.Linear.Deserialize(buf); err != nil {
		return err
	}
	if err = m.Angular.Deserialize(buf); err != nil {
		return err
	}
	return err
}
