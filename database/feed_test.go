package database

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"

	"github.com/crazcalm/go-rss-reader/file"
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

func getTestFile(path []string) []byte {
	pathString := strings.Join(path, string(os.PathSeparator))
	file, err := ioutil.ReadFile(pathString)
	if err != nil {
		log.Fatal(fmt.Errorf("(Test File Error) Occured when reading %s: %s", pathString, err.Error()))
	}
	return file
}

//getTestFeed -- Used to fetch locally stored raw rss file
//and turn them into gofeed.Feed so that I can use them in my tests
func getTestFeed(name string) *gofeed.Feed {
	path := []string{"test_data", "rss_raw", name}
	file := getTestFile(path)

	feedParser := gofeed.NewParser()
	feed, err := feedParser.Parse(bytes.NewReader(file))
	if err != nil {
		log.Fatal(fmt.Errorf("(Test file Error) Occured when parsing %s: %s", path, err.Error()))
	}
	return feed
}

func TestGuiItemsData(t *testing.T) {
	var noData gofeed.Feed
	feedTitle := "Title"

	tests := []struct {
		Feed          *Feed
		ExpectedTitle []string
		ExpectedDate  []string
	}{
		{&Feed{1, URLs[snXML], feedTitle, Tags[snXML], getTestFeed(snXML)}, []string{"SN 649: Meltdown & Spectre Emerge", "SN 648: Post Spectre?", "SN 647: The Dark Caracal", "SN 646: The InSpectre", "SN 645: The Speculation Meltdown", "SN 644: NSA Fingerprints", "SN 643: The Story of Bitcoin", "SN 642: BGP", "SN 641: The iOS Security Trade-off", "SN 640: More News & Feedback"}, []string{"Feb 06", "Jan 30", "Jan 23", "Jan 16", "Jan 09", "Jan 02", "Dec 26", "Dec 19", "Dec 12", "Dec 05"}},
		{&Feed{2, "NO URL", feedTitle, []string{}, &noData}, []string{}, []string{}},
	}

	for _, test := range tests {
		data := test.Feed.GuiItemsData()

		if len(data) != len(test.ExpectedTitle) || len(data) != len(test.ExpectedDate) {
			t.Errorf("%s does not have the expected number of episodes. Expected %d, but got %d", test.Feed.Title, len(test.ExpectedDate), len(data))
		}

		for index, item := range data {
			if !strings.EqualFold(item.Title, test.ExpectedTitle[index]) {
				t.Errorf("For feed %s, expected %s episode title, but got %s", test.Feed.Title, test.ExpectedTitle[index], item.Title)
			}

			if !strings.EqualFold(item.Date, test.ExpectedDate[index]) {
				t.Errorf("For feed %s, expected %s episode date, but got %s", test.Feed.Title, test.ExpectedDate[index], item.Date)
			}
		}
	}
}

func TestEpisodeTotal(t *testing.T) {
	var noData gofeed.Feed
	feedTitle := "Title"

	tests := []struct {
		Feed          *Feed
		ExpectedCount int
	}{
		{&Feed{1, URLs[snXML], feedTitle, Tags[snXML], getTestFeed(snXML)}, 10},
		{&Feed{2, URLs[alertXML], feedTitle, Tags[alertXML], getTestFeed(alertXML)}, 10},
		{&Feed{3, URLs[xkcdXML], feedTitle, Tags[xkcdXML], getTestFeed(xkcdXML)}, 4},
		{&Feed{4, "NO URL", feedTitle, []string{}, &noData}, 0},
	}

	for _, test := range tests {
		if test.Feed.EpisodeTotal() != test.ExpectedCount {
			t.Errorf("Expected %s to have %d episodes, but got %d", test.Feed.URL, test.ExpectedCount, test.Feed.EpisodeTotal())
		}
	}
}

func TestNewFeed(t *testing.T) {
	//FileData variables
	goodData := file.Data{"http://www.leoville.tv/podcasts/sn.xml", []string{"security", "favorite", "audio"}}
	badData := file.Data{"http://www.linux-magazine.com/rs/feed/lmi_full", []string{"error"}}
	noData := file.Data{"", []string{}}

	//Error
	errorStr := "Error occured while trying to parse feed"

	tests := []struct {
		Data      file.Data
		ExpectErr bool
		Err       string
		Title     string
	}{
		{goodData, false, "", "Security Now (MP3)"},
		{badData, true, errorStr, ""},
		{noData, true, errorStr, ""},
	}

	for _, test := range tests {
		feed, err := NewFeed(-1, test.Data)

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

func TestGetEpisodeFeed(t *testing.T) {
	feedTitle := "Title"

	tests := []struct {
		Feed          *Feed
		EpisodeNum    int
		ExpectedTitle string
		ExpectError   bool
	}{
		{&Feed{1, URLs[snXML], feedTitle, Tags[snXML], getTestFeed(snXML)}, 0, "SN 649: Meltdown & Spectre Emerge", false},
		{&Feed{2, URLs[snXML], feedTitle, Tags[snXML], getTestFeed(snXML)}, 5, "SN 644: NSA Fingerprints", false},
		{&Feed{3, URLs[snXML], feedTitle, Tags[snXML], getTestFeed(snXML)}, 12, "None", true},
		{&Feed{4, URLs[snXML], feedTitle, Tags[snXML], getTestFeed(snXML)}, -3, "None", true},
	}

	for _, test := range tests {
		episode, err := test.Feed.GetEpisode(test.EpisodeNum)

		if err == nil && test.ExpectError {
			t.Errorf("Expected an error, but none was received.")
		}

		if err != nil && !test.ExpectError {
			t.Errorf("Got an unexpected err: %s", err.Error())
		}

		//We expected and error and we got an error.
		//Test case passed
		if err != nil && test.ExpectError {
			continue
		}

		if !strings.EqualFold(episode.Data.Title, test.ExpectedTitle) {
			t.Errorf("Expected the title to be %s, but got %s instead.", test.ExpectedTitle, episode.Data.Title)
		}

	}
}

func TestGetEpisodes(t *testing.T) {
	var noData gofeed.Feed
	feedTitle := "Title"

	tests := []struct {
		Name        string
		Feed        *Feed
		ExpectedNum int
	}{
		{"snXML", &Feed{1, URLs[snXML], feedTitle, Tags[snXML], getTestFeed(snXML)}, 10},
		{"alertXML", &Feed{2, URLs[alertXML], feedTitle, Tags[alertXML], getTestFeed(alertXML)}, 10},
		{"xkcdXML", &Feed{3, URLs[xkcdXML], feedTitle, Tags[xkcdXML], getTestFeed(xkcdXML)}, 4},
		{"No Url", &Feed{4, "NO URL", feedTitle, []string{}, &noData}, 0},
	}

	for _, test := range tests {
		if len(test.Feed.GetEpisodes()) != test.ExpectedNum {
			t.Errorf("For %s we expected %d episodes but got %d instead.", test.Name, test.ExpectedNum, len(test.Feed.GetEpisodes()))
		}
	}

}
