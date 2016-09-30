package main;

import (
	"os";
	"bufio";
	"encoding/json";
	// "errors";
	"github.com/taskcluster/file-browser-go/browser";
	//"gopkg.in/vmihailenco/msgpack.v2";
)

func main(){
	reader := bufio.NewReader(os.Stdin);
	decoder := json.NewDecoder(reader);
	encoder := json.NewEncoder(os.Stdout);
	out := make(chan interface{});
	var cmd browser.Command;
	var err error = nil;

	go func(){
		for{
			res := <-out;
			if res.(*browser.ResultSet).Cmd == "Exit" {
				os.Exit(0);
				break;
			}
			err = encoder.Encode(res);
		}
	}();

	for err == nil || err.Error() == "EOF"{
		err = decoder.Decode(&cmd);
		browser.Run(cmd, out);
	}
	close(out);
}
