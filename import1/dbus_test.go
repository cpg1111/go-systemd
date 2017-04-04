package import1

import (
	"testing"
)

func TestNewImportd(t *testing.T) {
	conn, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	if conn == nil {
		t.Errorf("conn is nil \n")
		return
	}
}
