package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	pb "github.com/deepakkamesh/sonny/sonny"
	google_pb "github.com/golang/protobuf/ptypes/empty"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:2233", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to rpc endpoint: %v", err)
	}
	defer conn.Close()
	ctrl := pb.NewDevicesRPCClient(conn)

	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:    "PingController",
			Aliases: []string{"ping"},
			Usage:   "Ping the controller.",
			Action: func(c *cli.Context) error {
				if _, err := ctrl.Ping(context.Background(), &google_pb.Empty{}); err != nil {
					return err
				}
				log.Printf("Pinging controller successful")
				return nil
			},
		},
		{
			Name:    "PIRDetect",
			Aliases: []string{"pir"},
			Usage:   "PIR Sensor",
			Action: func(c *cli.Context) error {
				h, err := ctrl.PIRDetect(context.Background(), &google_pb.Empty{})
				if err != nil {
					return err
				}
				log.Printf("PIR Sensor %v", h.On)
				return nil
			},
		},
		{
			Name:    "ForwardSweep",
			Aliases: []string{"sweep"},
			Usage:   "Forward Sweep ultrasonic sensor",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "angle,a",
					Usage: "Increment Angle",
				},
			},
			Action: func(c *cli.Context) error {
				d, err := ctrl.ForwardSweep(context.Background(), &pb.SweepReq{
					Angle: int32(c.Int("angle")),
				})
				if err != nil {
					return err
				}
				log.Printf("Sweep %v", d.Distance)
				return nil
			},
		},
		{
			Name:    "Distance",
			Aliases: []string{"dist"},
			Usage:   "Distance from ultrasonic sensor",
			Action: func(c *cli.Context) error {
				d, err := ctrl.Distance(context.Background(), &google_pb.Empty{})
				if err != nil {
					return err
				}
				log.Printf("Distance %v", d.Distance)
				return nil
			},
		},
		{
			Name:    "Heading",
			Aliases: []string{"head"},
			Usage:   "Magnetic Heading",
			Action: func(c *cli.Context) error {
				h, err := ctrl.Heading(context.Background(), &google_pb.Empty{})
				if err != nil {
					return err
				}
				log.Printf("Heading %v", h.Heading)
				return nil
			},
		},
		{
			Name:    "ServoRotate",
			Aliases: []string{"servo"},
			Usage:   "Rotate the servo.",
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "servo,s",
					Usage: "Select servo 1 (left-right) or 2 (top-bottom) ",
				},
				cli.UintFlag{
					Name:  "angle, a",
					Usage: "Angle of Rotation",
				},
			},
			Action: func(c *cli.Context) error {
				if _, err := ctrl.ServoRotate(context.Background(), &pb.ServoReq{
					Servo: uint32(c.Uint("servo")),
					Angle: uint32(c.Uint("angle")),
				}); err != nil {
					log.Printf("ServoRotate failed %v", err)
				}
				return nil
			},
		},
		{
			Name:    "LEDOn",
			Aliases: []string{"led"},
			Usage:   "Turn on the LED.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "on, o",
					Usage: "Turn on/off LED",
				},
			},
			Action: func(c *cli.Context) error {
				if _, err := ctrl.LEDOn(context.Background(), &pb.LEDOnReq{On: c.Bool("on")}); err != nil {
					log.Printf("LEDOn failed %v", err)
				}
				return nil
			},
		},
		{
			Name:    "LEDBlink",
			Aliases: []string{"blink"},
			Usage:   "Blink LED.",
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "duration, d",
					Value: 1000,
					Usage: "value for blink duration (1000ms)",
				},
				cli.UintFlag{
					Name:  "times, t",
					Value: 10,
					Usage: "no of times to blink LED",
				},
			},
			Action: func(c *cli.Context) error {
				if _, err := ctrl.LEDBlink(context.Background(), &pb.LEDBlinkReq{
					Duration: uint32(c.Uint("duration")),
					Times:    uint32(c.Uint("times")),
				}); err != nil {
					log.Printf("LEDBlink failed %v", err)
				}
				return nil
			},
		},
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "lang, l",
					Value:  "english",
					Usage:  "language for the greeting",
					EnvVar: "LEGACY_COMPAT_LANG,APP_LANG,LANG",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
	}
	app.Run(os.Args)

}