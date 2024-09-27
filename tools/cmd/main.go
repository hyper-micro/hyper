package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hyper-micro/hyper/tools/command"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "hyper",
		Version: "0.0.1",
		Usage:   "A cli tool to generate project code",
		Action: func(ctx *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize new project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "project",
						Aliases:  []string{"p"},
						Usage:    "Project name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "mod",
						Aliases:  []string{"m"},
						Usage:    "Go module name",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					cancel := waitDot()

					err := command.NewInitCommand(
						command.InitCommandArgs{
							ProjectName: ctx.String("project"),
							Mod:         ctx.String("mod"),
						},
					).Run()

					cancel()

					if err != nil {
						return fmt.Errorf("init project fail: %v", err)
					}

					printSuccess("Program initialization successful")

					return nil
				},
			},
			{
				Name:    "replace",
				Aliases: []string{"r"},
				Usage:   "Replace project template",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "project",
						Aliases:  []string{"p"},
						Usage:    "Project name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "mod",
						Aliases:  []string{"m"},
						Usage:    "Go module name",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					cancel := waitDot()
					err := command.NewInitCommand(
						command.InitCommandArgs{
							ProjectName: ctx.String("project"),
							Mod:         ctx.String("mod"),
						},
					).Replace()
					cancel()
					if err != nil {
						return fmt.Errorf("replace project fail: %v", err)
					}
					printSuccess("Program replace successful")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		printError(err.Error())
		os.Exit(1)
	}
}

func waitDot() func() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			fmt.Printf("..")
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(600 * time.Millisecond)
			}
		}
	}()
	return func() {
		cancel()
	}
}

func printSuccess(msg string) {
	fmt.Printf("\n\033[0;32;40m%s\033[0m %s\n", "[OK]", msg)
}

func printError(msg string) {
	fmt.Printf("\n\033[0;31;40m%s\033[0m %s\n\n", "[ERROR]", msg)
}
