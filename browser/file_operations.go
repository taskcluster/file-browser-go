package browser;

import (
	"os";
	"encoding/json";
	// "fmt";
	// "gopkg.in/vmihailenco/msgpack.v2";
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
	written to stdout in msgpack format.
*/

func GetFileDiv (file *os.File) int64 {
	finfo, _ := file.Stat();
	if finfo.Size() % CHUNKSIZE == 0 {
		return finfo.Size() / CHUNKSIZE;
	}
	return finfo.Size() / CHUNKSIZE + 1;
}

func GetFile (path string ) *ResultSet{
	if !ValidateDirPath(&path) || IsDir(path) {
		return FailedResultSet("GetFile", path, "Not a valid path.");
	}
	file, err := os.Open(path);
	if err != nil {
		return FailedResultSet("GetFile",path, err.Error());
	}
	maxdiv := GetFileDiv(file);
	defer WritePiecesToStdout(path, file, maxdiv);
	return &ResultSet{
		Cmd: "GetFile",
		Path: path,
		Data: &FileData{
			TotalPieces: maxdiv,
			CurrentPiece: 0,
			Data: []byte{},
		},
	}
}

func WritePiecesToStdout(path string, file *os.File, maxdiv int64){
	enc := json.NewEncoder(os.Stdout);
	buff := make([]byte, CHUNKSIZE);
	var i int64;
	go func(){
		defer file.Close();
		for i = 1; i <= maxdiv; i++ {
			_, _ = file.Read(buff);
			res := &ResultSet{
				Cmd: "GetFile",
				Path: path,
				Data: &FileData{
					TotalPieces: maxdiv,
					CurrentPiece: i,
					Data: buff,
				},
			}
			enc.Encode(res);
		}

	}();
}

/*
PutFile:
Receive command of the form
{
	"Cmd": "PutFile",
	"Args": [<path>],
	"Data": <bytes>,
}
If directory does not exist return empty result set.
If file does not exist create it and write bytes from data.
else append bytes to end of file.
Return resultset
*/

func PutFile(path string, data []byte) *ResultSet {
	file, err := os.OpenFile(path, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0666);
	if err != nil {
		// fmt.Println(err.Error());
		return FailedResultSet("PutFile",path, err.Error());
	}
	defer file.Close();
	_,_ = file.Write(data);
	return &ResultSet{
		Cmd: "PutFile",
		Path: path,
	};
}

