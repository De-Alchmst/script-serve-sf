package main

import (
	"os"
	"syscall"
	"time"
)


type CacheEntry struct {
	InputBuffer  *[]byte
	OutputBuffer  []byte
}


var cache map[Pid]map[Fid]CacheEntry

func initCache() {
	cache = make(map[Pid]map[Fid]CacheEntry)
}


// remove entries by dead processes
func CleanCache() {
	for {
		// clear cache every five minutes
		time.Sleep(time.Minute * 5)
		
		for k := range cache {
			if !PidExists(k) {
				delete(cache, k)
			}
		}
	}
}


// dafaq! what is wrong with you, UNIX? U ok?

// Source - https://stackoverflow.com/a/59459658
// Posted by Paul
// Retrieved 2026-09-14, License - CC BY-SA 4.0

func PidExists(pidPropa Pid) bool {
	pid := int32(pidPropa)
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}
	if err.Error() == "os: process already finished" {
		return false
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	switch errno {
	case syscall.ESRCH:
		return false
	case syscall.EPERM:
		return true
	}
	return false
}

