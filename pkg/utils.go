package pkg

import (
	"log"
	"os"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		log.Printf("failed to stat file %s, %v", path, err)
		return false
	}
	return s.IsDir()
}

func cleanByteSlice(buf *[]byte, len int) {
	// This is a time-consuming op, but I cannot find another way?
	for i := 0; i < len; i++ {
		(*buf)[i] = 0
	}
}
