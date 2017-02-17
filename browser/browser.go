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
