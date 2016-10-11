package browser;

import (
	"os";
	"io";
	"path/filepath";
	"encoding/json";
	"container/list";
)

func Move (oldpath, newpath string) interface{} {
	OpAdd();
	if IsLocked(oldpath) {
		return FailedResultSet("Move", oldpath, "Path locked for another operation.");
	}
	LockPath(oldpath);
	LockPath(newpath);

	defer func (){
		UnlockPath(oldpath);
		UnlockPath(newpath);
		OpDone();
	}();

	err := os.Rename(oldpath, newpath);
	if err != nil {
		return FailedResultSet("Move", oldpath, err.Error());
	}
	return &ResultSet{
		Cmd: "Move",
		Path: newpath,
	}
}

func Remove (path string) interface{} {
	OpAdd();
	if IsLocked(path) {
		return FailedResultSet("Move", path, "Path locked for another operation.");
	}
	LockPath(path);
	defer func() {
		UnlockPath(path);
		OpDone();
	}();
	err := os.RemoveAll(path);
	if err != nil {
		return FailedResultSet("Remove", path, err.Error());
	}
	return &ResultSet{
		Cmd: "Remove",
		Path: path,
	};
}

// Function for copying file/dirs

func Copy (oldpath, newpath string, out io.Writer) interface{} {
	file, err := os.Open(oldpath);
	if err != nil {
		return FailedResultSet("Copy", oldpath, err.Error());
	}
	file.Close();

	finfo, err := os.Stat(newpath);
	if err != nil || !finfo.IsDir() {
		return FailedResultSet("Copy", newpath, "Destination not valid.");
	}

	// Append the filename to the new path
	_, f := filepath.Split(oldpath);
	newpath = filepath.Join(newpath,f);

	// Add to the wait group before the go routine
	// to avoid a race condition
	OpAdd();
	// BFS Copying method
	go func (oldpath, newpath string) {
		// Release the lock after the goroutine completes
		defer OpDone();

		enc := json.NewEncoder(out);
		queue := list.New();
		lockedPaths := make([]string,0);
		queue.PushBack(oldpath);
		errStr := "";

		for queue.Len() > 0 {
			path := queue.Front().Value.(string);
			queue.Remove(queue.Front());

			file, err := os.Open(path);
			if err != nil {
				errStr += err.Error() + "\n";
				continue;
			}
			lockedPaths = append(lockedPaths, path);
			LockPath(path);
			finfo, err := file.Stat();
			if err != nil {
				errStr += err.Error() + "\n";
				continue;
			}
			npath := newpath + path[len(oldpath):]
			if finfo.IsDir() {
				err = os.Mkdir(npath, finfo.Mode().Perm());
				if err != nil {
					errStr += err.Error() + "\n";
					continue;
				}
				sub, err := file.Readdirnames(-1);
				if err != nil {
					errStr += err.Error() + "\n";
					continue;
				}
				for _, name := range sub {
					queue.PushBack(filepath.Join(path,name));
				}
			}else{
				nfile, err := os.OpenFile(npath, os.O_CREATE | os.O_WRONLY, 0777);
				if err != nil {
					errStr += err.Error() + "\n";
					continue;
				}
				_, err = io.Copy(nfile,file);
				if err != nil {
					errStr += err.Error() + "\n";
					continue;
				}
				_ = nfile.Chmod(finfo.Mode());
				nfile.Close();
			}
			file.Close();
		}
		for _, p := range lockedPaths {
			UnlockPath(p);
		}

		res := &ResultSet{
			Cmd: "Copy Complete",
			Path: newpath,
			Err: errStr,
		}
		enc.Encode(res);

	}(oldpath, newpath);

	return &ResultSet{
		Cmd: "Copy Started",
		Path: newpath,
	}
}
