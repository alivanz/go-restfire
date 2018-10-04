package restfire

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type SSEListener func(string, json.RawMessage) error

func SSEHandle(reader io.Reader, listener SSEListener) error {
	var event string
	r := bufio.NewReader(reader)
	sep := []byte(": ")
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return err
		}
		line = line[:len(line)-1]
		spl := bytes.SplitN(line, sep, 2)
		if len(spl) == 1 {
			continue
		}
		if len(spl) != 2 {
			return fmt.Errorf("SSE format error")
		}
		switch string(spl[0]) {
		case "event":
			event = string(spl[1])
		case "data":
			err = listener(event, spl[1])
			if err != nil {
				return err
			}
		}
	}
}
