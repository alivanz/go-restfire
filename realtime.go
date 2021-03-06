package restfire

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	utils "github.com/alivanz/go-utils"
)

const (
	EventPut       = "put"
	EventPatch     = "patch"
	EventKeepAlive = "keep-alive"
	EventCancel    = "cancel"
	EventRevoked   = "revoked"
)

var (
	DataUntouched = errors.New("data untouched")
)

type rtdb_event_data struct {
	Path string           `json:"path"`
	Data *json.RawMessage `json:"data"`
}

func (x *rtdb) watchRequest(path string) (*http.Request, error) {
	uri := x.url + path + ".json"
	if x.refresher != nil {
		if x.refresher.Token() != nil {
			uri = uri + "?" + x.authParam()
		}
	}
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "text/event-stream")
	return request, nil
}
func (x *rtdb) Watch(path string, listener RealtimeDatabaseListener) error {
	client := &http.Client{}
	// client.Timeout = 1 * time.Minute
	request, err := x.watchRequest(path)
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("(watch %s) %v", path, err)
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

func WatchRequest(url string, token TokenInfo) *http.Request {
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	if !strings.HasSuffix(url, ".json") {
		url = url + ".json"
	}
	if token != nil {
		url = url + "?" + tokenAuth(token)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Accept", "text/event-stream")
	return request
}
func tokenAuth(token TokenInfo) string {
	return "auth=" + token.IDToken()
}

func WatchAndListen(url string, token TokenInfo) (<-chan utils.SSEEvent, io.Closer, error) {
	request := WatchRequest(url, token)
	body, err := utils.HttpDo(
		&http.Client{},
		request,
		utils.MustHTTP200,
		utils.MustHTTPContentType("text/event-stream"),
	)
	if err != nil {
		return nil, nil, err
	}
	c := make(chan utils.SSEEvent, 8)
	go ssePump(c, body)
	return c, body, nil
}

func ssePump(c chan<- utils.SSEEvent, body io.ReadCloser) {
	defer close(c)
	defer body.Close()
	sse := utils.NewSSE(body)
	for {
		event, err := sse.ReadEvent()
		if err != nil {
			break
		}
		c <- *event
	}
}

type SSEData struct {
	Path string          `json:"path"`
	Data json.RawMessage `json:"data"`
}

func inList(test string, list []string) bool {
	for _, elem := range list {
		if elem == test {
			return true
		}
	}
	return false
}

func EventParse(event utils.SSEEvent, v interface{}) (string, error) {
	var sseData SSEData
	if inList(event.Event, []string{
		EventPut,
		EventPatch,
	}) {
		if err := json.Unmarshal(event.Data, &sseData); err != nil {
			return "", err
		}
	}
	switch event.Event {
	case EventPut:
		nulldata := JSONExpand(sseData.Path, []byte("null"))
		// log.Print(string(nulldata))
		if err := json.Unmarshal(nulldata, v); err != nil {
			return "", err
		}
		fallthrough
	case EventPatch:
		return sseData.Path, json.Unmarshal(JSONExpand(sseData.Path, sseData.Data), v)
	}
	return "", DataUntouched
}

func HasMutualPrefix(a, b string) bool {
	return strings.HasPrefix(a, b) || strings.HasPrefix(b, a)
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
