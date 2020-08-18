package scoirparser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//UserName in the structure of the output json
type UserName struct {
	First  string `json:"first"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last"`
}

//User in the structure of the output json
type User struct {
	Id    int32    `json:"id"`
	Name  UserName `json:"name"`
	Phone string   `json:"string"`
}

//FileError to record the errors in the file for writing later
type FileError struct {
	Line     int
	ErrorMsg string
}

//ParseCsvToJsonBytes will takes a reader so should be used for any file loading situiation.
func ParseCsvToJsonBytes(data io.Reader) ([]byte, []FileError, error) {
	rows, err := mapColumns(data)
	if err != nil {
		return nil, nil, err
	}
	return makeJson(rows)
}

// INTERNAL_ID : 8 digit positive integer. Cannot be empty.
// FIRST_NAME : 15 character max string. Cannot be empty.
// MIDDLE_NAME : 15 character max string. Can be empty.
// LAST_NAME : 15 character max string. Cannot be empty.
// PHONE_NUM : string that matches this pattern ###-###-####. Cannot be empty.
func getUserId(id string, i int) (int32, *FileError) {
	if len(id) == 0 {
		reqErr := FileError{
			Line:     i,
			ErrorMsg: "user id is required",
		}
		return 0, &reqErr
	}
	if len(id) != 8 {
		lenErr := FileError{
			Line:     i,
			ErrorMsg: fmt.Sprintf("invalid Id length lenght=%d", len(id)),
		}
		return 0, &lenErr
	}
	num, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		errRow := FileError{
			Line:     i,
			ErrorMsg: "error confiming user id is int",
		}

		return 0, &errRow
	}
	return int32(num), nil
}

func nameToLongOrEmpty(key string, name string, index int) *FileError {
	if len(name) == 0 && key != "MIDDLE_NAME" {
		emptyErr := FileError{
			Line:     index,
			ErrorMsg: fmt.Sprintf("%s cannot be empty", key),
		}
		return &emptyErr
	}
	if len(name) > 15 {
		longErr := FileError{
			Line:     index,
			ErrorMsg: fmt.Sprintf("%s:%sis longer then 15 chars", key, name),
		}
		return &longErr
	}
	return nil
}

func validatePhoneNumber(phoneNumber string, index int) *FileError {
	if len(phoneNumber) == 0 {
		emptyErr := FileError{
			Line:     index,
			ErrorMsg: "PHONE_NUM cannot be empty",
		}
		return &emptyErr
	}
	numSlice := strings.Split(phoneNumber, "-")
	valError := FileError{
		Line:     index,
		ErrorMsg: fmt.Sprintf("PHONE_NUM:%s is an invald phone number", phoneNumber),
	}
	if len(numSlice) != 3 {
		return &valError
	}
	for i, n := range numSlice {
		_, err := strconv.Atoi(n)
		if err != nil {
			return &valError
		}
		if i <= 1 {
			if len(n) != 3 {
				return &valError
			}
		} else {
			if len(n) != 4 {
				return &valError
			}
		}
	}
	return nil
}

func getUsersName(key string, name string, index int, user *User) *FileError {
	switch key {
	case "FIRST_NAME":
		err := nameToLongOrEmpty(key, name, index)
		if err != nil {
			return err
		}
		user.Name.First = name
	case "MIDDLE_NAME":
		err := nameToLongOrEmpty(key, name, index)
		if err != nil {
			return err
		}
		user.Name.Middle = name
	case "LAST_NAME":
		err := nameToLongOrEmpty(key, name, index)
		if err != nil {
			return err
		}
		user.Name.Last = name
	}
	return nil
}

func makeJson(rows []map[string]string) ([]byte, []FileError, error) {
	var users []User
	var errors []FileError
	for i, r := range rows {
		i += 2
		user := User{}
		skip := false
		for key, item := range r {
			switch key {
			case "INTERNAL_ID":
				id, err := getUserId(item, i)
				if err != nil {
					errors = append(errors, *err)
					skip = true
					break
				}
				user.Id = id
			case "PHONE_NUM":
				err := validatePhoneNumber(item, i)
				if err != nil {
					errors = append(errors, *err)
					skip = true
					break
				}
				user.Phone = item
			// just fall throught the names to use the default
			case "FIRST_NAME":
				fallthrough
			case "MIDDLE_NAME":
				fallthrough
			case "LAST_NAME":
				fallthrough
			default:
				// we are going to create the rest of the user in the getUserName func
				err := getUsersName(key, item, i, &user)
				if err != nil {
					errors = append(errors, *err)
					skip = true
					break
				}
			}
		}
		if !skip {
			users = append(users, user)
		}
	}
	data, err := json.Marshal(users)
	if err != nil {
		return data, errors, err
	}
	return data, errors, nil
}

func mapColumns(data io.Reader) ([]map[string]string, error) {
	r := csv.NewReader(data)
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows, nil
}
