package main

import (
	"bazil.org/fuse"
)

type Fid uint64


type FileNode struct {
	Dirent fuse.Dirent;
	FullPath string
	Executable bool
	Children []Fid;
}


var fileMap map[Fid]FileNode
var sourceDir string
var nextFid Fid = 1
