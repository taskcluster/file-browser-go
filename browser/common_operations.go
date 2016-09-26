package browser;

import (
	"os";
	"io";
	"path/filepath";
	"encoding/json";
)

func Move (oldpath, newpath string) *ResultSet {
	if IsLocked[oldpath] {
		return FailedResultSet("Move", oldpath, "Path locked for another operation.");
	}
	LockPath(oldpath);
	LockPath(newpath);

	defer func (){
		UnlockPath(oldpath);
		UnlockPath(newpath);
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

func Remove (path string) *ResultSet {
	if IsLocked[path] {
		return FailedResultSet("Move", path, "Path locked for another operation.");
	}
	LockPath(path);
	defer UnlockPath(path);
	err := os.RemoveAll(path);
	if err != nil {
		return FailedResultSet("Remove", path, err.Error());
	}
	return &ResultSet{
		Cmd: "Remove",
		Path: path,
	};
}

func Copy (oldpath, newpath string) *ResultSet {
	file, err := os.Open(oldpath);
	if err != nil {
		return FailedResultSet("Copy", oldpath, err.Error());
	}
	file.Close();
	go CopyUtil(oldpath, newpath);
	return &ResultSet{
		Cmd: "Copy Started",
		Path: newpath,
	}
}

// BFS copying method
func CopyUtil (oldpath, newpath string) {
	enc := json.NewEncoder(os.Stdout);
	queue := make([]string,0);
	failedFiles := make([]string, 0);
	failedDirs := make([]string, 0);
	lockedPaths := make([]string,0);
	queue = append(queue, oldpath);
	errStr := "";
	
	for len(queue) != 0 {
		path := queue[0];
		queue = queue[1:];

		file, err := os.Open(path);
		if err != nil {
			failedFiles = append(failedFiles, path);
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
				failedDirs = append(failedDirs, path);
				errStr += err.Error() + "\n";
				continue;
			}
			sub, err := file.Readdirnames(-1);
			if err != nil {
				failedDirs = append(failedDirs, path);	
				errStr += err.Error() + "\n";
				continue;
			}
			for _, name := range sub {
				queue = append(queue, filepath.Join(path, name));
			}
		}else{
			nfile, err := os.OpenFile(npath, os.O_CREATE | os.O_WRONLY, 0777);
			if err != nil {
				failedFiles = append(failedFiles, path);	
				errStr += err.Error() + "\n";
				continue;
			}
			_, err = io.Copy(nfile,file);
			if err != nil {
				failedFiles = append(failedFiles, path);	
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
		Files: failedFiles,
		Dirs: failedDirs,
		Err: errStr,
	}
	enc.Encode(res);
}
