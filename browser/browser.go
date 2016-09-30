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
	return &ResultSet{
		Cmd:"Exit",
	}
}

func Run (cmd Command, out chan interface{} ) {
	if cmd.Cmd == "Exit" {
		out <- ExitBrowser();
		return;
	}
	if len(cmd.Args) == 0 {
		out <-  FailedResultSet(cmd.Cmd,"", "Not enough arguments.");
		return;
	}
	switch cmd.Cmd{
	case "List":
		out <-  List(cmd.Args[0]);
		return;

	case "GetFile":
		GetFile(cmd.Args[0], out);
		return;

	case "PutFile":
		out <-  PutFile(cmd.Args[0], cmd.Data);
		return;

	case "Move":
		if len(cmd.Args) < 1 {
			out <-  FailedResultSet(cmd.Cmd,"","Not enough arguments.");
			return;
		}
		out <-  Move(cmd.Args[0], cmd.Args[1]);
		return;
	case "Copy":
		if len(cmd.Args) < 1 {
			out <-  FailedResultSet(cmd.Cmd,"","Not enough arguments.");
			return;
		}
		out <-  Copy(cmd.Args[0], cmd.Args[1]);
		return;
	case "Remove":
		out <-  Remove(cmd.Args[0]);
		return;
	case "MakeDir":
		out <-  MakeDirectory(cmd.Args[0]);
		return;

	}
	out <-  FailedResultSet("","","No command specified.");
}
