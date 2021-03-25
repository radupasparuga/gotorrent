package main

import (
	"bytes"
	"fmt"
)

type metainfo struct {
}

func parser(data *bytes.Reader) string {
	ch, err := data.ReadByte()
	if err != nil {
		fmt.Println(nil, err)
	} else {
		switch ch {
		case 'i':
			n, err := data.ReadByte()
			if err != nil {
				return "File not encoded properly"
			} else {
				print(string(n))
			}
			end, endErr := data.ReadByte()
			if endErr != nil {
				return "File not encoded properly"
			} else if string(end) != "e" {
				return "File not encoded properly"
			} else {
				parser(data)
			}
		case 'l':
			fmt.Printf("l")
			parser(data)
		case 'd':
			fmt.Printf("d")
			parser(data)
		default:
			print(string(ch))
			parser(data)
		}
	}

	return "File encoded properly"
}
