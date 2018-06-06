package gui

import (
	"strings"
)

const (
	//HeaderText -- Generic header used for all screens
	HeaderText = "go-rss-reader"
	//QuitText -- Instructions on quiting the application
	QuitText = "Quit: Ctrl-C"
	//BackText -- Instructions on travelling to a previous screen
	BackText = "Back: Ctrl-B"
	//SelectText -- Instructions on selecting
	SelectText = "Select: Enter"
	//RefreshAllText -- Instructions on Refreshing all feeds
	RefreshAllText = "Refresh All: Ctrl-A"
	//RefreshOneText -- Instrcutions on Refreshing one feed
	RefreshOneText = "Refresh One: Ctrl-R"
)

var (
	//FeedsFooterText -- Text used for the Feeds footer
	FeedsFooterText = strings.Join([]string{QuitText, SelectText, RefreshAllText, RefreshOneText}, ", ")
	//EpisodesFooterText -- Text used for Episodes footer
	EpisodesFooterText = strings.Join([]string{QuitText, BackText, SelectText}, ", ")
	//PagerFooterText -- Text used for the Pager Footer
	PagerFooterText = strings.Join([]string{QuitText, BackText}, ", ")
)
