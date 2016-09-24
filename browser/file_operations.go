package browser;

import (
	"io/ioutil";
	"os";
	"path/filepath";
	"gopkg.in/vmihailenco/msgpack.v2";
)

// TODO: List, GetFile, PutFile, RunFile

const CHUNKSIZE = 2048

func List(dir string) *ResultSet {
	dirs, files := []string{}, []string{};
	if !ValidateDirPath(&dir)|| !IsDir(dir) {
		return EmptyResultSet();
	}
	finfo, err := ioutil.ReadDir(dir);
	if err != nil {
		return EmptyResultSet();
	}
	for _, f := range finfo {
		if f.IsDir() {
			dirs = append(dirs, f.Name());
		}else{
			files = append(files, f.Name());
		}
	}
	return NewListResultSet(dirs, files, dir, nil);
}

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

func GetFile (dir string ) *ResultSet{
	if !ValidateDirPath(&dir) || IsDir(dir) {
		return EmptyResultSet();
	}
	file, err := os.Open(dir);
	if err != nil {
		return EmptyResultSet();
	}
	maxdiv := GetFileDiv(file);
	defer WritePiecesToStdout(dir, file, maxdiv);
	return &ResultSet{
		Cmd: "GetFile",
		Path: dir,
		Data: &FileData{
			TotalPieces: maxdiv,
			CurrentPiece: 0,
			Data: []byte{},
		},
	}
}

func WritePiecesToStdout(dir string, file *os.File, maxdiv int64){
	enc := msgpack.NewEncoder(os.Stdout);
	buff := make([]byte, CHUNKSIZE);
	var i int64;
	go func(){
		for i = 1; i <= maxdiv; i++ {
			_, _ = file.Read(buff);
			res := &ResultSet{
				Cmd: "GetFile",
				Path: dir,
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
	if !IsDir(filepath.Dir(path)) {
		return EmptyResultSet();
	}
	file, err := os.OpenFile(path, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0666);
	if err != nil {
		return EmptyResultSet();
	}
	defer file.Close();
	_,_ = file.Write(data);
	return &ResultSet{
		Cmd: "PutFile",
		Path: path,
	};
}
