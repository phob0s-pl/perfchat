package tests

import (
	"net/http"
	"testing"
	"time"

	api "github.com/phob0s-pl/perfchat/apiv1"
)

func TestAddUserNoAuth(t *testing.T) {
	var (
		address = "localhost:9001"
		server  = NewServer(address)
		client  = api.NewClient(dummyuser, address)
		done    = make(chan bool)
	)

	go func() {
		if err := server.Srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server failed, err=%s", err)
		}
		done <- true
	}()
	time.Sleep(time.Millisecond * 10)

	if err := client.AddUser(admin); err == nil {
		t.Errorf("Adding user with no auth should fail")
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}

	if l := server.API.Engine.UsersCount(); l != 0 {
		t.Errorf("Expected no users in engine, but got %d", l)
	}

	<-done
}

func TestAddSingleUser(t *testing.T) {
	var (
		address = "localhost:9002"
		server  = NewServer(address)
		client  = api.NewClient(admin, address)
		done    = make(chan bool)
	)
	server.API.Engine.AddUser(admin)

	go func() {

		if err := server.Srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Server failed, err=%s", err)
		}
		done <- true
	}()
	time.Sleep(time.Millisecond * 10)

	if err := client.AddUser(dummyuser); err != nil {
		t.Errorf("Adding user %s failed, err=%s", dummyuser.Name, err)
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}
	<-done

	if l := server.API.Engine.UsersCount(); l != 2 {
		t.Errorf("Expected one user in engine, but got %d", l)
	}

	user, err := server.API.Engine.GetUserByName(dummyuser.Name)
	if err != nil {
		t.Fatalf("User %s was not addded to engine", dummyuser.Name)
	}

	if dummyuser.Name != user.Name {
		t.Fatalf("User name expected %q, got %q", dummyuser.Name, user.Name)
	}

	if dummyuser.AuthID != user.AuthID {
		t.Fatalf("User AuthID expected %q, got %q", dummyuser.AuthID, user.AuthID)
	}

	if dummyuser.Role != user.Role {
		t.Fatalf("User Role expected %q, got %q", dummyuser.Role, user.Role)
	}

	if dummyuser.Token != user.Token {
		t.Fatalf("User name expected %q, got %q", dummyuser.Token, user.Token)
	}

}
