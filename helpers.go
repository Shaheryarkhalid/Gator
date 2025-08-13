package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/Shaheryarkhalid/Gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		if strings.Contains(fmt.Sprintf("%v", err), "no rows") {
			return fmt.Errorf("No feeds to scrape.")
		}
		return scrapeFeeds(s)
	}
	fmt.Printf("Fetchting '%v' feed.\n", nextFeed.Name)
	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		fmt.Println("Error trying to mark feed fetched: %w", err)
	}
	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Printf("Error trying to fetch feed: %v\n", err)
		return nil
	}
	for _, rssItem := range rssFeed.Channel.Item {
		pubAt, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", rssItem.PubDate)
		post := database.Post{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			Title:       rssItem.Title,
			Url:         rssItem.Link,
			Description: rssItem.Description,
			PublishedAt: pubAt,
			FeedID:      nextFeed.ID,
		}
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams(post))
		if err != nil {
			if !strings.Contains(err.Error(), "") {
				fmt.Printf("Error while saving the posts: %v", err)
				continue
			}
		}
	}
	return nil
}

func fetchFeed(ctx context.Context, feedUrl string) (rssFeed *RSSFeed, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return rssFeed, fmt.Errorf("Error creating the request: %v", err)
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return rssFeed, fmt.Errorf("Error trying to get response from url: %v", err)
	}
	defer resp.Body.Close()
	err = xml.NewDecoder(resp.Body).Decode(&rssFeed)
	if err != nil {
		return rssFeed, fmt.Errorf("Error happened while decoding the response: %v", err)
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for idx, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[idx] = item
	}
	return rssFeed, nil
}
