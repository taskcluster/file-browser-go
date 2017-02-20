package browser

/*
FileData is used to send files over stdout.
The chunks are sent in order so that they may be
easily reassembled at the other end.
CurrentPiece -> [1, TotalPieces]
When CurrentPiece == TotalPieces, all the pieces have
been sent
*/
type FileData struct {
	TotalPieces  int64  `msgpack:"totalPieces"`
	CurrentPiece int64  `msgpack:"currentPiece"`
	Data         []byte `msgpack:"data"`
}

/*
FileInfo is used to send the result of a List operation
*/
type FileInfo struct {
	Name string `msgpack:"name"`
	Size int64  `msgpack:"size"`
	Dir  bool   `msgpack:"dir"`
}

/*
The ResultSet struct which is written over stdout
*/
type ResultSet struct {
	Id    string     `msgpack:"id"`
	Err   string     `msgpack:"error,omitempty"`
	Files []FileInfo `msgpack:"files,omitempty"`
	Data  *FileData  `msgpack:"fileData,omitempty"`
}

/*
Creates a ResultSet with Id = id and Err = err
*/
func FailedResultSet(id, err string) *ResultSet {
	return &ResultSet{
		Id:  id,
		Err: err,
	}
}
