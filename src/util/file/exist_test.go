package file

import "testing"

func TestFileExist(t *testing.T) {
	tmp := ""
	exis := FileExist(tmp)
	if !exis {
		t.Errorf("file actual exist but get not\n")
	}
}
