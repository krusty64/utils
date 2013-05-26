package json

import (
	"io/ioutil"
	"log"
	"net/http"
)

func FromUrl(url string, response interface{}) error {
	log.Println(url)
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}

	return unmarshal(data, response)
}
