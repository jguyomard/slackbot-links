package links

import "testing"

func TestGetLink(t *testing.T) {
	newlink := NewLink(testURL)
	newlink.Save()

	link, found := Get(newlink.GetID())
	link.Delete()

	if !found {
		t.Fatal("TestGetLink: link.Get() error (not found)")
	}
	if link.GetID() != newlink.GetID() {
		t.Fatal("TestGetLink: link.Get() error (ID doesn't match)")
	}
}

func TestGetUnknowLink(t *testing.T) {
	_, found := Get("unknow")
	if found {
		t.Fatal("TestGetUnknowLink: link.Get() error (found)")
	}
}
