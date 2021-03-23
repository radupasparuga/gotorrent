package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	pathPtr := flag.String("path", "", "Path to .torrent file")
	flag.Parse()

	data, err := ioutil.ReadFile(*pathPtr)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	parser(bytes.NewReader((data)))
}
