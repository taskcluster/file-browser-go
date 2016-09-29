package browser;

import "os";

var OutputFile *os.File = os.Stdout;

type Command struct {
	Cmd string `json:"cmd"`
	Args []string `json:"args"`
	Data []byte `json:"data"`
}

func ExitBrowser () interface{} {
	WaitForOperationsToComplete();
	return nil;
}

func Run (cmd Command, OutFile *os.File) interface{} {
	if OutFile == nil {
		OutputFile = os.Stdout;
	}else{
		OutputFile = OutFile;
	}
	if cmd.Cmd == "Exit" {
		return ExitBrowser();
	}
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

	}
	return FailedResultSet("","","No command specified.");
}
