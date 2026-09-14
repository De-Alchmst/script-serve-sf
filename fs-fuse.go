package main

import (
	"context"
	"os"
	"syscall"
	"errors"
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
)

type FS struct{}


// add PID to fuse context
// useful for keeping track of users based on the PID of the process
// serving them
func withContext(ctx context.Context, req fuse.Request) context.Context {
	return context.WithValue(ctx, "PID", req.Hdr().Pid)
}


func (FS) Root() (fs.Node, error) {
	return Fid(fileMap[rootFid].Dirent.Inode), nil
}


func (f Fid) Attr(ctx context.Context, a *fuse.Attr) error {
	node, validFid := fileMap[f]
	if !validFid { return errors.New("invalid Fid") }

	a.Inode = uint64(f)
	if node.Dirent.Type == fuse.DT_Dir {
		a.Mode = os.ModeDir | 0o777 // dr-xr-xr-x

	} else {
		if node.IsExecutable {
			a.Mode = 0o666 // .rw-rw-rw-
			a.Size = 0 // must set `fuse.OpenDirectIo` to read file of size 0
		} else {
			a.Mode = 0o444 // .r--r--r--
			a.Size = uint64(node.Size)
		}
	}

	return nil
}


func (f Fid) Lookup(ctx context.Context, name string) (fs.Node, error) {
	node, validFid := fileMap[f]
	if !validFid || node.Dirent.Type == fuse.DT_File {
		return nil, errors.New("invalid Fid")
	}

	for _, fid := range *node.Children {
		if fileMap[fid].Dirent.Name == name {
			return fid, nil
		}
	}

	return nil, syscall.ENOENT
}


func (f Fid) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error)  {
	node, validFid := fileMap[f]
	if !validFid || node.Dirent.Type == fuse.DT_File {
		return nil, errors.New("invalid Fid")
	}

	dirents := []fuse.Dirent{}
	for _, fid := range *node.Children {
		dirents = append(dirents, fileMap[fid].Dirent)
	}

	return dirents, nil
}


func (f Fid) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	node, validFid := fileMap[f]
	if !validFid {
		return nil, errors.New("invalid Fid")
	}
	log.Println("|| ", node.FullPath)

	// if executable, set `fuse.OpenDirectIo` so that it is not limited
	// by it's reported size
	if node.IsExecutable && node.Dirent.Type == fuse.DT_File {
		resp.Flags |= fuse.OpenDirectIO
	}

	return f, nil
}


func (f Fid) ReadAll(ctx context.Context) ([]byte, error) {
	node, validFid := fileMap[f]
	if !validFid || node.Dirent.Type == fuse.DT_Dir {
		return nil, errors.New("invalid Fid")
	}
	log.Println(">> ", node.FullPath)

	return node.ReadFile(ctx.Value("PID").(uint32))
}


func (f Fid) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	node, validFid := fileMap[f]
	if !validFid || node.Dirent.Type == fuse.DT_Dir {
		return errors.New("invalid Fid")
	}
	log.Println("OO ", node.FullPath)

	data, err := node.ReadFile(ctx.Value("PID").(uint32))
	if err != nil { return err }

	fuseutil.HandleRead(req, resp, data)
	return nil
}

