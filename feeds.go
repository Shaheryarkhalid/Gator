package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Shaheryarkhalid/Gator/internal/databse"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title string `xml:"title"`
		Link string	 `xml:"link"`
		Description string `xml:"description"` 
		Item []RSSItem `xml:"item"`
	} `xml:"channel"`
}
type RSSItem struct{
	Title string `xml:"title"`
	Link string	 `xml:"link"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"`
} 

func fetchFeed(ctx context.Context, feedUrl string)(*RSSFeed, error){
	feed := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil  )
	if err  != nil {
		return &feed, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err  != nil {
		return &feed, err
	}
	defer resp.Body.Close()
	data , err := io.ReadAll(resp.Body)
	if err != nil {
		return &feed, fmt.Errorf("Error: Trying to read resp body:\n%v", err)
	 }
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return &feed, fmt.Errorf("Error: Trying to Unmarshal resp body:\n%v", err)
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	return &feed, nil
}


func scrapeFeeds(s *state){
	nextFeed, err :=  s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("Error: Trying to get next feed to fetch.")
		fmt.Println(err)
		os.Exit(1)
	}
	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		fmt.Println("Error: Trying to mark feed fetched.")
		fmt.Println(err)
		os.Exit(1)
	}
	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Println("Error: Trying to fetch next feed.")
		fmt.Println(err)
		os.Exit(1)
	}
		for _, item := range rssFeed.Channel.Item{
			pubAt, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
			createPostParams := databse.CreatePostParams{
				ID: uuid.New(),
				Title: item.Title,
				Url: item.Link,
				Description: item.Description,
				PublishedAt: pubAt,
				FeedID: nextFeed.ID,
			}
			createdPost, err := s.db.CreatePost(context.Background(), createPostParams)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(createdPost.Title)
	}
}
