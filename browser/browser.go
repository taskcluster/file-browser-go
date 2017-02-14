package browser

import (
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
	"io"
	"os"
)

type Command struct {
	Id   string   `msgpack:"id"`
	Cmd  string   `msgpack:"cmd"`
	Args []string `msgpack:"args"`
	Data []byte   `msgpack:"data"`
}

func RunCmd(cmd Command, out io.Writer) {
	encoder := msgpack.NewEncoder(out)
	if cmd.Id == "" {
		encoder.Encode(FailedResultSet("", "No Id supplied."))
		return
	}
	if len(cmd.Args) == 0 {
		res := FailedResultSet(cmd.Id, "Not enough arguments.")
		WriteOut(encoder, res)
		return
	}

	id := cmd.Id

	switch cmd.Cmd {

	case "ls":
		res := List(cmd.Id, cmd.Args[0])
		WriteOut(encoder, res)
		return

	case "getfile":
		GetFile(cmd.Id, cmd.Args[0], out)
		return

	case "putfile":
		res := PutFile(cmd.Id, cmd.Args[0], cmd.Data)
		WriteOut(encoder, res)
		return

	case "mv":
		if len(cmd.Args) < 1 {
			res := FailedResultSet(id, "Not enough arguments.")
			WriteOut(encoder, res)
			return
		}
		res := Move(cmd.Id, cmd.Args[0], cmd.Args[1])
		WriteOut(encoder, res)
		return
	case "cp":
		if len(cmd.Args) < 1 {
			res := FailedResultSet(id, "Not enough arguments.")
			WriteOut(encoder, res)
			return
		}
		res := Copy(cmd.Id, cmd.Args[0], cmd.Args[1], out)
		if res != nil {
			WriteOut(encoder, res)
		}
		return
	case "rm":
		res := Remove(cmd.Id, cmd.Args[0])
		WriteOut(encoder, res)
		return

	case "mkdir":
		res := MakeDirectory(cmd.Id, cmd.Args[0])
		WriteOut(encoder, res)
		return

	}
	res := FailedResultSet(id, "No command specified.")
	WriteOut(encoder, res)
}

func Run(in *os.File, out *os.File) {
	decoder := msgpack.NewDecoder(in)
	var cmd Command
	var err error = nil
	for {
		err = decoder.Decode(&cmd)
		if err != nil {
			fmt.Print(err.Error())
			break
		}
		RunCmd(cmd, out)
	}
	if err == io.EOF {
		return
	}
}
