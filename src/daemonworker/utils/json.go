package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func UnmarshalToJSON(in io.ReadCloser, target interface{}) error {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &target)
}

func PrettyPrintJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err.Error()
	}
	return string(data)
}
