package scoirparser

import (
	"encoding/json"
	"strings"
	"testing"
)

const good = `INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM
12345678,Bobby,,Tables,555-555-5555`

const badNames = `INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM
12345678,Bobby,,,555-555-5555
12345678,,,Fish,555-555-5555`

const badNumber = `INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM
12345678,Bobby,,Tables,555-55-5555`

const skipOne = `INTERNAL_ID,FIRST_NAME,MIDDLE_NAME,LAST_NAME,PHONE_NUM
12345678,Bobby,,Tables,555-555-5555
123,Bobby,,Tables,555-555-5555`

func TestParserGood(t *testing.T) {
	reader := strings.NewReader(good)
	var userData []User
	b, elist, err := ParseCsvToJsonBytes(reader)
	if err != nil {
		t.Errorf("should not have faild got an error: %s", err)
	}
	if len(elist) != 0 {
		t.Errorf("should have no errors got: %d", len(elist))
	}
	json.Unmarshal(b, &userData)
	if len(userData) != 1 {
		t.Errorf("userData should be ony one in length: %d", len(userData))
	}
}
func TestParserBadNames(t *testing.T) {
	reader := strings.NewReader(badNames)
	var userData []User
	b, elist, err := ParseCsvToJsonBytes(reader)
	if err != nil {
		t.Errorf("should not have faild got an error: %s", err)
	}
	if len(elist) != 2 {
		t.Errorf("should have one errors got: %d", len(elist))
	}
	json.Unmarshal(b, &userData)
	if len(userData) != 0 {
		t.Errorf("userData should be empty in length: %d", len(userData))
	}
}
func TestParserBadNumber(t *testing.T) {
	reader := strings.NewReader(badNumber)
	var userData []User
	b, elist, err := ParseCsvToJsonBytes(reader)
	if err != nil {
		t.Errorf("should not have faild got an error: %s", err)
	}
	if len(elist) != 1 {
		t.Errorf("should have one errors got: %d", len(elist))
	}
	json.Unmarshal(b, &userData)
	if len(userData) != 0 {
		t.Errorf("userData should be empty in length: %d", len(userData))
	}
}

func TestParserSkipOne(t *testing.T) {
	reader := strings.NewReader(skipOne)
	var userData []User
	b, elist, err := ParseCsvToJsonBytes(reader)
	if err != nil {
		t.Errorf("should not have faild got an error: %s", err)
	}
	if len(elist) != 1 {
		t.Errorf("should have one errors got: %d", len(elist))
	}
	json.Unmarshal(b, &userData)
	if len(userData) != 1 {
		t.Errorf("userData should be one in length: %d", len(userData))
	}
}
