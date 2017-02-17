package browser

import (
	"io"
	"os"
)

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

func GetFile(id string, out chan *ResultSet, path string) {
	OpAdd()
	defer OpDone()
	if !ValidateDirPath(&path) || IsDir(path) {
		res := FailedResultSet(id, "Not a valid path.")
		out <- res
		return
	}
	file, err := os.Open(path)
	if err != nil {
		res := FailedResultSet(id, err.Error())
		out <- res
		return
	}
	maxdiv := GetFileDiv(file)

	buff := make([]byte, CHUNKSIZE)
	var i int64
	defer file.Close()
	for i = 1; i <= maxdiv; i++ {
		var res *ResultSet = nil
		n, err := file.Read(buff)

		if err != nil && err != io.EOF {
			out <- &ResultSet{
				Id:  id,
				Err: err.Error(),
			}
			break
		}

		res = &ResultSet{
			Id: id,
			Data: &FileData{
				TotalPieces:  maxdiv,
				CurrentPiece: i,
				Data:         buff[:n],
			},
		}
		out <- res
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

func PutFile(id string, outChan chan *ResultSet, path string, data []byte) {
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
