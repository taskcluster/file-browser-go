package browser;

import (
	"os";
	"io";
	"path/filepath";
	"container/list";
	"gopkg.in/vmihailenco/msgpack.v2";
)

func Move (id, oldpath, newpath string) interface{} {
	OpAdd();
	if IsLocked(oldpath) {
		return FailedResultSet(id, "mv", oldpath, "Path locked for another operation.");
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
		return FailedResultSet(id, "mv", oldpath, err.Error());
	}
	return &ResultSet{
		Id : id,
		Cmd: "mv",
		Path: newpath,
	}
}

func Remove (id, path string) interface{} {
	OpAdd();
	if IsLocked(path) {
		return FailedResultSet(id, "mv", path, "Path locked for another operation.");
	}
	LockPath(path);
	defer func() {
		UnlockPath(path);
		OpDone();
	}();
	err := os.RemoveAll(path);
	if err != nil {
		return FailedResultSet(id, "rm", path, err.Error());
	}
	return &ResultSet{
		Id : id,
		Cmd: "rm",
		Path: path,
	};
}

// Function for copying file/dirs

func Copy (id, oldpath, newpath string, out io.Writer) interface{} {
	file, err := os.Open(oldpath);
	if err != nil {
		return FailedResultSet(id, "cp", oldpath, err.Error());
	}
	file.Close();

	finfo, err := os.Stat(newpath);
	if err != nil || !finfo.IsDir() {
		return FailedResultSet(id, "cp", newpath, "Destination not valid.");
	}

	// Append the filename to the new path
	_, f := filepath.Split(oldpath);
	newpath = filepath.Join(newpath,f);

	// Add to the wait group before the go routine
	// to avoid a race condition
	OpAdd();
	// BFS Copying method
	go func (id, oldpath, newpath string) {
		// Release the lock after the goroutine completes
		defer OpDone();

		enc := msgpack.NewEncoder(out);
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
			Id : id,
			Cmd: "cp",
			Path: newpath,
			Err: errStr,
		}
		WriteOut(enc, res);

	}(id, oldpath, newpath);

	return nil;

}
