// Automatically generated from the message definition "tf2_msgs/TFMessage.msg"
package tf2_msgs

import (
	"bytes"
	"encoding/binary"
	"geometry_msgs"

	"github.com/akio/rosgo/ros"
)

type _MsgTFMessage struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgTFMessage) Text() string {
	return t.text
}

func (t *_MsgTFMessage) Name() string {
	return t.name
}

func (t *_MsgTFMessage) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgTFMessage) NewMessage() ros.Message {
	m := new(TFMessage)
	m.Transforms = []geometry_msgs.TransformStamped{}
	return m
}

var (
	MsgTFMessage = &_MsgTFMessage{
		`geometry_msgs/TransformStamped[] transforms
`,
		"tf2_msgs/TFMessage",
		"94810edda583a504dfda3829e70d7eec",
	}
)

type TFMessage struct {
	Transforms []geometry_msgs.TransformStamped `rosmsg:"transforms:TransformStamped[]"`
}

func (m *TFMessage) Type() ros.MessageType {
	return MsgTFMessage
}

func (m *TFMessage) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	binary.Write(buf, binary.LittleEndian, uint32(len(m.Transforms)))
	for _, e := range m.Transforms {
		if err = e.Serialize(buf); err != nil {
			return err
		}
	}
	return err
}

func (m *TFMessage) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	{
		var size uint32
		if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
			return err
		}
		m.Transforms = make([]geometry_msgs.TransformStamped, int(size))
		for i := 0; i < int(size); i++ {
			if err = m.Transforms[i].Deserialize(buf); err != nil {
				return err
			}
		}
	}
	return err
}
