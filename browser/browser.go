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
/*
func RunCmd(cmd Command, outChan chan interface{}) {
	if cmd.Id == "" {
		outChan <- FailedResultSet("", "No Id supplied.")
		return
	}
	if len(cmd.Args) == 0 {
		res := FailedResultSet(cmd.Id, "Not enough arguments.")
		outChan <- res
		return
	}

	id := cmd.Id

	switch cmd.Cmd {

	case "ls":
		List(cmd.Id, outChan, cmd.Args[0])
		return

	case "getfile":
		GetFile(cmd.Id, outChan, cmd.Args[0])
		return

	case "putfile":
		PutFile(cmd.Id, outChan, cmd.Args[0], cmd.Data)
		return

	case "mv":
		if len(cmd.Args) < 1 {
			outChan <- FailedResultSet(id, "Not enough arguments.")
			return
		}
		Move(cmd.Id, outChan, cmd.Args[0], cmd.Args[1])
		return
	case "cp":
		if len(cmd.Args) < 1 {
			outChan <- FailedResultSet(id, "Not enough arguments.")
			return
		}
		Copy(cmd.Id, outChan, cmd.Args[0], cmd.Args[1])
		return
	case "rm":
		Remove(cmd.Id, outChan, cmd.Args[0])
		return

	case "mkdir":
		MakeDirectory(cmd.Id, outChan, cmd.Args[0])
		return

	}
	res := FailedResultSet(id, "No command specified.")
	outChan <- res
}
*/

func Run(in *os.File, out *os.File) {
	decoder := msgpack.NewDecoder(in)
	encoder := msgpack.NewEncoder(out)
	var cmd Command
	var err error = nil
	InitRegistry()
	outChan := make(chan interface{})
	go func() {
		for {
			encoder.Encode(<-outChan)
		}
	}()
	for {
		err = decoder.Decode(&cmd)
		if err != nil {
			fmt.Print(err.Error())
			break
		}
		go RunCommand(cmd.Cmd)(cmd, outChan)
	}
	if err == io.EOF {
		return
	}
}
