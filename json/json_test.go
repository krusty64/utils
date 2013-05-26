package json

import (
	"strconv"
	"strings"
	"testing"
)

const input1 = `
{
  "a": "a",
  "b": "b",
  "array":["abc", "123", "asdf"],
	"struct":{
    "a": "a",
    "b": "b"
	},
  "n": 100
}
`
const syntaxError1 = `
{
	"a": "a",
	"b": "b",
}
`

const syntaxError2 = `
{
	"a": "a"
	"b": "b"
}
`

type Input1 struct {
	A      string
	B      string
	Array  []string
	Struct struct {
		A string
		B string
	}
	N int
}

func testUnmarshal(input string, i interface{}, t *testing.T) {
	err := unmarshal([]byte(input), i)
	if err != nil {
		t.Error(i, err)
	}
}

func testSyntaxError(input string, errorline int, i interface{}, t *testing.T) {
	err := unmarshal([]byte(input), i)
	if err == nil ||
		!strings.Contains(err.Error(), "Error in line "+strconv.Itoa(errorline)) {
		t.Error(i, err)
	}
}

func TestParsing(t *testing.T) {
	i := Input1{}
	testUnmarshal(input1, &i, t)
	testSyntaxError(syntaxError1, 5, &i, t)
	testSyntaxError(syntaxError2, 4, &i, t)
}
