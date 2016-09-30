package browser;

type FileData struct {
	TotalPieces int64 `json:"totalPieces"`
	CurrentPiece int64 `json:"currentPieces"`
	Data []byte `json:"data"`
}

type FileInfo struct {
	Name string  `json:"name"`
	Size int64 `json:"size"`
	Dir bool `json:"dir"`
}

type ResultSet struct {
	Cmd string `json:"cmd"`
	Path string `json:"path"`
	Err string `json:"error"`
	Files []FileInfo `json:"files"`
	Data *FileData `json:"fileData"`
}

func FailedResultSet (cmd, path string, err string) *ResultSet {
	return &ResultSet{
		Cmd: cmd,
		Path: path,
		Err: err,
	}
}
