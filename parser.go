package main

import (
	"bytes"
	"fmt"
)

func parser(data *bytes.Reader) {
	ch, err := data.ReadByte()
	if err != nil {
		fmt.Println(nil, err)
	}
	switch ch {
	case 'i':
		fmt.Printf("i")
		parser(data)
	case 'l':
		fmt.Printf("l")
		parser(data)
	case 'd':
		fmt.Printf("d")
		parser(data)
	default:
		test1, err1 := data.ReadByte()
		if err1 != nil {
			fmt.Println(nil, err1)
		}
		test2, err2 := data.ReadByte()
		if err2 != nil {
			fmt.Println(nil, err2)
		}
		if test1 == 'e' && test2 == 'e' {
			fmt.Println("end")
		} else {
			fmt.Printf(string(ch))
			parser(data)
		}
	}
}
