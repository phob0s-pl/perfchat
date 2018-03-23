package tests

import (
	"net/http"
	"testing"
	"time"

	api "github.com/phob0s-pl/perfchat/apiv1"
)

func TestGetUsersByNotAddedUser(t *testing.T) {
	var (
		address = "localhost:9021"
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

	if _, err := client.GetUsers(); err == nil {
		t.Errorf("Getting users with user not in engine should fail")
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}
	<-done
}

func TestGetUsers(t *testing.T) {
	var (
		address = "localhost:9022"
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

	users, err := client.GetUsers()
	if err != nil {
		t.Errorf("Getting users failed, err=%s", err)
	}

	if err := server.Srv.Close(); err != nil {
		t.Errorf("Closing server failed, err=%s", err)
	}
	<-done

	engineUsers := server.API.Engine.ListUsers()
	if len(users) != len(engineUsers) {
		t.Fatalf("Expected %d users, got %d", len(engineUsers), len(users))
	}

	for _, user := range users {
		// for security reasons
		if user.AuthID != "" || user.Token != "" {
			t.Errorf("User %s shouldn't have token or auth set", user.Name)
		}
		engineUser, err := server.API.Engine.GetUserByName(user.Name)
		if err != nil {
			t.Errorf("Returned user %s, which is not in engine", user.Name)
			continue
		}
		if engineUser.Role != user.Role {
			t.Errorf("User %s, role mismatch, got %s, expected %s", user.Name, user.Role, engineUser.Role)
		}
	}
}
