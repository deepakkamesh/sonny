////go:generate gengo msg
//go:generate gengo msg geometry_msgs/Quaternion
//go:generate gengo msg geometry_msgs/Point
//go:generate gengo msg geometry_msgs/Pose
//go:generate gengo msg geometry_msgs/PoseWithCovariance
//go:generate gengo msg geometry_msgs/TwistWithCovariance
//go:generate gengo msg geometry_msgs/Vector3
//go:generate gengo msg geometry_msgs/Twist
//go:generate gengo msg nav_msgs/Odometry
// Package uses akio/rosgo. But it needs a patch to ensure fixed length arrays
// are handled properly. Apply https://github.com/akio/rosgo/pull/18.
package ros

import (
	"fmt"
	"math"
	"os"
	"time"

	"geometry_msgs"
	"nav_msgs"
	"std_msgs"
	"tf2_msgs"

	roslib "github.com/akio/rosgo/ros"
	"github.com/deepakkamesh/go-roomba/constants"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/golang/glog"
	"github.com/westphae/quaternion"
	"gonum.org/v1/gonum/mat"
)

type Pose struct {
	x, y, yaw  float64
	covariance [9]float64
}
type Vel = Pose

const (
	EPS    float64 = 0.0001
	TWO_PI float64 = 6.28318
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

	node, e := roslib.NewNode(name, os.Args)
	if e != nil {
		return e
	}
	n.node = node
	node.Logger().SetSeverity(roslib.LogLevelInfo)

	// Subscribers.
	n.node.NewSubscriber("/cmd_vel", geometry_msgs.MsgTwist, n.Twist)

	// Publishers.
	pubOdom := node.NewPublisher("/odom", nav_msgs.MsgOdometry)
	pubTf := node.NewPublisher("/tf", tf2_msgs.MsgTFMessage)
	pubTest := node.NewPublisher("/chatter", std_msgs.MsgString)
	n.publishOdom(pubOdom, pubTf, pubTest)
	return nil

}

