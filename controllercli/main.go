package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/deepakkamesh/sonny/sonny"
	google_pb "github.com/golang/protobuf/ptypes/empty"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
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
				fmt.Println("Pinging controller: ", c.Args().First())
				if _, err := ctrl.Ping(context.Background(), &google_pb.Empty{}); err != nil {
					log.Printf("Ping to controller failed: %v", err)
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
