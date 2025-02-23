package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sync"
	"time"
)

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
	db.Exec("DELETE FROM files")

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
	// batch := make(chan FileInfo, 100)
	// go sqlWorker(batch, &wg, db)

	//───────────────────────────────────────────────────────────────────────────────────
	// Create worker pool to throttle goroutines
	//───────────────────────────────────────────────────────────────────────────────────
	jobs := make(chan Context, options.channelSize)
	for w := 0; w < options.nbWorkers; w++ {
		wg.Add(1)
		go workerPool(jobs, &wg, &mutex, &options)
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
		jobs <- Context{path, entry, db}

		return nil
	})

	close(jobs)

	log.Printf("%d files processed", fileCount)

	if err != nil {
		fmt.Printf("filepath.WalkDir() returned %v\n", err)
	}

	wg.Wait()

	// finally, add remaining rows left by any worker
	log.Printf("# of remaining rows: %d", len(rows))
	insertRows(db, rows)

	// get stats
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}
