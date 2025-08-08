# Gator RSS Feed Aggregator

A command-line RSS feed aggregator built in Go that allows users to manage and follow RSS feeds, with automatic feed scraping and post browsing capabilities.

## Features

- **User Management**: Register, login, and manage multiple users. Users are not password protected, knowing the username is enough to act as that user.
- **Feed Management**: Add, list, and manage RSS feeds
- **Feed Following**: Follow/unfollow feeds on a per-user basis
- **Automatic Aggregation**: Continuously scrape RSS feeds at configurable intervals
- **Post Browsing**: Browse posts from your followed feeds
- **PostgreSQL Storage**: Persistent data storage with SQLC-generated queries

## Prerequisites

- PostgreSQL database
- Git (optional, for cloning)

## Quick Start

### Option 1: Download Pre-built Binary

1. Download the latest release from the [releases page](https://github.com/tbirddv/gator/releases)
2. Make it executable:
```bash
chmod +x gator
```
3. Optionally, move to your PATH:
```bash
sudo mv gator /usr/local/bin/
```

### Option 2: Clone and Use

1. Clone the repository:
```bash
git clone https://github.com/tbirddv/gator.git
cd gator
```

2. The pre-built executable is already included in the repository.

### Setup

1. Set up your PostgreSQL database and note the connection string.

2. Create a configuration file `~/.gatorconfig.json`:
```json
{
  "db_url": "postgres://username:password@localhost/gator?sslmode=disable",
  "current_user_name": ""
}
```

3. Ensure your PostgreSQL database exists and is accessible.

4. You're ready to use Gator!

## Building from Source (Optional)

If you prefer to build from source:

**Prerequisites:**
- Go 1.19 or higher

**Steps:**
```bash
git clone https://github.com/tbirddv/gator.git
cd gator
go mod download
go build -o gator
```

## Usage

### User Management

**Register a new user:**
```bash
./gator register <username>
```

**Login as an existing user:**
```bash
./gator login <username>
```

**List all users:**
```bash
./gator users
```

**Reset all users:**
```bash
./gator reset
```

### Feed Management

**Add a new RSS feed:**
```bash
./gator addfeed <feed_name> <feed_url>
```
*Note: This automatically follows the feed for the current user*

**List all feeds:**
```bash
./gator feeds
```

**Follow an existing feed:**
```bash
./gator follow <feed_url>
```

**Unfollow a feed:**
```bash
./gator unfollow <feed_url>
```

**List feeds you're following:**
```bash
./gator following
```

### Feed Aggregation

**Start the RSS aggregator:**
```bash
./gator agg <time_between_requests>
```

Examples:
- `./gator agg 1m` - Scrape feeds every minute
- `./gator agg 30s` - Scrape feeds every 30 seconds
- `./gator agg 1h` - Scrape feeds every hour

*Note: This runs continuously until stopped with Ctrl+C*

### Browse Posts

**Browse recent posts from your followed feeds:**
```bash
./gator browse [limit]
```

Examples:
- `./gator browse` - Browse 2 posts (default)
- `./gator browse 10` - Browse 10 most recent posts

### Help

**Get help for all commands:**
```bash
./gator help
```

**Get help for a specific command:**
```bash
./gator help <command_name>
```

## Project Structure

```
.
├── gator                   # Pre-built executable
├── main.go                 # Entry point and command routing
├── commands.go            # Command definitions and initialization
├── handlers.go            # Command handler implementations
├── utils.go               # Utility functions for feeds and users
├── internal/
│   ├── config/
│   │   └── config.go      # Configuration management
│   ├── database/
│   │   ├── db.go          # Database connection
│   │   ├── models.go      # Generated database models
│   │   └── *.sql.go       # Generated SQLC queries
│   └── rssfeed/
│       └── rssfeed.go     # RSS feed fetching and parsing
└── sql/
    ├── queries/           # SQL queries for SQLC
    └── schema/           # Database migration files
```

## Database Schema

The application uses the following main tables:
- `users` - User accounts
- `feeds` - RSS feed definitions
- `feed_follows` - Many-to-many relationship between users and feeds
- `posts` - Individual RSS feed items/posts

## Technologies Used

- **Go** - Programming language
- **PostgreSQL** - Database
- **SQLC** - SQL query code generation
- **github.com/lib/pq** - PostgreSQL driver
- **github.com/google/uuid** - UUID generation

## Example Workflow

1. Register and login:
```bash
./gator register alice
./gator login alice
```

2. Add and follow some feeds:
```bash
./gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
./gator addfeed "Go Blog" "https://go.dev/blog/feed.atom"
```

3. Start aggregating feeds:
```bash
./gator agg 2m
```

4. In another terminal, browse posts:
```bash
./gator browse 5
```

## Configuration

The application uses a JSON configuration file stored at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://user:password@localhost/dbname?sslmode=disable",
  "current_user_name": "current_logged_in_user"
}
```

## Troubleshooting

- **Permission denied**: Make sure the executable has proper permissions (`chmod +x gator`)
- **Database connection issues**: Verify your PostgreSQL server is running and the connection string is correct
- **Command not found**: Ensure you're running `./gator` from the correct directory or have added it to your PATH

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is part of the Boot.dev backend development course.