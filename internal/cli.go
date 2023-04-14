package internal

import (
	"fmt"
	"github.com/UnnoTed/horizontal"
	"github.com/lunabrain-ai/lunapipe/internal/config"
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
	logConfig config.LogConfig,
) *cli.App {
	setupLogging(logConfig.Level)

	flagsFromCtx := func(context *cli.Context) Flags {
		sync := context.Bool("sync")
		quiet := context.Bool("quiet")
		return Flags{
			Sync:  sync,
			Quiet: quiet,
		}
	}

	return &cli.App{
		Name:  "lunapipe",
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
			&cli.BoolFlag{
				Name:    "interact",
				Aliases: []string{"i"},
				Usage:   "For a template, interactively prompt for parameters",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "chat",
				Description: "Chat with GPT.",
				Action: func(context *cli.Context) error {
					flags := flagsFromCtx(context)
					stream := !flags.Sync
					println("Starting chat, close with ctrl+D...")
					return client.Chat(stream)
				},
			},
			{
				Name:        "configure",
				Description: "Configure the CLI.",
				Action: func(context *cli.Context) error {
					fmt.Printf("Enter your API key: ")
					var apiKey string
					_, err := fmt.Scanf("%s", &apiKey)
					println()
					if err != nil {
						return err
					}
					return config.NewConfigurator(apiKey)
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
