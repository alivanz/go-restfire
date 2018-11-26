package restfire

import (
	"log"
	"testing"
)

func TestExpand(t *testing.T) {
	log.Print(string(JSONExpand("/abc/def", []byte("123"))))
}
