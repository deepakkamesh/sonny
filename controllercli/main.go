package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:    "PingController",
			Aliases: []string{"ping"},
			Usage:   "Ping the controller.",
			Action: func(c *cli.Context) error {
				fmt.Println("Pinging controller: ", c.Args().First())
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
	/*
		conn, err := grpc.Dial("localhost:2233", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to connect to rpc endpoint: %v", err)
		}
		defer conn.Close()

		c := pb.NewDevicesRPCClient(conn)
		if _, err := c.Ping(context.Background(), &google_pb.Empty{}); err != nil {
			log.Printf("Ping to controller failed: %v", err)
		}
	*/

}
