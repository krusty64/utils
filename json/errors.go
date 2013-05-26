package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

func filter_harmless(err error, data []byte) error {
	if err == nil {
		return nil
	}

	// This error is relatively minor, usually data that doesn't match
	// the interface, but parsing still finished.
	_, ok := err.(*json.UnmarshalTypeError)
	if ok {
		log.Println("Warning:", err)
		return nil
	}

	syntax, ok := err.(*json.SyntaxError)
	if !ok {
		// must be a real error
		return err
	}

	js := string(data)
	start, end := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
	if idx := strings.Index(js[start:], "\n"); idx >= 0 {
		end = start + idx
	}

	line, pos := strings.Count(js[:start], "\n"), int(syntax.Offset)-start-1

	new_error := fmt.Sprintf("Error in line %d: %s \n", line+1, err) +
		fmt.Sprintf("%s\n%s^", js[start:end], strings.Repeat(" ", pos))

	return errors.New(new_error)
}
