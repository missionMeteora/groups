package groups

import (
	"os"
	"testing"
)

func TestGroups(t *testing.T) {
	var (
		g   *Groups
		gs  []string
		err error
	)

	if g, err = New("./_testdata"); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("./_testdata")

	if gs, err = g.Set("0", "users"); err != nil {
		t.Fatalf("Error encountered while setting: %v", err)
	}

	if len(gs) != 1 {
		t.Fatalf("Invalid number of groups, expected %v and received %v", 1, len(gs))
	}

	if gs[0] != "users" {
		t.Fatalf("Invalid group name: expected %v and received %v", "users", gs[0])
	}

	if err = g.Close(); err != nil {
		t.Fatal(err)
	}

	if g, err = New("./_testdata"); err != nil {
		t.Fatal(err)
	}

	return
}
