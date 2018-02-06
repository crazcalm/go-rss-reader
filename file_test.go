package rss

import (
	"path/filepath"
	"strings"
	"testing"
)

/*
ExtractFileContent is responsible for extracting the
content from the rss url file. It does not validate this data.
*/
func TestExtractFileContent(t *testing.T) {
	//urls
	sn := "http://www.leoville.tv/podcasts/sn.xml"
	alerts := "https://www.us-cert.gov/ncas/alerts.xml"
	bulletins := "https://www.us-cert.gov/ncas/bulletins.xml"
	activity := "https://www.us-cert.gov/ncas/current-activity.xml"
	npr := "https://www.npr.org/rss/rss.php?id=1001"
	fbi := "https://www.fbi.gov/feeds/fbi-this-week-podcast/rss.xml"
	linuxSecurity := "http://feeds.feedburner.com/LinuxJournalSecurity"
	linuxSysadmin := "http://feeds.feedburner.com/LinuxJournalSysadmin"
	linuxMag := "http://www.linux-magazine.com/rss/feed/lmi_full"
	sed := "http://softwareengineeringdaily.com/feed/podcast/"
	xkcd := "https://www.xkcd.com/rss.xml"
	commitStrip := "http://www.commitstrip.com/en/feed/"
	rssUrls := []string{sn, alerts, bulletins, activity, npr, fbi, linuxSecurity, linuxSysadmin, linuxMag, sed, xkcd, commitStrip}

	//urls to tags map
	urlToTags := map[string][]string{
		sn:            []string{"security", "favorite", "audio"},
		alerts:        []string{"security"},
		bulletins:     []string{"security"},
		activity:      []string{"security"},
		npr:           []string{"news"},
		fbi:           []string{"audio", "FBI"},
		linuxSecurity: []string{"security"},
		linuxSysadmin: []string{"devops"},
		linuxMag:      []string{"linux"},
		sed:           []string{"audio"},
		xkcd:          []string{"comic"},
		commitStrip:   []string{"comic"},
	}

	//files
	dir := "test_data"
	empty := filepath.Join(dir, "empty_rss")
	oneRss := filepath.Join(dir, "one_rss")
	oneRssWithTags := filepath.Join(dir, "one_rss_with_tags")
	multipleRss := filepath.Join(dir, "multiple_rss")
	multipleRssWithTags := filepath.Join(dir, "multiple_rss_with_tags")

	tests := []struct {
		File string
		Urls []string
		Tags map[string][]string
	}{
		{empty, []string{}, urlToTags},
		{oneRss, []string{sn}, urlToTags},
		{oneRssWithTags, []string{sn}, urlToTags},
		{multipleRss, rssUrls, urlToTags},
		{multipleRssWithTags, rssUrls, urlToTags},
	}

	for _, test := range tests {
		results := ExtractFileContent(test.File)

		//Check number of items
		if len(test.Urls) != len(results) {
			t.Errorf("The number of results do not match. Expected %d, but got %d", len(test.Urls), len(results))
		}

		//Check to see I got the expected urls
		for index, fileData := range results {
			url := test.Urls[index]

			if test.Urls[index] != fileData.URL {
				t.Errorf("Expected %s, but got %s", url, fileData.URL)
			}

			//Check for the expected tags
			if len(fileData.Tags) > 0 {
				for index, tag := range fileData.Tags {
					if !strings.Contains(tag, test.Tags[url][index]) {
						t.Errorf("Tags for url %s don't match. Expected %s, but got %s", url, tag, test.Tags[url][index])
					}
				}
			}
		}
	}
}

/*
CheckFile is responsible for making sure that
the file exist, which includes

- making sure that the file path exists
- making sure that the file path is not a directory
*/
func TestCheckFile(t *testing.T) {
	file := filepath.Join("test_data", "one_rss")
	dir := filepath.Join("test_data")
	notAFile := filepath.Join("test_data", "not_a_file")

	tests := []struct {
		File      string
		ExpectErr bool
		Err       string
	}{
		{"", true, "file cannot be an empty string"},
		{dir, true, "test_data is not a file"},
		{file, false, "none"},
		{notAFile, true, "file test_data/not_a_file does not exist"},
	}

	for _, test := range tests {
		err := CheckFile(test.File)

		//Check error case
		if err != nil && test.ExpectErr {
			if !strings.EqualFold(err.Error(), test.Err) {
				t.Errorf("Error %s != %s", err.Error(), test.Err)
			}

			//Error case is working as expected
			continue
		}

		//Check expected error but did not get one case
		if err == nil && test.ExpectErr {
			t.Errorf("Expected err %s, but got nothing", test.Err)
		}

		//Check for unexpected error
		if err != nil && !test.ExpectErr {
			t.Errorf("Unexpected err %s", err.Error())
		}

		//Got an unexpected Error
		if err != nil {
			t.Errorf("Got an unexpected error: %s", err.Error())
		}
	}

}
