package tests

import (
	"net/http"
	"testing"
	"time"

	api "github.com/phob0s-pl/perfchat/apiv1"
)

func TestRoomCreateNoUser(t *testing.T) {
	var (
		address  = "localhost:9041"
		server   = NewServer(address)
		client   = api.NewClient(dummyuser, address)
		done     = make(chan bool)
		roomsPre = server.API.Engine.RoomsCount()
	)

	go func() {
		if err := server.Srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server failed, err=%s", err)
		}
		done <- true
	}()
	time.Sleep(time.Millisecond * 10)

	if err := client.RoomCreate(roomDummy.Name); err == nil {
		t.Errorf("Adding room with no auth should fail")
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}

	if l := server.API.Engine.RoomsCount(); l != roomsPre {
		t.Errorf("Expected %d rooms in engine, but got %d", roomsPre, l)
	}

	<-done
}

func TestRoomCreateNoName(t *testing.T) {
	var (
		address  = "localhost:9042"
		server   = NewServer(address)
		client   = api.NewClient(dummyuser, address)
		done     = make(chan bool)
		roomsPre = server.API.Engine.RoomsCount()
	)
	server.API.Engine.AddUser(dummyuser)

	go func() {
		if err := server.Srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server failed, err=%s", err)
		}
		done <- true
	}()
	time.Sleep(time.Millisecond * 10)

	if err := client.RoomCreate(""); err == nil {
		t.Errorf("Adding room with no name should fail")
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}

	if l := server.API.Engine.RoomsCount(); l != roomsPre {
		t.Errorf("Expected %d rooms in engine, but got %d", roomsPre, l)
	}

	<-done
}

func TestRoomCreateDouble(t *testing.T) {
	var (
		address  = "localhost:9043"
		server   = NewServer(address)
		client   = api.NewClient(dummyuser, address)
		done     = make(chan bool)
		roomsPre = server.API.Engine.RoomsCount()
	)
	server.API.Engine.AddUser(dummyuser)

	go func() {
		if err := server.Srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server failed, err=%s", err)
		}
		done <- true
	}()
	time.Sleep(time.Millisecond * 10)

	if err := client.RoomCreate(roomDummy.Name); err != nil {
		t.Errorf("Adding room for a first time failed, err=%s", err)
	}

	if err := client.RoomCreate(roomDummy.Name); err == nil {
		t.Errorf("Adding room for a second should fail")
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}

	if l := server.API.Engine.RoomsCount(); l != roomsPre+1 {
		t.Errorf("Expected %d rooms in engine, but got %d", roomsPre+1, l)
	}

	<-done
}
