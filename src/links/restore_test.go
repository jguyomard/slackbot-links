package links

import "testing"

func ExampleRestore() {
	Restore("testdata/links.json")
	// Output:
	//  - Link test-restore (https://ilonet.fr/welcome-to-hugo/) restored.
	//  - Link test-restore (url=https://ilonet.fr/welcome-to-hugo/) already analysed. Continue.
	//  - Link test-restore (https://ilonet.fr/welcome-to-hugo/) deletion restored.
	// ES Error: elastic: Error 404 (Not Found)
	//  - Error deleting Link test-restore (https://ilonet.fr/welcome-to-hugo/). Continue.
	// jsonDecode ERROR (Unmarshal): invalid character 'i' looking for beginning of value
}

func TestRestore_badFile(t *testing.T) {
	if Restore("testdata/notfound.json") != false {
		t.Fatal("TestRestore_badFile: Restore() with bad file return true")
	}
}
