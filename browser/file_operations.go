package browser;

import (
	"os";
	"io";
	"io/ioutil";
	"encoding/json";
	"path/filepath";
)

// TODO: GetFile, PutFile, RunFile

const CHUNKSIZE = 2048

/* GetFile
	GetFile will write the content of files to stdout.
	GetFile2 initially returns a pointer to a ResultSet:
	The first result set contains a the path of the file,
	total number of pieces, and is numbered as piece 0.
	Use this for assembling the file.
	&ResultSet {
		Cmd: "cat",
		Path: <path of the file>,
		Data: &FileData{
			TotalPieces: <Total number of pieces>,
			CurrentPiece: 0,
			Data: nil,
		},
	}

	The file is then read in chunks of size CHUNKSIZE and
	written to stdout in json format.
*/

func GetFileDiv (file *os.File) int64 {
	finfo, _ := file.Stat();
	if finfo.Size() == 0 {
		return 1;
	}
	if finfo.Size() % CHUNKSIZE == 0 {
		return finfo.Size() / CHUNKSIZE;
	}
	return finfo.Size() / CHUNKSIZE + 1;
}

func GetFile (id, path string, out io.Writer) {
	OpAdd();
	defer OpDone();
	encoder := json.NewEncoder(out);
	if !ValidateDirPath(&path) || IsDir(path) {
		res := FailedResultSet(id, "getfile", path, "Not a valid path.");
		WriteJson(encoder, res);
		return;
	}
	file, err := os.Open(path);
	if err != nil {
		res := FailedResultSet(id, "getfile",path, err.Error());
		WriteJson(encoder, res);
		return;
	}
	maxdiv := GetFileDiv(file);
	res := &ResultSet{
		Id:  id,
		Cmd: "getfile",
		Path: path,
		Data: &FileData{
			TotalPieces: maxdiv,
			CurrentPiece: 0,
			Data: []byte{},
		},
	}
	WriteJson(encoder, res);

	buff := make([]byte, CHUNKSIZE);
	var i int64;
	defer file.Close();
	for i=1; i <= maxdiv; i++ {
		n, _ := file.Read(buff);
		res := &ResultSet{
			Id:  id,
			Cmd: "getfile",
			Path: path,
			Data: &FileData{
				TotalPieces: maxdiv,
				CurrentPiece: i,
				Data: buff[:n],
			},
		}
		WriteJson(encoder, res);
	}
}
/*
PutFile:
Receive command of the form
{
	"Cmd": "putfile",
	"Args": [<path>],
	"Data": <bytes>,
}
If directory does not exist return empty result set.
If file does not exist create it and write bytes from data.
else append bytes to end of file.
Return resultset
*/


func PutFile(id, path string, data []byte) interface{} {
	OpAdd();
	defer OpDone();
	file, err := os.OpenFile(path, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0666);
	if err != nil {
		// fmt.Println(err.Error());
		return FailedResultSet(id, "putfile",path, err.Error());
	}
	defer file.Close();
	_,_ = file.Write(data);
	return &ResultSet{
		Id : id,
		Cmd: "putfile",
		Path: path,
	};
}

/*
putFile2
Use instead of PutFile
params: 
	path : string, path where the file is to be written
	data : []byte, data to be written to the file

*/

var tempPath map[string]string = make(map[string]string);

func WriteToTemp(path string, data []byte) bool {
	file, err := os.OpenFile(path, os.O_APPEND | os.O_WRONLY, 0666);
	if err != nil {
		return false;
	}
	_, err = file.Write(data);
	if err != nil {
		return false;
	}
	file.Close();
	return true;
}

func PutFile2 (id, path string, data []byte) interface{} {
	OpAdd();
	defer OpDone();
	if tempPath[path] == "" {
		finfo, err := os.Stat(filepath.Dir(path));
		if err != nil || !finfo.IsDir() {
			return FailedResultSet(id, "putfile", path, "Path not valid.");
		}
		tf, err := ioutil.TempFile("", "putfile");
		if err != nil {
			return FailedResultSet(id, "putfile", path, err.Error());
		}
		tempPath[path] = tf.Name();
		tf.Close();
		if WriteToTemp(tempPath[path], data) == false {
			return FailedResultSet(id, "putfile", path, "Unable to write to temp file.");
		}
		LockPath(tempPath[path]);
		return &ResultSet{
			Id : id,
			Cmd: "putfile",
			Path: path,
		}
	}

	if len(data) == 0 {
		err := os.Rename(tempPath[path], path);
		tempPath[path] = "";
		UnlockPath(tempPath[path]);
		if err != nil {
			return FailedResultSet(id, "putfile", path, "Unable to move file to desired location.");
		}
		return &ResultSet{
			Id: id,
			Cmd: "putfile",
			Path: path,
		}
	}

	WriteToTemp(tempPath[path], data);

	return &ResultSet{
		Id : id,
		Cmd: "putfile",
		Path: path,
	};
};
