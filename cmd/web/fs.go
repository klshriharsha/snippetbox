package main

import (
	"net/http"
	"path/filepath"
)

// staticFileSystem implements a custom file system to ensure that
// directories are not served from a file server
type staticFileSystem struct {
	fs http.FileSystem
}

// Open ensures that the requested file is not a directory.
// if it is a directory, it ensures that it has index.html file inside it
func (sfs staticFileSystem) Open(name string) (http.File, error) {
	// attempt to open the requested file
	f, err := sfs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	// check if the requested file is a directory
	s, _ := f.Stat()
	if s.IsDir() {
		// check if the directory has an index html file
		index := filepath.Join(name, "index.html")
		if _, err := sfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
