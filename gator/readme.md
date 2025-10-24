# Gator

Gator is a CLI program that aggregates RSS feeds. It allows users to perform the following:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

## Prerequisites

### Required Software

Before running Gator, you'll need to have the following installed on your system:

- **Go 1.25 or higher** - Required to build and run the application
  - Download from [https://golang.org/dl/](https://golang.org/dl/)
  - Follow the installation instructions for your operating system
- **PostgreSQL** - Database server for storing RSS feeds and posts
  - Download from [https://www.postgresql.org/download/](https://www.postgresql.org/download/)
  - Make sure the PostgreSQL service is running and you have connection credentials

### Development Tools (Optional)

These tools are only needed if you plan to modify the database schema or queries:

- Goose CLI tool for database migrations (https://github.com/pressly/goose/)
- SQLC (https://docs.sqlc.dev/en/latest/overview/install.html)

## Installation

Install the Gator CLI using Go's built-in package manager:

```bash
go install github.com/jfcisco/bootdev-projects/gator@latest
```

This will install the `gator` binary to your `$GOPATH/bin` directory. Make sure this directory is in your system's `PATH` so you can run the `gator` command from anywhere.

> **Note**: Since this project is located in a subdirectory of the repository, the install path includes the `/gator` subdirectory path.

## Configuration

Before using Gator, you need to create a configuration file in your home directory:

1. Create a file named `.gatorconfig.json` in your home directory
2. Add your PostgreSQL database connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username`, `password`, and database details with your actual PostgreSQL credentials. The `current_user_name` field will be automatically populated when you register or login.

## Getting Started

### Basic Commands

Once installed and configured, you can use the following commands:

#### User Management
- `gator register <username>` - Create a new user account
- `gator login <username>` - Login as an existing user
- `gator users` - List all registered users

#### Feed Management
- `gator addfeed <name> <url>` - Add a new RSS feed to the system
- `gator feeds` - List all available RSS feeds
- `gator follow <url>` - Follow an RSS feed
- `gator following` - List feeds you're currently following
- `gator unfollow <url>` - Unfollow an RSS feed

#### Content Browsing
- `gator agg` - Aggregate (fetch) new posts from all feeds
- `gator browse [limit]` - Browse recent posts from your followed feeds

### Example Usage

1. **Register a new user:**
   ```bash
   gator register myusername
   ```

2. **Add a new RSS feed:**
   ```bash
   gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
   ```

3. **Follow the feed:**
   ```bash
   gator follow "https://techcrunch.com/feed/"
   ```

4. **Fetch new posts:**
   ```bash
   gator agg
   ```

5. **Browse your posts:**
   ```bash
   gator browse 10
   ```
