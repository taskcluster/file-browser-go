package browser

import (
	"io"
	"os"
)

func init() {
	RegisterCommand("getfile", OnePathWrapper(GetFile))
	RegisterCommand("putfile", PutFileWrapper())
}

const CHUNKSIZE = 2048

/* GetFile
GetFile will write the content of files to stdout.
GetFile2 initially returns a pointer to a ResultSet:
The first result set contains a the path of the file,
total number of pieces, and is numbered as piece 0.
Use this for assembling the file.

The file is then read in chunks of size CHUNKSIZE and
written to stdout in msgpack format.
*/

func GetFileDiv(file *os.File) int64 {
	finfo, _ := file.Stat()
	if finfo.Size() == 0 {
		return 1
	}
	if finfo.Size()%CHUNKSIZE == 0 {
		return finfo.Size() / CHUNKSIZE
	}
	return finfo.Size()/CHUNKSIZE + 1
}

func GetFile(id string, out chan<- *ResultSet, path string) {
	OpAdd()
	defer OpDone()
	path, valid := ValidateDirPath(path)

	if !valid || IsDir(path) {
		out <- FailedResultSet(id, "Not a valid path.")
		return
	}
	file, err := os.Open(path)
	if err != nil {
		out <- FailedResultSet(id, err.Error())
		return
	}
	maxdiv := GetFileDiv(file)

	buff := make([]byte, CHUNKSIZE)

	// Declared because GetMaxDiv returns int64
	var i int64
	defer file.Close()
	for i = 1; i <= maxdiv; i++ {
		n, err := file.Read(buff)

		if err != nil && err != io.EOF {
			out <- &ResultSet{
				Id:  id,
				Err: err.Error(),
			}
			break
		}

		if err == io.EOF {
			out <- &ResultSet{
				Id: id,
				Data: &FileData{
					CurrentPiece: i,
					Data:         []byte{},
				},
			}
			break
		}

		out <- &ResultSet{
			Id: id,
			Data: &FileData{
				TotalPieces:  maxdiv,
				CurrentPiece: i,
				Data:         buff[:n],
			},
		}
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

func PutFile(id string, outChan chan<- *ResultSet, path string, data []byte) {
	OpAdd()
	defer OpDone()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}
	defer file.Close()
	_, err = file.Write(data)

	if err != nil {
		outChan <- FailedResultSet(id, err.Error())
		return
	}

	outChan <- &ResultSet{
		Id: id,
	}
}
