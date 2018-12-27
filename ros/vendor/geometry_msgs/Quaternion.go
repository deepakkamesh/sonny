
// Automatically generated from the message definition "geometry_msgs/Quaternion.msg"
package geometry_msgs
import (
    "bytes"
    "encoding/binary"
    "github.com/akio/rosgo/ros"
)


type _MsgQuaternion struct {
    text string
    name string
    md5sum string
}

func (t *_MsgQuaternion) Text() string {
    return t.text
}

func (t *_MsgQuaternion) Name() string {
    return t.name
}

func (t *_MsgQuaternion) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgQuaternion) NewMessage() ros.Message {
    m := new(Quaternion)
	m.X = 0.0
	m.Y = 0.0
	m.Z = 0.0
	m.W = 0.0
    return m
}

var (
    MsgQuaternion = &_MsgQuaternion {
        `# This represents an orientation in free space in quaternion form.

float64 x
float64 y
float64 z
float64 w
`,
        "geometry_msgs/Quaternion",
        "a779879fadf0160734f906b8c19c7004",
    }
)

type Quaternion struct {
	X float64 `rosmsg:"x:float64"`
	Y float64 `rosmsg:"y:float64"`
	Z float64 `rosmsg:"z:float64"`
	W float64 `rosmsg:"w:float64"`
}

func (m *Quaternion) Type() ros.MessageType {
	return MsgQuaternion
}

func (m *Quaternion) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    binary.Write(buf, binary.LittleEndian, m.X)
    binary.Write(buf, binary.LittleEndian, m.Y)
    binary.Write(buf, binary.LittleEndian, m.Z)
    binary.Write(buf, binary.LittleEndian, m.W)
    return err
}


func (m *Quaternion) Deserialize(buf *bytes.Reader) error {
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
    if err = binary.Read(buf, binary.LittleEndian, &m.W); err != nil {
        return err
    }
    return err
}
