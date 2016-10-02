package main;

import (
	"os";
	"github.com/taskcluster/file-browser-go/browser";
)

func main(){
	browser.Run(os.Stdin, os.Stdout);
}
