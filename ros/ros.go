//go:generate gengo msg geometry_msgs/Vector3
//go:generate gengo msg geometry_msgs/Twist

package ros

import (
	"geometry_msgs"
	"os"

	roslib "github.com/akio/rosgo/ros"
	"github.com/deepakkamesh/go-roomba/constants"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/golang/glog"
)

type ROSConn struct {
	node  roslib.Node
	rover devices.Platform
}

func NewRos(rover devices.Platform) *ROSConn {
	return &ROSConn{
		rover: rover,
	}
}

func (n *ROSConn) Shutdown() {
	n.node.Shutdown()
}

func (n *ROSConn) StartNode(name string) error {

	node, err := roslib.NewNode(name, os.Args)
	if err != nil {
		return err
	}
	n.node = node
	node.Logger().SetSeverity(roslib.LogLevelInfo)
	return nil
}

func (n *ROSConn) Spinup() {
	go func() {
		for n.node.OK() {
			n.node.SpinOnce()
		}
	}()
}

func (n *ROSConn) Twist(msg *geometry_msgs.Twist) {
	glog.V(2).Infof("Twist msg: %+v", *msg)
	// linear vel = m/s | angular vel = rad/s
	// ref: https://snapcraft.io/blog/your-first-robot-the-driver-4-5
	lX := msg.Linear.X
	aZ := msg.Angular.Z

	vL := lX - aZ*constants.RoombaRadius
	vR := lX + aZ*constants.RoombaRadius
	if err := n.rover.DirectDrive(int16(vR*1000), int16(vL*1000)); err != nil {
		glog.Warningf("Failed to drive roomba:%v", err)
	}
}

// ListenCmdVel listens to the cmd_vel message and controls the rover.
func (n *ROSConn) ListenCmdVel() {
	n.node.NewSubscriber("/cmd_vel", geometry_msgs.MsgTwist, n.Twist)
}
