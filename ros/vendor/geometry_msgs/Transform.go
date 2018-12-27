// Automatically generated from the message definition "geometry_msgs/Transform.msg"
package geometry_msgs

import (
	"bytes"

	"github.com/akio/rosgo/ros"
)

type _MsgTransform struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgTransform) Text() string {
	return t.text
}

func (t *_MsgTransform) Name() string {
	return t.name
}

func (t *_MsgTransform) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgTransform) NewMessage() ros.Message {
	m := new(Transform)
	m.Translation = Vector3{}
	m.Rotation = Quaternion{}
	return m
}

var (
	MsgTransform = &_MsgTransform{
		`# This represents the transform between two coordinate frames in free space.

Vector3 translation
Quaternion rotation
`,
		"geometry_msgs/Transform",
		"ac9eff44abf714214112b05d54a3cf9b",
	}
)

type Transform struct {
	Translation Vector3    `rosmsg:"translation:Vector3"`
	Rotation    Quaternion `rosmsg:"rotation:Quaternion"`
}

func (m *Transform) Type() ros.MessageType {
	return MsgTransform
}

func (m *Transform) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	if err = m.Translation.Serialize(buf); err != nil {
		return err
	}
	if err = m.Rotation.Serialize(buf); err != nil {
		return err
	}
	return err
}

func (m *Transform) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = m.Translation.Deserialize(buf); err != nil {
		return err
	}
	if err = m.Rotation.Deserialize(buf); err != nil {
		return err
	}
	return err
}
