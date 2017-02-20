package main

import (
	"./browser"
	"os"
)

func main() {
	browser.Run(os.Stdin, os.Stdout)
}
