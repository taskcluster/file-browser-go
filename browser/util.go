package browser;

import (
	"sync";
	"os";
	"path/filepath";
)


// Locking paths to make sure there's no interference

var IsLocked map[string]bool = make(map[string]bool);

func LockPath (path string) {
	IsLocked[path] = true;
}

func UnlockPath (path string) {
	IsLocked[path] = false;
}

// On exit make browser wait until every copy operation
// is complete

var copy_running sync.WaitGroup;

func CopyAdd () {
	copy_running.Add(1);
}

func CopyDone () {
	copy_running.Done();
}

// GetFile Lock

var get_file_running sync.WaitGroup;

func LockFile(){
	get_file_running.Add(1);
}

func UnlockFile(){
	get_file_running.Done();
}

// Utility functions
func WaitForOperationsToComplete() {
	copy_running.Wait();
	get_file_running.Wait();
}

func IsDir (dir string) bool {
	file, err := os.Open(dir);
	defer file.Close();
	if err != nil {
		return false;
	}
	finfo, err := file.Stat();
	if err != nil {
		return false;
	}
	return finfo.IsDir();
}

func ValidateDirPath (dir *string) bool {
	*dir = filepath.Clean(*dir);
	if !filepath.IsAbs(*dir) {
		return false;
	}
	return true;
}
