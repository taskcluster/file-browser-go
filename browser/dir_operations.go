package browser

import (
	"io/ioutil"
	"os"
)

func init() {
	registerCommand("ls", onePathWrapper(List))
	registerCommand("mkdir", onePathWrapper(MakeDirectory))
}

func List(id string, outChan chan<- *ResultSet, path string) {
	OpAdd()
	defer OpDone()
	path, valid := validateDirPath(path)
	if !valid || !isDir(path) {
		outChan <- FailedResultSet(id, "Not a directory.")
		return
	}
	finfo, err := ioutil.ReadDir(path)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}
	files := []FileInfo{}
	for _, f := range finfo {
		files = append(files, FileInfo{
			Name: f.Name(),
			Size: f.Size(),
			Dir:  f.IsDir(),
		})
	}
	outChan <- &ResultSet{
		Id:    id,
		Files: files,
	}
}

func MakeDirectory(id string, outChan chan<- *ResultSet, path string) {
	OpAdd()
	defer OpDone()
	path, valid := validateDirPath(path)
	if !valid {
		outChan <- FailedResultSet(id, "Not a valid path.")
		return
	}
	err := os.Mkdir(path, 0777)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}
	outChan <- &ResultSet{
		Id: id,
	}
}
