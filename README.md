# gator

A command-line RSS feed aggregator that helps you follow and manage your favorite RSS feeds.

## Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.23 or later
- PostgreSQL database

## Installation

Install the gator CLI using Go:

```bash
go install github.com/jasonwashburn/gator@latest
```

## Configuration

1. Create a configuration file at `~/.gatorconfig.json` with your PostgreSQL database URL:

```json
{
    "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
    "current_user_name": ""
}
```

Replace the database URL with your actual PostgreSQL connection string.

## Usage

### User Management

1. Register a new user:
```bash
gator register your_username
```

2. Login as a user:
```bash
gator login your_username
```

3. List all users:
```bash
gator users
```

### Feed Management

1. Add a new feed (example using the Go Blog):
```bash
gator addfeed "Go Blog" "https://go.dev/blog/feed.atom"
```

2. List all available feeds:
```bash
gator feeds
```

3. Follow a feed (must be added first):
```bash
gator follow "https://go.dev/blog/feed.atom"
```

4. List feeds you're following:
```bash
gator following
```

5. Unfollow a feed:
```bash
gator unfollow "https://go.dev/blog/feed.atom"
```

### Reading Posts

1. Browse your feed posts:
```bash
gator browse [limit]
```
The optional `limit` parameter specifies how many posts to display (default is 2).

### Feed Aggregation

Start the feed aggregator to fetch new posts:
```bash
gator agg "5m"
```
This will fetch new posts every 5 minutes. You can adjust the interval as needed (e.g., "1h" for hourly updates).

## Notes

- You must be logged in to use feed management commands (addfeed, follow, unfollow, browse)
- The feed must be added to the system before you can follow it
- The aggregator will fetch posts from all feeds in the system, not just the ones you follow