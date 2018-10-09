package restfire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
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
	client.Timeout = 1 * time.Minute
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("(%s) %v", url, err)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("(%s) %v", url, err)
	}
	defer resp.Body.Close()
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
