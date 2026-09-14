package main

import (
	"errors"
	"os"
	"os/exec"
	"strconv"

	"bazil.org/fuse"
)

type Fid uint64


type FileNode struct {
	Dirent fuse.Dirent;
	FullPath string
	IsExecutable bool
	Size int64 // if not executable...
	// needs to be pointer, as cannot take reference to map item
	// alternative would be havind fileMap `map[Fid]*FileNode`, but that would be
	// even more annoying
	Children *[]Fid;
}


var fileMap map[Fid]FileNode
var sourceDir string
var rootFid Fid = 1
var nextFid = rootFid


func (n FileNode) ReadFile(pid uint32) ([]byte, error){
	if n.Dirent.Type == fuse.DT_Dir {
		return nil, errors.New("reading directory, lel")
	}

	// if executable, execute with PID as arg
	if n.IsExecutable {
		return exec.Command(n.FullPath, strconv.FormatUint(uint64(pid), 10)).Output()

	// else just read
	} else {
		return os.ReadFile(n.FullPath)
	}
}
