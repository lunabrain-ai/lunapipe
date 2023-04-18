package cli

import (
	"fmt"
	"github.com/UnnoTed/horizontal"
	"github.com/lunabrain-ai/lunapipe/internal/config"
	logcfg "github.com/lunabrain-ai/lunapipe/internal/log"
	"github.com/lunabrain-ai/lunapipe/internal/openai"
	"github.com/lunabrain-ai/lunapipe/internal/prompt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

// TODO breadchris this should be a provided dependency
func setupLogging(level string) {
	logLevel := zerolog.InfoLevel
	if level == "debug" {
		logLevel = zerolog.DebugLevel
	}
	log.Logger = zerolog.New(horizontal.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger().Level(logLevel)
}

func New(
	openaiConfig openai.Config,
	logConfig logcfg.Config,
	cfg config.Configurator,
) *cli.App {
	setupLogging(logConfig.Level)

	type commonFlags struct {
		sync  bool
		quiet bool
		model string
	}

	flagsFromCtx := func(context *cli.Context) commonFlags {
		sync := context.Bool("sync")
		quiet := context.Bool("quiet")
		model := context.String("model")
		return commonFlags{
			sync:  sync,
			quiet: quiet,
			model: model,
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
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "model to use to handle prompt",
			},
		},
		Action: func(ctx *cli.Context) error {
			flags := flagsFromCtx(ctx)

			// TODO breadchris duplicate code
			if flags.model != "" {
				openaiConfig.Model = flags.model
			}

			argPrompt := ctx.Args().First()

			loader := prompt.NewLoader(
				argPrompt,
				ctx.StringSlice("param"),
				ctx.String("template"),
				ctx.String("prompts"),
				flags.sync,
				flags.quiet,
				ctx.Bool("interact"),
			)

			createdPrompt, err := loader.Create()
			if err != nil {
				return err
			}

			log.Debug().Str("prompt", createdPrompt).Msg("sending prompt")
			stream := !flags.sync

			client, err := openai.NewOpenAIQAClient(openaiConfig)
			if err != nil {
				return err
			}

			resp, err := client.Ask(createdPrompt, stream)
			if err != nil {
				return err
			}
			if flags.sync {
				fmt.Println(resp)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:        "chat",
				Description: "Chat with GPT.",
				Action: func(context *cli.Context) error {
					flags := flagsFromCtx(context)
					stream := !flags.sync

					// TODO breadchris duplicate code
					if flags.model != "" {
						openaiConfig.Model = flags.model
					}

					client, err := openai.NewOpenAIQAClient(openaiConfig)
					if err != nil {
						return err
					}

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
					path, err := cfg.Configure(apiKey)
					if err != nil {
						return err
					}
					fmt.Printf("Configuration saved to %s\n", path)
					return nil
				},
			},
		},
	}
}
