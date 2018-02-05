package rss

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//FileData -- Container to keep the url and tags of a feed
type FileData struct {
	URL  string
	Tags []string
}

/*
ExtractFileContent is responsible for extracting the
content from the rss url file. It does not validate this data.
*/
func ExtractFileContent(f string) (results []FileData) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("Error when reading file: %s", err.Error())
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		l := strings.TrimSpace(line)

		//Check to see if the line is empty
		if strings.EqualFold(l, "") {
			continue
		}

		//Collect urls and tags
		items := strings.Split(l, " ")
		if len(items) == 1 {
			results = append(results, FileData{items[0], []string{}})
		} else if len(items) > 1 {
			results = append(results, FileData{items[0], items[1:]})
		}
	}
	return results
}

/*
CheckFile is responsible for making sure that
the file exist, which includes

- making sure that the file path exists
- making sure that the file path is not a directory
*/
func CheckFile(file string) (err error) {
	fileInfo, err := os.Stat(file)

	//File does not exist
	if os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", file)
	}

	//Path exists, but it is a directory and not a file
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is not a file", file)
	}

	return nil
}
