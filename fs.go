package main

import (
	"errors"
	"os"
	"os/exec"
	"strconv"

	"bazil.org/fuse"
)

type Fid uint64
type Pid uint32


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


func (n FileNode) ReadFile(pid Pid) ([]byte, error){
	if n.Dirent.Type == fuse.DT_Dir {
		return nil, errors.New("reading directory, lel")
	}

	// if executable, execute with PID as arg
	if n.IsExecutable {
		// if cached result of 'POST'
		writeResult, exists := getPidWriteResult(pid, Fid(n.Dirent.Inode))
		if exists { return writeResult, nil }

		// else
		return exec.Command(n.FullPath, strconv.FormatUint(uint64(pid), 10)).Output()

	// else just read
	} else {
		return os.ReadFile(n.FullPath)
	}
}


// if `Fid` under `Pid` has occupied `OutputBuffer`, set it to `nil` and return
func getPidWriteResult(pid Pid, fid Fid) ([]byte, bool) {
	userCache, exists := cache[pid]
	if exists {
		fidData, exists := userCache[fid]
		if exists {
			if fidData.OutputBuffer != nil {
				data := fidData.OutputBuffer
				fidData.OutputBuffer = nil
				return data, true
			}
		}
	}

	return nil, false
}


// returns `Fid` `InputBuffer` under `Pid`, creates if needed
// code not optimal, but readable!
func getPidWriteBuffer(pid Pid, fid Fid) (*[]byte) {
	userCache, exists := cache[pid]
	if !exists {
		userCache = make(map[Fid]CacheEntry)
		cache[pid] = userCache
	}

	fidData, exists := userCache[fid]
	if !exists {
		fidData = CacheEntry{
			InputBuffer: nil,
			OutputBuffer: nil,
		}
		userCache[fid] = fidData
	}

	if fidData.InputBuffer == nil {
		fidData.InputBuffer = &[]byte{}
		userCache[fid] = fidData
	}

	return fidData.InputBuffer
}


// expects that there is valid `InputBuffer`, will overwrite unused `OutputBuffer`
func executeCachedFile(pid Pid, fid Fid) error {
	node, validFid := fileMap[fid]
	if !validFid || node.Dirent.Type == fuse.DT_Dir {
		return errors.New("invalid Fid")
	}

	userCache := cache[pid]

	input := userCache[fid]
	output, err := exec.Command(node.FullPath, strconv.FormatUint(uint64(pid), 10), string(*input.InputBuffer)).Output()
	if err != nil { return err }

	userCache[fid] = CacheEntry{
		InputBuffer: nil,
		OutputBuffer: output,
	}

	return nil
}
