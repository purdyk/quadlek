package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/jirwin/quadlek/plugins/echo"
	"github.com/jirwin/quadlek/plugins/karma"
	"github.com/jirwin/quadlek/plugins/random"
	"github.com/jirwin/quadlek/quadlek"
	"github.com/urfave/cli"
	"github.com/jirwin/quadlek/plugins/spotify"
	"github.com/jirwin/quadlek/plugins/nextep"
)

const Version = "0.0.1"

func run(c *cli.Context) error {
	var apiToken string
	if c.IsSet("api-key") {
		apiToken = c.String("api-key")
	} else {
		cli.ShowAppHelp(c)
		return cli.NewExitError("Missing --api-key arg.", 1)
	}

	var verificationToken string
	if c.IsSet("verification-token") {
		verificationToken = c.String("verification-token")
	} else {
		cli.ShowAppHelp(c)
		return cli.NewExitError("Missing --verification-token arg.", 1)
	}

	dbPath := c.String("db-path")

	bot, err := quadlek.NewBot(context.Background(), apiToken, verificationToken, dbPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error creating bot")
		return nil
	}

	err = bot.RegisterPlugin(echo.Register())
	if err != nil {
		fmt.Printf("error registering echo plugin: %s", err.Error())
		return nil
	}

	err = bot.RegisterPlugin(karma.Register())
	if err != nil {
		fmt.Printf("error registering karma plugin: %s", err.Error())
		return nil
	}

	err = bot.RegisterPlugin(random.Register())
	if err != nil {
		fmt.Printf("error registering random plugin: %s", err.Error())
		return nil
	}

	err = bot.RegisterPlugin(spotify.Register())
	if err != nil {
		fmt.Printf("error registering spotify plugin: %s", err.Error())
		return nil
	}

	if c.IsSet("tvdb-key") {
		tvdbKey := c.String("tvdb-key")

		err = bot.RegisterPlugin(nextep.Register(tvdbKey))
		if err != nil {
			fmt.Printf("error registering nextep plugin: %s", err.Error())
			return nil
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	bot.Start()
	<-signals
	bot.Stop()

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "quadlek"
	app.Version = Version
	app.Usage = "a slack bot"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key",
			Usage:  "The slack api token for the bot",
			EnvVar: "QUADLEK_API_TOKEN",
		},
		cli.StringFlag{
			Name:   "verification-token",
			Usage:  "The slack webhook verification token.",
			EnvVar: "QUADLEK_VERIFICATION_TOKEN",
		},
		cli.StringFlag{
			Name:   "db-path",
			Usage:  "The path where the database is stored.",
			Value:  "quadlek.db",
			EnvVar: "QUADLEK_DB_PATH",
		},
		cli.StringFlag{
			Name: "tvdb-key",
			Usage: "The TVDB api key for the bot, used by the nextep command",
			EnvVar: "QUADLEK_TVDB_TOKEN",
		},
	}

	app.Run(os.Args)
}
