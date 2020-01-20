package graylogger

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type graylogHelpersSuite struct {
	suite.Suite
}

func (s graylogHelpersSuite) TestCleanString() {
	// 1. Removing whitespaces from a string
	text := `
		a,
		b,
		c
	`

	str := cleanString(text)
	s.Equal("a, b, c", str)

	// 2. Removing indentation from JSON
	// {
	// 	"A": "a",
	// 	"B": "b",
	// 	"C": "c"
	// }
	type JS struct {
		A string
		B string
		C string
	}

	js := JS{
		A: "a",
		B: "b",
		C: "c",
	}

	obj, err := json.MarshalIndent(js, "", "  ")
	s.Equal(nil, err)

	str = cleanString(string(obj))
	s.Equal(`{"A":"a","B":"b","C":"c"}`, str)
}

func (s graylogHelpersSuite) TestPrettifyObject() {
	// 1. From a struct ...
	type S struct {
		A string
		B string
		C string
	}

	structure := S{
		A: "a",
		B: "b",
		C: "c",
	}

	pretty := cleanString(prettifyObject(structure))
	s.Equal(`{"A":"a","B":"b","C":"c"}`, pretty)

	// 2. From a slice ...
	slice := []interface{}{1, "2", true}
	pretty = cleanString(prettifyObject(slice))
	s.Equal(`[1,"2",true]`, pretty)

	// 3. From a map ...
	stringMap := map[string]string{
		"A": "a",
		"B": "b",
		"C": "c",
	}

	pretty = cleanString(prettifyObject(stringMap))
	s.Equal(`{"A":"a","B":"b","C":"c"}`, pretty)

	// 4. From a string ...
	str := `
		{
		  "A":"a"
		}
	`

	pretty = cleanString(prettifyObject(str))
	s.Equal(`{"A":"a"}`, pretty)

	// 4. From a string ...
	var any interface{} = 101

	pretty = cleanString(prettifyObject(any))
	s.Equal("101", pretty)
}

func TestGraylogHelpersSuite(t *testing.T) {
	suite.Run(t, new(graylogHelpersSuite))
}
