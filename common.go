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

var (
	Client = http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 0 * time.Second,
		},
		// Timeout: 1 * time.Minute,
	}
)

func requestdata(method string, url string, data interface{}, out interface{}, errframe error) error {
	var body io.Reader
	if data != nil {
		raw, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(raw)
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("(%s) %v", url, err)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := Client.Do(request)
	if err != nil {
		return fmt.Errorf("(%s) %v", url, err)
	}
	defer resp.Body.Close()
	return finalize(resp, out, errframe)
}
func finalize(resp *http.Response, out interface{}, errframe error) error {
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err = json.Unmarshal(raw, errframe)
		if err != nil {
			return err
		}
		return errframe
	}
	if out == nil {
		return nil
	}
	err = json.Unmarshal(raw, out)
	if err != nil {
		return err
	}
	return nil
}
