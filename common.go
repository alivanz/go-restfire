package restfire

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type firebaseAuth struct {
	key string
}

func requestdata(method string, url string, data interface{}, out interface{}, errframe error) error {
	var body io.Reader
	if data != nil {
		raw, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(raw)
	}
	client := &http.Client{}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	return finalize(resp, out, errframe)
}
func finalize(resp *http.Response, out interface{}, errframe error) error {
	if resp.StatusCode != 200 {
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw, errframe)
		if err != nil {
			return err
		}
		return errframe
	}
	if out == nil {
		return nil
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, out)
	if err != nil {
		return err
	}
	return nil
}
