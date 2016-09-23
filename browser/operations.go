package browser;

import (
	"io/ioutil";
	"os";
	"gopkg.in/vmihailenco/msgpack.v2";
)

const PART_SIZE = 2048

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

func Cat(dir string) *ResultSet{
	if !ValidateDirPath(&dir) || IsDir(dir) {
		return EmptyResultSet();
	}
	file, err := os.Open(dir);
	if err != nil {
		return EmptyResultSet();
	}

	_, err = ioutil.ReadAll(file);
	if err != nil {
		return EmptyResultSet();
	}

	return NewCatResultSet(dir);
}

// Cat utility functions

func CatDiv (file *os.File) int64 {
	finfo, _ := file.Stat();
	if finfo.Size() % PART_SIZE == 0 {
		return finfo.Size() / PART_SIZE;
	}
	return finfo.Size() / PART_SIZE + 1;
}

func Cat2 (dir string ) *ResultSet{
	enc := msgpack.NewEncoder(os.Stdout);
	if !ValidateDirPath(&dir) || IsDir(dir) {
		return EmptyResultSet();
	}
	file, _ := os.Open(dir);
	maxdiv := CatDiv(file);
	defer WritePiecesToStdout(dir, file, maxdiv, enc);
	return &ResultSet{
		Cmd: "cat",
		Path: dir,
		Data: &FileData{
			TotalPieces: maxdiv,
			CurrentPiece: 0,
			Data: []byte{},
		},
	}
}

func WritePiecesToStdout(dir string, file *os.File, maxdiv int64, enc *msgpack.Encoder){
	buff := make([]byte, 2048);
	var i int64;
	go func(){
		for i = 1; i <= maxdiv; i++ {
			_, _ = file.Read(buff);
			res := &ResultSet{
				Cmd: "cat",
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
