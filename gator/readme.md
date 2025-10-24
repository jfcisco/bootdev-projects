# Gator

Gator is a CLI program that aggregates RSS feeds. It allows users to perform the following:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

## Prerequisites

- Go 1.25 or higher
- PostgreSQL
- Goose CLI tool for database migrations (https://github.com/pressly/goose/)
- SQLC (https://docs.sqlc.dev/en/latest/overview/install.html)
