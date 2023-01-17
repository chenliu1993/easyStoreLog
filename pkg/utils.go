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
