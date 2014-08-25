package parser

import (
	"bytes"
	"encoding/xml"
	"strings"
	"time"
)

type rss1Feed struct {
	XMLName xml.Name    `xml:"RDF"`
	Channel rss1Channel `xml:"channel"`
	Items   []rssItem   `xml:"item"`
	Image   rssImage    `xml:"image"`
}

type rss1Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Image       rssImage `xml:"image"`
	TTL         int      `xml:"ttl"`
	SkipHours   []int    `xml:"skipHours>hour"`
	SkipDays    []string `xml:"skipDays>day"`
}

func ParseRss1(b []byte) (Feed, error) {
	var f Feed
	var rss rss1Feed

	decoder := xml.NewDecoder(bytes.NewReader(b))
	decoder.DefaultSpace = "parserfeed"

	if err := decoder.Decode(&rss); err != nil {
		return f, err
	}

	var image = rss.Image
	if image == (rssImage{}) {
		image = rss.Channel.Image
	}

	f = Feed{
		Title:       rss.Channel.Title,
		Description: rss.Channel.Description,
		SiteLink:    rss.Channel.Link,
		Image: Image{
			image.Title, image.Url,
			image.Width, image.Height},
	}

	if rss.Channel.TTL != 0 {
		f.TTL = time.Duration(rss.Channel.TTL) * time.Minute
	}

	f.SkipHours = make(map[int]bool, len(rss.Channel.SkipHours))
	for _, v := range rss.Channel.SkipHours {
		f.SkipHours[v] = true
	}

	f.SkipDays = make(map[string]bool, len(rss.Channel.SkipDays))
	for _, v := range rss.Channel.SkipDays {
		f.SkipDays[strings.Title(v)] = true
	}

	for _, i := range rss.Items {
		article := Article{Title: i.Title, Description: i.Description, Link: i.Link}
		if i.Id == "" {
			article.Id = i.Link
		} else {
			article.Id = i.Id
		}

		var err error
		if i.PubDate != "" {
			if article.Date, err = parseDate(i.PubDate); err != nil {
				return f, err
			}
		} else if i.Date != "" {
			if article.Date, err = parseDate(i.Date); err != nil {
				return f, err
			}
		}
		f.Articles = append(f.Articles, article)
	}

	f.HubLink = getHubLink(b)

	return f, nil
}
