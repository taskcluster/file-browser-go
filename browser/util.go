package browser

import (
	"os"
	"path/filepath"
	"sync"
)

// Locking paths to make sure there's no interference

var pathLock sync.Mutex

var isLocked map[string]bool = make(map[string]bool)

func LockPath(path string) {
	pathLock.Lock()
	isLocked[path] = true
	pathLock.Unlock()
}

func UnlockPath(path string) {
	pathLock.Lock()
	isLocked[path] = false
	pathLock.Unlock()
}

func IsLocked(path string) bool {
	pathLock.Lock()
	defer pathLock.Unlock()
	if isLocked[path] {
		return true
	}
	dir, f := filepath.Split(path)
	for f != "" {
		if isLocked[dir] {
			return true
		}
		dir, f = filepath.Split(dir)
	}
	return false
}

// On exit make browser wait until every copy operation
// is complete
var op sync.WaitGroup

func OpAdd() {
	op.Add(1)
}

func OpDone() {
	op.Done()
}

// Utility functions
func WaitForOperationsToComplete() {
	op.Wait()
}

func IsDir(dir string) bool {
	file, err := os.Open(dir)
	defer file.Close()
	if err != nil {
		return false
	}
	finfo, err := file.Stat()
	if err != nil {
		return false
	}
	return finfo.IsDir()
}

func ValidateDirPath(dir *string) bool {
	*dir = filepath.Clean(*dir)
	if !filepath.IsAbs(*dir) {
		return false
	}
	return true
}

// Helper functions to wrap methods

func OnePathWrapper(fun func(string, chan<- *ResultSet, string)) func(Command, chan<- *ResultSet) {
	return func(cmd Command, outChan chan<- *ResultSet) {
		if len(cmd.Args) == 0 {
			outChan <- FailedResultSet(cmd.Id, "No path specified")
			return
		}
		fun(cmd.Id, outChan, cmd.Args[0])
	}
}

func TwoPathWrapper(fun func(string, chan<- *ResultSet, string, string)) func(Command, chan<- *ResultSet) {
	return func(cmd Command, outChan chan<- *ResultSet) {
		if len(cmd.Args) < 2 {
			outChan <- FailedResultSet(cmd.Id, "Not enough arguments")
			return
		}
		fun(cmd.Id, outChan, cmd.Args[0], cmd.Args[1])
	}
}

func PutFileWrapper() func(Command, chan<- *ResultSet) {
	return func(cmd Command, outChan chan<- *ResultSet) {
		if len(cmd.Args) == 0 {
			outChan <- FailedResultSet(cmd.Id, "Not enough arguments")
			return
		}
		PutFile(cmd.Id, outChan, cmd.Args[0], cmd.Data)
	}
}
