package rss

import (
	"ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
)

const (
	alertXML = "alerts.xml"
	linuxXML = "LinuxJournalSecurity.xml"
	nprXML = "npr.xml"
	snXML = "sn.xml"
	thoughtWorksXML = "thoughtworks.xml"
	xkcdXML = "xkcd.xml"
)

var (
	URLs = map[string]string{
		alertXML: "https://www.us-cert.gov/ncas/alerts.xml",
		linuxXML: "http://feeds.feedburner.com/LinuxJournalSecurity",
		nprXML: "https://www.npr.org/rss/rss.php?id=1001",
		snXML: "http://www.leoville.tv/podcasts/sn.xml",
		thoughtWorksXML: "http://feeds.soundcloud.com/users/soundcloud:users:94605026/sounds.rss",
		xkcdXML: "https://www.xkcd.com/rss.xml",
	}
	Tags = map[string][]string{
		alertXML: []string{"security"},
		linuxXML: []string{"security"},
		nprXML: []string{"news"},
		snXML: []string{"security", "favorite", "audio"},
		thoughtWorksXML: []string{"audio"},
		xkcdXML: []string{"comic"},
	}
)

func checkForInternet() bool {
	/*
		Use: https://stackoverflow.com/questions/10040954/alternative-to-google-finance-api

		That post has a link to a stock market api call that returns json. Given
		that it is stock market data it will, hopefully, always be online and reachable
		from all countries.

		I can check the response from that call and decided from that whether or not
		I have internet.
	*/
	//TODO: implement later
	return true
}

//getTestFeed -- Used to fetch locally stored raw rss file
//and turn them into gofeed.Feed so that I can use them in my tests
func getTestFeed(name string) *gofeed.Feed {
	path := filepath.Join("test_data", "raw_rss", name)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Errorf("(Test File Error) Occured when reading %s: %s", path, err.Error()))
	}
	
	feedParser := gofeed.NewParser()
	feed, err := feedParser.Parse(file)
	if err != nil {
		log.Fatal(fmt.Errorf("(Test file Error) Occured when parsing %s: %s", path, err.Error()))
	}
	return feed
}

func TestTitle(t *testing.T) {
	tests := []struct{
		Feed *Feed
		ExpectedTitle string
	}{
		{&Feed{URLs[snXML], Tags[snXML], getTestFeed(snXML)}},
	}

	for _, test := range tests {
		if !strings.EqualFold(test.Feed.Title(), test.ExpectedTitle){
			//write error
		}
	}
}

func TestNewFeed(t *testing.T) {
	//FileData variables
	goodData := FileData{"http://www.leoville.tv/podcasts/sn.xml", []string{"security", "favorite", "audio"}}
	badData := FileData{"http://www.linux-magazine.com/rs/feed/lmi_full", []string{"error"}}
	noData := FileData{"", []string{}}

	//Error
	errorStr := "Error occured while trying to parse feed"

	tests := []struct {
		Data      FileData
		ExpectErr bool
		Err       string
		Title     string
	}{
		{goodData, false, "", "Security Now (MP3)"},
		{badData, true, errorStr, ""},
		{noData, true, errorStr, ""},
	}

	for _, test := range tests {
		feed, err := NewFeed(test.Data)

		//Case: Received and error while expecting an error
		if err != nil && test.ExpectErr {
			if !strings.EqualFold(err.Error(), test.Err) {
				t.Errorf("Expected error: %s, received err -- %s", test.Err, err.Error())
			}
			continue
		}

		//Case: Received and error while not expecting an error
		if err != nil && !test.ExpectErr {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		//Case: Expected an error, but did not receive an error
		if err == nil && test.ExpectErr {
			t.Errorf("Expected error: %s, but no error was received", err.Error())
		}

		//Case: Did not expect an error, and I did not receive an error
		if err == nil && !test.ExpectErr {
			if !strings.EqualFold(feed.Data.Title, test.Title) {
				t.Errorf("Expected title: %s, but got %s", test.Title, feed.Data.Title)
			}
			continue
		}
	}
}
