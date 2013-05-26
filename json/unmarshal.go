package json

import (
	"encoding/json"
)

func unmarshal(data []byte, response interface{}) error {
	err := json.Unmarshal(data, response)
	if err != nil {
		return filter_harmless(err, data)
	}

	return nil
}
