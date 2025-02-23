package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
	"gorm.io/gorm"
)

// this will hold all information we want about files we meet
type FileInfo struct {
	gorm.Model
	Name     string // name including full path
	Size     int64  // file size in bytes
	Hash     string // sha256 sum if any
	Basename string // file name without full path
	Ext      string // file extension
	Kind     string // file magic number
}

// build a new info from path and dir entry
func NewFileInfo(path string, d fs.DirEntry, opts *CliOptions) *FileInfo {
	// get hash if any
	var hash string

	// compute hash if requested
	if opts.computeHash {
		hash = path_sha256(path)
	}

	return &FileInfo{
		Name:     path,
		Size:     getMetaData(d),
		Hash:     hash,
		Basename: filepath.Base(path),
		Ext:      filepath.Ext(path),
	}
}

func (f FileInfo) String() string {
	return fmt.Sprintf("name='%s' base='%s' extension=%s size=%d hash=%s magic=%s", f.Name, f.Basename, f.Ext, f.Size, f.Hash, f.Kind)
}

// get dir entry metadata
func getMetaData(d fs.DirEntry) int64 {
	info, err := d.Info()
	if err != nil {
		log.Printf("error <%s> getting info for file <%s>", err, info.Name())
	}

	return info.Size()

}

// get file magic number and hash if asked
func (f *FileInfo) getMagicOrHash(wantHash, wantMagic bool) {
	// for the hash, we have to read the whole file
	if wantHash {
		data, err := os.ReadFile(f.Name)
		if err != nil {
			log.Printf("error <%s> opening file <%s>", err, f.Name)
		}

		// calcultate hash
		hash := sha256.Sum256(data)
		f.Hash = hex.EncodeToString(hash[:])

		// if in addition we want magic
		if wantMagic {
			kind, _ := filetype.Match(data)
			f.Kind = kind.MIME.Value
		}

	} else {
		// we only need to read the first 261 bytes to get the magic number
		file, err := os.Open(f.Name)
		if err != nil {
			log.Printf("error <%s> opening file <%s>", err, f.Name)
		}

		head := make([]byte, 261)
		file.Read(head)
		kind, _ := filetype.Match(head)
		f.Kind = kind.MIME.Value
	}
}

// this is used by GORM to set the table name
func (FileInfo) TableName() string {
	return "files"
}
