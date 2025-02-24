// the worker pool is use to manage the number of go routines created
// without such a control, system resource could be exhausted
package main

import (
	"io/fs"
	"log"
	"sync"

	"gorm.io/gorm"
)

// data of this type will hold the dir entry and db connexion
type JobData struct {
	path  string
	entry fs.DirEntry
}

// those data will not change once initialized
type Context struct {
	wg    *sync.WaitGroup
	db    *gorm.DB
	mutex *sync.Mutex
	batch chan FileInfo
}

// worker used to get file details
func worker2(entries <-chan JobData, ctx *Context, opts *CliOptions) {
	defer ctx.wg.Done()

	// the chanel received the dir entries from the WalkDir()
	for d := range entries {
		// build the file info structure
		info := NewFileInfo(d.path, d.entry, opts)

		// manage options
		info.getMagicOrHash(opts.computeHash, opts.magic)

		// send info structure to specialized goroutine
		ctx.batch <- *info
	}
}

// func worker(entries <-chan JobData, ctx *Context, opts *CliOptions) {
// 	log.Printf("running goroutine %d\n", opts.id)
// 	defer ctx.wg.Done()

// 	// the chanel received the dir entries from the WalkDir()
// 	for d := range entries {
// 		// build the file info structure
// 		info := NewFileInfo(d.path, d.entry, opts)

// 		// manage options
// 		info.getMagicOrHash(opts.computeHash, opts.magic)

// 		// update mina slice: critical section
// 		ctx.mutex.Lock()
// 		rows = append(rows, *info)

// 		if len(rows) >= opts.batchSize {
// 			insertRows(ctx.db, rows)
// 			log.Printf("inside goroutine %d", opts.id)

// 			// reset slice
// 			rows = nil
// 		}
// 		ctx.mutex.Unlock()
// 		// end of critical section
// 	}
// }

// worker to insert data in batch using Gorm features
func sqlWorker(dataToInsert <-chan FileInfo, ctx *Context, opts *CliOptions) {
	defer ctx.wg.Done()

	for row := range dataToInsert {
		rows = append(rows, row)

		if len(rows) >= opts.batchSize {
			ctx.wg.Add(1)
			go insertRows(ctx, rows)
			log.Printf("inside goroutine %d", opts.id)

			// reset slice
			rows = nil
		}
	}

}
