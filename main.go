package main;

import (
	"fmt";
	"os";
	"bufio";
	"encoding/json";

	"github.com/taskcluster/file-browser-go/browser";
	"gopkg.in/vmihailenco/msgpack.v2";
)

func main(){
	reader := bufio.NewReader(os.Stdin);
	decoder := json.NewDecoder(reader);
	encoder := msgpack.NewEncoder(os.Stdout);
	var cmd browser.Command;
	var err error;
	for cmd.Cmd != "exit" && err == nil {
		err = decoder.Decode(&cmd);
		if err != nil {
			fmt.Println(err.Error());
		}else{
			res := browser.Run(cmd);
			if res == nil {
				break;
			}
			err = encoder.Encode(res);
		}
	}
}
