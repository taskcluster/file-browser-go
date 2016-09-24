package browser;

import (
	"os";
	"path/filepath";
)

type Command struct {
	Cmd string;
	Args []string;
	Data []byte;
}

// Utility functions

func IsDir (dir string) bool {
	file, err := os.Open(dir);
	if err != nil {
		return false;
	}
	finfo, err := file.Stat();
	if err != nil {
		return false;
	}
	return finfo.IsDir();
}

func ValidateDirPath (dir *string) bool {
	*dir = filepath.Clean(*dir);
	if !filepath.IsAbs(*dir) {
		return false;
	}
	return true;
}

func Run (cmd Command) *ResultSet{
	if len(cmd.Args) == 0 {
		return EmptyResultSet();
	}
	switch cmd.Cmd{
	case "List":
		return List(cmd.Args[0]);

	case "GetFile":
		return GetFile(cmd.Args[0]);

	case "PutFile":
		return PutFile(cmd.Args[0], cmd.Data);

	case "Exit":
		return nil;
	}
	return EmptyResultSet();
}
