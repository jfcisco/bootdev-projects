package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jfcisco/gator/internal/database"
	"github.com/jfcisco/gator/internal/rss"
	"github.com/lib/pq"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("handlerLogin: no username provided")
	}

	name := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user %s does not exist", name)
	} else if err != nil {
		return fmt.Errorf("handlerLogin: %w", err)
	}

	if err := s.config.SetUser(user.Name); err != nil {
		return fmt.Errorf("handlerLogin: %w", err)
	}

	fmt.Printf("Logged in as %s\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("handlerRegister: no username provided")
	}

	name := cmd.args[0]
	createParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	user, err := s.db.CreateUser(context.Background(), createParams)
	if err != nil {
		return err
	}

	s.config.SetUser(user.Name)
	fmt.Printf("registered new user: %v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("handlerReset: %w", err)
	}

	if err := s.config.SetUser(""); err != nil {
		return fmt.Errorf("handlerReset: %w", err)
	}

	fmt.Println("successfully reset all users")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("handlerUsers: %w", err)
	}

	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	s.db.MarkFeedFetched(context.Background(), feed.ID)

	fmt.Printf("fetching %s at %v\n", feed.Url, time.Now())

	content, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("fetched %v items:\n", len(content.Channel.Items))
	for _, item := range content.Channel.Items {
		newPost := database.CreatePostParams{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullString{String: item.PubDate, Valid: item.PubDate != ""},
			FeedID:      feed.ID,
		}

		_, err := s.db.CreatePost(context.Background(), newPost)

		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code.Name() == "unique_violation" {
			// Unique constraint violation, skip this post
			continue
		} else if err != nil {
			return err
		}
	}
	fmt.Println("fetch complete")
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("please enter a duration string")
	}

	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("agg: fetching every %v\n", duration)

	ticker := time.NewTicker(duration)
	for {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
		<-ticker.C
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("please enter a name and url")
	}

	name, feedUrl := cmd.args[0], cmd.args[1]

	if _, err := url.ParseRequestURI(feedUrl); err != nil {
		return fmt.Errorf("feed url is not valid")
	}

	createArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       feedUrl,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), createArgs)

	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}
	fmt.Printf("new feed created: %v\n", feed)

	// Attempt to follow the newly-created feed
	if _, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feedsWithUser, err := s.db.GetFeedsAndUser(context.Background())
	if err != nil {
		return fmt.Errorf("handlerFeeds: %w", err)
	}

	for _, fwu := range feedsWithUser {
		fmt.Printf("%s (%s) | added by %s\n", fwu.Feed.Name, fwu.Feed.Url, fwu.User.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("no feed url specified to follow")
	}

	feedUrl := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("no feed exists with url '%s'", feedUrl)
	} else if err != nil {
		return fmt.Errorf("handlerFollow: %w", err)
	}

	// Check that given feed isn't already being followed yet by this user
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return fmt.Errorf("handlerFollow: %w", err)
	}

	for _, follow := range follows {
		if follow.Feed.Url == feedUrl {
			return fmt.Errorf("given feed is already followed")
		}
	}

	params := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Printf("%s successfully followed the feed %s\n", follow.User.Name, follow.Feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		user.Name,
	)

	if err != nil {
		return fmt.Errorf("handlerFollowing: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("(No feeds followed)")
		return nil
	}

	for _, follow := range follows {
		fmt.Println(follow.Feed.Name)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("handlerUnfollow: please specify a feed URL to unfollow")
	}

	deleteArgs := database.DeleteFeedFollowParams{UserID: user.ID, Url: cmd.args[0]}
	err := s.db.DeleteFeedFollow(context.Background(), deleteArgs)
	if err != nil {
		return fmt.Errorf("handlerUnfollow: %w", err)
	}

	fmt.Println("Successfully unfollowed feed")
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.args) > 0 {
		result, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			fmt.Println(fmt.Errorf("invalid limit: %w, defaulting to 2", err))
			limit = 2
		} else {
			limit = int32(result)
		}
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("handlerBrowse: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("* %s (%s)\n", post.Title, post.Url)
	}

	return nil
}
