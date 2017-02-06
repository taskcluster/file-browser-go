package main

import (
	"./browser"
	"os"
  // "fmt"
)

func main() {
	browser.Run(os.Stdin, os.Stdout)
  // fmt.Print("Exit");
}
