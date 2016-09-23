package browser;

import (
	"os";
	"path/filepath";
)

type Command struct {
	Cmd string;
	Args []string;
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
	switch cmd.Cmd{
	case "ls":
		if len(cmd.Args) == 0{
			return EmptyResultSet();
		}
		return List(cmd.Args[0]);

	case "cat":
		if len(cmd.Args) == 0{
			return EmptyResultSet();
		}
		return Cat2(cmd.Args[0]);

	case "exit":
		return nil;
	}
	return nil;
}
