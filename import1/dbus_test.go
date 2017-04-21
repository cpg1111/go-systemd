package import1

import (
	"fmt"
	"os"
	"testing"
)

func TestNewImportd(t *testing.T) {
	conn, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if conn == nil {
		t.Errorf("conn is nil \n")
	}
}

func TestImportTar(t *testing.T) {
	conn, err := New()
	if err != nil {
		t.Fatal(err)
	}
	gopath := os.Getenv("GOPATH")
	imgPath := fmt.Sprintf("%ssrc/github.com/coreos/go-systemd/test.tar", gopath)
	img, err := os.Open(imgPath)
	if err != nil {
		t.Fatal(err)
	}
	fd := img.Fd()
	t.Log(fd)
	transfer, err := conn.ImportTar(int(fd), "test-tar-image", false, false)
	if err != nil {
		t.Error(err)
	}
	t.Log(transfer)
	if transfer == nil {
		t.Error("did not transfer")
	}
}

func TestImportRaw(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	imgPath := fmt.Sprintf("%s/src/github.com/coreos/go-systemd/test.iso", gopath)
	conn, err := New()
	if err != nil {
		t.Fatal(err)
	}
	img, err := os.Open(imgPath)
	if err != nil {
		t.Fatal(err)
	}
	fd := img.Fd()
	transfer, err := conn.ImportRaw(int(fd), "test-iso-image", false, false)
	if err != nil {
		t.Error(err)
	}
	t.Log(transfer)
	if transfer == nil {
		t.Error("did not transfer")
	}
}
