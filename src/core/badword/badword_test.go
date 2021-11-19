package badword

import "testing"

func TestBadword(t *testing.T) {
	t.Log(BadWordSearch("you mother fucker"))
	t.Log(BadWordReplace("you mother fucker"))
}
