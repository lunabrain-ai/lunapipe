package internal

import (
	"fmt"
	"github.com/UnnoTed/horizontal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

type Flags struct {
	Sync  bool
	Quiet bool
}

// TODO breadchris this should be a provided dependency
func setupLogging(level string) {
	logLevel := zerolog.InfoLevel
	if level == "debug" {
		logLevel = zerolog.DebugLevel
	}
	log.Logger = zerolog.New(horizontal.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger().Level(logLevel)
}

func NewCLI(
	client QAClient,
	config LogConfig,
) *cli.App {
	setupLogging(config.Level)

	flagsFromCtx := func(context *cli.Context) Flags {
		sync := context.Bool("sync")
		quiet := context.Bool("quiet")
		return Flags{
			Sync:  sync,
			Quiet: quiet,
		}
	}

	return &cli.App{
		Name:  "aicli",
		Usage: "AI for your CLI!",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "sync",
				Aliases: []string{"s"},
				Usage:   "do not stream output",
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "do not print prompt",
			},
			&cli.StringFlag{
				Name:  "prompts",
				Usage: "Directory containing additional prompt templates.",
			},
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "use a template to generate the prompt",
			},
			&cli.StringSliceFlag{
				Name:    "param",
				Aliases: []string{"p"},
				Usage:   "parameters used when formatting template",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "chat",
				Action: func(context *cli.Context) error {
					flags := flagsFromCtx(context)
					stream := !flags.Sync
					println("Starting chat, close with ctrl+D...")
					return client.Chat(stream)
				},
			},
		},
		Action: func(context *cli.Context) error {
			flags := flagsFromCtx(context)
			prompt, err := getPrompt(context, flags)
			if err != nil {
				return err
			}

			log.Debug().Str("prompt", prompt).Msg("sending prompt")
			stream := !flags.Sync
			resp, err := client.Ask(prompt, stream)
			if err != nil {
				return err
			}
			if flags.Sync {
				fmt.Println(resp)
			}
			return nil
		},
	}
}
