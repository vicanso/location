package router

import "testing"

func TestNewGroup(t *testing.T) {
	g := NewGroup("/xx")
	groups := GetGroups()
	if groups[len(groups)-1] != g {
		t.Fatalf("new group fail")
	}
}
