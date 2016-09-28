package browser;

type FileData struct {
	TotalPieces int64 `json:"totalPieces"`
	CurrentPiece int64 `json:"currentPieces"`
	Data []byte `json:"data"`
}

type ResultSet struct {
	Cmd string `json:"cmd"`
	Path string `json:"path"`
	Dirs []string `json:"directories"`
	Files []string `json:"files"`
	Err string `json:"error"`
	Data *FileData `json:"fileData"`
}

func FailedResultSet (cmd, path string, err string) *ResultSet {
	return &ResultSet{
		Cmd: cmd,
		Path: path,
		Err: err,
	}
}

func (r *ResultSet) GetDirs() []string {
	return r.Dirs;
}

func (r *ResultSet) GetFiles() []string {
	return r.Files;
}

func (r *ResultSet) GetError() string {
	return r.Err;
}
