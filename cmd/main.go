package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	asn1parser "github.com/dutchsec/asn1-scheme-parser"
	"github.com/fatih/color"

	cli "github.com/urfave/cli"
)

var Version = "0.1"

var helpTemplate = `NAME:
{{.Name}} - {{.Usage}}

DESCRIPTION:
{{.Description}}

USAGE:
{{.Name}} {{if .Flags}}[flags] {{end}}command{{if .Flags}}{{end}} [arguments...]

COMMANDS:
{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
{{end}}{{if .Flags}}
FLAGS:
{{range .Flags}}{{.}}
{{end}}{{end}}
VERSION:
` + Version +
	`{{ "\n"}}`

func main() {
	app := New()
	app.RunAndExitOnError()
}

var globalFlags = []cli.Flag{}

type cmd struct {
	*cli.App
}

func VersionAction(c *cli.Context) {
	fmt.Println(color.YellowString(fmt.Sprintf("asn1-scheme-parser")))
}

func New() *cmd {
	app := cli.NewApp()
	app.Name = "asn1-scheme-parser"
	app.Author = "DutchSec"
	app.Usage = ""
	app.Description = ``
	app.Flags = globalFlags
	app.CustomAppHelpTemplate = helpTemplate
	app.Commands = []cli.Command{
		{
			Name:   "version",
			Action: VersionAction,
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Action = func(c *cli.Context) error {
		if args := c.Args(); len(args) == 0 {
			cli.ShowAppHelp(c)

			ec := cli.NewExitError(color.RedString("[!] No asn1 scheme set"), 1)
			return ec
		}

		r, err := os.Open(c.Args().First())
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		parser := asn1parser.NewParser(r)

		definition, err := parser.Parse()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Println(definition.Name)
		fmt.Println(strings.Repeat("=", len(definition.Name)))

		for _, t := range definition.Types {
			fmt.Printf("%s: %s\n", t.Name(), reflect.TypeOf(t))
		}

		return nil
	}
	return &cmd{
		App: app,
	}
}
