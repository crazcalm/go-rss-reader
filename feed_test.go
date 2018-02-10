package rss

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
)

const (
	alertXML        = "alerts.xml"
	linuxXML        = "LinuxJournalSecurity.xml"
	nprXML          = "npr.xml"
	snXML           = "sn.xml"
	thoughtWorksXML = "thoughtworks.xml"
	xkcdXML         = "xkcd.xml"
)

var (
	URLs = map[string]string{
		alertXML:        "https://www.us-cert.gov/ncas/alerts.xml",
		linuxXML:        "http://feeds.feedburner.com/LinuxJournalSecurity",
		nprXML:          "https://www.npr.org/rss/rss.php?id=1001",
		snXML:           "http://www.leoville.tv/podcasts/sn.xml",
		thoughtWorksXML: "http://feeds.soundcloud.com/users/soundcloud:users:94605026/sounds.rss",
		xkcdXML:         "https://www.xkcd.com/rss.xml",
	}
	Tags = map[string][]string{
		alertXML:        []string{"security"},
		linuxXML:        []string{"security"},
		nprXML:          []string{"news"},
		snXML:           []string{"security", "favorite", "audio"},
		thoughtWorksXML: []string{"audio"},
		xkcdXML:         []string{"comic"},
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
	path := filepath.Join("test_data", "rss_raw", name)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Errorf("(Test File Error) Occured when reading %s: %s", path, err.Error()))
	}

	feedParser := gofeed.NewParser()
	feed, err := feedParser.Parse(bytes.NewReader(file))
	if err != nil {
		log.Fatal(fmt.Errorf("(Test file Error) Occured when parsing %s: %s", path, err.Error()))
	}
	return feed
}

func TestGuiItemsData(t *testing.T) {
	var noData gofeed.Feed

	tests := []struct {
		Feed          *Feed
		ExpectedTitle []string
		ExpectedDate  []string
	}{
		{&Feed{URLs[snXML], Tags[snXML], getTestFeed(snXML)}, []string{"SN 649: Meltdown & Spectre Emerge", "SN 648: Post Spectre?", "SN 647: The Dark Caracal", "SN 646: The InSpectre", "SN 645: The Speculation Meltdown", "SN 644: NSA Fingerprints", "SN 643: The Story of Bitcoin", "SN 642: BGP", "SN 641: The iOS Security Trade-off", "SN 640: More News & Feedback"}, []string{"Feb 06", "Jan 30", "Jan 23", "Jan 16", "Jan 09", "Jan 02", "Dec 26", "Dec 19", "Dec 12", "Dec 05"}},
		{&Feed{"NO URL", []string{}, &noData}, []string{}, []string{}},
	}

	for _, test := range tests {
		data := test.Feed.GuiItemsData()

		if len(data) != len(test.ExpectedTitle) || len(data) != len(test.ExpectedDate) {
			t.Errorf("%s does not have the expected number of episodes. Expected %d, but got %d", test.Feed.Title(), len(test.ExpectedDate), len(data))
		}

		for index, item := range data {
			if !strings.EqualFold(item.Title, test.ExpectedTitle[index]) {
				t.Errorf("For feed %s, expected %s episode title, but got %s", test.Feed.Title(), test.ExpectedTitle[index], item.Title)
			}

			if !strings.EqualFold(item.Date, test.ExpectedDate[index]) {
				t.Errorf("For feed %s, expected %s episode date, but got %s", test.Feed.Title(), test.ExpectedDate[index], item.Date)
			}
		}
	}
}

func TestEpisodeTotal(t *testing.T) {
	var noData gofeed.Feed

	tests := []struct {
		Feed          *Feed
		ExpectedCount int
	}{
		{&Feed{URLs[snXML], Tags[snXML], getTestFeed(snXML)}, 10},
		{&Feed{URLs[alertXML], Tags[alertXML], getTestFeed(alertXML)}, 10},
		{&Feed{URLs[xkcdXML], Tags[xkcdXML], getTestFeed(xkcdXML)}, 4},
		{&Feed{"NO URL", []string{}, &noData}, 0},
	}

	for _, test := range tests {
		if test.Feed.EpisodeTotal() != test.ExpectedCount {
			t.Errorf("Expected %s to have %d episodes, but got %d", test.Feed.URL, test.ExpectedCount, test.Feed.EpisodeTotal())
		}
	}
}

func TestTitle(t *testing.T) {
	var noData gofeed.Feed

	tests := []struct {
		Feed          *Feed
		ExpectedTitle string
	}{
		{&Feed{URLs[snXML], Tags[snXML], getTestFeed(snXML)}, "Security Now (MP3)"},
		{&Feed{URLs[alertXML], Tags[alertXML], getTestFeed(alertXML)}, "US-CERT Alerts"},
		{&Feed{URLs[xkcdXML], Tags[xkcdXML], getTestFeed(xkcdXML)}, "xkcd.com"},
		{&Feed{"NO URL", []string{}, &noData}, "No URL"},
	}

	for _, test := range tests {
		if !strings.EqualFold(test.Feed.Title(), test.ExpectedTitle) {
			t.Errorf("Expected %s, but got %s", test.ExpectedTitle, test.Feed.Title())
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
