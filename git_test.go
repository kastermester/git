package git

import (
	"testing"
)

func TestStuff(t *testing.T) {
	git, err := New()
	if err != nil {
		t.Fatal(err)
	}
	err = git.SyncRepositoryToRemoteBranch("./test/something", "https://github.com/kastermester/jobs", "ce5a422")
	if err != nil {
		t.Fatal(err)
	}
}
