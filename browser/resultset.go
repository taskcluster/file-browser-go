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
	Err error `json:"Error"`
	Data *FileData `json:"FileData"`
}

func ErrorResultSet (e error) *ResultSet {
	return &ResultSet{
		Err: e,
	}
}

func NewListResultSet (d, f []string, path string, e error) *ResultSet {
	return &ResultSet{
		Cmd: "ls",
		Path: path,
		Dirs: d,
		Files: f,
		Err: e,
	}
}

func NewCatResultSet(path string) *ResultSet {
	return &ResultSet{
		Cmd: "cat",
		Path: path,
	}
}

func EmptyResultSet () *ResultSet {
	return &ResultSet{};
}

func (r *ResultSet) GetDirs() []string {
	return r.Dirs;
}

func (r *ResultSet) GetFiles() []string {
	return r.Files;
}

func (r *ResultSet) GetError() error {
	return r.Err;
}
