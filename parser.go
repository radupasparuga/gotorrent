package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type metainfo struct {
}

func parser(data *bytes.Reader) (string, error) {
	ch, err := data.ReadByte()
	if err != nil {
		return "", err
	} else {
		switch ch {
		case 'i': // todo: negative integers
			n, err := data.ReadByte()
			if err != nil {
				return "", err
			} else {
				print(string(n))
			}
			end, endErr := data.ReadByte()
			if endErr != nil {
				return "", endErr
			} else if string(end) != "e" {
				return "", errors.New("file not encoded properly")
			} else {
				parser(data)
			}
		case 'l':
			var list [2]string // todo figure out how to handle array size
			var index int = 0
			for {
				check, checkErr := data.ReadByte()
				if checkErr != nil {
					return "", checkErr
				} else if check == 'e' { // todo handle lists/dictionaries/etc inside a list
					break
				}
				data.UnreadByte()
				str, strErr := parseString(data)
				if strErr != nil {
					return "", strErr
				}
				list[index] = str
				index = index + 1
			}
			fmt.Println(list)
			parser(data)
		case 'd':
			fmt.Printf("d")
			parser(data)
		default:
			data.UnreadByte()
			str, strErr := parseString(data)
			if strErr != nil {
				return "", strErr
			}
			fmt.Println(str)
			parser(data)
		}
	}

	return "File encoded properly", nil
}

func parseString(data *bytes.Reader) (string, error) {
	len, lenErr := data.ReadByte()
	if lenErr != nil {
		return "", errors.New("string doesn't begin with string length")
	}
	colon, colonErr := data.ReadByte()
	if colonErr != nil || string(colon) != ":" {
		return "", errors.New("string not formatted properly")
	}

	intLen, _ := strconv.Atoi(string(len))
	buf := make([]byte, intLen)
	if _, sErr := io.ReadFull(data, buf); sErr != nil {
		return "", sErr
	}
	return string(buf), nil
}
