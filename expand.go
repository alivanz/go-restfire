package restfire

import (
	"encoding/json"
	"strings"
)

func JSONExpand(path string, data []byte) []byte {
	path = strings.TrimPrefix(path, "/")
	return jsonExpand(path, data)
}
func jsonExpand(path string, data []byte) []byte {
	if len(path) == 0 {
		return data
	}
	subs := strings.Split(path, "/")
	return jsonExpandPaths(subs, json.RawMessage(data))
}
func jsonExpandPaths(paths []string, data []byte) []byte {
	if len(paths) > 0 {
		return jsonExpandPathOnce(paths[0], jsonExpandPaths(paths[1:], data))
	}
	return data
}
func jsonExpandPathOnce(path string, data []byte) []byte {
	b, err := json.Marshal(map[string]json.RawMessage{
		path: data,
	})
	if err != nil {
		panic(err)
	}
	return b
}
