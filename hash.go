package main

import (
	"crypto/sha256"
	"encoding/hex"
	// "fmt"
	"log"
	"os"
	// "github.com/h2non/filetype"
)

func path_sha256(path string) string {
	// read whole file
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("error <%s> opening file <%s>", err, path)
	}

	// calculate hash
	hash := sha256.Sum256(data)

	// calculate file magic number
	// kind, _ := filetype.Match(data)
	// if kind == filetype.Unknown {
	// 	fmt.Println("Unknown file type")
	// }

	// fmt.Printf("File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)

	return hex.EncodeToString(hash[:])
}
