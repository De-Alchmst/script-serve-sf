package main

import (
	"os"
	"log"
	"fmt"
)


func usage() {
	log.Fatal("Usage: ", os.Args[0], " <source> <mountpoint>")
}


func main() {
	if len(os.Args) != 3 {
		usage()
	}

	sourceDir   = os.Args[1] // global defined in fs.go
	// mountpoint := os.Args[2]

	err := discoverSource()
	if err != nil {
		log.Fatal("Cannot scan ", sourceDir)
	}

	var i Fid
	for i = 1; i < nextFid; i++ {
		fmt.Println(fileMap[i].FullPath, fileMap[i].Executable)
	}
}
