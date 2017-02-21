package browser

import (
	"os"
	"path/filepath"
	"sync"
)

// Locking paths to make sure there's no interference

var pathLock sync.Mutex

var lck = make(map[string]bool)

func lockPath(path string) {
	pathLock.Lock()
	lck[path] = true
	pathLock.Unlock()
}

func unlockPath(path string) {
	pathLock.Lock()
	lck[path] = false
	pathLock.Unlock()
}

func isLocked(path string) bool {
	pathLock.Lock()
	defer pathLock.Unlock()
	if lck[path] {
		return true
	}
	dir, f := filepath.Split(path)
	for f != "" {
		if lck[dir] {
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

func isDir(dir string) bool {
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

func validateDirPath(dir string) (string, bool) {
	cleanDir := filepath.Clean(dir)
	if !filepath.IsAbs(cleanDir) {
		return cleanDir, false
	}
	return cleanDir, true
}

// Helper functions to wrap methods

func onePathWrapper(fun func(string, chan<- *ResultSet, string)) func(Command, chan<- *ResultSet) {
	return func(cmd Command, outChan chan<- *ResultSet) {
		if len(cmd.Args) == 0 {
			outChan <- FailedResultSet(cmd.Id, "No path specified")
			return
		}
		fun(cmd.Id, outChan, cmd.Args[0])
	}
}

func twoPathWrapper(fun func(string, chan<- *ResultSet, string, string)) func(Command, chan<- *ResultSet) {
	return func(cmd Command, outChan chan<- *ResultSet) {
		if len(cmd.Args) < 2 {
			outChan <- FailedResultSet(cmd.Id, "Not enough arguments")
			return
		}
		fun(cmd.Id, outChan, cmd.Args[0], cmd.Args[1])
	}
}

func putFileWrapper() func(Command, chan<- *ResultSet) {
	return func(cmd Command, outChan chan<- *ResultSet) {
		if len(cmd.Args) == 0 {
			outChan <- FailedResultSet(cmd.Id, "Not enough arguments")
			return
		}
		PutFile(cmd.Id, outChan, cmd.Args[0], cmd.Data)
	}
}
