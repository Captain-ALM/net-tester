package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
)

var (
	buildName    = ""
	buildVersion = "develop"
	buildDate    = ""
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&CommandStruct{service: Echo, netType: Server}, "")
	subcommands.Register(&CommandStruct{service: Echo, netType: Client}, "")
	subcommands.Register(&CommandStruct{service: Pattern, netType: Server}, "")
	subcommands.Register(&CommandStruct{service: Pattern, netType: Client}, "")
	subcommands.Register(&Info{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
