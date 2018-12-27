// Automatically generated from the message definition "geometry_msgs/TransformStamped.msg"
package geometry_msgs

import (
	"bytes"
	"encoding/binary"
	"std_msgs"

	"github.com/akio/rosgo/ros"
)

type _MsgTransformStamped struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgTransformStamped) Text() string {
	return t.text
}

func (t *_MsgTransformStamped) Name() string {
	return t.name
}

func (t *_MsgTransformStamped) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgTransformStamped) NewMessage() ros.Message {
	m := new(TransformStamped)
	m.Header = std_msgs.Header{}
	m.ChildFrameId = ""
	m.Transform = Transform{}
	return m
}

var (
	MsgTransformStamped = &_MsgTransformStamped{
		`# This expresses a transform from coordinate frame header.frame_id
# to the coordinate frame child_frame_id
#
# This message is mostly used by the 
# <a href="http://wiki.ros.org/tf">tf</a> package. 
# See its documentation for more information.

Header header
string child_frame_id # the frame id of the child frame
Transform transform
`,
		"geometry_msgs/TransformStamped",
		"b5764a33bfeb3588febc2682852579b0",
	}
)

type TransformStamped struct {
	Header       std_msgs.Header `rosmsg:"header:Header"`
	ChildFrameId string          `rosmsg:"child_frame_id:string"`
	Transform    Transform       `rosmsg:"transform:Transform"`
}

func (m *TransformStamped) Type() ros.MessageType {
	return MsgTransformStamped
}

func (m *TransformStamped) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	if err = m.Header.Serialize(buf); err != nil {
		return err
	}
	binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.ChildFrameId))))
	buf.Write([]byte(m.ChildFrameId))
	if err = m.Transform.Serialize(buf); err != nil {
		return err
	}
	return err
}

func (m *TransformStamped) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = m.Header.Deserialize(buf); err != nil {
		return err
	}
	{
		var size uint32
		if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
			return err
		}
		data := make([]byte, int(size))
		if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
			return err
		}
		m.ChildFrameId = string(data)
	}
	if err = m.Transform.Deserialize(buf); err != nil {
		return err
	}
	return err
}
