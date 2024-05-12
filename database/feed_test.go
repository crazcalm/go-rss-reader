package database

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetFeedDataFromSiteV2(t *testing.T) {
	t.Parallel()

	response_message := "hello world"

	good_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, response_message)
	}))

	bad_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "TeaPot Error", 418)
	}))
	not_up_server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s", r.Proto)
	}))

	t.Cleanup(func() {
		bad_server.Close()
		good_server.Close()
	})

	tcs := []struct {
		name   string
		url    string
		reader func(io.Reader) ([]byte, error)
		err    error
	}{
		{"Not Okay Status Code", bad_server.URL, io.ReadAll, errors.New("returned a status code of 418")},
		{"Read Response Error", good_server.URL, func(_ io.Reader) ([]byte, error) {
			return nil, errors.New("Read Error")
		}, errors.New("Unable to read response body: Read Error")},
		{"No Server Running", not_up_server.URL, io.ReadAll, errors.New("Error trying to get the raw feed data from : Get \"\": unsupported protocol scheme \"\"")},
		{"Happy Path", good_server.URL, io.ReadAll, nil},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetFeedDataFromSiteV2(tc.url, tc.reader)
			if err != nil {
				if strings.Contains(err.Error(), tc.err.Error()) == false {
					t.Fatalf("Expected %q, but got %q", tc.err.Error(), err.Error())
				}
			} else {
				if strings.Compare(result, response_message) != 0 {
					t.Fatalf("Expected %q, but got %q", result, response_message)
				}
			}
		})
	}

}

func TestGetFeedDataFromSite(t *testing.T) {
	tests := []string{
		"http://www.leoville.tv/podcasts/sn.xml",
		"https://www.xkcd.com/rss.xml",
		"https://latenightlinux.com/feed/mp3",
		"http://feeds.feedburner.com/InternetHistoryPodcast",
		"https://feeds.mozilla-podcasts.org/irl",
	}

	for _, test := range tests {
		body, err := GetFeedDataFromSite(test)
		if err != nil {
			t.Errorf("Error for url (%s): %s", test, err.Error())
		}

		if strings.EqualFold(body, "") {
			t.Errorf("The body was empty for url %s", test)
		}
	}
}
