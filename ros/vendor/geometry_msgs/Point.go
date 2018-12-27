
// Automatically generated from the message definition "geometry_msgs/Point.msg"
package geometry_msgs
import (
    "bytes"
    "encoding/binary"
    "github.com/akio/rosgo/ros"
)


type _MsgPoint struct {
    text string
    name string
    md5sum string
}

func (t *_MsgPoint) Text() string {
    return t.text
}

func (t *_MsgPoint) Name() string {
    return t.name
}

func (t *_MsgPoint) MD5Sum() string {
    return t.md5sum
}

func (t *_MsgPoint) NewMessage() ros.Message {
    m := new(Point)
	m.X = 0.0
	m.Y = 0.0
	m.Z = 0.0
    return m
}

var (
    MsgPoint = &_MsgPoint {
        `# This contains the position of a point in free space
float64 x
float64 y
float64 z
`,
        "geometry_msgs/Point",
        "4a842b65f413084dc2b10fb484ea7f17",
    }
)

type Point struct {
	X float64 `rosmsg:"x:float64"`
	Y float64 `rosmsg:"y:float64"`
	Z float64 `rosmsg:"z:float64"`
}

func (m *Point) Type() ros.MessageType {
	return MsgPoint
}

func (m *Point) Serialize(buf *bytes.Buffer) error {
    var err error = nil
    binary.Write(buf, binary.LittleEndian, m.X)
    binary.Write(buf, binary.LittleEndian, m.Y)
    binary.Write(buf, binary.LittleEndian, m.Z)
    return err
}


func (m *Point) Deserialize(buf *bytes.Reader) error {
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
