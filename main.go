package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pivotal-cf/jhanda"
)

type cmdArgs struct {
	CalendarID string `long:"calendar-id" short:"c" default:"primary" description:"The ID of the Google Calendar you wish to purge"`
	StartDate  string `long:"start-date"                              description:"The date when the tool should start purging at 12:01 AM local time in YYYY-MM-DD format (default: 1970-01-01)"`
	EndDate    string `long:"end-date"                                description:"The date when the tool should end purging at 11:59 PM local time in YYYY-MM-DD format (default: today)"`
	DryRun     bool   `long:"dry-run"                                 description:"Will print what will happen to stderr, but won't actually do it"`
	IsVerbose  bool   `long:"verbose" short:"v"                       description:"Enable more verbose logging"`
	Help       bool   `long:"help" short:"h"                          description:"Print this message and quit"`
}

type options struct {
	cmdArgs
	stdout *log.Logger
	stderr *log.Logger
	debug  *log.Logger
	client *http.Client
}

const helpMsg = `
Clean up all your old Google Calendar entries!

Options: 
	%s
`

func main() {
	var opts options
	opts.stdout = log.New(os.Stdout, "", 0)
	opts.stderr = log.New(os.Stderr, "", 0)

	var args cmdArgs

	_, err := jhanda.Parse(&args, os.Args[1:])
	if err != nil {
		opts.stderr.Fatal(err)
	}

	opts.cmdArgs = args

	if opts.Help {
		usage, err := jhanda.PrintUsage(args)
		if err != nil {
			opts.stderr.Fatal(err)
		}

		opts.stdout.Printf(helpMsg, strings.ReplaceAll(usage, "\n", "\n\t"))
		os.Exit(2)
	}

	if opts.IsVerbose {
		opts.debug = log.New(os.Stderr, "", 0)
		opts.debug.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	}

	getClient(&opts)
	processEvents(opts)
}
