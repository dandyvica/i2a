package main

import (
	"flag"
	// "log"
	// "os"
	"regexp"
	"runtime"
)

const CONFIG_FILE = "settings.yml"

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
	dryRun      bool           // if set, don't insert data, just trace goroutines
	id          int            // ID of the running goroutine
	trace       string         // name of the trace file if any
	configFile  string         // name of the YAML config file if defined
}

func CliArgs() CliOptions {
	// init struct
	var options CliOptions

	// load settings from configuration file if it exists
	// if _, err := os.Stat(options.configFile); err == nil {
	// 	cfg, err = ReadYAMLConfig(options.configFile)
	// 	if err != nil {
	// 		log.Fatalf("unable to open file '%s'", options.configFile)
	// 	}
	  
	// }	

	// if set, we want the line number from the file
	flag.StringVar(&options.startDir, "dir", ".", "name of the starting directory")
	flag.StringVar(&options.dbUrl, "db", "iaa.db", "url database connexion strings, or DB file for Sqlite3.")
	flag.StringVar(&options.matchString, "match", "", "insert only files matching the regular expression")
	flag.StringVar(&options.trace, "trace", "trace.out", "name of the trace file")
	flag.StringVar(&options.configFile, "cfg", "settings.yml", "name of the YAML config file if defined")

	flag.BoolVar(&options.computeHash, "sha256", false, "compute the sha256 sums")
	flag.BoolVar(&options.verbose, "v", false, "debugging mode")
	flag.BoolVar(&options.magic, "magic", false, "compute the file supposed MIME type")
	flag.BoolVar(&options.truncate, "truncate", false, "if set, delete all rows from the table first")
	flag.BoolVar(&options.dryRun, "dry", false, "if set, don't insert data, just trace goroutines")

	flag.IntVar(&options.nbWorkers, "workers", runtime.NumCPU(), "number of workers to start. Default to the number of cores")
	flag.IntVar(&options.batchSize, "sqlbatch", 1000, "SQL insert batch size")
	flag.IntVar(&options.channelSize, "chan", 100, "chanel size for workers")

	flag.Parse()

	// compile regexp if any
	if options.matchString != "" {
		options.match = regexp.MustCompile(options.matchString)
	}

	return options
}
