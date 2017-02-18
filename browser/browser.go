package browser

import (
	"gopkg.in/vmihailenco/msgpack.v2"
	"os"
)

type Command struct {
	Id   string   `msgpack:"id"`
	Cmd  string   `msgpack:"cmd"`
	Args []string `msgpack:"args"`
	Data []byte   `msgpack:"data,omitempty"`
}

func Run(in *os.File, out *os.File) {
	// decoder := msgpack.NewDecoder(in)
	// var err error = nil
	defer func() {
		in.Close()
		out.Close()
	}()
	outChan := make(chan *ResultSet)
	// inChan := make(chan Command)

	go EncodeCompressWrite(outChan, out)
	// go DecompressDecode(inChan, in)
	for {
		var cmd Command
		msgpack.NewDecoder(in).Decode(&cmd)
		go RunCommand(cmd, outChan)
	}

}
