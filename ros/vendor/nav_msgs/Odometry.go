// Automatically generated from the message definition "nav_msgs/Odometry.msg"
package nav_msgs

import (
	"bytes"
	"encoding/binary"
	"geometry_msgs"
	"std_msgs"

	"github.com/akio/rosgo/ros"
)

type _MsgOdometry struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgOdometry) Text() string {
	return t.text
}

func (t *_MsgOdometry) Name() string {
	return t.name
}

func (t *_MsgOdometry) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgOdometry) NewMessage() ros.Message {
	m := new(Odometry)
	m.Header = std_msgs.Header{}
	m.ChildFrameId = ""
	m.Pose = geometry_msgs.PoseWithCovariance{}
	m.Twist = geometry_msgs.TwistWithCovariance{}
	return m
}

var (
	MsgOdometry = &_MsgOdometry{
		`# This represents an estimate of a position and velocity in free space.  
# The pose in this message should be specified in the coordinate frame given by header.frame_id.
# The twist in this message should be specified in the coordinate frame given by the child_frame_id
Header header
string child_frame_id
geometry_msgs/PoseWithCovariance pose
geometry_msgs/TwistWithCovariance twist
`,
		"nav_msgs/Odometry",
		"cd5e73d190d741a2f92e81eda573aca7",
	}
)

type Odometry struct {
	Header       std_msgs.Header                   `rosmsg:"header:Header"`
	ChildFrameId string                            `rosmsg:"child_frame_id:string"`
	Pose         geometry_msgs.PoseWithCovariance  `rosmsg:"pose:PoseWithCovariance"`
	Twist        geometry_msgs.TwistWithCovariance `rosmsg:"twist:TwistWithCovariance"`
}

func (m *Odometry) Type() ros.MessageType {
	return MsgOdometry
}

func (m *Odometry) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	if err = m.Header.Serialize(buf); err != nil {
		return err
	}
	binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.ChildFrameId))))
	buf.Write([]byte(m.ChildFrameId))
	if err = m.Pose.Serialize(buf); err != nil {
		return err
	}
	if err = m.Twist.Serialize(buf); err != nil {
		return err
	}
	return err
}

func (m *Odometry) Deserialize(buf *bytes.Reader) error {
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
	if err = m.Pose.Deserialize(buf); err != nil {
		return err
	}
	if err = m.Twist.Deserialize(buf); err != nil {
		return err
	}
	return err
}
