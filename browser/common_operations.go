package browser

import (
	"container/list"
	"io"
	"os"
	"path/filepath"
)

func Move(id string, outChan chan *ResultSet, oldpath, newpath string) {
	OpAdd()
	if IsLocked(oldpath) {
		outChan <- FailedResultSet(id, "Path locked for another operation.")
		return
	}
	LockPath(oldpath)
	LockPath(newpath)

	defer func() {
		UnlockPath(oldpath)
		UnlockPath(newpath)
		OpDone()
	}()

	err := os.Rename(oldpath, newpath)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}
	outChan <- &ResultSet{
		Id: id,
	}
}

func Remove(id string, outChan chan *ResultSet, path string) {
	OpAdd()
	if IsLocked(path) {
		outChan <- FailedResultSet(id, "Path locked for another operation.")
		return
	}
	LockPath(path)
	defer func() {
		UnlockPath(path)
		OpDone()
	}()
	err := os.RemoveAll(path)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}
	outChan <- &ResultSet{
		Id: id,
	}
}

// Function for copying file/dirs

func Copy(id string, outChan chan *ResultSet, oldpath, newpath string) {
	file, err := os.Open(oldpath)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}
	file.Close()

	finfo, err := os.Stat(newpath)
	if err != nil || !finfo.IsDir() {
		outChan <- FailedResultSet(id, "Destination not valid.")
		return
	}

	// Append the filename to the new path
	_, f := filepath.Split(oldpath)
	newpath = filepath.Join(newpath, f)

	// Add to the wait group before the go routine
	// to avoid a race condition
	OpAdd()
	// BFS Copying method
	// Was initially a separate goroutine but is now a function
	// since the whole method is invoked as a goroutine
	func(id, oldpath, newpath string) {
		// Release the lock after the goroutine completes
		defer OpDone()

		queue := list.New()
		lockedPaths := make([]string, 0)
		queue.PushBack(oldpath)
		errStr := ""

		for queue.Len() > 0 {
			path := queue.Front().Value.(string)
			queue.Remove(queue.Front())

			file, err := os.Open(path)
			if err != nil {
				errStr += err.Error() + "\n"
				continue
			}
			lockedPaths = append(lockedPaths, path)
			LockPath(path)
			finfo, err := file.Stat()
			if err != nil {
				errStr += err.Error() + "\n"
				continue
			}
			npath := newpath + path[len(oldpath):]
			if finfo.IsDir() {
				err = os.Mkdir(npath, finfo.Mode().Perm())
				if err != nil {
					errStr += err.Error() + "\n"
					continue
				}
				sub, err := file.Readdirnames(-1)
				if err != nil {
					errStr += err.Error() + "\n"
					continue
				}
				for _, name := range sub {
					queue.PushBack(filepath.Join(path, name))
				}
			} else {
				nfile, err := os.OpenFile(npath, os.O_CREATE|os.O_WRONLY, 0777)
				if err != nil {
					errStr += err.Error() + "\n"
					continue
				}
				_, err = io.Copy(nfile, file)
				if err != nil {
					errStr += err.Error() + "\n"
					continue
				}
				_ = nfile.Chmod(finfo.Mode())
				nfile.Close()
			}
			file.Close()
		}
		for _, p := range lockedPaths {
			UnlockPath(p)
		}

		res := &ResultSet{
			Id:  id,
			Err: errStr,
		}
		outChan <- res

	}(id, oldpath, newpath)
}
