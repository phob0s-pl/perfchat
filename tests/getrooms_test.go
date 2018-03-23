package tests

import (
	"net/http"
	"testing"
	"time"

	api "github.com/phob0s-pl/perfchat/apiv1"
)

func TestGetRoomsByNotAddedUser(t *testing.T) {
	var (
		address = "localhost:9011"
		server  = NewServer(address)
		client  = api.NewClient(dummyuser, address)
		done    = make(chan bool)
	)
	server.API.Engine.AddUser(useralpha)
	server.API.Engine.AddUser(userbeta)

	go func() {
		if err := server.Srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server failed, err=%s", err)
		}
		done <- true
	}()
	time.Sleep(time.Millisecond * 10)

	if _, err := client.GetRooms(); err == nil {
		t.Errorf("Getting rooms with user not in engine should fail")
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}
	<-done
}

func TestGetRoomsSimple(t *testing.T) {
	var (
		address = "localhost:9012"
		server  = NewServer(address)
		client  = api.NewClient(admin, address)
		done    = make(chan bool)
	)
	server.API.Engine.AddUser(useralpha)
	server.API.Engine.AddUser(userbeta)
	server.API.Engine.AddUser(admin)

	go func() {
		if err := server.Srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server failed, err=%s", err)
		}
		done <- true
	}()
	time.Sleep(time.Millisecond * 10)

	rooms, err := client.GetRooms()
	if err != nil {
		t.Errorf("Getting rooms failed, err=%s", err)
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}
	<-done

	engineRooms := server.API.Engine.ListRooms()
	if len(rooms) != len(engineRooms) {
		t.Fatalf("Expected %d rooms, got %d", len(engineRooms), len(rooms))
	}
}
