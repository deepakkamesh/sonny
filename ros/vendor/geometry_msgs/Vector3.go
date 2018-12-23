
// Automatically generated from the message definition "geometry_msgs/Vector3.msg"
package geometry_msgs
import (
    "bytes"
    "encoding/binary"
    "github.com/akio/rosgo/ros"
)


type _MsgVector3 struct {
    text string
    name string
    md5sum string
}

func (t *_MsgVector3) Text() string {
    return t.text
}

func (t *_MsgVector3) Name() string {
    return t.name
}

func (t *_MsgVector3) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgVector3) NewMessage() ros.Message {
    m := new(Vector3)
	m.X = 0.0
	m.Y = 0.0
	m.Z = 0.0
    return m
}

var (
    MsgVector3 = &_MsgVector3 {
        `# This represents a vector in free space. 
# It is only meant to represent a direction. Therefore, it does not
# make sense to apply a translation to it (e.g., when applying a 
# generic rigid transformation to a Vector3, tf2 will only apply the
# rotation). If you want your data to be translatable too, use the
# geometry_msgs/Point message instead.

float64 x
float64 y
float64 z`,
        "geometry_msgs/Vector3",
        "4a842b65f413084dc2b10fb484ea7f17",
    }
)

type Vector3 struct {
	X float64 `rosmsg:"x:float64"`
	Y float64 `rosmsg:"y:float64"`
	Z float64 `rosmsg:"z:float64"`
}

func (m *Vector3) Type() ros.MessageType {
	return MsgVector3
}

func (m *Vector3) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    binary.Write(buf, binary.LittleEndian, m.X)
    binary.Write(buf, binary.LittleEndian, m.Y)
    binary.Write(buf, binary.LittleEndian, m.Z)
    return err
}


func (m *Vector3) Deserialize(buf *bytes.Reader) error {
    var err error = nil
    if err = binary.Read(buf, binary.LittleEndian, &m.X); err != nil {
        return err
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.Y); err != nil {
        return err
    }
    if err = binary.Read(buf, binary.LittleEndian, &m.Z); err != nil {
        return err
    }
    return err
}
