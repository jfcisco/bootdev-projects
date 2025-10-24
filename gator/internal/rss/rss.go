package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

// Internal struct to help with parsing feed XML.
type rssFeedXML struct {
	Channel struct {
		Title string `xml:"title"`
		Links []struct {
			Name xml.Name
			Data string `xml:",chardata"`
		} `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func (f *rssFeedXML) toRSSFeed() *RSSFeed {
	// Handle feeds with multiple <link> elements. We want the one with an empty namespace.
	var linkData string
	for _, link := range f.Channel.Links {
		if link.Name.Space == "" {
			linkData = link.Data
			break
		}
	}

	return &RSSFeed{
		Channel: RSSChannel{
			Title:       f.Channel.Title,
			Description: f.Channel.Description,
			Link:        linkData,
			Items:       f.Channel.Items,
		},
	}
}

type RSSChannel struct {
	Title       string
	Link        string
	Description string
	Items       []RSSItem
}

type RSSFeed struct {
	Channel RSSChannel
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (f *RSSFeed) unescape() {
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)

	for i, item := range f.Channel.Items {
		f.Channel.Items[i].Title = html.UnescapeString(item.Title)
		f.Channel.Items[i].Description = html.UnescapeString(item.Description)
	}
}

// FetchFeed fetches and parses the RSS feed from the given URL.
func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	req.Header.Set("User-Agent", "gator")
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	defer res.Body.Close()

	rawFeed, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %w", err)
	}

	ifeed := &rssFeedXML{}
	if err := xml.Unmarshal(rawFeed, ifeed); err != nil {
		return nil, fmt.Errorf("feed unmarshal error: %w", err)
	}
	feed := ifeed.toRSSFeed()
	feed.unescape()
	return feed, nil
}
