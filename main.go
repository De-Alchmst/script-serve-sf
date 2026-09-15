package main

import (
	"os"
	"log"
	"fmt"
	"syscall"
	"os/signal"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)


func usage() {
	log.Fatal("Usage: ", os.Args[0], " <source> <mountpoint>")
}


func main() {
	// parse args
	if len(os.Args) != 3 {
		usage()
	}

	sourceDir  := os.Args[1]
	mountpoint := os.Args[2]

	// analyze the source dir
	err := discoverSource(sourceDir)
	if err != nil { log.Fatal("Cannot scan ", sourceDir) }

	// initialise cache
	initCache()
	
	// mount fuse
	con, err := fuse.Mount(
		mountpoint,
		fuse.FSName("fuse/SFTH-api-test"),
		fuse.Subtype("apifs"),
		fuse.AllowOther(),
	)
	if err != nil { log.Fatal("cannot mount") }

	
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// start serving the fs
	serveDone := make(chan error, 1)
	go func() {
		server := fs.New(con, &fs.Config {
			Debug: fuse.Debug,        // do nothing
			WithContext: withContext, // store PID in context
		})
		serveDone <- server.Serve(FS{})
	}()

	fmt.Println("Serving at", mountpoint)

	select {
		case <-serveDone:
		case <-sigChan:
	}

	err = fuse.Unmount(mountpoint)
	con.Close()

}
