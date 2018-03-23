package tests

import (
	"net/http"
	"testing"
	"time"

	api "github.com/phob0s-pl/perfchat/apiv1"
)

func TestPing(t *testing.T) {
	var (
		address = "localhost:9093"
		server  = NewServer(address)
		client  = api.NewClient(admin, address)
		done    = make(chan bool)
	)

	go func() {
		if err := server.Srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server failed, err=%s", err)
		}
		done <- true
	}()
	time.Sleep(time.Millisecond * 10)

	if err := client.Ping(); err != nil {
		t.Errorf("Sending PING failed, err=%s", err)
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}

	<-done
}
