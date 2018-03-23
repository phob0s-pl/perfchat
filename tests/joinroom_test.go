package tests

import (
	"net/http"
	"testing"
	"time"

	api "github.com/phob0s-pl/perfchat/apiv1"
)

func TestRoomJoinNoUser(t *testing.T) {
	var (
		address  = "localhost:9033"
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

	if err := client.RoomJoin(roomDummy.Name); err == nil {
		t.Errorf("Joining room with no auth should fail")
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}

	if l := server.API.Engine.RoomsCount(); l != roomsPre {
		t.Errorf("Expected %d rooms in engine, but got %d", roomsPre, l)
	}

	<-done
}

func TestRoomJoinNoRoom(t *testing.T) {
	var (
		address  = "localhost:9033"
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

	if err := client.RoomJoin(roomDummy.Name); err == nil {
		t.Errorf("Joining nonexist room should fail")
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}

	if l := server.API.Engine.RoomsCount(); l != roomsPre {
		t.Errorf("Expected %d rooms in engine, but got %d", roomsPre, l)
	}

	<-done
}

func TestRoomJoin(t *testing.T) {
	var (
		address  = "localhost:9034"
		server   = NewServer(address)
		clientA  = api.NewClient(useralpha, address)
		clientB  = api.NewClient(userbeta, address)
		done     = make(chan bool)
		roomsPre = server.API.Engine.RoomsCount()
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

	if err := clientA.RoomCreate(roomAlpha.Name); err != nil {
		t.Errorf("RoomCreate(%s) failed, err=%s", roomAlpha.Name, err)
	}

	if err := clientB.RoomJoin(roomAlpha.Name); err != nil {
		t.Errorf("RoomJoin(%s) failed, err=%s", roomAlpha.Name, err)
	}

	if err := clientB.RoomJoin(roomAlpha.Name); err == nil {
		t.Errorf("RoomJoin(%s) for second time should fail", roomAlpha.Name)
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}

	if l := server.API.Engine.RoomsCount(); l != roomsPre+1 {
		t.Errorf("Expected %d rooms in engine, but got %d", roomsPre+1, l)
	}

	room, err := server.API.Engine.GetRoomByName(roomAlpha.Name)
	if err != nil {
		t.Errorf("GetRoomByName(%s)  failed, err=%s", roomAlpha.Name, err)
	}

	if lu := len(room.Users); lu != 2 {
		t.Errorf("Expected %d users in room, but got %d", 2, lu)
	}

	<-done
}
