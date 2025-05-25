package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
)

type Info struct {
}

func (i *Info) Name() string {
	return "info"
}

func (i *Info) Synopsis() string {
	return "Shows software information."
}

func (i *Info) Usage() string {
	return i.Name()
}

func (i *Info) SetFlags(*flag.FlagSet) {
}

func (i *Info) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	fmt.Println("Name: ", buildName)
	fmt.Println("Version: ", buildVersion)
	fmt.Println("Build Date: ", buildDate)
	return subcommands.ExitSuccess
}
