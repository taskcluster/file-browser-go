package browser

import (
	"os"
	"sync"
)

// Use for supporting Open calls on a file/directory
type fileHandle struct {
	// handleID is only set when file is registered
	handleID uint64
	dir      bool
	// Flags used while opening the file
	flags int
	// Permissions
	perm os.FileMode
	// This file descriptor will store the read and write
	// positions in the file
	file *os.File
	// This list stores indices of duplicated handles
	dup []uint64
}

type fileRegistry struct {
	openFiles    map[uint64]*fileHandle
	nextHandleID uint64
	sync.Mutex
}

func (f *fileRegistry) openFile(path string, flags int, perm os.FileMode) (uint64, error) {
	file, err := os.OpenFile(path, flags, perm)
	if err != nil {
		return 0, err
	}
	finfo, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return 0, err
	}
	f.Lock()
	defer f.Unlock()
	handleID := f.nextHandleID
	f.nextHandleID++
	f.openFiles[handleID] = &fileHandle{
		handleID: handleID,
		perm:     perm,
		dir:      finfo.IsDir(),
		flags:    flags,
		file:     file,
		dup:      []uint64{},
	}
	return handleID, nil
}

func (f *fileRegistry) closeFile(handleID uint64) error {
	f.Lock()
	defer f.Unlock()
	fh, ok := f.openFiles[handleID]
	if !ok {
		return errInvalidHandle
	}
	err := fh.file.Close()
	delete(f.openFiles, handleID)
	return err
}

func (f *fileRegistry) get(handleID uint64) (*fileHandle, bool) {
	f.Lock()
	defer f.Unlock()
	fh, ok := f.openFiles[handleID]
	return fh, ok
}

func (f *fileRegistry) dup(handleID uint64) (uint64, error) {
	return 0, nil
}

func (f *fileRegistry) dup2(handleID uint64, handleID2 uint64) (uint64, error) {
	return 0, nil
}
