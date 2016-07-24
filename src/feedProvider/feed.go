package feedProvider

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type NewsFeed struct {
	Channel Channel `xml:"channel" json:"channel"`
}

type Channel struct {
	Title       string `xml:"title" json:"title"`
	Link        string `xml:"link" json:"link"`
	Description string `xml:"description" json:"description"`
	Items       []Item `xml:"item" json:"item"`
}

type Item struct {
	Title       string `xml:"title" json:"title"`
	Link        string `xml:"link" json:"link"`
	Description string `xml:"description" json:"description"`
	PubDate     string `xml:"pubDate" json:"pubDate"`
}

func ReadNewsFeed(url string) *NewsFeed {
	resultFeed := new(NewsFeed)

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(body, resultFeed)

	return resultFeed
}
