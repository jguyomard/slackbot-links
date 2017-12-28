package links

import (
	"testing"
	"time"
)

func TestSaveLink(t *testing.T) {
	link := NewLink(testURL)

	link.SetTitle(testTitle)
	if link.Title != testTitle {
		t.Fatal("TestSaveLink: link.SetTitle() error")
	}

	link.SetExcerpt(testExcerpt)
	if link.Excerpt != testExcerpt {
		t.Fatal("TestSaveLink: link.SetExcerpt() error")
	}

	link.SetImageURL(testImageURL)
	if link.ImageURL != testImageURL {
		t.Fatal("TestSaveLink: link.SetImageURL() error")
	}

	dateSharedAt, _ := time.Parse("2006-01-02T15:04:05.000Z", testSharedAt)
	link.SetSharedAt(&dateSharedAt)
	if link.SharedAt.String() != dateSharedAt.String() {
		t.Fatal("TestSaveLink: link.SetSharedAt() error")
	}

	link.SetSharedBy("1", testSharedBy)
	if link.SharedBy.Name != testSharedBy {
		t.Fatal("TestSaveLink: link.SetSharedBy() error")
	}

	link.SetSharedOn("1", testSharedOn)
	if link.SharedOn.Name != testSharedOn {
		t.Fatal("TestSaveLink: link.SetSharedOn() error")
	}

	saved := link.Save()
	if !saved {
		t.Fatal("TestSaveLink: link.Save() error")
	}

	// Save link ID for other tests
	linkID = link.GetID()
}

// This test should returns one duplicate
func TestFindDuplicate(t *testing.T) {
	link := NewLink(testURL)
	duplicates := link.FindDuplicates()

	if duplicates.GetTotal() != 1 {
		t.Fatalf("TestFindDuplicate: link.FindDuplicates() error, %d duplicate(s) link(s) found instead of 1", duplicates.GetTotal())
	}
}

func TestDeleteLink(t *testing.T) {
	link, _ := Get(linkID)
	deleted := link.Delete()
	if !deleted {
		t.Fatal("TestDeleteLink: link.Delete() error, deleted != true")
	}

	deletedAgain := link.Delete()
	if deletedAgain {
		t.Fatal("TestDeleteLink: link.Delete() error, deletedAgain != false")
	}
}

// This test should returns no duplicate
func TestFindNoDuplicate(t *testing.T) {
	link := NewLink(testURL)
	duplicates := link.FindDuplicates()

	if duplicates.GetTotal() > 0 {
		t.Fatalf("TestFindNoDuplicate: link.FindDuplicates() error, %d duplicate(s) link(s) found instead of 0", duplicates.GetTotal())
	}
}
