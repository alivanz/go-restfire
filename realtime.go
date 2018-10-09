package restfire

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type rtdb_event_data struct {
	Path string           `json:"path"`
	Data *json.RawMessage `json:"data"`
}

func (x *rtdb) watchRequest(path string) (*http.Request, error) {
	uri := x.url + path + ".json?" + x.authParam()
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "text/event-stream")
	return request, nil
}
func (x *rtdb) Watch(path string, listener RealtimeDatabaseListener) error {
	client := &http.Client{}
	client.Timeout = 1 * time.Minute
	request, err := x.watchRequest(path)
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("(watch %s) %v", path, err)
	}
	if resp.StatusCode != 401 {
		resp.Body.Close()
		// Check if refresher available
		if x.refresher == nil {
			return fmt.Errorf("return code %d", resp.StatusCode)
		}
		// Refresh auth
		if err = x.refresher.AuthRefresh(); err != nil {
			return err
		}
		// retry
		request, err = x.watchRequest(path)
		resp, err = client.Do(request)
		if err != nil {
			return fmt.Errorf("(watch %s) %v", path, err)
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("(watch %s) return code %d", path, resp.StatusCode)
	}
	return SSEHandle(resp.Body, func(event string, data json.RawMessage) error {
		var parsed rtdb_event_data
		switch event {
		case "put":
			fallthrough
		case "patch":
			err := json.Unmarshal(data, &parsed)
			if err != nil {
				return err
			}
		}
		switch event {
		case "put":
			if parsed.Data == nil {
				listener.OnDelete(parsed.Path)
			} else {
				listener.OnPut(parsed.Path, *parsed.Data)
			}
		case "patch":
			listener.OnPatch(parsed.Path, *parsed.Data)
		case "keep-alive":
			return nil
		case "cancel":
			listener.OnCancel()
		case "auth_revoked":
			listener.OnAuthRevoked()
		}
		return nil
	})
}

type SilentRealtimeDatabaseListener struct{}

func (SilentRealtimeDatabaseListener) OnPut(path string, data json.RawMessage)   {}
func (SilentRealtimeDatabaseListener) OnPatch(path string, data json.RawMessage) {}
func (SilentRealtimeDatabaseListener) OnDelete(path string)                      {}
func (SilentRealtimeDatabaseListener) OnCancel()                                 {}
func (SilentRealtimeDatabaseListener) OnAuthRevoked()                            {}

type LogRealtimeDatabaseListener struct{}

func (LogRealtimeDatabaseListener) OnPut(path string, data json.RawMessage) {
	log.Printf("put: %s %v", path, string(data))
}
func (LogRealtimeDatabaseListener) OnPatch(path string, data json.RawMessage) {
	log.Printf("patch: %s %v", path, string(data))
}
func (LogRealtimeDatabaseListener) OnDelete(path string) {
	log.Printf("delete: %s", path)
}
func (LogRealtimeDatabaseListener) OnCancel() {
	log.Printf("cancel")
}
func (LogRealtimeDatabaseListener) OnAuthRevoked() {
	log.Printf("revoked")
}
