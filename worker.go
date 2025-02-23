// the worker pool is use to manage the number of go routines created
// without such a control, system resource could be exhausted
package main

import (
	"io/fs"
	"sync"

	"gorm.io/gorm"
)

// data of this type will hold the dir entry and db connexion
type Context struct {
	path  string
	entry fs.DirEntry
	db    *gorm.DB
}

// worker used to get file details
func workerPool(entries <-chan Context, wg *sync.WaitGroup, mutex *sync.Mutex, opts *CliOptions) {
	defer wg.Done()

	// the chanel received the dir entries from the WalkDir()
	for d := range entries {
		// build the file info structure
		info := NewFileInfo(d.path, d.entry, opts)

		// manage options
		info.getMagicOrHash(opts.computeHash, opts.magic)

		// update mina slice: critical section
		mutex.Lock()
		rows = append(rows, *info)

		if len(rows) >= opts.batchSize {
			insertRows(d.db, rows)

			// reset slice
			rows = nil
		}
		mutex.Unlock()
		// end of critical section
	}
}

// worker to insert data in batch using Gorm features
func sqlWorker(data <-chan FileInfo, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()

	for row := range data {
		rows = append(rows, row)

		if len(rows) > 100 {
			db.Create(rows)
		}
	}

}
