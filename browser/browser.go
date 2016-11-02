package browser;

import (
	"os";
	"encoding/json";
	"io";
)

type Command struct {
	Cmd string `json:"cmd"`
	Args []string `json:"args"`
	Data []byte `json:"data"`
}

func RunCmd (cmd Command, out io.Writer) {
	encoder := json.NewEncoder(out);
	if len(cmd.Args) == 0 {
		res := FailedResultSet(cmd.Cmd,"", "Not enough arguments.");
		encoder.Encode(res);
		return;
	}
	switch cmd.Cmd{

	case "ls":
		res := List(cmd.Args[0]);
		encoder.Encode(res);
		return;

	case "getfile":
		GetFile(cmd.Args[0], out);
		return;

	case "putfile":
		res := PutFile2(cmd.Args[0], cmd.Data);
		encoder.Encode(res);
		return;

	case "mv":
		if len(cmd.Args) < 1 {
			res := FailedResultSet(cmd.Cmd,"","Not enough arguments.");
			encoder.Encode(res);
			return;
		}
		res := Move(cmd.Args[0], cmd.Args[1]);
		encoder.Encode(res);
		return;
	case "cp":
		if len(cmd.Args) < 1 {
			res := FailedResultSet(cmd.Cmd,"","Not enough arguments.");
			encoder.Encode(res);
			return;
		}
		res := Copy(cmd.Args[0], cmd.Args[1], out);
		if res != nil {
			encoder.Encode(res);
		}
		return;
	case "rm":
		res := Remove(cmd.Args[0]);
		encoder.Encode(res);
		return;

	case "mkdir":
		res := MakeDirectory(cmd.Args[0]);
		encoder.Encode(res);
		return;

	}
	res := FailedResultSet("","","No command specified.");
	encoder.Encode(res);
}

func Run(in *os.File, out *os.File) {
	decoder := json.NewDecoder(in);
	var cmd Command;
	var err error = nil;
	for {
		err = decoder.Decode(&cmd);
		if err != nil {
			break;
		}
		RunCmd(cmd, out);
	}
	if err == io.EOF {
		WaitForOperationsToComplete();
	}
}
