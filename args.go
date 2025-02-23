package main

import (
	"flag"
	"regexp"
	"runtime"
)

// this will hold all options
type CliOptions struct {
	startDir    string         // starting directory
	computeHash bool           // true if we want to compute SHA256 hash for each file
	nbWorkers   int            // number of workers to start
	dbUrl       string         // DB connection url
	batchSize   int            // GROM batch insert size
	channelSize int            // buffered or unbuffered chanel
	verbose     bool           // print out list of files
	matchString string         // regular expression to filter files to insert
	match       *regexp.Regexp // the corresponding re compiled
	magic       bool           // true if we want to compute the file supposed MIME type
	truncate    bool           // if set, delete all rows from the table first
}

func CliArgs() CliOptions {
	// init struct
	var options CliOptions

	// if set, we want the line number from the file
	flag.StringVar(&options.startDir, "dir", ".", "name of the starting directory")
	flag.StringVar(&options.dbUrl, "db", "iaa.db", "url database connexion strings, or DB file for Sqlite3.")
	flag.StringVar(&options.matchString, "match", "", "insert only files matching the regular expression")

	flag.BoolVar(&options.computeHash, "sha256", false, "compute the sha256 sums")
	flag.BoolVar(&options.verbose, "v", false, "debugging mode")
	flag.BoolVar(&options.magic, "magic", false, "compute the file supposed MIME type")
	flag.BoolVar(&options.truncate, "truncate", false, "if set, delete all rows from the table first")

	flag.IntVar(&options.nbWorkers, "workers", runtime.NumCPU(), "number of workers to start. Default to the number of cores")
	flag.IntVar(&options.batchSize, "batch", 1000, "SQL insert batch size")
	flag.IntVar(&options.channelSize, "chan", 100, "chanel size for workers")

	flag.Parse()

	// compile regexp if any
	if options.matchString != "" {
		options.match = regexp.MustCompile(options.matchString)
	}

	return options
}
