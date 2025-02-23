// all methods for inserting data into SQL
package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// create DB connexion
func getDB(opts *CliOptions) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(opts.dbUrl), &gorm.Config{
		CreateBatchSize: opts.batchSize,
	})

	if err != nil {
		panic("failed to connect database")
	}

	// create table if any
	db.AutoMigrate(&FileInfo{})

	return db
}

// insert row into DB
func insertRows(db *gorm.DB, rows []FileInfo) {
	result := db.Create(&rows)

	if result.Error == nil {
		log.Printf("==============> %d inserted\n", len(rows))
	} else {
		log.Printf("error %v, '%v' when inserting rows\n", len(rows), result)
	}
}
