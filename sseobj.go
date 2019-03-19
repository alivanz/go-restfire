package restfire

import (
	"io"

	utils "github.com/alivanz/go-utils"
)

type FirebaseWatcher struct {
	c      <-chan utils.SSEEvent
	closer io.Closer
}

func (f *FirebaseWatcher) Close() error {
	f.c = nil
	if f.closer == nil {
		return nil
	}
	return f.closer.Close()
}
func (f *FirebaseWatcher) CEvent() <-chan utils.SSEEvent {
	return f.c
}
func (f *FirebaseWatcher) HandleOpen(url string) error {
	if f.c == nil {
		sse, closer, err := WatchAndListen(url, nil)
		if err != nil {
			return err
		}
		f.c = sse
		f.closer = closer
	}
	return nil
}
func (f *FirebaseWatcher) CEventHandle(event utils.SSEEvent, ok bool, v interface{}, onerror func(error), onmodify func()) {
	if !ok {
		f.c = nil
	}
	modified, err := DecodeSSE(event, v)
	if err != nil {
		f.Close()
		onerror(err)
		return
	}
	if modified {
		onmodify()
	}
}

// DecodeSSE decode SSE. Returns true if data modified.
func DecodeSSE(event utils.SSEEvent, v interface{}) (bool, error) {
	// logFirebase.Print(event)
	_, err := EventParse(event, v)
	if err == DataUntouched {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
