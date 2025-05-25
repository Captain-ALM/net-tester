package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	net2 "gitcove.com/alfred/net-tester/net"
	"gitcove.com/alfred/net-tester/services"
	"gitcove.com/alfred/net-tester/updates"
	"github.com/google/subcommands"
	"net"
	"os"
	"os/signal"
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
	outputInterval time.Duration
	netType        NetType
	service        ServiceType
	pattern        []byte
	patternMin     uint
	patternMax     uint
	outputStats    bool
	bufferSize     uint
	timeout        time.Duration
}

func (c *CommandStruct) Name() string {
	return string(c.netType) + "-" + string(c.service)
}

func (c *CommandStruct) Synopsis() string {
	toReturn := ""
	if c.netType == Server {
		toReturn = "Listen for a connection on a specified address - "
	} else if c.netType == Client {
		toReturn = "Connect to a specified address - "
	}
	if c.service == Echo {
		toReturn += "Echo received data back to the sender."
	} else if c.service == Pattern {
		toReturn += "Send a pattern, receive a specific pattern."
	}
	return toReturn
}

func (c *CommandStruct) Usage() string {
	toReturn := c.Name() + " [-stats] [-interval <duration>] [-timeout <duration>] [-buffer <size>]"
	if c.service == Pattern {
		toReturn += " [-start <length>] [-end <length>] [-hex]"
	}
	toReturn += " <addr-family> <addr>"
	if c.service == Pattern {
		toReturn += " <pattern>"
	}
	return toReturn + "\n  Run the " + string(c.service) + " service as a " + string(c.netType) +
		`.
  addr-family : The address family of the socket
  addr : The address of the socket (Can include port).`
}

func (c *CommandStruct) SetFlags(set *flag.FlagSet) {
	set.BoolVar(&c.outputStats, "stats", false, "- : Output final statistics on exit.")
	set.DurationVar(&c.outputInterval, "interval", time.Second*15, "duration : The interval between STDOUT writes, 0 to disable.")
	set.DurationVar(&c.timeout, "timeout", 0, "duration : The timeout for sending or receiving.")
	set.UintVar(&c.bufferSize, "buffer", 1024, "size : The size of the receive buffer in bytes.")
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

	quitter := updates.NewQuitter()
	var lConn net.Listener = nil
	go func() {
		<-sigs
		quitter.Quit()
		if lConn != nil {
			_ = lConn.Close()
		}
	}()
	if len(f.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "addr-family is required")
		return subcommands.ExitUsageError
	} else {
		c.addrFamily = f.Args()[0]
	}
	if len(f.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "addr is required")
		return subcommands.ExitUsageError
	} else {
		c.addr = f.Args()[1]
	}
	if c.service == Pattern && len(f.Args()) < 3 {
		fmt.Fprintln(os.Stderr, "pattern is required")
		return subcommands.ExitUsageError
	} else if c.service == Pattern {
		if c.hexPattern {
			var err error
			c.pattern, err = hex.DecodeString(f.Args()[2])
			if err != nil {
				fmt.Fprintln(os.Stderr, "pattern is invalid")
				return subcommands.ExitUsageError
			}
		} else {
			c.pattern = []byte(f.Args()[2])
		}
		if len(c.pattern) < 1 {
			fmt.Fprintln(os.Stderr, "pattern is empty")
			return subcommands.ExitUsageError
		}
	}
	if c.service == Pattern && c.patternMin > c.patternMax {
		fmt.Fprintln(os.Stderr, "start larger than end")
		return subcommands.ExitUsageError
	}

	var err error
	var conn net.Conn = nil
	if c.netType == Server {
		fmt.Fprintln(os.Stderr, "Listening for a connection on ", c.addrFamily, c.addr)
		lConn, err = net.Listen(c.addrFamily, c.addr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return subcommands.ExitFailure
		}
		err = errors.New("")
		for err != nil && conn == nil && quitter.Active() {
			conn, err = lConn.Accept()
		}
		_ = lConn.Close()
		lConn = nil
		if !quitter.Active() || conn == nil {
			return subcommands.ExitFailure
		}
	} else if c.netType == Client {
		fmt.Fprintln(os.Stderr, "Connecting to ", c.addrFamily, c.addr)
		conn, err = net.Dial(c.addrFamily, c.addr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return subcommands.ExitFailure
		}
		if conn == nil {
			return subcommands.ExitFailure
		}
	} else {
		fmt.Fprintln(os.Stderr, "netType is invalid")
		return subcommands.ExitUsageError
	}
	fmt.Fprintln(os.Stderr, "Connected to : ", conn.RemoteAddr().String())
	fmt.Fprintln(os.Stderr, "Connected from : ", conn.LocalAddr().String())
	updater := &updates.Update{StartTime: time.Now()}
	net2.RunClient(conn, c.GetService(), quitter, updater, c.bufferSize, c.timeout)
	for quitter.Active() {
		if c.outputInterval > 0 {
			select {
			case <-time.After(c.outputInterval):
				c.statsOut(updater)
			case <-quitter.Quitter():
				break
			}
		} else {
			select {
			case <-quitter.Quitter():
				break
			}
		}
	}
	if c.outputStats {
		c.statsOut(updater)
	}
	return subcommands.ExitSuccess
}

func (c *CommandStruct) GetService() services.Service {
	if c.service == Pattern {
		return &services.PatternService{Pattern: c.pattern, MinLength: c.patternMin, MaxLength: c.patternMax}
	} else {
		return &services.EchoService{}
	}
}

func (c *CommandStruct) statsOut(update *updates.Update) {
	if c.service == Pattern {
		fmt.Println("Pattern Length : IN|OUT :", update.PatternLengthIn, "|", update.PatternLengthOut)
	}
	fmt.Println("Bytes : IN|OUT :", update.BytesReceived, "|", update.BytesSent)
	timeTaken := uint64(time.Now().Sub(update.StartTime) / time.Second)
	if timeTaken <= 0 {
		timeTaken = 1
	}
	fmt.Println("Speed (B/s) : IN|OUT :", update.BytesReceived/timeTaken, "|", update.BytesSent/timeTaken)
}
