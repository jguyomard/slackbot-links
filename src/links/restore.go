package links

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type logData struct {
	Action string
	Level  string
	Msg    string
	Time   *time.Time
	Link   Link
}

var (
	linksIdsAnalysed = []string{}
)

// Restore links from json file
func Restore(filepath string) bool {

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening links-file:", err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		restoreLink(scanner.Bytes())
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return true
}

func restoreLink(lineStr []byte) bool {
	var line logData

	// json decode
	if err := json.Unmarshal(lineStr, &line); err != nil {
		fmt.Printf("jsonDecode ERROR (Unmarshal): %v\n", err)
		return false
	}

	switch line.Action {

	// Save Link (insert or update)
	case "save":
		if linkAlreadyAnalysed(line.Link) {
			fmt.Printf(" - Link %s (url=%s) already analysed. Continue.\n", line.Link.ID, line.Link.URL)
			return false
		}
		if line.Link.Save() {
			markLinkAsAnalysed(line.Link)
			fmt.Printf(" - Link %s (%s) restored.\n", line.Link.ID, line.Link.URL)
			return true
		}
		fmt.Printf(" - Error saving Link %s (%s). Continue.\n", line.Link.ID, line.Link.URL)
		return false

	// Delete Link
	case "delete":
		if line.Link.Delete() {
			fmt.Printf(" - Link %s (%s) deletion restored.\n", line.Link.ID, line.Link.URL)
			return true
		}
		fmt.Printf(" - Error deleting Link %s (%s). Continue.\n", line.Link.ID, line.Link.URL)
		return false

	}

	return false
}

func markLinkAsAnalysed(link Link) {
	linksIdsAnalysed = append(linksIdsAnalysed, link.ID)
}

func linkAlreadyAnalysed(link Link) bool {
	return inArr(link.ID, linksIdsAnalysed)
}

func inArr(str string, arr []string) bool {
	for _, val := range arr {
		if str == val {
			return true
		}
	}
	return false
}
