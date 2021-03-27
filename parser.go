package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

func decode(data *bytes.Reader) (map[string]interface{}, error) {
	metadata := map[string]interface{}{
		"key": "value",
	}

	parser(data, metadata)

	return metadata, nil
}

func parser(data *bytes.Reader, metadata map[string]interface{}) (string, error) { // todo handle multiple return data types
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
				parser(data, metadata)
			}
		case 'l':
			array, _ := parseList(data)
			fmt.Println(array)
			parser(data, metadata)
		case 'd':
			dict, _ := parseDict(data)
			parser(data, metadata)
		default:
			data.UnreadByte()
			str, strErr := parseString(data)
			if strErr != nil {
				return "", strErr
			}
			fmt.Println(str)
			parser(data, metadata)
		}
	}

	return "File encoded properly", nil
}

func parseList(data *bytes.Reader) ([2]interface{}, error) {
	var list [2]interface{} // todo figure out how to handle array size
	var index int = 0
	for {
		check, checkErr := data.ReadByte()
		if checkErr != nil {
			return list, checkErr
		} else if check == 'e' {
			break
		} else if check == 'l' {
			list[index], _ = parseList(data)
			index = index + 1
		} else {
			data.UnreadByte()
			str, strErr := parseString(data)
			if strErr != nil {
				return list, strErr
			}
			list[index] = str
			index = index + 1
		}
	}
	return list, nil
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
