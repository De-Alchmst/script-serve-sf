package main

import (
	"bazil.org/fuse"
	"os"
	"path/filepath"
	"errors"
  "github.com/wneessen/go-fileperm"

)


func discoverSource() error {
	fileMap = make(map[Fid]FileNode)

	// populate root
	root := FileNode{
		Dirent: fuse.Dirent{
			Inode: uint64(nextFid),
			Name: "/",
			Type: fuse.DT_Dir,
		},
		Children: []Fid{},
		FullPath: sourceDir,
		Executable: true,
	}
	fileMap[nextFid] = root

	nextFid += 1

	// check if pointing at dir
	dirInfo, err := os.Stat(sourceDir)
	if err != nil { return err }

	// check if dir
	if !dirInfo.IsDir() {
		return errors.New("source dir is not a directory")
	}

	// search
	return searchDir(root, sourceDir)
}


func searchDir(self FileNode, path string) (error) {
	children, err := os.ReadDir(path)
	if err != nil { return err }

	for _, child := range children {

		// create child node
		childNode := FileNode{
			Dirent: fuse.Dirent{
				Inode: uint64(nextFid),
				Name : child.Name(),
				Type : fuse.DT_File,
			},
			FullPath: filepath.Join(path, child.Name()),
			Children: nil,
		}

		// check if executable
		p, err := fileperm.New(childNode.FullPath)
		if err != nil { return err }
		childNode.Executable = p.UserExecutable()
		
		// add child
		self.Children = append(self.Children, nextFid)
		fileMap[nextFid] = childNode

		nextFid += 1

		// recurse
		if child.IsDir() {
			childNode.Children = []Fid{}
			childNode.Dirent.Type = fuse.DT_Dir
			searchDir(childNode, childNode.FullPath)
		}
	}

	return nil
}
