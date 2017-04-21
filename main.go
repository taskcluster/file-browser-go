package main

import (
	"github.com/taskcluster/file-browser-go/browser"
	"os"
)

func main() {
	browser.Run(os.Stdin, os.Stdout)
}
