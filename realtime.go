package restfire

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type rtdb_event_data struct {
	Path string           `json:"path"`
	Data *json.RawMessage `json:"data"`
}

func (x *rtdb) Watch(path string, listener RealtimeDatabaseListener) error {
	client := &http.Client{}
	uri := x.url + path + ".json?" + x.authParam()
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "text/event-stream")
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("return code %d", resp.StatusCode)
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
