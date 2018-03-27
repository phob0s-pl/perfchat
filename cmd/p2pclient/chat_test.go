package main

import "testing"

func TestGetRandomRoom(t *testing.T) {
	chat := NewChat("")
	if name := chat.GetRandomRoom(); name != "" {
		t.Errorf("GetRandomRoom for empty chat returned %s", name)
	}

	chat.joinedGroups["test"] = &Group{Name: "test"}
	if name := chat.GetRandomRoom(); name != "test" {
		t.Errorf("GetRandomRoom for nonempty chat returned %s, expected %s", name, "test")
	}
}

func TestRemoveGroup(t *testing.T) {
	chat := NewChat("")
	chat.joinedGroups["test"] = &Group{Name: "test"}
	chat.RemoveGroup("invalid")
	if len(chat.joinedGroups) != 1 {
		t.Errorf("Removing nonexist room shouldnt influence room counter")
	}

	chat.RemoveGroup("test")
	if len(chat.joinedGroups) != 0 {
		t.Errorf("Removing room should decrease room count")
	}

}
