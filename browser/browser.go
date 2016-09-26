package browser;

import (
	"os";
	"path/filepath";
)

var IsLocked map[string]bool = make(map[string]bool);

func LockPath (path string) {
	IsLocked[path] = true;
}

func UnlockPath (path string) {
	IsLocked[path] = false;
}

type Command struct {
	Cmd string;
	Args []string;
	Data []byte;
}

// Utility functions

func IsDir (dir string) bool {
	file, err := os.Open(dir);
	defer file.Close();
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
		return FailedResultSet(cmd.Cmd,"", "Not enough arguments.");
	}
	switch cmd.Cmd{
	case "List":
		return List(cmd.Args[0]);

	case "GetFile":
		return GetFile(cmd.Args[0]);

	case "PutFile":
		return PutFile(cmd.Args[0], cmd.Data);

	case "Move":
		if len(cmd.Args) < 1 {
			return FailedResultSet(cmd.Cmd,"","Not enough arguments.");
		}
		return Move(cmd.Args[0], cmd.Args[1]);
	case "Copy":
		if len(cmd.Args) < 1 {
			return FailedResultSet(cmd.Cmd,"","Not enough arguments.");
		}
		return Copy(cmd.Args[0], cmd.Args[1]);
	case "Remove":
		return Remove(cmd.Args[0]);
	case "MakeDir":
		return MakeDirectory(cmd.Args[0]);

	case "Exit":
		return nil;
	}
	return FailedResultSet("","","No command specified.");
}
