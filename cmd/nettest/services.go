package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type NetType string

var Client NetType = "client"
var Server NetType = "server"

type ServiceType string

var Echo ServiceType = "echo"
var Pattern ServiceType = "pattern"

type CommandStruct struct {
	hexPattern     bool
	addrFamily     string
	addr           string
	port           uint16
	outputInterval time.Duration
	netType        NetType
	service        ServiceType
	pattern        string
	patternMin     uint
	patternMax     uint
	outputStats    bool
}

func (c *CommandStruct) Name() string {
	return string(c.netType) + "-" + string(c.service)
}

func (c *CommandStruct) Synopsis() string {
	toReturn := ""
	if c.netType == Server {
		toReturn = "Listen for a connection on a specified address and port - "
	} else if c.netType == Client {
		toReturn = "Connect to a specified address and port - "
	}
	if c.service == Echo {
		toReturn += "Echo received data back to the sender."
	} else if c.service == Pattern {
		toReturn += "Send a pattern, receive a specific pattern."
	}
	return toReturn
}

func (c *CommandStruct) Usage() string {
	toReturn := c.Name() + " <addr-family> <addr> <port>"
	if c.service == Pattern {
		toReturn += " <pattern>"
	}
	toReturn += " [-interval <duration>]"
	if c.service == Pattern {
		toReturn += " [-start <length>] [-end <length>]"
	}
	return toReturn + "\n  Run the " + string(c.service) + " service as a " + string(c.netType) +
		`.
  addr-family : The address family of the socket
  addr : The address of the socket.
  port : THe port of the socket.`
}

func (c *CommandStruct) SetFlags(set *flag.FlagSet) {
	set.BoolVar(&c.outputStats, "stats", false, "- : Output final statistics on exit.")
	set.DurationVar(&c.outputInterval, "interval", time.Second*15, "duration : The interval between STDOUT writes, 0 to disable.")
	if c.service == Pattern {
		set.UintVar(&c.patternMin, "start", 1, "start : Pattern starting length.")
		set.UintVar(&c.patternMax, "end", 1024, "end : Pattern max length.")
		set.BoolVar(&c.hexPattern, "hex", false, "- : The pattern is in hex format with no deliminators.")
	}
}

func (c *CommandStruct) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	// Safe shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	termChan := make(chan struct{})
	go func() {
		<-sigs
		termChan <- struct{}{}
	}()
	if len(f.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "addr-family is required")
		return subcommands.ExitUsageError
	} else {
		c.addr = f.Args()[0]
	}
	if len(f.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "addr is required")
		return subcommands.ExitUsageError
	} else {
		c.addr = f.Args()[1]
	}
	if len(f.Args()) < 3 {
		fmt.Fprintln(os.Stderr, "port is required")
		return subcommands.ExitUsageError
	} else {
		cPort, err := strconv.ParseUint(f.Args()[2], 10, 16)
		if err != nil {
			fmt.Fprintln(os.Stderr, "port is invalid")
			return subcommands.ExitUsageError
		}
		c.port = uint16(cPort)
	}
	if c.service == Pattern && len(f.Args()) < 4 {
		fmt.Fprintln(os.Stderr, "pattern is required")
		return subcommands.ExitUsageError
	} else {
		c.pattern = f.Args()[3]
		if c.hexPattern {
			dPattern, err := hex.DecodeString(c.pattern)
			if err != nil {
				fmt.Fprintln(os.Stderr, "pattern is invalid")
				return subcommands.ExitUsageError
			}
			c.pattern = string(dPattern)
		}
		if len(c.pattern) < 1 {
			fmt.Fprintln(os.Stderr, "pattern is empty")
			return subcommands.ExitUsageError
		}
	}

	//TODO implement me
	panic("implement me")

	return subcommands.ExitSuccess
}
