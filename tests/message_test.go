package tests

import (
	"net/http"
	"testing"
	"time"

	api "github.com/phob0s-pl/perfchat/apiv1"
)

func TestMessage(t *testing.T) {
	var (
		address     = "localhost:9064"
		server      = NewServer(address)
		clientAdmin = api.NewClient(admin, address)
		clientA     = api.NewClient(useralpha, address)
		clientB     = api.NewClient(userbeta, address)
		clientD     = api.NewClient(dummyuser, address)
		done        = make(chan bool)
	)
	server.API.Engine.AddUser(admin)

	go func() {
		if err := server.Srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server failed, err=%s", err)
		}
		done <- true
	}()
	time.Sleep(time.Millisecond * 10)

	if err := clientAdmin.AddUser(useralpha); err != nil {
		t.Errorf("Adduser(%s) failed, err=%s", useralpha.Name, err)
	}

	if err := clientAdmin.AddUser(userbeta); err != nil {
		t.Errorf("Adduser(%s) failed, err=%s", userbeta.Name, err)
	}

	if err := clientAdmin.AddUser(dummyuser); err != nil {
		t.Errorf("Adduser(%s) failed, err=%s", dummyuser.Name, err)
	}

	if err := clientA.RoomCreate(roomAlpha.Name); err != nil {
		t.Errorf("RoomCreate(%s) failed, err=%s", roomAlpha.Name, err)
	}

	if err := clientB.RoomJoin(roomAlpha.Name); err != nil {
		t.Errorf("RoomJoin(%s) failed, err=%s", roomAlpha.Name, err)
	}

	if err := clientB.SendMessage(&api.Message{Content: "hello", Room: roomAlpha.Name}); err != nil {
		t.Errorf("SendMessage(%s) failed, err=%s", roomAlpha.Name, err)
	}
	if err := clientB.SendMessage(&api.Message{Content: "hello", Room: roomAlpha.Name}); err != nil {
		t.Errorf("SendMessage(%s) failed, err=%s", roomAlpha.Name, err)
	}
	time.Sleep(time.Millisecond * 10)

	msgs, err := clientA.ReceiveMessage()
	if err != nil {
		t.Errorf("ReceiveMessage(%s) failed, err=%s", roomAlpha.Name, err)
	}

	dmsg, err := clientD.ReceiveMessage()
	if err != nil {
		t.Errorf("ReceiveMessage(%s) failed, err=%s", roomAlpha.Name, err)
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}
	<-done

	if len(dmsg) != 0 {
		t.Errorf("User not in room shouldnt receive messages, but got %d", len(dmsg))
	}

	if len(msgs) != 2 {
		t.Fatalf("Expected one message, got %d", len(msgs))
	}

	if msgs[0].User != userbeta.Name {
		t.Errorf("msgs.user, got %s, expected %s", msgs[0].User, userbeta.Name)
	}

	if msgs[0].Room != roomAlpha.Name {
		t.Errorf("msgs.room, got %s, expected %s", msgs[0].Room, roomAlpha.Name)
	}

	if msgs[0].Content != "hello" {
		t.Errorf("msgs.content, got %s, expected %s", msgs[0].Content, "hello")
	}

}
