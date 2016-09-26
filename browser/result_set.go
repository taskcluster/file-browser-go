package browser;

type FileData struct {
	TotalPieces int64 `json:"TotalPieces"`
	CurrentPiece int64 `json:"CurrentPieces"`
	Data []byte `json:"Data"`
}

type ResultSet struct {
	Cmd string `json:"Cmd"`
	Path string `json:"Path"`
	Dirs []string `json:"Directories"`
	Files []string `json:"Files"`
	Err string `json:"Error"`
	Data *FileData `json:"FileData"`
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
