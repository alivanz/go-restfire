package restfire

import (
	"log"
	"testing"
	"time"
)

func TestRealtime(t *testing.T) {
	c, closer, err := WatchAndListen("https://monadminer.firebaseio.com/test", nil)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	go func() {
		time.Sleep(10 * time.Second)
		closer.Close()
	}()
	for event := range c {
		log.Print(event)
	}
}
