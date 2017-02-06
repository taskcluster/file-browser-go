package browser;

import (
	"os";
	"encoding/json";
	"io";
  "fmt";
)

type Command struct {
	Id string `json:"id"`
	Cmd string `json:"cmd"`
	Args []string `json:"args"`
	Data []byte `json:"data"`
}

func RunCmd (cmd Command, out io.Writer) {
	encoder := json.NewEncoder(out);
	if cmd.Id == "" {
		encoder.Encode(FailedResultSet("null", cmd.Cmd, "", "No Id supplied."));
		return;
	}
	if len(cmd.Args) == 0 {
		res := FailedResultSet(cmd.Id, cmd.Cmd,"", "Not enough arguments.");
		WriteJson(encoder, res);
		return;
	}

	id := cmd.Id;

	switch cmd.Cmd{

	case "ls":
		res := List(cmd.Id, cmd.Args[0]);
		WriteJson(encoder, res);
		return;

	case "getfile":
		GetFile(cmd.Id, cmd.Args[0], out);
		return;

	case "putfile":
		res := PutFile2(cmd.Id, cmd.Args[0], cmd.Data);
		WriteJson(encoder, res);
		return;

	case "mv":
		if len(cmd.Args) < 1 {
			res := FailedResultSet(id, cmd.Cmd,"","Not enough arguments.");
			WriteJson(encoder, res);
			return;
		}
		res := Move(cmd.Id, cmd.Args[0], cmd.Args[1]);
		WriteJson(encoder, res);
		return;
	case "cp":
		if len(cmd.Args) < 1 {
			res := FailedResultSet(id, cmd.Cmd,"","Not enough arguments.");
			WriteJson(encoder, res);
			return;
		}
		res := Copy(cmd.Id, cmd.Args[0], cmd.Args[1], out);
		if res != nil {
			WriteJson(encoder, res);
		}
		return;
	case "rm":
		res := Remove(cmd.Id, cmd.Args[0]);
		WriteJson(encoder, res);
		return;

	case "mkdir":
		res := MakeDirectory(cmd.Id, cmd.Args[0]);
		WriteJson(encoder, res);
		return;

	}
	res := FailedResultSet(id, "","","No command specified.");
	WriteJson(encoder, res);
}

func Run(in *os.File, out *os.File) {
	decoder := json.NewDecoder(in);
	var cmd Command;
	var err error = nil;
	for {
		err = decoder.Decode(&cmd);
		if err != nil {
      fmt.Print(err.Error());
			break;
		}
		RunCmd(cmd, out);
	}
	if err == io.EOF {
		WaitForOperationsToComplete();
	}
}
