package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	db "github.com/KaePee/go-rssag/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func startScraper(db *db.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}

		wg.Wait()

	}
}

func scrapeFeed(queries *db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := queries.MarkFeedAsFetched(ctx, feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url.String)
	if err != nil {
		log.Println("error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		log.Println("Found post with title:", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Items))
	for _, item := range rssFeed.Channel.Items {
		description := pgtype.Text{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123, item.PubDate)
		// log.Printf("Logging original pubDate of feed: %v", item.PubDate)
		// log.Printf("Logging parsed published at Time: %v", pubAt)
		
		if err != nil {
			log.Printf("couldn't parse date %v with err %v, using current UTC time", item.PubDate, err)
			//pubAt = time.Now().UTC()
			continue
		}
		_, err = queries.CreatePost(ctx, db.CreatePostParams{
			ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
			CreatedAt:   pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
			UpdatedAt:   pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
			Title:       item.Title,
			Description: description,
			PublishedAt: pgtype.Timestamp{Time: pubAt, Valid: true},
			Url: item.Link,
			FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue 
			}
			log.Println("failed to create post:", err)
		}
	}

}
