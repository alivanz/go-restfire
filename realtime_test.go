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
		time.Sleep(40 * time.Second)
		closer.Close()
	}()
	var data interface{}
	for event := range c {
		EventParse(event, &data)
		if err == DataUntouched {
			log.Print(event)
			continue
		} else if err != nil {
			t.Log(err)
			t.Fail()
		}
		log.Print(data)
	}
}
