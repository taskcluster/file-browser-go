package browser;

import (
	"io/ioutil";
	"os";
)

func List(path string) interface{} {
	OpAdd();
	defer OpDone();
	if !ValidateDirPath(&path)|| !IsDir(path) {
		return FailedResultSet("ls",path, "Not a directory.");
	}
	finfo, err := ioutil.ReadDir(path);
	if err != nil {
		return FailedResultSet("ls",path, err.Error());
	}
	files := []FileInfo{};
	for _, f := range finfo {
		files = append(files, FileInfo{
			Name: f.Name(),
			Size: f.Size(),
			Dir: f.IsDir(),
		});
	}
	return &ResultSet{
		Cmd: "ls",
		Path: path,
		Files: files,
	}
}

func MakeDirectory (path string) interface{} {
	OpAdd();
	defer OpDone();
	if !ValidateDirPath(&path) {
		return FailedResultSet("mkdir",path, "Not a valid path.");
	}
	err := os.Mkdir(path, 0777);
	if err != nil {
		return FailedResultSet("mkdir",path, err.Error());
	}
	return &ResultSet{
		Cmd: "mkdir",
		Path: path,
	}
}
