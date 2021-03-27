package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

func decode(data *bytes.Reader) (map[string]interface{}, error) {
	metadata := map[string]interface{}{}

	parsed, parsedErr := parser(data, metadata)
	if parsedErr != nil {
		fmt.Println(parsedErr)
	} else {
		fmt.Println(parsed)
	}
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
			dict, dictErr := parseDict(data)
			if dictErr != nil {
				return "", dictErr
			}
			fmt.Println(dict)
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

func parseString(data *bytes.Reader) (string, error) { // todo fix case where length is >= 10
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

func parseList(data *bytes.Reader) ([5]interface{}, error) {
	var list [5]interface{} // todo figure out how to handle array size
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
		} else if check == 'd' {
			parsedDictionary, parsedDictionaryErr := parseDict(data)
			if parsedDictionaryErr != nil {
				return list, parsedDictionaryErr
			}
			list[index] = parsedDictionary
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

func parseDict(data *bytes.Reader) (map[string]interface{}, error) {
	dictionary := map[string]interface{}{}
	var bounce int = 1
	var key string = ""
	var value [1]interface{} = [1]interface{}{""}
	for {
		check, checkErr := data.ReadByte()
		if checkErr != nil {
			return dictionary, checkErr
		} else if check == 'e' {
			break
		} else if check == 'l' {
			parsedList, parsedListErr := parseList(data)
			if parsedListErr != nil {
				return dictionary, parsedListErr
			}
			if bounce == 1 {
				return dictionary, errors.New("List can't be key in dictionary")
			} else {
				value[0] = parsedList
				dictionary[key] = value[0]
				bounce = 1
			}
		} else if check == 'd' {
			parsedDictionary, parsedDictionaryErr := parseDict(data)
			if parsedDictionaryErr != nil {
				return dictionary, parsedDictionaryErr
			}
			if bounce == 1 {
				return dictionary, errors.New("Dictionaries can't be key in dictionary")
			} else {
				value[0] = parsedDictionary
				dictionary[key] = value[0]
				bounce = 1
			}
		} else {
			data.UnreadByte()
			str, strErr := parseString(data)
			if strErr != nil {
				return dictionary, strErr
			}
			if bounce == 1 {
				key = str
				bounce = 0
			} else {
				value[0] = str
				dictionary[key] = value[0]
				bounce = 1
			}
		}
	}

	return dictionary, nil
}