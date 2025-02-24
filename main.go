package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"runtime/trace"
)

const FILE_TABLE = "files"

// this will hold all FileInfo data structures for batch SQL insert
var rows []FileInfo

func main() {
	start := time.Now()

	//───────────────────────────────────────────────────────────────────────────────────
	// get cli options
	//───────────────────────────────────────────────────────────────────────────────────
	options := CliArgs()
	log.Printf("options: %+v", options)

	//───────────────────────────────────────────────────────────────────────────────────
	// connect to DB
	//───────────────────────────────────────────────────────────────────────────────────
	db := getDB(&options)

	//───────────────────────────────────────────────────────────────────────────────────
	// delete all rows if requested
	//───────────────────────────────────────────────────────────────────────────────────
	db.Exec(fmt.Sprintf("DELETE FROM %s", FILE_TABLE))

	//───────────────────────────────────────────────────────────────────────────────────
	// if -trace was requested
	//───────────────────────────────────────────────────────────────────────────────────
	if options.trace != "" {
		f, err := os.OpenFile(options.trace, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			panic(err)
		}
		trace.Start(f)
		defer trace.Stop()
	}

	//───────────────────────────────────────────────────────────────────────────────────
	// to synchronize goroutines
	//───────────────────────────────────────────────────────────────────────────────────
	var wg sync.WaitGroup
	mutex := sync.Mutex{}

	//───────────────────────────────────────────────────────────────────────────────────
	// this will hold the rows to insert
	//───────────────────────────────────────────────────────────────────────────────────
	rows = make([]FileInfo, 0)

	//───────────────────────────────────────────────────────────────────────────────────
	// Create SQL pool to batch insert data
	//───────────────────────────────────────────────────────────────────────────────────
	batch := make(chan FileInfo, 100)
	context := Context{&wg, db, &mutex, batch}

	go sqlWorker(batch, &context, &options)

	//───────────────────────────────────────────────────────────────────────────────────
	// Create worker pool to throttle goroutines
	//───────────────────────────────────────────────────────────────────────────────────
	jobs := make(chan JobData, options.channelSize)
	for w := 0; w < options.nbWorkers; w++ {
		wg.Add(1)
		options.id = w
		log.Printf("starting goroutine %d\n", w)
		go worker2(jobs, &context, &options)
	}

	fileCount := 0

	//───────────────────────────────────────────────────────────────────────────────────
	// walk through the given directory, sending the context for every file and directory
	//───────────────────────────────────────────────────────────────────────────────────
	err := filepath.WalkDir(options.startDir, func(path string, entry fs.DirEntry, err error) error {
		// if user entered a regexp, match the file against the re
		if options.match != nil && !options.match.MatchString(path) {
			return nil
		}

		fileCount += 1

		if options.verbose {
			log.Printf("processing file '%s'\n", path)
		}

		// is this is a directory, don't calculate its hash
		if entry.IsDir() {
			return nil
		}

		// send job to any worker
		jobs <- JobData{path, entry}

		return nil
	})

	close(jobs)

	if err != nil {
		fmt.Printf("filepath.WalkDir() returned %v\n", err)
	}

	wg.Wait()

	// finally, add remaining rows left by any worker
	log.Printf("# of remaining rows: %d", len(rows))
	context.db.Create(&rows)
	//insertRows(&context, rows)
	log.Printf("%d files processed", fileCount)

	// get stats
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}
