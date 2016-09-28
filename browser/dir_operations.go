package browser;

import (
	"io/ioutil";
	"os";
)

func List(path string) interface{} {
	dirs, files := []string{}, []string{};
	if !ValidateDirPath(&path)|| !IsDir(path) {
		return FailedResultSet("List",path, "Not a directory.");
	}
	finfo, err := ioutil.ReadDir(path);
	if err != nil {
		return FailedResultSet("List",path, err.Error());
	}
	for _, f := range finfo {
		if f.IsDir() {
			dirs = append(dirs, f.Name());
		}else{
			files = append(files, f.Name());
		}
	}
	return &ResultSet{
		Dirs: dirs,
		Files: files,
		Cmd: "List",
		Path: path,
	}
}

func MakeDirectory (path string) interface{} {
	if !ValidateDirPath(&path) {
		return FailedResultSet("MakeDir",path, "Not a valid path.");
	}
	err := os.Mkdir(path, 0777);
	if err != nil {
		return FailedResultSet("MakeDir",path, err.Error());
	}
	return &ResultSet{
		Cmd: "MakeDir",
		Path: path,
	}
}