func (n *ROSConn) publishOdom(pubOdom roslib.Publisher, pubTf roslib.Publisher, pubTest roslib.Publisher) {
	data, _ := n.rover.GetRoombaTelemetry()

	prevTicksLeft := int(data[constants.SENSOR_LEFT_ENCODER])
	prevTicksRight := int(data[constants.SENSOR_RIGHT_ENCODER])
	prevTmstmp := time.Now()
	pose := Pose{}
	vel := Vel{}
	seq := uint32(0)
	var deltaX, deltaY, totalLeftDist, totalRightDist float64
	poseCovar := mat.NewDense(3, 3, nil)

	// Start update loop.
	go func() {
		for n.node.OK() {
			time.Sleep(100 * time.Millisecond) // 10 Hz.

			n.node.SpinOnce()

			// Test
			var msg std_msgs.String
			msg.Data = fmt.Sprintf("hello %s", time.Now().String())
			pubTest.Publish(&msg)

			// Get current time.
			tmstmp := time.Now()
			dt := tmstmp.Sub(prevTmstmp).Seconds()

			// Get cumulative ticks (wraps around at 65535).
			data, _ := n.rover.GetRoombaTelemetry()
			totalTicksLeft := int(data[constants.SENSOR_LEFT_ENCODER])
			totalTicksRight := int(data[constants.SENSOR_RIGHT_ENCODER])

			// Compute ticks since last update
			ticksLeft := totalTicksLeft - prevTicksLeft
			ticksRight := totalTicksRight - prevTicksRight
			prevTicksLeft = totalTicksLeft
			prevTicksRight = totalTicksRight

			// Handle wrap around
			if math.Abs(float64(ticksLeft)) > 0.9*constants.MAX_ENCODER_TICKS {
				ticksLeft = (ticksLeft % constants.MAX_ENCODER_TICKS) + 1
			}
			if math.Abs(float64(ticksRight)) > 0.9*constants.MAX_ENCODER_TICKS {
				ticksRight = (ticksRight % constants.MAX_ENCODER_TICKS) + 1
			}

			// Compute distance travelled by each wheel (in meters).
			leftWheelDist := float64(ticksLeft) * math.Pi * (constants.WheelDia / 1000) / constants.TicksPerRevolution
			rightWheelDist := float64(ticksRight) * math.Pi * (constants.WheelDia / 1000) / constants.TicksPerRevolution
			deltaDist := (rightWheelDist + leftWheelDist) / 2.0

			wheelDistDiff := rightWheelDist - leftWheelDist
			deltaYaw := wheelDistDiff / (constants.AxleDist / 1000)

			//	measuredLeftVel := leftWheelDist / dt
			//	measuredRightVel := rightWheelDist / dt

			// Moving straight.
			if math.Abs(wheelDistDiff) < EPS {
				deltaX = deltaDist * math.Cos(pose.yaw)
				deltaY = deltaDist * math.Sin(pose.yaw)
			} else {
				turnRadius := (constants.RoombaRadius / 1000) * (leftWheelDist + rightWheelDist) / wheelDistDiff
				deltaX = turnRadius * (math.Sin(pose.yaw+deltaYaw) - math.Sin(pose.yaw))
				deltaY = -turnRadius * (math.Cos(pose.yaw+deltaYaw) - math.Cos(pose.yaw))
			}

			totalLeftDist += leftWheelDist
			totalRightDist += rightWheelDist

			if math.Abs(dt) > EPS {
				vel.x = deltaDist / dt
				vel.y = 0
				vel.yaw = deltaYaw / dt
			} else {
				vel.x = 0
				vel.y = 0
				vel.yaw = 0
			}

			// Update covariances.
			// Ref: "Introduction to Autonomous Mobile Robots" (Siegwart 2004, page 189).
			// TODO: Perform experiments to find these nondeterministic parameters.
			kr := 1.0
			kl := 1.0
			cosYawAndHalfDelta := math.Cos(pose.yaw + (deltaYaw / 2.0)) // deltaX?
			sinYawAndHalfDelta := math.Sin(pose.yaw + (deltaYaw / 2.0)) // deltaY?
			distOverTwoWB := deltaDist / ((constants.AxleDist / 1000) * 2.0)

			invCovar := mat.NewDense(2, 2, nil)
			invCovar.Set(0, 0, kr*math.Abs(rightWheelDist))
			invCovar.Set(0, 1, 0)
			invCovar.Set(1, 0, 0)
			invCovar.Set(1, 1, kl*math.Abs(leftWheelDist))

			Finc := mat.NewDense(3, 2, nil)
			Finc.Set(0, 0, (cosYawAndHalfDelta/2.0)-(distOverTwoWB*sinYawAndHalfDelta))
			Finc.Set(0, 1, (cosYawAndHalfDelta/2.0)+(distOverTwoWB*sinYawAndHalfDelta))
			Finc.Set(1, 0, (sinYawAndHalfDelta/2.0)+(distOverTwoWB*cosYawAndHalfDelta))
			Finc.Set(1, 1, (sinYawAndHalfDelta/2.0)-(distOverTwoWB*cosYawAndHalfDelta))
			Finc.Set(2, 0, (1.0 / (constants.AxleDist / 1000)))
			Finc.Set(2, 1, (-1.0 / (constants.AxleDist / 1000)))
			FincT := Finc.T() // TODO: Verify is this is boost::numeric::ublas::trans().

			Fp := mat.NewDense(3, 3, nil)
			Fp.Set(0, 0, 1.0)
			Fp.Set(0, 1, 0.0)
			Fp.Set(0, 2, (-deltaDist)*sinYawAndHalfDelta)
			Fp.Set(1, 0, 0.0)
			Fp.Set(1, 1, 1.0)
			Fp.Set(1, 2, deltaDist*cosYawAndHalfDelta)
			Fp.Set(2, 0, 0.0)
			Fp.Set(2, 1, 0.0)
			Fp.Set(2, 2, 1.0)
			FpT := Fp.T()

			var velCovar1 mat.Dense
			velCovar1.Mul(invCovar, FincT)
			var velCovar mat.Dense
			velCovar.Mul(Finc, &velCovar1)

			vel.covariance[0] = velCovar.At(0, 0)
			vel.covariance[1] = velCovar.At(0, 1)
			vel.covariance[2] = velCovar.At(0, 2)
			vel.covariance[3] = velCovar.At(1, 0)
			vel.covariance[4] = velCovar.At(1, 1)
			vel.covariance[5] = velCovar.At(1, 2)
			vel.covariance[6] = velCovar.At(2, 0)
			vel.covariance[7] = velCovar.At(2, 1)
			vel.covariance[8] = velCovar.At(2, 2)

			var poseCovarTmp mat.Dense
			poseCovarTmp.Mul(poseCovar, FpT)
			poseCovarTmp.Mul(Fp, &poseCovarTmp)
			poseCovar.Add(&poseCovarTmp, &velCovar)

			pose.covariance[0] = poseCovar.At(0, 0)
			pose.covariance[1] = poseCovar.At(0, 1)
			pose.covariance[2] = poseCovar.At(0, 2)
			pose.covariance[3] = poseCovar.At(1, 0)
			pose.covariance[4] = poseCovar.At(1, 1)
			pose.covariance[5] = poseCovar.At(1, 2)
			pose.covariance[6] = poseCovar.At(2, 0)
			pose.covariance[7] = poseCovar.At(2, 1)
			pose.covariance[8] = poseCovar.At(2, 2)

			// Update pose.
			pose.x += deltaX
			pose.y += deltaY
			pose.yaw = normalizeAngle(pose.yaw + deltaYaw)

			/************* Publish odometry ************/
			quat := quaternion.FromEuler(0, 0, pose.yaw)

			// Populate odomMsg.
			odomMsg := nav_msgs.Odometry{
				Header: std_msgs.Header{
					Seq:     seq,
					Stamp:   roslib.NewTime(uint32(time.Now().Unix()), uint32(time.Now().UnixNano())),
					FrameId: "odom_frame",
				},
				ChildFrameId: "base_link",
				Pose: geometry_msgs.PoseWithCovariance{
					Pose: geometry_msgs.Pose{
						Position: geometry_msgs.Point{
							X: pose.x,
							Y: pose.y,
						},
						Orientation: geometry_msgs.Quaternion{
							X: quat.X,
							Y: quat.Y,
							Z: quat.Z,
							W: quat.W,
						},
					},
				},
				Twist: geometry_msgs.TwistWithCovariance{
					Twist: geometry_msgs.Twist{
						Linear: geometry_msgs.Vector3{
							X: vel.x,
							Y: vel.y,
						},
						Angular: geometry_msgs.Vector3{
							Z: vel.yaw,
						},
					},
				},
			}

			//	glog.Infof("DeltaDist %+v , %v, %v", odomMsg, vel.y, vel.yaw)
			// Update covariances.
			odomMsg.Pose.Covariance[0] = pose.covariance[0]
			odomMsg.Pose.Covariance[1] = pose.covariance[1]
			odomMsg.Pose.Covariance[5] = pose.covariance[2]
			odomMsg.Pose.Covariance[6] = pose.covariance[3]
			odomMsg.Pose.Covariance[7] = pose.covariance[4]
			odomMsg.Pose.Covariance[11] = pose.covariance[5]
			odomMsg.Pose.Covariance[30] = pose.covariance[6]
			odomMsg.Pose.Covariance[31] = pose.covariance[7]
			odomMsg.Pose.Covariance[35] = pose.covariance[8]
			odomMsg.Twist.Covariance[0] = vel.covariance[0]
			odomMsg.Twist.Covariance[1] = vel.covariance[1]
			odomMsg.Twist.Covariance[5] = vel.covariance[2]
			odomMsg.Twist.Covariance[6] = vel.covariance[3]
			odomMsg.Twist.Covariance[7] = vel.covariance[4]
			odomMsg.Twist.Covariance[11] = vel.covariance[5]
			odomMsg.Twist.Covariance[30] = vel.covariance[6]
			odomMsg.Twist.Covariance[31] = vel.covariance[7]
			odomMsg.Twist.Covariance[35] = vel.covariance[8]

			// Publish odom Msg.
			pubOdom.Publish(&odomMsg)

			// Publish Odom TF.
			odomTF := geometry_msgs.TransformStamped{
				Header: std_msgs.Header{
					Seq:     seq,
					Stamp:   roslib.NewTime(uint32(time.Now().Unix()), uint32(time.Now().UnixNano())),
					FrameId: "odom_frame",
				},
				ChildFrameId: "base_link",
				Transform: geometry_msgs.Transform{
					Translation: geometry_msgs.Vector3{
						X: pose.x,
						Y: pose.y,
					},
					Rotation: geometry_msgs.Quaternion{
						X: quat.X,
						Y: quat.Y,
						Z: quat.Z,
						W: quat.W,
					},
				},
			}
			TFMsg := tf2_msgs.TFMessage{
				Transforms: []geometry_msgs.TransformStamped{odomTF},
			}
			pubTf.Publish(&TFMsg)

			// The END.
			prevTmstmp = tmstmp
			seq++
			n.node.SpinOnce()
		}
	}()

}

func normalizeAngle(angle float64) float64 {
	a := angle
	for a < -math.Pi {
		a += TWO_PI
	}
	for a > math.Pi {
		a -= TWO_PI
	}
	return a
}

func (n *ROSConn) Twist(msg *geometry_msgs.Twist) {
	// linear vel = m/s | angular vel = rad/s
	// ref: https://snapcraft.io/blog/your-first-robot-the-driver-4-5
	lX := msg.Linear.X
	aZ := msg.Angular.Z

	vL := lX - aZ*(constants.RoombaRadius/1000)
	vR := lX + aZ*(constants.RoombaRadius/1000)

	glog.V(2).Infof("TwistMsg: %+v => Speed Left:%.2f Right:%.2f", msg, vL, vR)
	if err := n.rover.DirectDrive(int16(vR*1000), int16(vL*1000)); err != nil {
		glog.Warningf("Failed to drive roomba: %v", err)
	}
}
