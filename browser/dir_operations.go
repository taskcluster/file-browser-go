package browser;

import (
	"io/ioutil";
	"os";
)

func List(id, path string) interface{} {
	OpAdd();
	defer OpDone();
	if !ValidateDirPath(&path)|| !IsDir(path) {
		return FailedResultSet(id, "Not a directory.");
	}
	finfo, err := ioutil.ReadDir(path);
	if err != nil {
		return FailedResultSet(id, err.Error());
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
		Id : id,
		// Cmd: "ls",
		// Path: path,
		Files: files,
	}
}

func MakeDirectory (id, path string) interface{} {
	OpAdd();
	defer OpDone();
	if !ValidateDirPath(&path) {
		return FailedResultSet(id, "Not a valid path.");
	}
	err := os.Mkdir(path, 0777);
	if err != nil {
		return FailedResultSet(id, err.Error());
	}
	return &ResultSet{
		Id : id,
		// Cmd: "mkdir",
		// Path: path,
	}
}
